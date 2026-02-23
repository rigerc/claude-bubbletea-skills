package adapter

import "context"

// KiloAdapter executes prompts via the kilo CLI and supports model selection
// via the --model flag.
type KiloAdapter struct {
	model string
}

// NewKiloAdapter returns a KiloAdapter configured for the given model.
// Pass an empty string to use the agent's default model.
func NewKiloAdapter(model string) *KiloAdapter {
	return &KiloAdapter{model: model}
}

func (a *KiloAdapter) Name() AgentType { return AgentKilo }

func (a *KiloAdapter) SupportsModelSelection() bool { return true }

func (a *KiloAdapter) Execute(ctx context.Context, prompt string, onOutput func(string)) error {
	return runProcess(ctx, AgentCommands[AgentKilo], prompt, a.model, onOutput)
}

// FetchModels returns the list of models available through kilo.
func (a *KiloAdapter) FetchModels(ctx context.Context) ([]string, error) {
	return FetchModels(ctx, AgentKilo)
}
