package adapter

import "context"

// OpencodeAdapter executes prompts via the opencode CLI and supports model
// selection via the --model flag.
type OpencodeAdapter struct {
	model string
}

// NewOpencodeAdapter returns an OpencodeAdapter configured for the given model.
// Pass an empty string to use the agent's default model.
func NewOpencodeAdapter(model string) *OpencodeAdapter {
	return &OpencodeAdapter{model: model}
}

func (a *OpencodeAdapter) Name() AgentType { return AgentOpencode }

func (a *OpencodeAdapter) SupportsModelSelection() bool { return true }

func (a *OpencodeAdapter) Execute(ctx context.Context, prompt string, onOutput func(string)) error {
	return runProcess(ctx, AgentCommands[AgentOpencode], prompt, a.model, onOutput)
}

// FetchModels returns the list of models available through opencode.
func (a *OpencodeAdapter) FetchModels(ctx context.Context) ([]string, error) {
	return FetchModels(ctx, AgentOpencode)
}
