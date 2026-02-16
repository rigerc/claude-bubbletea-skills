# BubbleTea v2 — API Reference

## Package Import

```go
import tea "charm.land/bubbletea/v2"
```

---

## Program Lifecycle

### `tea.NewProgram(model, opts...) *Program`
Creates a new program. Options are applied in order.

### `(*Program).Run() (Model, error)`
Starts the program, blocks until exit. Returns final model state.

### `(*Program).Send(msg Msg)`
Sends a message to the running program from an external goroutine.

---

## Program Options

| Option | Description |
|---|---|
| `WithContext(ctx)` | Cancel program via context |
| `WithOutput(io.Writer)` | Custom output writer |
| `WithInput(io.Reader)` | Custom input reader |
| `WithEnvironment([]string)` | Override environment |
| `WithFPS(fps int)` | Frame rate 1–120, default 60 |
| `WithColorProfile(p)` | Override color detection |
| `WithWindowSize(w, h int)` | Override terminal size (useful for testing) |
| `WithFilter(fn)` | Intercept/transform messages before delivery |
| `WithoutRenderer()` | Disable rendering (headless mode) |
| `WithoutSignalHandler()` | Disable default signal handling |
| `WithoutCatchPanics()` | Disable panic recovery |

---

## The View Type

```go
type View struct {
    Content                  string
    Cursor                   *Cursor
    BackgroundColor          color.Color
    ForegroundColor          color.Color
    WindowTitle              string
    ProgressBar              *ProgressBar
    AltScreen                bool
    ReportFocus              bool
    DisableBracketedPasteMode bool
    MouseMode                MouseMode
    KeyboardEnhancements     KeyboardEnhancements
    OnMouse                  func(MouseMsg) Cmd  // view-local mouse handler
}
```

### `tea.NewView(content string) View`
Create a View with the given content string.

### Cursor Shapes
| Constant | Description |
|---|---|
| `CursorBlock` | Block cursor |
| `CursorUnderline` | Underline cursor |
| `CursorBar` | I-beam / bar cursor |

### Mouse Modes
| Constant | Description |
|---|---|
| `MouseModeNone` | No mouse events |
| `MouseModeCellMotion` | Clicks, drags, wheel |
| `MouseModeAllMotion` | All movement events (high traffic) |

---

## Message Types

### Keyboard

| Type | Fields | Notes |
|---|---|---|
| `KeyPressMsg` | `Key() KeyEvent`, `String() string`, `IsRepeat bool` | Key down |
| `KeyReleaseMsg` | `Key() KeyEvent` | Key up (enhanced mode only) |
| `KeyMsg` | Interface | Both press and release |

**KeyEvent fields:**
```go
type KeyEvent struct {
    Code        rune     // special key constant or printable rune
    ShiftedCode rune     // shifted variant
    BaseCode    rune     // PC-101 base layout
    Mod         KeyMod   // bitmask of modifiers
    Text        string   // printable text representation
    IsRepeat    bool
}
```

**Key constants:** `KeyEnter`, `KeyEscape`, `KeyBackspace`, `KeyDelete`,
`KeyTab`, `KeySpace`, `KeyUp`, `KeyDown`, `KeyLeft`, `KeyRight`,
`KeyHome`, `KeyEnd`, `KeyPageUp`, `KeyPageDown`,
`KeyF1`–`KeyF63`, `KeyInsert`, `KeyPrintScreen`, `KeyScrollLock`, `KeyPause`

**Modifier constants:** `ModShift`, `ModAlt`, `ModCtrl`, `ModMeta`, `ModSuper`, `ModHyper`, `ModCapsLock`, `ModNumLock`

### Mouse

| Type | Fields |
|---|---|
| `MouseClickMsg` | `X, Y int`, `Button MouseButton`, `Mod KeyMod` |
| `MouseReleaseMsg` | `X, Y int`, `Button MouseButton`, `Mod KeyMod` |
| `MouseWheelMsg` | `X, Y int`, `Button MouseButton`, `Mod KeyMod` |
| `MouseMotionMsg` | `X, Y int`, `Button MouseButton`, `Mod KeyMod` |

All implement `MouseMsg` interface: `Mouse() MouseEvent`.

