package adapter

import "context"

// CodexAdapter executes prompts via the OpenAI Codex CLI.
type CodexAdapter struct{}

// NewCodexAdapter returns a CodexAdapter.
func NewCodexAdapter() *CodexAdapter {
	return &CodexAdapter{}
}

func (a *CodexAdapter) Name() AgentType { return AgentCodex }

func (a *CodexAdapter) SupportsModelSelection() bool { return false }

func (a *CodexAdapter) Execute(ctx context.Context, prompt string, onOutput func(string)) error {
	return runProcess(ctx, AgentCommands[AgentCodex], prompt, "", onOutput)
}
