package adapter

import "context"

// CursorAdapter executes prompts via the Cursor agent CLI.
type CursorAdapter struct{}

// NewCursorAdapter returns a CursorAdapter.
func NewCursorAdapter() *CursorAdapter {
	return &CursorAdapter{}
}

func (a *CursorAdapter) Name() AgentType { return AgentCursor }

func (a *CursorAdapter) SupportsModelSelection() bool { return false }

func (a *CursorAdapter) Execute(ctx context.Context, prompt string, onOutput func(string)) error {
	return runProcess(ctx, AgentCommands[AgentCursor], prompt, "", onOutput)
}
