package orchestrator

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"ralphio/internal/adapter"
	"ralphio/internal/plan"
	"ralphio/internal/prompt"
	"ralphio/internal/state"
)

const iterationDelay = 2000 * time.Millisecond

// Orchestrator drives the Ralph loop in a background goroutine. It communicates
// with the BubbleTea TUI exclusively via typed messages over channels: outbound
// messages go on msgCh, inbound commands arrive on cmdCh.
//
// In the agent-driven model, the orchestrator no longer selects tasks or runs
// validation. Instead, it builds a mode-based prompt from PROMPT_build.md or
// PROMPT_plan.md and passes it to the agent. The agent reads tasks.json,
// selects the highest-priority pending task, implements it, and updates
// tasks.json directly.
type Orchestrator struct {
	projectDir    string
	msgCh         chan<- tea.Msg
	cmdCh         <-chan any
	currentAdapter adapter.Adapter
	promptBuilder *prompt.Builder

	paused bool
}

// New returns an Orchestrator ready to be started with Run.
//
// msgCh should be a buffered channel (capacity >= 64) so that sends from the
// orchestrator goroutine do not block when the TUI is briefly busy.
// cmdCh carries commands from the TUI (TogglePauseCmd, ChangeAdapterCmd, etc.).
func New(
	projectDir string,
	initialAdapter adapter.Adapter,
	msgCh chan<- tea.Msg,
	cmdCh <-chan any,
) *Orchestrator {
	return &Orchestrator{
		projectDir:    projectDir,
		msgCh:         msgCh,
		cmdCh:         cmdCh,
		currentAdapter: initialAdapter,
		promptBuilder: prompt.New(projectDir),
	}
}

