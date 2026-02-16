---
name: bubbles-v2
description: >-
  Generates, explains, and debugs Bubble Tea component code using Charm's
  Bubbles v2 library for Go. Use when building TUI apps that need spinners,
  text inputs, text areas, lists, tables, progress bars, viewports, file
  pickers, timers, stopwatches, help views, paginators, or keybindings.
  Covers the v2 API at charm.land/bubbles/v2, migration from v1, light/dark
  theme handling, and all component patterns.
argument-hint: [component|pattern|feature description]
metadata:
  author: skill-generator
  version: "1.0"
  category: go-cli
---

# Bubbles v2 — Pre-Built Components for Bubble Tea

Bubbles v2 (`charm.land/bubbles/v2`) provides ready-made TUI components for
Bubble Tea v2 applications. Requires Bubble Tea v2 + Lip Gloss v2.

## Quick Reference

- Full component API → [references/COMPONENTS.md](references/COMPONENTS.md)
- v1 → v2 migration guide → [references/MIGRATION.md](references/MIGRATION.md)

---

## Migrating from Bubbles v1?

If the codebase uses `github.com/charmbracelet/bubbles`, read the migration
guide before making any changes: **[references/MIGRATION.md](references/MIGRATION.md)**

Search-and-replace import paths:
```
github.com/charmbracelet/bubbles/  →  charm.land/bubbles/v2/
github.com/charmbracelet/bubbles   →  charm.land/bubbles/v2
```

Upgrade all three together:
```sh
go get charm.land/bubbletea/v2
go get charm.land/bubbles/v2
go get charm.land/lipgloss/v2
```

Key breaking changes at a glance:

| Area | v1 | v2 |
|---|---|---|
| Import prefix | `github.com/charmbracelet/bubbles` | `charm.land/bubbles/v2` |
| Key press type | `case tea.KeyMsg:` | `case tea.KeyPressMsg:` |
| Width/Height | exported fields (`m.Width = 40`) | getter/setter methods (`m.SetWidth(40)`) |
| DefaultKeyMap | variable | function: `DefaultKeyMap()` |
| NewModel | deprecated alias | removed — use `New()` |
| Adaptive colors | `AdaptiveColor` (auto) | `DefaultStyles(isDark bool)` (explicit) |
| Progress colors | `string` hex | `lipgloss.Color` / `color.Color` |
| Viewport ctor | `New(w, h int)` | `New(...Option)` |
| Stopwatch ctor | `NewWithInterval(d)` | `New(WithInterval(d))` |
| Timer ctor | `NewWithInterval(t,i)` | `New(t, WithInterval(i))` |
| Spinner tick | `spinner.Tick()` (pkg func) | `model.Tick()` (method) |

---

## Light and Dark Styles

**Critical:** Bubbles v2 removes automatic color adaptation. You must pass
`isDark bool` to style constructors. Always handle `tea.BackgroundColorMsg`:

```go
func (m model) Init() tea.Cmd {
    return tea.RequestBackgroundColor  // request terminal background
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.BackgroundColorMsg:
        isDark := msg.IsDark()
        m.help.Styles    = help.DefaultStyles(isDark)
        m.list.Styles    = list.DefaultStyles(isDark)
        m.textInput.SetStyles(textinput.DefaultStyles(isDark))
        m.textArea.SetStyles(textarea.DefaultStyles(isDark))
    }
    return m, nil
}
```

Quick alternative (blocks, not suitable for SSH/Wish):
```go
import "charm.land/lipgloss/v2/compat"
isDark := compat.HasDarkBackground()
```

---

## Component Overview

| Component | Constructor | Purpose |
|---|---|---|
| `spinner` | `spinner.New(opts...)` | Animated loading indicator |
| `textinput` | `textinput.New()` | Single-line text input |
| `textarea` | `textarea.New()` | Multi-line text input |
| `list` | `list.New(items, delegate, w, h)` | Browseable item list with filter |
| `table` | `table.New(opts...)` | Scrollable tabular data |
| `progress` | `progress.New(opts...)` | Progress bar (animated or static) |
| `viewport` | `viewport.New(opts...)` | Scrollable content area |
| `filepicker` | `filepicker.New()` | File system navigator |
| `paginator` | `paginator.New(opts...)` | Pagination logic + display |
| `help` | `help.New()` | Auto-generated keybinding help |
| `key` | `key.NewBinding(opts...)` | Keybinding definitions |
| `timer` | `timer.New(timeout, opts...)` | Countdown timer |
| `stopwatch` | `stopwatch.New(opts...)` | Count-up stopwatch |
| `cursor` | (embedded in inputs) | Text cursor behavior |

---

## Spinner

```go
import "charm.land/bubbles/v2/spinner"

sp := spinner.New(
    spinner.WithSpinner(spinner.Dot),    // Line, Dot, MiniDot, Jump, Pulse, Points, Globe, Moon, Monkey, Meter, Hamburger, Ellipsis
    spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
)

// In Init():
return sp.Tick()   // method, NOT package function (v2 change)

// In Update():
case spinner.TickMsg:
    sp, cmd = sp.Update(msg)
    return m, cmd

// In View():
sp.View()
```

