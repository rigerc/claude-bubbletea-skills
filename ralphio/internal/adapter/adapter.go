// Package adapter defines the Adapter interface and agent type constants for
// the ralphio multi-client execution system.
package adapter

import "context"

// AgentType identifies a supported AI coding agent.
type AgentType string

const (
	AgentClaude   AgentType = "claude"
	AgentCursor   AgentType = "cursor"
	AgentCodex    AgentType = "codex"
	AgentOpencode AgentType = "opencode"
	AgentKilo     AgentType = "kilo"
	AgentPi       AgentType = "pi"
)

// ValidAgents is the ordered list of all supported agent types.
var ValidAgents = []AgentType{
	AgentCursor,
	AgentClaude,
	AgentCodex,
	AgentOpencode,
	AgentKilo,
	AgentPi,
}

// AgentsSupportingModel is the subset of agents that accept a --model flag
// and expose a model listing command.
var AgentsSupportingModel = []AgentType{AgentOpencode, AgentKilo, AgentPi}

// Adapter executes AI agent prompts and streams output back to the caller.
type Adapter interface {
	// Name returns the agent type identifier.
	Name() AgentType

	// Execute runs the agent with the given prompt. Each chunk of displayable
	// text is delivered to onOutput as it arrives (streaming). Execute blocks
	// until the agent process exits.
	Execute(ctx context.Context, prompt string, onOutput func(text string)) error

	// SupportsModelSelection reports whether this adapter accepts a model flag.
	SupportsModelSelection() bool
}

// ModelFetcher is implemented by adapters that can enumerate available models.
// Agents in AgentsSupportingModel implement both Adapter and ModelFetcher.
type ModelFetcher interface {
	FetchModels(ctx context.Context) ([]string, error)
}

// NewAdapter returns the concrete Adapter for the given agent type and model.
// model is only meaningful for agents that support model selection
// (AgentOpencode, AgentKilo, AgentPi); it is ignored for others.
// An unknown agent type falls back to the Claude adapter.
func NewAdapter(agent AgentType, model string) Adapter {
	switch agent {
	case AgentCursor:
		return NewCursorAdapter()
	case AgentCodex:
		return NewCodexAdapter()
	case AgentOpencode:
		return NewOpencodeAdapter(model)
	case AgentKilo:
		return NewKiloAdapter(model)
	case AgentPi:
		return NewPiAdapter(model)
	default:
		// AgentClaude and any unknown value fall back to Claude.
		return NewClaudeAdapter()
	}
}
