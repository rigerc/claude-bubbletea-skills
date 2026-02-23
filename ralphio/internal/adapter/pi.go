package adapter

import "context"

// PiAdapter executes prompts via the pi CLI and supports model selection via
// the --model flag.
type PiAdapter struct {
	model string
}

// NewPiAdapter returns a PiAdapter configured for the given model.
// Pass an empty string to use the agent's default model.
func NewPiAdapter(model string) *PiAdapter {
	return &PiAdapter{model: model}
}

func (a *PiAdapter) Name() AgentType { return AgentPi }

func (a *PiAdapter) SupportsModelSelection() bool { return true }

func (a *PiAdapter) Execute(ctx context.Context, prompt string, onOutput func(string)) error {
	return runProcess(ctx, AgentCommands[AgentPi], prompt, a.model, onOutput)
}

// FetchModels returns the list of models available through pi.
func (a *PiAdapter) FetchModels(ctx context.Context) ([]string, error) {
	return FetchModels(ctx, AgentPi)
}
