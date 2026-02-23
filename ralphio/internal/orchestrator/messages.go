// Package orchestrator runs the Ralph loop in a background goroutine and
// communicates with the BubbleTea TUI via typed messages over channels.
package orchestrator

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"ralphio/internal/adapter"
	"ralphio/internal/plan"
)

// --- Messages sent from orchestrator → TUI ---

// LoopStateMsg is a full snapshot of the orchestrator state, sent after every
// significant change. The TUI should replace its cached state on receipt.
type LoopStateMsg struct {
	Iteration      int
	TotalTasks     int
	CompletedTasks int
	CurrentTask    *plan.Task  // nil if idle
	Tasks          []plan.Task // full task list (copy)
	Status         string      // "running", "paused", "stopped", "error"
	LoopMode       plan.LoopMode
	ActiveAgent    adapter.AgentType
	ActiveModel    string
}

// AgentOutputMsg carries a single chunk of streaming text from the running
// agent. The TUI appends each chunk to its output viewport as it arrives.
type AgentOutputMsg struct {
	Text string
}

// IterationStartMsg fires when the orchestrator begins an agent invocation.
// In the agent-driven model, TaskID and TaskTitle are empty — the agent
// selects the task autonomously.
type IterationStartMsg struct {
	Iteration int
	TaskID    string
	TaskTitle string
}

// IterationCompleteMsg fires when an iteration finishes (pass or fail).
// Passed is determined by output signal detection, not an explicit validation
// command run by the orchestrator.
type IterationCompleteMsg struct {
	Iteration int
	TaskID    string
	Passed    bool
	Duration  time.Duration
}

// LoopDoneMsg fires when all tasks are complete or the loop exits cleanly with
// no pending work remaining.
type LoopDoneMsg struct{}

// LoopErrorMsg fires when the orchestrator encounters a fatal error that
// prevents the loop from continuing.
type LoopErrorMsg struct {
	Err error
}

// LoopPausedMsg fires when the loop transitions to the paused state.
type LoopPausedMsg struct{}

// LoopResumedMsg fires when the loop resumes after being paused.
type LoopResumedMsg struct{}

// LoopModeChangedMsg fires when the orchestrator automatically switches modes.
type LoopModeChangedMsg struct {
	Mode   plan.LoopMode
	Reason string
}

// --- Commands sent from TUI → orchestrator ---

// RetryCmd is kept for backward compatibility. The orchestrator no longer
// handles per-task retries; the agent manages its own retry strategy.
type RetryCmd struct{}

// SkipCmd is kept for backward compatibility. The orchestrator no longer
// handles per-task skips; use ChangeModeCmd or StopCmd instead.
type SkipCmd struct{}

// TogglePauseCmd pauses the loop if it is running, or resumes it if paused.
type TogglePauseCmd struct{}

// ChangeAdapterCmd asks the orchestrator to switch to a different agent and
// optional model for subsequent iterations.
type ChangeAdapterCmd struct {
	Agent adapter.AgentType
	Model string
}

// ChangeModeCmd asks the orchestrator to switch between planning and building
// modes. The change takes effect on the next iteration.
type ChangeModeCmd struct {
	Mode plan.LoopMode
}

// StopCmd asks the orchestrator to stop the loop cleanly. The caller is
// responsible for cancelling the context that was passed to Run.
type StopCmd struct{}

// Ensure the message types satisfy tea.Msg at compile time so that callers
// get a clear error if the interface changes.
var (
	_ tea.Msg = LoopStateMsg{}
	_ tea.Msg = AgentOutputMsg{}
	_ tea.Msg = IterationStartMsg{}
	_ tea.Msg = IterationCompleteMsg{}
	_ tea.Msg = LoopDoneMsg{}
	_ tea.Msg = LoopErrorMsg{}
	_ tea.Msg = LoopPausedMsg{}
	_ tea.Msg = LoopResumedMsg{}
	_ tea.Msg = LoopModeChangedMsg{}
)
