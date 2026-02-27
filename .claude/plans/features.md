# Implementation Plan
## Async Tasks · Context · Crash Recovery · First Run · Modal

---

## 1. Async Task Infrastructure + Context Propagation

### The Problem
`ui.Run` calls `tea.NewProgram(m).Run()` with no context. Background goroutines have no cancellation signal, so navigating away or quitting mid-task leaks them silently.

### New Files
- `internal/task/task.go` — `Run`, `RunWithTimeout`, message types
- `internal/task/msgs.go` — `StartedMsg`, `DoneMsg[T]`, `ErrMsg`, `ProgressMsg`
- `internal/ui/spinner/spinner.go` — thin wrapper around `bubbles/spinner`

### Changes to Existing Files

**`internal/ui/ui.go`** — accept a `context.Context`:
```go
func Run(ctx context.Context, m rootModel) error {
    _, err := tea.NewProgram(m, tea.WithContext(ctx)).Run()
    return err
}
```

**`main.go`** — create root context with cancel, wire it through:
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
if err := ui.Run(ctx, ui.New(*cfg, configPath)); err != nil { ... }
```

**`internal/ui/model.go`** — store `ctx` and `cancel` on `rootModel`:
```go
type rootModel struct {
    ctx    context.Context
    cancel context.CancelFunc
    // ...existing fields
}
```
Handle `task.ErrMsg` at root level → route to `status.SetError`.
Handle `tea.Quit` from task cancellation on `ctx.Done()`.

### The Task Package

```go
// task/task.go

type Result[T any] struct {
    Value T
    Err   error
    Label string
}

// Run executes fn in a goroutine, sends DoneMsg[T] when complete.
// Respects ctx cancellation — sends ErrMsg(context.Canceled) if ctx fires first.
func Run[T any](ctx context.Context, label string, fn func(ctx context.Context) (T, error)) tea.Cmd {
    return func() tea.Msg {
        done := make(chan Result[T], 1)
        go func() {
            v, err := fn(ctx)
            done <- Result[T]{Value: v, Err: err, Label: label}
        }()
        select {
        case r := <-done:
            if r.Err != nil {
                return ErrMsg{Label: label, Err: r.Err}
            }
            return DoneMsg[T]{Label: label, Value: r.Value}
        case <-ctx.Done():
            return ErrMsg{Label: label, Err: ctx.Err()}
        }
    }
}
```

### Spinner Integration

Each screen that initiates a task holds a `spinner.Model` and a `loading bool`. Pattern:
```go
// In a screen's Update:
case task.ErrMsg:
    s.loading = false
    return s, status.SetError(msg.Err.Error(), 4*time.Second)
case task.DoneMsg[MyResult]:
    s.loading = false
    s.data = msg.Value
    return s, status.SetSuccess("Loaded", 2*time.Second)
```

The screen's `Body()` returns either the spinner view or the real content based on `s.loading`.

---

## 2. Panic / Crash Recovery

### The Problem
An unhandled panic in any `Update` or `View` tears the terminal raw without restoring it, leaving the shell unusable. The log file gets nothing useful.

### Changes: `main.go`

Wrap the `ui.Run` call in a deferred recovery function. The key constraint: bubbletea's `Program.Run()` already has an internal recover, but it re-panics — so we catch at the outermost layer.

```go
func main() {
    // ... cobra, config, logging setup unchanged ...

    defer func() {
        if r := recover(); r != nil {
            // Terminal was already restored by tea's internal cleanup.
            // Log the stack trace and print a clean exit message.
            buf := make([]byte, 8192)
            n := runtime.Stack(buf, false)
            applogger.Fatal().
                Str("panic", fmt.Sprintf("%v", r)).
                Str("stack", string(buf[:n])).
                Msg("panic: unrecovered")
            fmt.Fprintf(os.Stderr, "\n[scaffold] crashed — see debug.log for details\n")
            os.Exit(2)
        }
    }()

    if err := ui.Run(ctx, ui.New(*cfg, configPath)); err != nil {
        applogger.Fatal().Err(err).Msg("UI failed")
    }
}
```

### Changes: `internal/ui/ui.go`

Enable bubbletea's own `WithReportPanic` option so it restores the terminal before re-panicking up to our recover:

```go
func Run(ctx context.Context, m rootModel) error {
    _, err := tea.NewProgram(m,
        tea.WithContext(ctx),
        tea.WithANSICompressor(),
    ).Run()
    return err
}
```

> **Note:** Debug mode should always write a log file (`debug.log`) so the stack trace isn't lost. Adjust `setupLogOutput` to always open the file when `cfg.Debug` is true — which it already does. For release builds consider always writing crash logs to `~/.appname/crash.log` regardless of debug flag.

---

## 3. First-Run Detection

### The Problem
There's no way to know if this is the first launch (no config file written yet), or if a config was written by an older version that's missing new fields.

### Changes: `config/config.go` — add `ConfigVersion`

```go
type Config struct {
    ConfigVersion int    `json:"configVersion" koanf:"configVersion" cfg_exclude:"true"`
    // ...existing fields unchanged
}

