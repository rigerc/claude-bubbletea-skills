# Plan: Refactor ralphio to Align with Ralph Workflow

## Context

The current `ralphio` implementation is a **task runner** that deterministically executes pre-defined tasks from `tasks.json`. The Ralph workflow (per the playbook) is fundamentally different: an agent autonomously manages its own plan, selects tasks, and evolves the implementation through iteration.

This plan refactors ralphio to align with the Ralph methodology while preserving the multi-adapter architecture and TUI.

---

## Key Differences (Current vs Ralph)

| Aspect | Current ralphio | Ralph Playbook |
|--------|-----------------|----------------|
| Task selection | Orchestrator picks by priority | Agent picks "most important" |
| Modes | BUILD only | PLANNING + BUILDING |
| Plan format | `tasks.json` (JSON) | `IMPLEMENTATION_PLAN.md` (markdown) |
| Context per loop | Single task fields | Full specs + plan + AGENTS.md |
| Plan mutation | Status only | Add/remove tasks, discoveries, notes |

---

## Phase 1: File Structure Alignment

### 1a: New Required Files

Create project structure expected by Ralph:

```
project-root/
├── PROMPT_build.md          # Build mode instructions
├── PROMPT_plan.md           # Plan mode instructions  
├── AGENTS.md                # Operational guide (build/test commands)
├── IMPLEMENTATION_PLAN.md   # Agent-maintained task list
└── specs/                   # Requirement specifications
    └── *.md
```

### 1b: Deprecate `tasks.json`

- Keep for backward compatibility (optional migration)
- Agent reads/writes `IMPLEMENTATION_PLAN.md` instead
- TUI displays markdown plan content

### 1c: State File Update

Update `.ralph/state.json`:

```json
{
  "currentIteration": 14,
  "currentTaskId": "T-14",
  "loopStatus": "running",
  "loopMode": "building",
  "activeAdapter": "claude",
  "lastUpdated": "..."
}
```

Add `loopMode`: `"planning"` or `"building"`.

---

## Phase 2: Loop Mode System

### 2a: Mode Definition — `internal/plan/mode.go`

```go
type LoopMode string

const (
    ModePlanning LoopMode = "planning"
    ModeBuilding LoopMode = "building"
)
```

### 2b: Mode-Aware Orchestrator

Update `orchestrator.go` to:

1. Check current mode from state
2. Pass appropriate prompt file to adapter
3. PLANNING mode: agent generates/updates plan only
4. BUILDING mode: agent implements from plan

### 2c: Mode Selection in TUI

- Add `m` key binding to toggle mode
- Dashboard shows current mode indicator
- Mode change persists to state file

---

## Phase 3: Agent-Driven Task Selection

### 3a: Remove Orchestrator Task Selection

Current (remove):

```go
next := plan.NextTask(tasks)
if next == nil {
    o.send(LoopDoneMsg{})
    return
}
```

New approach:

```go
// Agent receives full context and decides what to do
prompt := o.buildPrompt()
execErr := o.currentAdapter.Execute(ctx, prompt, func(text string) {
    o.send(AgentOutputMsg{Text: text})
})
```

### 3b: Prompt Composition — `internal/prompt/builder.go`

Build prompt from project files each iteration:

```go
type PromptBuilder struct {
    projectDir string
    mode       plan.LoopMode
}

func (b *PromptBuilder) Build() (string, error) {
    // 1. Read PROMPT_build.md or PROMPT_plan.md based on mode
    // 2. Return full prompt text
    return promptContent, nil
}
```

The prompt files themselves instruct the agent to:
- Study `specs/*`
- Read `IMPLEMENTATION_PLAN.md`
- Study `src/lib/*`
- Select most important task
- Implement or plan accordingly

### 3c: Agent Commits, Not Orchestrator

The agent (via prompt instructions) handles:
- Updating `IMPLEMENTATION_PLAN.md`
- Running validation via `AGENTS.md` commands
- Git commit when done

Orchestrator just:
1. Builds prompt
2. Executes adapter
3. Streams output to TUI
4. Waits for completion
5. Loops

---

## Phase 4: Plan Manager Refactor

### 4a: Markdown Plan Parser — `internal/plan/parser.go`

Parse `IMPLEMENTATION_PLAN.md` for TUI display:

```go
type ParsedPlan struct {
    Tasks     []MarkdownTask
    RawContent string
}

type MarkdownTask struct {
    ID          string
    Title       string
    Status      string // derived from checkbox [x] or [ ]
    Priority    int    // derived from order
    RawMarkdown string
}

func ParsePlan(content string) (*ParsedPlan, error)
```

### 4b: Plan Watcher

