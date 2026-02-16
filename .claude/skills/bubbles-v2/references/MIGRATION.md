# Migrating from Bubbles v1 to v2

Source: `charm.land/bubbles/v2`
Previous: `github.com/charmbracelet/bubbles`

> Companion upgrades required: Bubble Tea v2 + Lip Gloss v2 must be upgraded at the same time.

```sh
go get charm.land/bubbletea/v2
go get charm.land/bubbles/v2
go get charm.land/lipgloss/v2
```

---

## 1. Import Paths

```go
// Before
import "github.com/charmbracelet/bubbles/spinner"
// After
import "charm.land/bubbles/v2/spinner"
```

Search-and-replace:
```
github.com/charmbracelet/bubbles/  →  charm.land/bubbles/v2/
github.com/charmbracelet/bubbles   →  charm.land/bubbles/v2
```

> Note: `runeutil` and `memoization` are now internal packages — not importable.

---

## 2. Global Patterns

### 2a. `tea.KeyMsg` → `tea.KeyPressMsg`

```go
// Before
case tea.KeyMsg:
// After
case tea.KeyPressMsg:
```

### 2b. Width/Height: Fields → Methods

Affected: `filepicker`, `help`, `progress`, `table`, `textinput`, `viewport`

```go
// Before
m.Width = 40; m.Height = 20
fmt.Println(m.Width, m.Height)
// After
m.SetWidth(40); m.SetHeight(20)
fmt.Println(m.Width(), m.Height())
```

### 2c. `DefaultKeyMap` Variables → Functions

Affected: `paginator`, `textarea`, `textinput`

```go
// Before: km := textinput.DefaultKeyMap
// After:  km := textinput.DefaultKeyMap()
```

### 2d. `AdaptiveColor` → Explicit `isDark bool`

Lip Gloss v2 removes `AdaptiveColor`. All `DefaultStyles()` now require `isDark bool`.

```go
// Before: h := help.New()  // auto-adapted
// After:
h := help.New()
h.Styles = help.DefaultStyles(isDark)   // must pass isDark
```

### 2e. `NewModel` Aliases Removed

All `NewModel` aliases removed. Use `New()` directly.

Affected: `help`, `list`, `paginator`, `spinner`, `textinput`

---

## 3. Per-Component Breaking Changes

### cursor

| v1 | v2 |
|---|---|
| `model.Blink` | `model.IsBlinked` |
| `model.BlinkCmd()` | `model.Blink()` |

### filepicker

| v1 | v2 |
|---|---|
| `DefaultStylesWithRenderer(r)` | `DefaultStyles()` |
| `model.Height = 10` | `model.SetHeight(10)` |
| `_ = model.Height` | `_ = model.Height()` |

### help

| v1 | v2 |
|---|---|
| `model.Width = 80` | `model.SetWidth(80)` |
| `_ = model.Width` | `_ = model.Width()` |
| `NewModel()` | `New()` |
| `DefaultStyles()` | `DefaultStyles(isDark)` |

New: `DefaultDarkStyles()`, `DefaultLightStyles()`

### list

| v1 | v2 |
|---|---|
| `DefaultStyles()` | `DefaultStyles(isDark)` |
| `NewDefaultItemStyles()` | `NewDefaultItemStyles(isDark)` |
| `styles.FilterPrompt` | `styles.Filter.Focused.Prompt` |
| `styles.FilterCursor` | `styles.Filter.Cursor` |
| `NewModel(...)` | `New(...)` |

### paginator

| v1 | v2 |
|---|---|
| `DefaultKeyMap` (var) | `DefaultKeyMap()` (func) |
| `NewModel(...)` | `New(...)` |
| `model.UsePgUpPgDownKeys` | Removed — customize `KeyMap` directly |
| `model.UseLeftRightKeys` | Removed |
| `model.UseUpDownKeys` | Removed |
| `model.UseHLKeys` | Removed |
| `model.UseJKKeys` | Removed |

### progress

