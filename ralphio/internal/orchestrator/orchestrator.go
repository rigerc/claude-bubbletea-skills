package orchestrator

import (
	"context"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"ralphio/internal/adapter"
	"ralphio/internal/plan"
	"ralphio/internal/state"
	"ralphio/internal/validator"
)

const iterationDelay = 2000 * time.Millisecond

// Orchestrator drives the Ralph loop in a background goroutine. It communicates
// with the BubbleTea TUI exclusively via typed messages over channels: outbound
// messages go on msgCh, inbound commands arrive on cmdCh.
//
// All mutable state (paused, skipCurrent, retryCurrent, currentAdapter) is
// accessed only from the goroutine that calls Run — no mutex is required.
type Orchestrator struct {
	projectDir     string
	msgCh          chan<- tea.Msg
	cmdCh          <-chan any
	currentAdapter adapter.Adapter

	// iteration-scoped flags; reset at the start of each iteration.
	paused       bool
	skipCurrent  bool
	retryCurrent bool
}

// New returns an Orchestrator ready to be started with Run.
//
// msgCh should be a buffered channel (capacity ≥ 64) so that sends from the
// orchestrator goroutine do not block when the TUI is briefly busy.
// cmdCh carries commands from the TUI (RetryCmd, SkipCmd, etc.).
func New(
	projectDir string,
	initialAdapter adapter.Adapter,
	msgCh chan<- tea.Msg,
	cmdCh <-chan any,
) *Orchestrator {
	return &Orchestrator{
		projectDir:     projectDir,
		msgCh:          msgCh,
		cmdCh:          cmdCh,
		currentAdapter: initialAdapter,
	}
}

// Run starts the Ralph loop. It blocks until ctx is cancelled or all tasks are
// complete. Run is intended to be called in a dedicated goroutine.
func (o *Orchestrator) Run(ctx context.Context) {
	// Load initial state and tasks so the TUI has something to display
	// immediately on startup.
	st, err := state.Load(o.projectDir)
	if err != nil {
		o.send(LoopErrorMsg{Err: fmt.Errorf("loading state: %w", err)})
		return
	}
	st.LoopStatus = state.StatusRunning
	st.ActiveAdapter = string(o.currentAdapter.Name())

	tasks, err := o.loadTasks()
	if err != nil {
		o.send(LoopErrorMsg{Err: err})
		return
	}

	o.send(o.snapshot(tasks, st))

	planMgr := plan.NewManager(o.projectDir)

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

		// Reload tasks from disk so any external edits are picked up.
		tasks, err = planMgr.LoadTasks()
		if err != nil {
			o.send(LoopErrorMsg{Err: fmt.Errorf("loading tasks: %w", err)})
			return
		}

		next := plan.NextTask(tasks)
		if next == nil {
			o.send(LoopDoneMsg{})
			return
		}

		if err := o.runIteration(ctx, st, planMgr, tasks, next); err != nil {
			// A non-nil error from runIteration means the context was cancelled
			// or a fatal I/O error occurred — either way, stop the loop.
			o.shutdown(tasks, st)
			return
		}

		// Reload tasks after the iteration (runIteration may have saved them).
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

// runIteration executes a single task: mark in_progress, run the agent,
// validate, update the task status, and persist everything. Returns a non-nil
// error only on unrecoverable failures (context cancellation or disk I/O).
func (o *Orchestrator) runIteration(
	ctx context.Context,
	st *state.State,
	planMgr *plan.Manager,
	tasks []plan.Task,
	task *plan.Task,
) error {
	st.CurrentIteration++
	st.CurrentTaskID = task.ID
	st.ActiveAdapter = string(o.currentAdapter.Name())

	// Mark task as in-progress.
	plan.UpdateTask(tasks, task.ID, plan.StatusInProgress, task.RetryCount)
	if err := planMgr.SaveTasks(tasks); err != nil {
		o.send(LoopErrorMsg{Err: fmt.Errorf("saving tasks (in_progress): %w", err)})
		return err
	}

	o.send(IterationStartMsg{
		Iteration: st.CurrentIteration,
		TaskID:    task.ID,
		TaskTitle: task.Title,
	})
	o.send(o.snapshot(tasks, st))

	start := time.Now()

	prompt := fmt.Sprintf(
		"Task ID: %s\nTitle: %s\n\nDescription:\n%s\n\nValidation command: %s",
		task.ID,
		task.Title,
		task.Description,
		task.ValidationCommand,
	)

	// Execute agent with streaming output forwarded to the TUI.
	execErr := o.currentAdapter.Execute(ctx, prompt, func(text string) {
		o.send(AgentOutputMsg{Text: text})
	})

	// Handle context cancellation from within Execute.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Run validation if a command is configured and the agent did not error.
	var validResult *validator.Result
	passed := execErr == nil

	if passed && task.ValidationCommand != "" {
		r := validator.Run(ctx, task.ValidationCommand)
		validResult = &r
		passed = r.Passed

		// Check again for context cancellation after the validation run.
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	duration := time.Since(start)

	// Update task status based on outcome.
	maxRetries := task.MaxRetries
	if maxRetries == 0 {
		maxRetries = 3 // sensible default when not specified
	}

	switch {
	case passed:
		plan.UpdateTask(tasks, task.ID, plan.StatusCompleted, 0)
	case task.RetryCount < maxRetries:
		// Increment retry count and reset to pending so the next loop
		// iteration picks it up again.
		plan.UpdateTask(tasks, task.ID, plan.StatusPending, task.RetryCount+1)
	default:
		// Retry budget exhausted.
		plan.UpdateTask(tasks, task.ID, plan.StatusSkipped, task.RetryCount)
	}

	o.send(IterationCompleteMsg{
		Iteration:        st.CurrentIteration,
		TaskID:           task.ID,
		ValidationResult: validResult,
		Passed:           passed,
		Duration:         duration,
	})

	if err := planMgr.SaveTasks(tasks); err != nil {
		o.send(LoopErrorMsg{Err: fmt.Errorf("saving tasks (post-iteration): %w", err)})
		return err
	}

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
	case RetryCmd:
		o.retryCurrent = true
	case SkipCmd:
		o.skipCurrent = true
	case TogglePauseCmd:
		o.paused = !o.paused
		if o.paused {
			o.send(LoopPausedMsg{})
		} else {
			o.send(LoopResumedMsg{})
		}
	case ChangeAdapterCmd:
		o.currentAdapter = adapter.NewAdapter(c.Agent, c.Model)
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
		ActiveAgent:    o.currentAdapter.Name(),
		ActiveModel:    st.ActiveModel,
	}
}

// loadTasks is a convenience wrapper used on first startup before the plan
// manager is constructed.
func (o *Orchestrator) loadTasks() ([]plan.Task, error) {
	mgr := plan.NewManager(o.projectDir)
	tasks, err := mgr.LoadTasks()
	if err != nil {
		return nil, fmt.Errorf("loading tasks: %w", err)
	}
	return tasks, nil
}