const CurrentConfigVersion = 1
```

### Changes: `config/defaults.go`

Set `ConfigVersion: CurrentConfigVersion` in `DefaultConfig()`.

### New Helper: `config/firstrun.go`

```go
package config

import "os"

// IsFirstRun returns true when no config file exists at the given path.
// A first run means the app has never written its config to disk.
func IsFirstRun(configPath string) bool {
    if configPath == "" {
        return false // no config file expected
    }
    _, err := os.Stat(configPath)
    return os.IsNotExist(err)
}

// NeedsUpgrade returns true when the loaded config's version is behind
// the current schema version. Callers should migrate and re-save.
func NeedsUpgrade(cfg *Config) bool {
    return cfg.ConfigVersion < CurrentConfigVersion
}
```

### Changes: `main.go` / `cmd/root.go`

After `loadConfig`, pass the first-run flag into the model:

```go
cfg, configPath := loadConfig()
firstRun := config.IsFirstRun(configPath)
if err := ui.Run(ctx, ui.New(*cfg, configPath, firstRun)); err != nil { ... }
```

### Changes: `internal/ui/ui.go` + `model.go`

`newRootModel` accepts `firstRun bool`. If true, it pushes `screens.NewWelcome()` as the initial screen instead of `screens.NewHome()`, or navigates to Welcome on the first `Init`:

```go
func (m rootModel) Init() tea.Cmd {
    cmds := tea.Batch(
        tea.RequestBackgroundColor,
        m.themeMgr.Init(m.cfg.UI.ThemeName, false, m.width),
    )
    if m.firstRun {
        return tea.Batch(cmds, func() tea.Msg {
            return NavigateMsg{Screen: screens.NewWelcome()}
        })
    }
    return cmds
}
```

### New File: `internal/ui/screens/welcome.go`

A simple screen that:
1. Displays a welcome message + feature overview (static text / lipgloss-styled)
2. Has a single "Get Started →" action that sends `BackMsg` (returns to Home) and triggers `config.Save` with defaults, marking config as written

```go
type welcomeKeyMap struct {
    Continue key.Binding
}

type Welcome struct {
    theme.ThemeAware
    keys  welcomeKeyMap
    width int
}
```

The `Continue` key dispatches `WelcomeDoneMsg{}` which root handles by calling `config.Save` and navigating back to Home.

---

## 4. Modal / Overlay

### The Problem
There's no way to ask "Are you sure?" or show an inline popup without navigating to a new screen and losing context.

### New File: `internal/ui/modal/modal.go`

A self-contained component rendered by `rootModel` on top of the current screen. It is **not** a `Screen` — it's a separate layer.

```go
package modal

// Kind controls which buttons/actions are available.
type Kind int
const (
    KindConfirm Kind = iota // Yes / No
    KindAlert               // OK only
    KindPrompt              // single-line text input + Submit / Cancel
)

type Model struct {
    kind    Kind
    title   string
    body    string
    input   textinput.Model  // only used for KindPrompt
    visible bool
    keys    keyMap
}

// ConfirmedMsg is sent when the user accepts a KindConfirm modal.
type ConfirmedMsg struct{ ID string }

// CancelledMsg is sent when the user dismisses any modal.
type CancelledMsg struct{ ID string }

// PromptSubmittedMsg is sent when the user submits a KindPrompt modal.
type PromptSubmittedMsg struct{ ID, Value string }

