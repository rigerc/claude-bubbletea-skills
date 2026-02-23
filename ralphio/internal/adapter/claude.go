package adapter

import "context"

// ClaudeAdapter executes prompts via the Claude CLI.
type ClaudeAdapter struct{}

// NewClaudeAdapter returns a ClaudeAdapter. Claude does not support model
// selection via a flag in the current command configuration.
func NewClaudeAdapter() *ClaudeAdapter {
	return &ClaudeAdapter{}
}

func (a *ClaudeAdapter) Name() AgentType { return AgentClaude }

func (a *ClaudeAdapter) SupportsModelSelection() bool { return false }

func (a *ClaudeAdapter) Execute(ctx context.Context, prompt string, onOutput func(string)) error {
	return runProcess(ctx, AgentCommands[AgentClaude], prompt, "", onOutput)
}
