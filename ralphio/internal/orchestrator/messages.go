// Package orchestrator runs the Ralph loop in a background goroutine and
// communicates with the BubbleTea TUI via typed messages over channels.
package orchestrator

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"ralphio/internal/adapter"
	"ralphio/internal/plan"
	"ralphio/internal/validator"
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
	ActiveAgent    adapter.AgentType
	ActiveModel    string
}

// AgentOutputMsg carries a single chunk of streaming text from the running
// agent. The TUI appends each chunk to its output viewport as it arrives.
type AgentOutputMsg struct {
	Text string
}

// IterationStartMsg fires when the orchestrator begins processing a task.
type IterationStartMsg struct {
	Iteration int
	TaskID    string
	TaskTitle string
}

// IterationCompleteMsg fires when an iteration finishes (pass or fail).
type IterationCompleteMsg struct {
	Iteration        int
	TaskID           string
	ValidationResult *validator.Result // nil when no validation command was set
	Passed           bool
	Duration         time.Duration
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

// --- Commands sent from TUI → orchestrator ---

// RetryCmd asks the orchestrator to retry the current task immediately,
// resetting its retry counter.
type RetryCmd struct{}

// SkipCmd asks the orchestrator to mark the current task as skipped and move
// on to the next pending task.
type SkipCmd struct{}

// TogglePauseCmd pauses the loop if it is running, or resumes it if paused.
type TogglePauseCmd struct{}

// ChangeAdapterCmd asks the orchestrator to switch to a different agent and
// optional model for subsequent iterations.
type ChangeAdapterCmd struct {
	Agent adapter.AgentType
	Model string
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
)