| v1 | v2 |
|---|---|
| `WithGradient(a, b string)` | `WithColors(colors ...color.Color)` |
| `WithDefaultGradient()` | `WithDefaultBlend()` |
| `WithScaledGradient(a, b)` | `WithColors(...) + WithScaled(true)` |
| `WithDefaultScaledGradient()` | `WithDefaultBlend() + WithScaled(true)` |
| `WithSolidFill(string)` | `WithColors(singleColor)` |
| `WithColorProfile(p)` | Removed (automatic) |
| `p.Width = 40` | `p.SetWidth(40)` |
| `p.FullColor = "#FF0000"` | `p.FullColor = lipgloss.Color("#FF0000")` |

New: `WithColorFunc(fn)`, `WithScaled(bool)`

### spinner

| v1 | v2 |
|---|---|
| `NewModel()` | `New()` |
| `spinner.Tick()` (package func) | `model.Tick()` (method) |

### stopwatch

| v1 | v2 |
|---|---|
| `NewWithInterval(d)` | `New(WithInterval(d))` |

### table

| v1 | v2 |
|---|---|
| `model.viewport.Width` (direct access) | `model.Width()` / `model.SetWidth(w)` |

### textarea

| v1 | v2 |
|---|---|
| `DefaultKeyMap` (var) | `DefaultKeyMap()` (func) |
| `Style` (type) | `StyleState` (type) |
| `model.FocusedStyle` | `model.Styles.Focused` |
| `model.BlurredStyle` | `model.Styles.Blurred` |
| `DefaultStyles()` | `DefaultStyles(isDark)` |
| `model.SetCursor(col)` | `model.SetCursorColumn(col)` |
| `model.Cursor` (cursor.Model) | `model.Cursor()` (func → *tea.Cursor) |

New: `Column()`, `ScrollYOffset()`, `MoveToBegin()`, `MoveToEnd()`, `PageUp()`, `PageDown()`

### textinput

| v1 | v2 |
|---|---|
| `DefaultKeyMap` (var) | `DefaultKeyMap()` (func) |
| `NewModel()` | `New()` |
| `ti.Width = 40` | `ti.SetWidth(40)` |
| `ti.PromptStyle` | `StyleState.Prompt` |
| `ti.TextStyle` | `StyleState.Text` |
| `ti.PlaceholderStyle` | `StyleState.Placeholder` |
| `ti.CompletionStyle` | `StyleState.Suggestion` |
| `ti.CursorStyle` | `Styles.Cursor` |
| `ti.Cursor` (cursor.Model) | `ti.Cursor()` (func → *tea.Cursor) |

Set styles via:
```go
s := textinput.DefaultStyles(isDark)
s.Focused.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
ti.SetStyles(s)
```

### timer

| v1 | v2 |
|---|---|
| `NewWithInterval(timeout, interval)` | `New(timeout, WithInterval(interval))` |

### viewport

| v1 | v2 |
|---|---|
| `New(w, h int)` | `New(...Option)` — `New(WithWidth(w), WithHeight(h))` |
| `vp.Width = 80` | `vp.SetWidth(80)` |
| `vp.Height = 24` | `vp.SetHeight(24)` |
| `vp.YOffset = 5` | `vp.SetYOffset(5)` |
| `HighPerformanceRendering` | Removed entirely |

New viewport features: `SoftWrap`, `LeftGutterFunc`, `SetHighlights`, `HighlightNext/Previous`, `ClearHighlights`, `SetContentLines`, `GetContent`, `FillHeight`, `StyleLineFunc`, horizontal scrolling.

---

## 4. Light and Dark Styles Pattern

```go
// Recommended: query via Bubble Tea (works for SSH/Wish too)
func (m model) Init() tea.Cmd {
    return tea.RequestBackgroundColor
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.BackgroundColorMsg:
        isDark := msg.IsDark()
        m.help.Styles = help.DefaultStyles(isDark)
        m.list.Styles = list.DefaultStyles(isDark)
        m.textInput.SetStyles(textinput.DefaultStyles(isDark))
        m.textArea.SetStyles(textarea.DefaultStyles(isDark))
    }
    return m, nil
}

// Quick alternative (blocks, no SSH support):
import "charm.land/lipgloss/v2/compat"
isDark := compat.HasDarkBackground()

// Force dark/light explicitly:
h.Styles = help.DefaultDarkStyles()
h.Styles = help.DefaultLightStyles()
```

---