// ShowConfirm returns a Cmd that triggers a confirm modal.
func ShowConfirm(id, title, body string) tea.Cmd { ... }

// ShowAlert returns a Cmd that triggers an alert modal.
func ShowAlert(id, title, body string) tea.Cmd { ... }
```

### New Message Type: `internal/ui/modal/msgs.go`

```go
// ShowMsg is dispatched via tea.Cmd to display a modal.
type ShowMsg struct {
    ID    string
    Kind  Kind
    Title string
    Body  string
}
```

### Changes: `internal/ui/model.go`

Add `modal modal.Model` to `rootModel`. In `Update`:

```go
case modal.ShowMsg:
    m.modal = modal.New(msg)
    return m, nil

case modal.ConfirmedMsg, modal.CancelledMsg, modal.PromptSubmittedMsg:
    m.modal = modal.Model{} // clear
    // delegate to current screen so it can react
    updated, cmd := m.current.Update(msg)
    // ...
    return m, cmd
```

When `m.modal.Visible()`, key events are routed to the modal **instead of** the current screen:

```go
case tea.KeyPressMsg:
    if m.modal.Visible() {
        var cmd tea.Cmd
        m.modal, cmd = m.modal.Update(msg)
        return m, cmd
    }
    // ...normal key handling
```

### Changes: `internal/ui/model.go` — `View`

Render the modal overlay on top if visible:

```go
func (m rootModel) View() tea.View {
    base := lipgloss.JoinVertical(...)  // existing layout
    if m.modal.Visible() {
        return tea.NewView(modal.Overlay(base, m.modal.View(), m.width, m.height))
    }
    return tea.NewView(m.styles.App.Render(base))
}
```

### New File: `internal/ui/modal/overlay.go`

`Overlay(base, popup string, w, h int) string` uses lipgloss `Place` to center the popup over the base content with a dimmed background (via ANSI `\033[2m` on the base, or a lipgloss overlay style).

### Usage from any Screen

```go
// In settings screen's Update, before destructive action:
case tea.KeyPressMsg:
    if key.Matches(msg, s.keys.Reset) {
        return s, modal.ShowConfirm("reset-settings",
            "Reset Settings",
            "This will restore all defaults. Continue?",
        )
    }

// Handle the response:
case modal.ConfirmedMsg:
    if msg.ID == "reset-settings" {
        return s, s.doReset()
    }
```

---

## Execution Order

The features have dependencies that dictate build order:

1. **Context** first — it's a structural change that touches `main.go`, `ui.go`, and `model.go`. Everything else runs inside the program.
2. **Crash recovery** — add the `defer/recover` in `main.go` at the same time as context; they're in the same function.
3. **Async tasks** — builds on the ctx already threaded through; add `task/` package and spinner.
4. **Modal** — self-contained UI layer; add to `model.go` view/update once context is stable.
5. **First run** — last, since it introduces a new screen (`welcome.go`) that may want to use the modal ("Skip setup?") and the task system (async config write).

---

## File Change Summary

| File | Change |
|---|---|
| `main.go` | Add `ctx/cancel`, defer panic recover, pass `firstRun` to `ui.New` |
| `internal/ui/ui.go` | Accept `ctx`, pass to `tea.NewProgram` |
| `internal/ui/model.go` | Add `ctx`, `cancel`, `modal`, `firstRun`; modal routing in Update/View |
| `config/config.go` | Add `ConfigVersion` field |
| `config/defaults.go` | Set `ConfigVersion: CurrentConfigVersion` |
| `config/firstrun.go` | **new** — `IsFirstRun`, `NeedsUpgrade` |
| `internal/task/task.go` | **new** — `Run[T]`, `RunWithTimeout[T]` |
| `internal/task/msgs.go` | **new** — `DoneMsg[T]`, `ErrMsg`, `ProgressMsg` |
| `internal/ui/spinner/spinner.go` | **new** — thin spinner wrapper |
| `internal/ui/modal/modal.go` | **new** — `Model`, `New`, `Overlay` |
| `internal/ui/modal/msgs.go` | **new** — `ShowMsg`, `ConfirmedMsg`, `CancelledMsg`, `PromptSubmittedMsg` |
| `internal/ui/screens/welcome.go` | **new** — first-run welcome screen |