// Run starts the Ralph loop. It blocks until ctx is cancelled or all tasks are
// complete. Run is intended to be called in a dedicated goroutine.
func (o *Orchestrator) Run(ctx context.Context) {
	st, err := state.Load(o.projectDir)
	if err != nil {
		o.send(LoopErrorMsg{Err: fmt.Errorf("loading state: %w", err)})
		return
	}
	st.LoopStatus = state.StatusRunning
	st.ActiveAdapter = string(o.currentAdapter.Name())
	if st.LoopMode == "" {
		st.LoopMode = state.ModeBuilding
	}

	planMgr := plan.NewManager(o.projectDir)

	tasks, err := planMgr.LoadTasks()
	if err != nil {
		o.send(LoopErrorMsg{Err: fmt.Errorf("loading tasks: %w", err)})
		return
	}

	// Create tasks.json with schema template if it doesn't exist.
	if len(tasks) == 0 {
		if err := planMgr.CreateInitialTasks(); err != nil {
			o.send(LoopErrorMsg{Err: fmt.Errorf("creating tasks.json: %w", err)})
			return
		}
		// Reload tasks after creation to get the template.
		tasks, err = planMgr.LoadTasks()
		if err != nil {
			o.send(LoopErrorMsg{Err: fmt.Errorf("loading tasks.json after creation: %w", err)})
			return
		}
		o.send(LoopModeChangedMsg{Mode: state.ModePlanning, Reason: "Created tasks.json with schema — planning mode"})
		st.LoopMode = state.ModePlanning
	}

	o.send(o.snapshot(tasks, st))

	for {
		// Drain any pending TUI commands before deciding what to do next.
		o.drainCommands()

		// Honour context cancellation between iterations.
		select {
		case <-ctx.Done():
			o.shutdown(tasks, st)
			return
		default:
		}

		// When paused, spin in a tight sleep loop while still draining
		// commands so the TUI stays responsive.
		if o.paused {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Reload tasks so external edits and agent changes are picked up.
		tasks, err = planMgr.LoadTasks()
		if err != nil {
			o.send(LoopErrorMsg{Err: fmt.Errorf("loading tasks: %w", err)})
			return
		}

		// Check if all tasks are complete — agent may have finished.
		if allCompleted(tasks) {
			o.send(LoopDoneMsg{})
			return
		}

		if err := o.runIteration(ctx, st, tasks); err != nil {
			// A non-nil error from runIteration means the context was cancelled
			// or a fatal I/O error occurred — either way, stop the loop.
			o.shutdown(tasks, st)
			return
		}

		// Reload after agent may have mutated tasks.json.
		tasks, err = planMgr.LoadTasks()
		if err != nil {
			o.send(LoopErrorMsg{Err: fmt.Errorf("reloading tasks after iteration: %w", err)})
			return
		}

		o.send(o.snapshot(tasks, st))

		// Inter-iteration delay; also check for cancellation.
		select {
		case <-ctx.Done():
			o.shutdown(tasks, st)
			return
		case <-time.After(iterationDelay):
		}
	}
}

// runIteration executes a single agent invocation: builds the mode-based
// prompt, runs the agent, and streams output to the TUI. The agent is
// responsible for selecting tasks, updating tasks.json, and running validation.
func (o *Orchestrator) runIteration(
	ctx context.Context,
	st *state.State,
	tasks []plan.Task,
) error {
	st.CurrentIteration++
	st.ActiveAdapter = string(o.currentAdapter.Name())

	// Build prompt from the mode-based prompt file; fall back to an inline
	// prompt if the file is missing so the loop can still make progress.
	mode := plan.LoopMode(st.LoopMode)
	agentPrompt, err := o.promptBuilder.Build(mode)
	if err != nil {
		agentPrompt = fallbackPrompt(mode)
	}

	o.send(IterationStartMsg{
		Iteration: st.CurrentIteration,
		// TaskID and TaskTitle are intentionally empty: the agent selects
		// the task autonomously by reading tasks.json.
	})
	o.send(o.snapshot(tasks, st))

	start := time.Now()
	var outputBuf strings.Builder

	execErr := o.currentAdapter.Execute(ctx, agentPrompt, func(text string) {
		outputBuf.WriteString(text)
		o.send(AgentOutputMsg{Text: text})
	})

	if ctx.Err() != nil {
		return ctx.Err()
	}

	duration := time.Since(start)
	passed := execErr == nil && detectValidation(outputBuf.String())

	// Update the current task ID from whatever the agent marked in_progress.
	if inProgress := findInProgress(tasks); inProgress != nil {
		st.CurrentTaskID = inProgress.ID
	}

	o.send(IterationCompleteMsg{
		Iteration: st.CurrentIteration,
		TaskID:    st.CurrentTaskID,
		Passed:    passed,
		Duration:  duration,
	})

	if err := state.Save(o.projectDir, st); err != nil {
		o.send(LoopErrorMsg{Err: fmt.Errorf("saving state: %w", err)})
		return err
	}

	return nil
}

// shutdown transitions state to stopped, sends a final snapshot, and returns.
func (o *Orchestrator) shutdown(tasks []plan.Task, st *state.State) {
	st.LoopStatus = state.StatusStopped
	// Best-effort save; ignore errors on shutdown path.
	_ = state.Save(o.projectDir, st)
	o.send(o.snapshot(tasks, st))
}

// drainCommands reads all currently queued commands from cmdCh without
// blocking. It is called at the top of every loop iteration.
func (o *Orchestrator) drainCommands() {
	for {
		select {
		case cmd := <-o.cmdCh:
			o.handleCommand(cmd)
		default:
			return
		}
	}
}

// handleCommand applies a single TUI command to the orchestrator's mutable
// state. It must only be called from the Run goroutine.
func (o *Orchestrator) handleCommand(cmd any) {
	switch c := cmd.(type) {
	case TogglePauseCmd:
		o.paused = !o.paused
		if o.paused {
			o.send(LoopPausedMsg{})
		} else {
			o.send(LoopResumedMsg{})
		}
	case ChangeAdapterCmd:
		o.currentAdapter = adapter.NewAdapter(c.Agent, c.Model)
	case ChangeModeCmd:
		// Mode is persisted in state; update will be reflected at the next
		// iteration when runIteration reads st.LoopMode.
		o.send(LoopPausedMsg{}) // brief pause while mode switches
	case StopCmd:
		// StopCmd is a signal for the caller to cancel the context. The
		// orchestrator itself does not exit here; it will exit on the next
		// ctx.Done() check in Run.
	}
}

// send delivers msg to the TUI channel without blocking. If the channel is
// full the message is dropped; the TUI will catch up via the next snapshot.
func (o *Orchestrator) send(msg tea.Msg) {
	select {
	case o.msgCh <- msg:
	default:
	}
}

// snapshot builds a LoopStateMsg from the current in-memory task list and
// state. Tasks are copied so the TUI holds an independent slice.
func (o *Orchestrator) snapshot(tasks []plan.Task, st *state.State) LoopStateMsg {
	tasksCopy := make([]plan.Task, len(tasks))
	copy(tasksCopy, tasks)

	total := len(tasks)
	completed := 0
	for i := range tasks {
		if tasks[i].Status == plan.StatusCompleted {
			completed++
		}
	}

	var currentTask *plan.Task
	if st.CurrentTaskID != "" {
		for i := range tasks {
			if tasks[i].ID == st.CurrentTaskID &&
				tasks[i].Status == plan.StatusInProgress {
				t := tasks[i]
				currentTask = &t
				break
			}
		}
	}

	return LoopStateMsg{
		Iteration:      st.CurrentIteration,
		TotalTasks:     total,
		CompletedTasks: completed,
		CurrentTask:    currentTask,
		Tasks:          tasksCopy,
		Status:         st.LoopStatus,
		LoopMode:       plan.LoopMode(st.LoopMode),
		ActiveAgent:    o.currentAdapter.Name(),
		ActiveModel:    st.ActiveModel,
	}
}

// allCompleted returns true if every task in the list is completed.
// An empty task list is never considered done.
func allCompleted(tasks []plan.Task) bool {
	if len(tasks) == 0 {
		return false
	}
	for i := range tasks {
		if tasks[i].Status != plan.StatusCompleted {
			return false
		}
	}
	return true
}

// findInProgress returns the first in_progress task or nil.
func findInProgress(tasks []plan.Task) *plan.Task {
	for i := range tasks {
		if tasks[i].Status == plan.StatusInProgress {
			return &tasks[i]
		}
	}
	return nil
}

// detectValidation scans agent output for pass/fail signals. Best-effort:
// the agent is authoritative; this is a display hint for the TUI only.
func detectValidation(output string) bool {
	lower := strings.ToLower(output)
	successPhrases := []string{
		"tests passed",
		"all tests pass",
		"build successful",
		"✓",
		"ok  ",
	}
	for _, phrase := range successPhrases {
		if strings.Contains(lower, phrase) {
			return true
		}
	}
	return false
}

// fallbackPrompt returns a minimal inline prompt when the prompt file is
// missing, so the loop can still make progress without PROMPT_*.md files.
func fallbackPrompt(mode plan.LoopMode) string {
	if mode == plan.ModePlanning {
		return "Study tasks.json and source code. Create or update tasks.json with prioritized tasks. Plan only, do not implement."
	}
	return "Study tasks.json. Select the highest-priority pending task, implement it, run tests, commit."
}
