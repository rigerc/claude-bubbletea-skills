package adapter

// CommandConfig holds the command template and any additional environment
// variables required to invoke a particular AI agent. Ported from
// ralph/src/lib/services/config/constants.ts.
type CommandConfig struct {
	// Command is the base argv slice (without the prompt argument).
	Command []string
	// Env contains extra environment variables added to the process environment.
	Env map[string]string
}

// AgentCommands maps each AgentType to its command configuration.
var AgentCommands = map[AgentType]CommandConfig{
	AgentCursor: {
		Command: []string{
			"agent", "-p", "--force",
			"--output-format", "stream-json",
			"--stream-partial-output",
		},
	},
	AgentClaude: {
		Command: []string{
			"claude", "-p",
			"--dangerously-skip-permissions",
			"--output-format", "stream-json",
			"--verbose",
		},
	},
	AgentCodex: {
		Command: []string{"codex", "exec", "--full-auto", "--json"},
	},
	AgentOpencode: {
		Command: []string{"opencode", "run", "--format", "json"},
		Env:     map[string]string{"OPENCODE_PERMISSION": `{"*":"allow"}`},
	},
	AgentKilo: {
		Command: []string{"kilo", "run", "--format", "json"},
		Env:     map[string]string{"KILO_PERMISSION": `{"*":"allow"}`},
	},
	AgentPi: {
		Command: []string{"pi", "--mode", "json", "-p"},
	},
}