**Button constants:** `MouseLeft`, `MouseRight`, `MouseMiddle`,
`MouseWheelUp`, `MouseWheelDown`, `MouseWheelLeft`, `MouseWheelRight`,
`MouseButton8`–`MouseButton11`, `MouseNone`

### Window & Focus

| Type | Fields |
|---|---|
| `WindowSizeMsg` | `Width, Height int` |
| `FocusMsg` | (no fields) |
| `BlurMsg` | (no fields) |

### Terminal Color Queries

| Type | Method / Fields |
|---|---|
| `BackgroundColorMsg` | `IsDark() bool`, `Color() color.Color` |
| `ForegroundColorMsg` | `Color() color.Color` |
| `CursorColorMsg` | `Color() color.Color` |
| `ColorProfileMsg` | `Profile colorprofile.Profile` |

### Paste

| Type | Description |
|---|---|
| `PasteMsg` | Text from bracketed paste |
| `PasteStartMsg` | Bracketed paste start |
| `PasteEndMsg` | Bracketed paste end |

### Clipboard (OSC52)

| Type | Description |
|---|---|
| `ClipboardMsg` | `string(msg)` = clipboard contents |
| `PrimaryClipboardMsg` | Primary selection contents |

### Program Control

| Type | Description |
|---|---|
| `QuitMsg` | Graceful exit |
| `InterruptMsg` | Ctrl+C |
| `SuspendMsg` | Ctrl+Z suspend |
| `ResumeMsg` | Resume after suspend |
| `BatchMsg` | Internal: concurrent commands |

---

## Commands (Cmd)

A `Cmd` is `func() Msg` — runs in a goroutine and sends its return value as a message.

### Lifecycle
```go
tea.Quit        // var Cmd — send QuitMsg
tea.Interrupt   // var Cmd — send InterruptMsg
tea.Suspend     // var Cmd — send SuspendMsg
```

### Composition
```go
tea.Batch(cmds ...Cmd) Cmd    // run all concurrently
tea.Sequence(cmds ...Cmd) Cmd // run in order, each waits for previous
```

### Timers
```go
tea.Tick(d time.Duration, fn func(time.Time) Msg) Cmd
tea.Every(d time.Duration, fn func(time.Time) Msg) Cmd  // aligned to clock
```

### Terminal Queries
```go
tea.RequestWindowSize() Msg
tea.RequestBackgroundColor() Msg
tea.RequestForegroundColor() Msg
tea.RequestCursorColor() Msg
tea.RequestCursorPosition() Msg
tea.RequestTerminalVersion() Msg
tea.RequestCapability(s string) Msg
```

### Output
```go
tea.Printf(format string, args ...any) Cmd  // print above program
tea.Println(args ...any) Cmd
tea.Raw(seq any) Cmd                        // raw ANSI escape sequence
```

### Clipboard
```go
tea.SetClipboard(s string) Cmd
tea.ReadClipboard() Msg                     // returns ClipboardMsg
tea.SetPrimaryClipboard(s string) Cmd
tea.ReadPrimaryClipboard() Msg
```

### External Processes
```go
tea.Exec(cmd ExecCommand, onExit func(error) Msg) Cmd
tea.ExecProcess(cmd *exec.Cmd, onExit func(error) Msg) Cmd
```

---

## Errors

```go
tea.ErrInterrupted    // Ctrl+C / InterruptMsg received
tea.ErrProgramKilled  // Program killed externally
tea.ErrProgramPanic   // Panic was caught
```

---

## Logging

```go
tea.LogToFile(path, prefix string) (io.Closer, error)
tea.LogToFileWith(path, prefix string, logger *slog.Logger) (io.Closer, error)
```

Environment variables:
- `TEA_TRACE=<path>` — write trace log to path
- `TEA_DEBUG=true` — enable debug output

---

## Related Packages

| Package | Purpose |
|---|---|
| `github.com/charmbracelet/lipgloss` | ANSI styles, colors, layout |
| `github.com/charmbracelet/bubbles` | Reusable components (spinner, list, input, viewport…) |
| `github.com/charmbracelet/x/ansi` | ANSI utilities |
| `github.com/charmbracelet/colorprofile` | Terminal color capability detection |
