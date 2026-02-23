# Ralphio Usage Guide

Ralphio is an autonomous task execution engine that orchestrates AI coding agents to implement tasks from a `tasks.json` plan. Agents autonomously select tasks, implement them, run validation, and update the plan. Ralphio provides a real-time TUI showing streaming agent output.

## Quick Start

```bash
# Build
cd ralphio && go build .

# Point at a project and run in building mode (default)
./ralphio --project-dir ./my-project

# Start in planning mode — agent generates/updates tasks.json
./ralphio --project-dir ./my-project --mode planning

# Run with a specific agent
./ralphio --agent cursor

# Run with model selection (opencode, kilo, pi only)
./ralphio --agent opencode --model anthropic/claude-sonnet-4-20250514

# Wipe state and start fresh
./ralphio --reset-state
```

## CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--project-dir` | `.` | Directory containing `tasks.json` and prompt files |
| `--agent` | `claude` | AI agent: `claude`, `cursor`, `codex`, `opencode`, `kilo`, `pi` |
| `--model` | `""` | Model for agents supporting selection (`opencode`, `kilo`, `pi`) |
| `--mode` | `building` | Initial loop mode: `planning` or `building` |
| `--reset-state` | `false` | Clear `.ralph/` state directory for a fresh start |
| `--config` | `$HOME/.ralphio.json` | Path to config file |
| `--debug` | `false` | Enable trace logging to `debug.log` |
| `--log-level` | `info` | Log verbosity: `trace`, `debug`, `info`, `warn`, `error`, `fatal` |

## Loop Modes

Ralphio operates in two modes, toggled with `m` or set via `--mode`.

### Planning mode

The agent reads `PRD.md` (or `specs/*` as fallback) and existing source code, then creates or updates `tasks.json` with a prioritized list of tasks. It does **not** implement anything.

Embedded prompt (override with `PROMPT_plan.md` for custom workflows).

### Building mode (default)

The agent reads `tasks.json`, selects the highest-priority pending task, implements it, runs tests, commits, and updates `tasks.json`. Ralphio loops until all tasks are completed or the loop is stopped.

## Project Structure

Ralphio uses the following files in `--project-dir`:

```
my-project/
├── AGENTS.md         # Build/test commands, project conventions (optional)
├── PRD.md            # Product Requirements Document (optional)
├── PROMPT_build.md   # Custom build mode prompt (optional, overrides embedded)
├── PROMPT_plan.md    # Custom plan mode prompt (optional, overrides embedded)
├── tasks.json        # Agent-maintained task list
└── specs/            # Requirement specifications (optional, fallback if no PRD.md)
    └── *.md
```

The TUI dashboard shows which of the key files (PRD.md, AGENTS.md, tasks.json) are present with a `✓`/`✗` indicator.

**Auto-creation:** If `tasks.json` doesn't exist on startup, ralphio creates one with a schema template example task and switches to planning mode so the agent can generate real tasks from `PRD.md`.

### Embedded Prompts

Ralphio includes embedded prompts for both modes, so `PROMPT_build.md` and `PROMPT_plan.md` are **optional**. If present, they override the embedded prompts.

**Embedded planning prompt:**
- Reads `PRD.md` (or `specs/*` as fallback)
- Generates or updates `tasks.json` from requirements
- Does NOT implement anything

**Embedded building prompt:**
- Studies `PRD.md` (if present) and `tasks.json`
- Implements tasks autonomously
- Runs tests, commits, and updates `tasks.json`

### Custom Prompts (Optional)

For project-specific workflows, create `PROMPT_build.md` or `PROMPT_plan.md` in your project directory to override the embedded prompts.

## Task File Format (`tasks.json`)

The agent reads and writes `tasks.json` directly. The orchestrator no longer selects or updates tasks — the agent manages the plan autonomously.

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

### Task fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique task identifier |
| `title` | string | Yes | Short task title |
| `description` | string | Yes | Detailed task description |
| `priority` | int | No | Lower number = higher priority. Default: 0 |
| `status` | string | No | `pending`, `in_progress`, `completed`, `failed` |
| `retryCount` | int | No | Current retry count. Default: 0 |
| `maxRetries` | int | No | Max retries before agent gives up. Default: 3 |
| `validationCommand` | string | No | Shell command the agent runs to verify completion |

## Supported Agents

| Agent | Key | Model Selection | Notes |
|-------|-----|-----------------|-------|
| Claude | `claude` | No | Anthropic's Claude CLI |
| Cursor | `cursor` | No | Cursor's agent mode |
| Codex | `codex` | No | OpenAI's Codex CLI |
| Opencode | `opencode` | Yes | `opencode run --format json` |
| Kilo | `kilo` | Yes | `kilo run --format json` |
| Pi | `pi` | Yes | `pi --mode json -p` |

## Key Bindings

### Dashboard

