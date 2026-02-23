# Ralphio Usage Guide

Ralphio is an autonomous task execution engine that orchestrates AI coding agents to complete tasks from a `tasks.json` plan. It features real-time streaming output, validation commands, and multi-agent support.

## Quick Start

```bash
# Build
cd ralphio && go build .

# Create a tasks.json file (see Task File Format below)
./ralphio --project-dir ./my-project

# Run with a specific agent
./ralphio --agent cursor

# Run with model selection (opencode, kilo, pi only)
./ralphio --agent opencode --model anthropic/claude-sonnet-4-20250514
```

## CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--project-dir` | `.` | Directory containing `tasks.json` |
| `--agent` | `claude` | AI agent: `claude`, `cursor`, `codex`, `opencode`, `kilo`, `pi` |
| `--model` | `""` | Model for agents supporting selection (`opencode`, `kilo`, `pi`) |
| `--config` | `$HOME/.ralphio.json` | Path to config file |
| `--debug` | `false` | Enable trace logging |
| `--log-level` | `info` | Log verbosity: `trace`, `debug`, `info`, `warn`, `error`, `fatal` |

## Subcommands

| Command | Description |
|---------|-------------|
| `ralphio` | Start the TUI (default) |
| `ralphio run` | Explicit run command (same as default) |
| `ralphio version` | Print version info |
| `ralphio completion [bash\|zsh\|fish]` | Generate shell completions |

## Task File Format

Create `tasks.json` in your project directory:

```json
[
  {
    "id": "task-1",
    "title": "Implement user authentication",
    "description": "Add login/logout functionality with JWT tokens",
    "priority": 1,
    "status": "pending",
    "retryCount": 0,
    "maxRetries": 3,
    "validationCommand": "go test ./auth/..."
  },
  {
    "id": "task-2",
    "title": "Add rate limiting",
    "description": "Implement API rate limiting middleware",
    "priority": 2,
    "status": "pending",
    "retryCount": 0,
    "maxRetries": 3,
    "validationCommand": "go test ./middleware/..."
  }
]
```

### Task Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique task identifier |
| `title` | string | Yes | Short task title |
| `description` | string | Yes | Detailed task description (sent to agent) |
| `priority` | int | No | Lower number = higher priority. Default: 0 |
| `status` | string | No | `pending`, `in_progress`, `completed`, `failed`, `skipped` |
| `retryCount` | int | No | Current retry count. Default: 0 |
| `maxRetries` | int | No | Max retries before skipping. Default: 3 |
| `validationCommand` | string | No | Shell command to validate task completion |

## Supported Agents

| Agent | Key | Model Selection | Notes |
|-------|-----|-----------------|-------|
| Claude | `claude` | No | Anthropic's Claude CLI |
| Cursor | `cursor` | No | Cursor's agent mode |
| Codex | `codex` | No | OpenAI's Codex CLI |
| Opencode | `opencode` | Yes | `opencode run --format json` |
| Kilo | `kilo` | Yes | `kilo run --format json` |
| Pi | `pi` | Yes | `pi --mode json -p` |

## Model Selection

Only `opencode`, `kilo`, and `pi` support the `--model` flag:

```bash
# List available models
opencode models
kilo models
pi --list-models

# Run with specific model
./ralphio --agent opencode --model anthropic/claude-sonnet-4-20250514
./ralphio --agent kilo --model openai/gpt-4o
./ralphio --agent pi --model anthropic/claude-3-5-sonnet
```

## Key Bindings

### Dashboard

| Key | Action |
|-----|--------|
| `r` | Retry current task |
| `s` | Skip current task |
| `v` | View task detail |
| `c` | Change adapter/client |
| `h` | View iteration history |
| `p` | Pause/Resume loop |
| `q` / `ctrl+c` | Quit |
| `?` | Toggle help |

### Navigation

| Key | Action |
|-----|--------|
| `esc` | Go back / close dialog |

## State Persistence

Ralphio maintains state in `.ralph/state.json` within your project directory:

```json
{
  "currentIteration": 5,
  "currentTaskId": "task-3",
  "loopStatus": "running",
  "activeAdapter": "claude",
  "activeModel": "",
  "lastUpdated": "2024-01-15T10:30:00Z"
}
```

This enables crash recovery — restart ralphio after an interruption and it resumes from the last known state.

## Configuration File

Create `~/.ralphio.json` for persistent settings:

```json
{
  "logLevel": "info",
  "debug": false,
  "ui": {
    "altScreen": false,
    "mouseEnabled": true,
    "themeName": "default"
  },
  "ralph": {
    "projectDir": ".",
    "agent": "claude",
    "agentModel": "",
    "maxRetries": 3,
    "retryDelayMs": 5000,
    "agentTimeoutMs": 1800000,
    "iterationDelayMs": 2000,
    "iterations": 10,
    "validation": {
      "enabled": false,
      "commands": [],
      "failOnWarning": false
    }
  }
}
```

## Execution Flow

1. **Load State**: Read `.ralph/state.json` for recovery info
2. **Load Tasks**: Read `tasks.json` from project directory
3. **Select Task**: Pick highest-priority `pending` task
4. **Execute Agent**: Run agent with task description, stream output to TUI
5. **Validate**: If `validationCommand` is set, run it
6. **Update Status**: Mark task `completed`, `failed`, or `skipped`
7. **Persist**: Save updated `tasks.json` and `state.json`
8. **Repeat**: Continue until no pending tasks remain

## Examples

### Basic Workflow

```bash
# 1. Create project with tasks
mkdir my-feature && cd my-feature
cat > tasks.json << 'EOF'
[
  {
    "id": "1",
    "title": "Setup project structure",
    "description": "Create directory layout and initial files",
    "priority": 1,
    "status": "pending"
  },
  {
    "id": "2",
    "title": "Implement core logic",
    "description": "Write the main processing functions",
    "priority": 2,
    "status": "pending",
    "validationCommand": "go build ./..."
  }
]
EOF

# 2. Run ralphio
ralphio --agent claude
```

### With Validation

```bash
# Tasks with validation commands
cat > tasks.json << 'EOF'
[
  {
    "id": "test-task",
    "title": "Add unit tests",
    "description": "Write comprehensive unit tests for the API",
    "priority": 1,
    "status": "pending",
    "validationCommand": "go test ./... -v"
  }
]
EOF

ralphio --agent cursor
```

### Pause and Resume

```bash
# Start ralphio
ralphio

# Press 'p' to pause at any time
# Press 'p' again to resume

# Or kill and restart - state is preserved
# Ctrl+C to quit, then:
ralphio  # resumes from last checkpoint
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Agent not found | Ensure the agent CLI is installed and in PATH |
| Tasks not updating | Check file permissions on `tasks.json` and `.ralph/` |
| No output visible | Agent may not support `stream-json` output format |
| State corruption | Delete `.ralph/` directory to reset |