Watch for file changes to update TUI in real-time:

```go
func (m *Manager) Watch(ctx context.Context, onChange func(content string)) error
```

Use `fsnotify` or polling.

### 4c: Remove `plan.Manager` CRUD

Agent manages plan directly via file writes. Remove:
- `SaveTasks()`
- `UpdateTask()`
- `NextTask()`

Keep:
- File watching
- Parsing for TUI display

---

## Phase 5: Prompt Templates

### 5a: `PROMPT_plan.md` Template

Port from playbook:

```markdown
0a. Study `specs/*` with up to 250 parallel Sonnet subagents to learn the application specifications.
0b. Study @IMPLEMENTATION_PLAN.md (if present) to understand the plan so far.
0c. Study `src/lib/*` with up to 250 parallel Sonnet subagents to understand shared utilities & components.
0d. For reference, the application source code is in `src/*`.

1. Study @IMPLEMENTATION_PLAN.md (if present; it may be incorrect) and use up to 500 Sonnet subagents to study existing source code in `src/*` and compare it against `specs/*`. Use an Opus subagent to analyze findings, prioritize tasks, and create/update @IMPLEMENTATION_PLAN.md as a bullet point list sorted in priority of items yet to be implemented. Ultrathink. Consider searching for TODO, minimal implementations, placeholders, skipped/flaky tests, and inconsistent patterns.

IMPORTANT: Plan only. Do NOT implement anything. Do NOT assume functionality is missing; confirm with code search first.
```

### 5b: `PROMPT_build.md` Template

Port from playbook:

```markdown
0a. Study `specs/*` with up to 500 parallel Sonnet subagents to learn the application specifications.
0b. Study @IMPLEMENTATION_PLAN.md.
0c. For reference, the application source code is in `src/*`.

1. Your task is to implement functionality per the specifications using parallel subagents. Follow @IMPLEMENTATION_PLAN.md and choose the most important item to address. Before making changes, search the codebase (don't assume not implemented) using Sonnet subagents.
2. After implementing functionality or resolving problems, run the tests for that unit of code that was improved.
3. When you discover issues, immediately update @IMPLEMENTATION_PLAN.md with your findings.
4. When the tests pass, update @IMPLEMENTATION_PLAN.md, then `git add -A` then `git commit` with a message describing the changes.

99999. Important: When authoring documentation, capture the why.
999999. Single sources of truth, no migrations/adapters.
9999999. As soon as there are no build or test errors create a git tag.
9999999999. Keep @IMPLEMENTATION_PLAN.md current with learnings.
```

### 5c: Configurable Prompt Paths

Allow users to specify custom prompt files:

```bash
ralphio --prompt-dir ./my-prompts/
ralphio --build-prompt ./PROMPT_custom.md
```

---

## Phase 6: AGENTS.md Integration

### 6a: AGENTS.md Discovery

Orchestrator checks for `AGENTS.md` existence:

```go
func (b *PromptBuilder) hasAgentsMD() bool
```

### 6b: Include in Context

The prompt instructs agent to study `@AGENTS.md` for:
- Build commands
- Test commands
- Project-specific patterns

### 6c: TUI Display

Show AGENTS.md status in dashboard:
- "AGENTS.md: ✓ found" or "✗ missing"
- Quick view on key press

---

## Phase 7: Validation Strategy

### 7a: Agent-Driven Validation

Remove orchestrator's `validationCommand` execution. Agent runs validation per prompt instructions.

### 7b: Validation Status Detection

Parse agent output to detect validation results:

```go
type ValidationSignal struct {
    Passed bool
    Output string
}

func DetectValidation(output string) *ValidationSignal
```

Look for patterns:
- "tests passed" / "all tests pass"
- Exit code indicators
- Build success/failure messages

### 7c: Backpressure Visualization

TUI shows validation status detected from agent stream:
- Green: validation passed
- Red: validation failed, agent retrying
- Yellow: validation running

---

## Phase 8: TUI Updates

### 8a: Dashboard Mode Indicator

```
╔══════════════════════════════════════════════════════╗
║ RALPHIO | Project: auth-service | Mode: BUILDING    ║
╠══════════════════════════════════════════════════════╣
```

### 8b: Plan View from Markdown

Parse and display `IMPLEMENTATION_PLAN.md`:

```
╠══════════════════════ IMPLEMENTATION PLAN ════════════════╣
║ [x] Setup project scaffold                              ║
║ [x] Add linting                                         ║
║ [ ] Implement auth middleware  ← CURRENT                ║
║ [ ] Add integration tests                               ║
╚══════════════════════════════════════════════════════════╝
```

### 8c: File Status Indicators

Show which Ralph files exist:

```
╠════════════════════════ RALPH FILES ═════════════════════╣
║ PROMPT_build.md ✓  PROMPT_plan.md ✓  AGENTS.md ✓       ║
║ IMPLEMENTATION_PLAN.md ✓  specs/ (3 files) ✓            ║
╚══════════════════════════════════════════════════════════╝
```

### 8d: Key Binding Updates

| Key | Action |
|-----|--------|
| `m` | Toggle mode (planning/building) |
| `e` | Edit IMPLEMENTATION_PLAN.md |
| `R` | Regenerate plan (switch to planning mode) |
| `a` | View AGENTS.md |

---

## Phase 9: Crash Recovery

### 9a: State Recovery

On startup, check:
1. `.ralph/state.json` exists
2. `IMPLEMENTATION_PLAN.md` exists
3. Git status (uncommitted changes?)

### 9b: Resume Logic

- If mode was `building` and task in progress → agent will see it in plan
- If mode was `planning` → resume planning iteration
- No manual task state management needed

### 9c: Clean State Option

```bash
ralphio --reset-state
```

Clears `.ralph/` directory for fresh start.

---

## Implementation Order

| Order | Phase | Files | Dependencies |
|-------|-------|-------|--------------|
| 1 | 1c | `internal/state/state.go` | None |
| 2 | 2a | `internal/plan/mode.go` | None |
| 3 | 5a, 5b | `assets/PROMPT_plan.md`, `assets/PROMPT_build.md` | None |
| 4 | 3b | `internal/prompt/builder.go` | Phase 2, 3 |
| 5 | 4a | `internal/plan/parser.go` | None |
| 6 | 3a | `internal/orchestrator/orchestrator.go` refactor | Phases 1-5 |
| 7 | 7a, 7b | Remove validation from orchestrator | Phase 6 |
| 8 | 8a-8d | TUI updates | Phases 4, 6 |
| 9 | 6a-6c | AGENTS.md integration | Phase 3 |
| 10 | 9a-9c | Crash recovery | Phase 1 |

---

## File Tree (changes)

```
internal/
  plan/
    mode.go           # NEW: LoopMode constants
    parser.go         # NEW: Markdown plan parser
    plan.go           # MODIFY: Remove NextTask, UpdateTask, SaveTasks
  prompt/
    builder.go        # NEW: Prompt composition from files
  orchestrator/
    orchestrator.go   # MAJOR REFACTOR: Agent-driven loop
    messages.go       # MODIFY: Add mode-related messages
  state/
    state.go          # MODIFY: Add LoopMode field
  validator/
    validator.go      # DEPRECATE or keep for optional use
  ui/
    screens/
      dashboard.go    # MODIFY: Mode indicator, plan view, file status
config/
  config.go           # MODIFY: Add promptDir, ralphMode config
  defaults.go         # MODIFY: Default prompt paths
assets/
  PROMPT_plan.md      # NEW: Default planning prompt
  PROMPT_build.md     # NEW: Default building prompt
cmd/
  root.go             # MODIFY: Add --mode, --reset-state flags
```

---

## Verification

1. **Mode switching**: `ralphio --mode planning` starts in planning mode
2. **Agent task selection**: Orchestrator passes full prompt, agent picks task
3. **Plan mutation**: Agent can add new tasks to IMPLEMENTATION_PLAN.md
4. **TUI reflects plan**: Markdown changes appear in dashboard
5. **Crash recovery**: Restart resumes from state, agent sees in-progress task
6. **Two modes work**: Planning generates plan, building implements from it

---

## Migration Path

For existing `tasks.json` users:

1. **Optional migration tool**: `ralphio migrate` converts tasks.json → IMPLEMENTATION_PLAN.md
2. **Deprecation warning**: Log message if tasks.json found
3. **Full removal**: After 2 releases, remove tasks.json support

---

## Risks

| Risk | Mitigation |
|------|------------|
| Agent doesn't follow prompt | Well-tested prompt templates; user can customize |
| Plan parsing fragile | Use robust markdown parsing; fallback to raw display |
| Validation detection unreliable | Best-effort detection; agent handles actual validation |
| Breaking change for users | Migration tool; deprecation period |

---

## Success Criteria

- [ ] Agent selects tasks autonomously from plan
- [ ] PLANNING mode generates/updates IMPLEMENTATION_PLAN.md
- [ ] BUILDING mode implements from plan
- [ ] TUI shows live plan content from markdown
- [ ] AGENTS.md loaded for project-specific commands
- [ ] Agent can add/remove tasks from plan
- [ ] State recovery works after crash