## 5. Complete Removed Symbols Reference

| Package | Removed | Replacement |
|---|---|---|
| `cursor` | `Model.Blink` | `Model.IsBlinked` |
| `cursor` | `Model.BlinkCmd()` | `Model.Blink()` |
| `filepicker` | `DefaultStylesWithRenderer(r)` | `DefaultStyles()` |
| `filepicker` | `Model.Height` (field) | `Model.SetHeight()` / `Model.Height()` |
| `help` | `NewModel` | `New()` |
| `help` | `Model.Width` (field) | `Model.SetWidth()` / `Model.Width()` |
| `list` | `NewModel` | `New()` |
| `list` | `DefaultStyles()` | `DefaultStyles(isDark)` |
| `list` | `NewDefaultItemStyles()` | `NewDefaultItemStyles(isDark)` |
| `list` | `Styles.FilterPrompt` | `Styles.Filter` (textinput.Styles) |
| `list` | `Styles.FilterCursor` | `Styles.Filter.Cursor` |
| `paginator` | `DefaultKeyMap` (var) | `DefaultKeyMap()` (func) |
| `paginator` | `NewModel` | `New()` |
| `paginator` | `UsePgUpPgDownKeys` etc. | Customize `KeyMap` directly |
| `progress` | `WithGradient(a, b)` | `WithColors(colors...)` |
| `progress` | `WithDefaultGradient()` | `WithDefaultBlend()` |
| `progress` | `WithScaledGradient(a, b)` | `WithColors(...) + WithScaled(true)` |
| `progress` | `WithDefaultScaledGradient()` | `WithDefaultBlend() + WithScaled(true)` |
| `progress` | `WithSolidFill(string)` | `WithColors(color)` |
| `progress` | `WithColorProfile(p)` | Removed (automatic) |
| `progress` | `Model.Width` (field) | `Model.SetWidth()` / `Model.Width()` |
| `spinner` | `NewModel` | `New()` |
| `spinner` | `Tick()` (package func) | `Model.Tick()` |
| `stopwatch` | `NewWithInterval(d)` | `New(WithInterval(d))` |
| `table` | `Model.Width` (field) | `Model.SetWidth()` / `Model.Width()` |
| `table` | `Model.Height` (field) | `Model.SetHeight()` / `Model.Height()` |
| `textarea` | `DefaultKeyMap` (var) | `DefaultKeyMap()` (func) |
| `textarea` | `Style` (type) | `StyleState` (type) |
| `textarea` | `Model.FocusedStyle` | `Model.Styles.Focused` |
| `textarea` | `Model.BlurredStyle` | `Model.Styles.Blurred` |
| `textarea` | `Model.SetCursor(col)` | `Model.SetCursorColumn(col)` |
| `textarea` | `DefaultStyles()` | `DefaultStyles(isDark)` |
| `textinput` | `DefaultKeyMap` (var) | `DefaultKeyMap()` (func) |
| `textinput` | `NewModel` | `New()` |
| `textinput` | `Model.Width` (field) | `Model.SetWidth()` / `Model.Width()` |
| `textinput` | `Model.PromptStyle` | `StyleState.Prompt` |
| `textinput` | `Model.TextStyle` | `StyleState.Text` |
| `textinput` | `Model.PlaceholderStyle` | `StyleState.Placeholder` |
| `textinput` | `Model.CompletionStyle` | `StyleState.Suggestion` |
| `textinput` | `Model.CursorStyle` | `Styles.Cursor` |
| `textinput` | `Model.Cursor` (cursor.Model) | `Model.Cursor()` (func → *tea.Cursor) |
| `timer` | `NewWithInterval(t, i)` | `New(t, WithInterval(i))` |
| `viewport` | `New(w, h int)` | `New(...Option)` |
| `viewport` | `Model.Width` (field) | `Model.SetWidth()` / `Model.Width()` |
| `viewport` | `Model.Height` (field) | `Model.SetHeight()` / `Model.Height()` |
| `viewport` | `Model.YOffset` (field) | `Model.SetYOffset()` / `Model.YOffset()` |
| `viewport` | `HighPerformanceRendering` | Removed |
| `runeutil` | Entire package | Moved to `internal/runeutil` (not importable) |