---

## Text Input

```go
import "charm.land/bubbles/v2/textinput"

ti := textinput.New()
ti.Placeholder = "Enter text..."
ti.SetWidth(40)                          // method, not field (v2 change)
ti.CharLimit = 156
ti.EchoMode = textinput.EchoPassword    // for password fields

// Styles (pass isDark from BackgroundColorMsg):
ti.SetStyles(textinput.DefaultStyles(isDark))

// Focus management:
cmd := ti.Focus()    // returns blink command
ti.Blur()
ti.Focused()         // bool

// In Update():
case tea.KeyPressMsg:                    // KeyPressMsg, not KeyMsg (v2 change)
    ti, cmd = ti.Update(msg)

// Value:
ti.Value()
ti.SetValue("initial")
ti.Reset()

// Suggestions (autocomplete):
ti.SetSuggestions([]string{"foo", "bar"})
ti.ShowSuggestions = true

// Cursor:
ti.SetVirtualCursor(false)  // use real terminal cursor
f.Cursor = ti.Cursor()      // add to tea.Frame
```

---

## Text Area

```go
import "charm.land/bubbles/v2/textarea"

ta := textarea.New()
ta.Placeholder = "Enter some text..."
ta.SetWidth(60)
ta.SetHeight(10)
ta.ShowLineNumbers = true

// Styles:
ta.SetStyles(textarea.DefaultStyles(isDark))

// Focus:
cmd := ta.Focus()
ta.Blur()

// In Update():
ta, cmd = ta.Update(msg)

// Value / navigation:
ta.Value()
ta.SetValue("initial content")
ta.Line()        // 0-indexed row
ta.Column()      // 0-indexed column
ta.MoveToBegin() // renamed from MoveToBeginning (v2)
ta.MoveToEnd()

// KeyMap (function in v2, not variable):
km := textarea.DefaultKeyMap()
km.InsertNewline.SetEnabled(false)
ta.KeyMap = km
```

---

## List

```go
import "charm.land/bubbles/v2/list"

// Items implement list.Item interface:
type item struct{ title, desc string }
func (i item) FilterValue() string { return i.title }
func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }

items := []list.Item{item{"Foo", "A thing"}, item{"Bar", "Another"}}

delegate := list.NewDefaultDelegate()
l := list.New(items, delegate, width, height)
l.Title = "My List"
l.Styles = list.DefaultStyles(isDark)                 // isDark required (v2)
delegate.Styles = list.NewDefaultItemStyles(isDark)   // isDark required (v2)

// In Update():
l, cmd = l.Update(msg)

// Selected item:
if i, ok := l.SelectedItem().(item); ok { ... }

// Dynamic items:
cmd = l.SetItems(newItems)
cmd = l.InsertItem(0, newItem)
l.RemoveItem(idx)

// Filter state:
l.FilterState()        // Unfiltered | Filtering | FilterApplied
l.SetFilteringEnabled(false)

// Visibility toggles:
l.SetShowTitle(false)
l.SetShowHelp(false)
l.SetShowStatusBar(false)
l.SetShowPagination(false)
l.DisableQuitKeybindings()
```

---

## Table

```go
import "charm.land/bubbles/v2/table"

cols := []table.Column{
    {Title: "Name", Width: 20},
    {Title: "Age",  Width: 5},
}
rows := []table.Row{{"Alice", "30"}, {"Bob", "25"}}

t := table.New(
    table.WithColumns(cols),
    table.WithRows(rows),
    table.WithFocused(true),
    table.WithHeight(10),
    table.WithWidth(40),
    table.WithStyles(table.DefaultStyles()),  // no isDark (table unchanged in v2)
)

t.SetStyles(table.Styles{
    Header:   lipgloss.NewStyle().Bold(true),
    Cell:     lipgloss.NewStyle(),
    Selected: lipgloss.NewStyle().Reverse(true),
})

// In Update():
t, cmd = t.Update(msg)

// Access:
t.SelectedRow()      // table.Row
t.Cursor()           // int index
t.Rows()
t.SetRows(newRows)
t.SetCursor(n)
t.GotoTop() / t.GotoBottom()
```

---

## Progress Bar

```go
import "charm.land/bubbles/v2/progress"

// Static rendering (no animation):
p := progress.New(
    progress.WithColors(lipgloss.Color("#5A56E0"), lipgloss.Color("#EE6FF8")),
    progress.WithoutPercentage(),
    progress.WithWidth(40),
)
view := p.ViewAs(0.7)  // 70%

// Animated rendering:
p := progress.New(progress.WithDefaultBlend())
// In Init(): return p.Init()
// In Update():
case progress.FrameMsg:
    pm, cmd := p.Update(msg)
    p = pm.(progress.Model)
    return m, cmd

// Trigger animation:
cmd = p.SetPercent(0.9)
cmd = p.IncrPercent(0.1)
cmd = p.DecrPercent(0.1)

// Dynamic color function:
p := progress.New(progress.WithColorFunc(func(total, current float64) color.Color {
    if total < 0.5 { return lipgloss.Color("#FF0000") }
    return lipgloss.Color("#00FF00")
}))
```