| Key | Action |
|-----|--------|
| `m` | Toggle mode (planning ↔ building) |
| `e` | Open `tasks.json` in `$EDITOR` |
| `R` | Switch to planning mode (regenerate plan) |
| `p` | Pause / resume loop |
| `v` | View current task detail |
| `c` | Change adapter / client |
| `h` | View iteration history |
| `q` / `ctrl+c` | Quit |
| `↑` / `pgup` | Scroll output up (disables auto-scroll) |
| `end` | Scroll to bottom (re-enables auto-scroll) |

### Navigation

| Key | Action |
|-----|--------|
| `esc` | Go back / close screen |

## TUI Layout

```
╔══════════════════════════════════════════════════════╗
║  ██████╗  █████╗ ██╗     ██████╗ ██╗  ██╗           ║
║  ██╔══██╗██╔══██╗██║     ██╔══██╗██║  ██║           ║
║  ██████╔╝███████║██║     ██████╔╝███████║           ║
║  ██╔══██╗██╔══██╗██║     ██╔══██╗██╔══██║           ║
║  ██║  ██║██║  ██║███████╗██║  ██║██║  ██║           ║
║  ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝           ║
║  Project: ./my-project  |  Mode: BUILDING  |  #14   ║
╠══════════════════════════════════════════════════════╣
║  PROMPT_build.md ✓  PROMPT_plan.md ✓  AGENTS.md ✓  ║
║  tasks.json ✓                                        ║
╠══════════════════════════════════════════════════════╣
║  [sys]  Loop started — iteration 14                  ║
║                                                      ║
║  [#14]  Reading specs/auth.md...                     ║
║         Implementing JWT middleware...               ║
║         Running go test ./...                        ║
║                                                      ║
║  [sys]  Iteration 14: PASS                           ║
║                                                      ║
║  [#15]  Selecting next task from tasks.json... ▌     ║
╠══════════════════════════════════════════════════════╣
║  m: mode  e: edit plan  R: regen  p: pause  q: quit ║
╚══════════════════════════════════════════════════════╝
```

Agent output is shown as numbered turns `[#N]` with a streaming cursor `▌` while the agent is running. System events (mode changes, iteration results, errors) appear as `[sys]` turns in italics. The viewport auto-scrolls to new content; press `↑` to scroll back.

## State Persistence

Ralphio maintains state in `.ralph/state.json` within your project directory:

```json
{
  "currentIteration": 5,
  "currentTaskId": "task-3",
  "loopStatus": "running",
  "loopMode": "building",
  "activeAdapter": "claude",
  "activeModel": "",
  "lastUpdated": "2024-01-15T10:30:00Z"
}
```

This enables crash recovery — restart ralphio after an interruption and it resumes from the last known state. The agent will see any `in_progress` task in `tasks.json` and handle it appropriately.

To start completely fresh:

```bash
./ralphio --reset-state
```

## Execution Flow

1. **Load state** — read `.ralph/state.json`; apply `--reset-state` or `--mode` if passed
2. **Load tasks** — read `tasks.json` from project directory
3. **Send initial snapshot** — TUI displays current task list and mode
4. **Loop**:
   - Build prompt from `PROMPT_build.md` or `PROMPT_plan.md` (based on mode)
   - Execute agent with full prompt; stream output to TUI chat panel
   - Agent autonomously: selects task, implements, runs validation, commits, updates `tasks.json`
   - Orchestrator reloads `tasks.json` after each iteration
   - If all tasks are `completed`, loop exits
5. **Persist state** — save updated `state.json` after each iteration

## Configuration File

Create `~/.ralphio.json` for persistent settings:

```json
{
  "logLevel": "info",
  "debug": false,
  "ui": {
    "altScreen": false,
    "mouseEnabled": true
  },
  "ralph": {
    "projectDir": ".",
    "agent": "claude",
    "agentModel": "",
    "maxRetries": 3,
    "iterationDelayMs": 2000
  }
}
```

## Examples

### New project: plan then build

```bash
mkdir my-feature && cd my-feature

# Create a PRD (Product Requirements Document)
cat > PRD.md << 'EOF'
# Auth Service

Implement JWT authentication with login/logout endpoints.

## Requirements
- JWT token generation and validation
- Login endpoint with username/password
- Logout endpoint that invalidates tokens
- Middleware to protect authenticated routes
EOF

# 1. Generate plan from PRD
ralphio --mode planning --agent claude

# 2. Review tasks.json, then switch to building
ralphio --mode building --agent claude
```

### Resume after interruption

```bash
# tasks.json and .ralph/state.json are preserved
ralphio --agent claude  # picks up where it left off
```

### Reset and restart

```bash
ralphio --reset-state --mode planning
```

### Toggle mode at runtime

Press `m` in the dashboard to switch between PLANNING and BUILDING modes without restarting. The new mode takes effect on the next iteration.

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Agent not found | Ensure the agent CLI is installed and in `PATH` |
| No output visible | Agent may not support streaming JSON output format |
| Tasks not updating | Check file permissions on `tasks.json` |
| State corruption | Run with `--reset-state` to clear `.ralph/` |
| Want a fresh plan | Press `R` in dashboard or run with `--mode planning --reset-state` |
| Custom prompts ignored | Ensure file is named exactly `PROMPT_plan.md` or `PROMPT_build.md` |