---

## Viewport

```go
import "charm.land/bubbles/v2/viewport"

// Constructor changed in v2:
vp := viewport.New(viewport.WithWidth(80), viewport.WithHeight(24))
// or:
vp := viewport.New()
vp.SetWidth(80)
vp.SetHeight(24)

vp.SetContent(longString)
vp.SetContentLines(lines)   // []string, with virtual soft-wrap support

// Features:
vp.SoftWrap = true
vp.MouseWheelEnabled = true

// Line numbers gutter:
vp.LeftGutterFunc = func(info viewport.GutterContext) string {
    if info.Soft { return "     │ " }
    if info.Index >= info.TotalLines { return "   ~ │ " }
    return fmt.Sprintf("%4d │ ", info.Index+1)
}

// Highlighting (search results):
vp.SetHighlights(re.FindAllStringIndex(vp.GetContent(), -1))
vp.HighlightNext()
vp.HighlightPrevious()
vp.ClearHighlights()

// Per-line styling:
vp.StyleLineFunc = func(i int) lipgloss.Style { ... }

// Scroll:
vp.GotoTop() / vp.GotoBottom()
vp.ScrollPercent()          // 0.0–1.0
vp.SetYOffset(n)            // method, not field (v2 change)
```

---

## Help and Keybindings

```go
import "charm.land/bubbles/v2/help"
import "charm.land/bubbles/v2/key"

// Define keybindings:
type keyMap struct {
    Up   key.Binding
    Down key.Binding
    Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding  { return []key.Binding{k.Up, k.Down, k.Quit} }
func (k keyMap) FullHelp() [][]key.Binding { return [][]key.Binding{{k.Up, k.Down}, {k.Quit}} }

var keys = keyMap{
    Up:   key.NewBinding(key.WithKeys("k", "up"),   key.WithHelp("↑/k", "move up")),
    Down: key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "move down")),
    Quit: key.NewBinding(key.WithKeys("q"),         key.WithHelp("q", "quit")),
}

// Handle key presses:
case tea.KeyPressMsg:
    switch {
    case key.Matches(msg, keys.Up):   ...
    case key.Matches(msg, keys.Down): ...
    case key.Matches(msg, keys.Quit): return m, tea.Quit
    }

// Render help:
h := help.New()
h.Styles = help.DefaultStyles(isDark)   // isDark required (v2)
h.SetWidth(width)
view := h.View(keys)                    // auto-generates from ShortHelp/FullHelp
```

---

## Paginator

```go
import "charm.land/bubbles/v2/paginator"

p := paginator.New(
    paginator.WithPerPage(10),
    paginator.WithTotalPages(5),
)
p.Type = paginator.Dots    // Arabic (default) or Dots

// Or compute pages:
p.SetTotalPages(len(items))   // automatically sets TotalPages

// Get slice bounds for current page:
start, end := p.GetSliceBounds(len(items))
pageItems := items[start:end]

// KeyMap (function in v2):
km := paginator.DefaultKeyMap()
p.KeyMap = km

// In Update():
p, cmd = p.Update(msg)

// Navigation:
p.OnFirstPage() / p.OnLastPage()
p.PrevPage() / p.NextPage()
```

---

## Timer, Stopwatch, File Picker

```go
// Timer (countdown) — constructor changed in v2:
t := timer.New(30*time.Second)
t := timer.New(30*time.Second, timer.WithInterval(100*time.Millisecond))  // was NewWithInterval

// Stopwatch (count up) — constructor changed in v2:
sw := stopwatch.New()
sw := stopwatch.New(stopwatch.WithInterval(500*time.Millisecond))  // was NewWithInterval

// File picker — Height is now a method:
fp := filepicker.New()
fp.AllowedTypes = []string{".go", ".md"}
fp.SetHeight(10)   // was fp.Height = 10 in v1
```

See full API in [references/COMPONENTS.md](references/COMPONENTS.md).

---

## Embedding Components

Route messages to sub-components in Update; compose views in View:

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd
    var cmd tea.Cmd

    m.spinner, cmd = m.spinner.Update(msg)
    cmds = append(cmds, cmd)

    if m.input.Focused() {
        m.input, cmd = m.input.Update(msg)
        cmds = append(cmds, cmd)
    }

    return m, tea.Batch(cmds...)
}

func (m model) View() string {
    return lipgloss.JoinVertical(lipgloss.Left,
        m.spinner.View(),
        m.input.View(),
        m.help.View(m.keys),
    )
}
```
