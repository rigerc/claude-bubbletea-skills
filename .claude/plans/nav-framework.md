# Plan: BubbleTea v2 Navigation/Routing Framework

## Context

`template-v2-enhanced/` is a minimal BubbleTea v2 skeleton with logging (zerolog), CLI (cobra),
and config (koanf). The `internal/ui/model.go` is a bare 4-field model with no navigation.
Previously deleted packages (`nav/`, `screens/`, `styles/`, `keys/`) are being replaced with a
clean, extensible design. The goal is a stack-based router driven by `charm.land/bubbles/v2/list`
as the primary navigation component, with adaptive lipgloss-v2 styling and the bubbles `help`/`key`
system for keybindings.

## Dependencies to Add

```sh
go get charm.land/bubbles/v2
go get charm.land/lipgloss/v2
```

Add to the `require` block in `template-v2-enhanced/go.mod`. These must match versions confirmed
working together (see examples go.mod for version pins). `go.sum` updates automatically.

---

## File Structure

```
internal/ui/
├── model.go              [MODIFY]  — Root router; owns stack, help, AltScreen
├── nav/
│   └── nav.go           [CREATE]  — Screen interface + navigation messages
├── styles/
│   └── styles.go        [CREATE]  — Adaptive lipgloss styles
├── keys/
│   └── keys.go          [CREATE]  — Global key bindings
└── screens/
    ├── menu.go          [CREATE]  — List-based navigation menu
    └── detail.go        [CREATE]  — Example scrollable content screen
```

**Build order:** `keys` → `nav` → `styles` → `screens/detail` → `screens/menu` → `model`

---

## File 1: `internal/ui/nav/nav.go`

**Package:** `nav`

### Types

```go
// Screen is implemented by every navigable view.
type Screen interface {
    Init() tea.Cmd
    Update(tea.Msg) (Screen, tea.Cmd)  // returns Screen, not tea.Model
    View() string                       // string, not tea.View — Model.View() builds the final View
}

// Themeable is optional. Router calls SetTheme on push and on BackgroundColorMsg.
type Themeable interface {
    SetTheme(isDark bool)
}

// ScreenKeyMap is optional. Router uses it to show screen-specific help keys.
type ScreenKeyMap interface {
    HelpKeys() []key.Binding
}

// Navigation message types:
type PushMsg    struct{ Screen Screen }
type PopMsg     struct{}
type ReplaceMsg struct{ Screen Screen }
```

### Functions

```go
// Push returns a Cmd that sends PushMsg.
func Push(s Screen) tea.Cmd { return func() tea.Msg { return PushMsg{Screen: s} } }

// Pop returns a Cmd that sends PopMsg.
func Pop() tea.Cmd { return func() tea.Msg { return PopMsg{} } }

// Replace returns a Cmd that sends ReplaceMsg.
func Replace(s Screen) tea.Cmd { return func() tea.Msg { return ReplaceMsg{Screen: s} } }
```

---

## File 2: `internal/ui/keys/keys.go`

**Package:** `keys`

Implements `help.KeyMap` interface so it can be passed directly to `help.View()`.

```go
type GlobalKeyMap struct {
    Back key.Binding  // "esc"       — go to previous screen
    Quit key.Binding  // "ctrl+c"    — always quit (avoids conflict with list filter "q")
    Help key.Binding  // "?"         — toggle help expansion
}

func New() GlobalKeyMap

// Implements help.KeyMap:
func (k GlobalKeyMap) ShortHelp() []key.Binding  // [Back, Quit]
func (k GlobalKeyMap) FullHelp() [][]key.Binding // [[Back, Help], [Quit]]
```

**Critical:** `"q"` is intentionally NOT in `Quit`. The list uses `"q"` for filtering when
`DisableQuitKeybindings()` is not called. After calling it, `"q"` types into filter input.
Only `ctrl+c` is the guaranteed quit path.

---

## File 3: `internal/ui/styles/styles.go`

**Package:** `styles`

```go
type Theme struct {
    App           lipgloss.Style   // outer container: Margin(1, 2)
    Title         lipgloss.Style   // list title bar (overrides list default)
    StatusMessage lipgloss.Style   // list status bar message
    HelpBar       lipgloss.Style   // help bar at bottom
    Detail        lipgloss.Style   // body text in detail screens
    Subtle        lipgloss.Style   // de-emphasized text
}

// New constructs a Theme using lipgloss.LightDark(isDark) for adaptive colors.
func New(isDark bool) Theme

// HelpBarHeight returns lines consumed by the help bar.
// Short mode = 1 line, full mode = 3 lines.
func HelpBarHeight(showAll bool) int
```

**Implementation note for `New`:**
```go
ld := lipgloss.LightDark(isDark)
Theme{
    App:   lipgloss.NewStyle().Margin(1, 2),
    Title: lipgloss.NewStyle().Bold(true).
               Foreground(lipgloss.Color("#FFFDF5")).
               Background(lipgloss.Color("#25A065")).
               Padding(0, 1),
    StatusMessage: lipgloss.NewStyle().
               Foreground(ld(lipgloss.Color("#04B575"), lipgloss.Color("#10CC85"))),
    HelpBar: lipgloss.NewStyle().
               Foreground(ld(lipgloss.Color("#626262"), lipgloss.Color("#9B9B9B"))),
    ...
}
```

---

## File 4: `internal/ui/screens/menu.go`

**Package:** `screens`

### MenuItem

```go
type MenuItem struct {
    title, description string
    action             tea.Cmd  // typically nav.Push(NewDetailScreen(...))
}

func (i MenuItem) Title() string       { return i.title }
func (i MenuItem) Description() string { return i.description }
func (i MenuItem) FilterValue() string { return i.title }

func NewMenuItem(title, description string, action tea.Cmd) MenuItem
```

`action` is a `tea.Cmd` capturing the target screen. When `enter` is pressed, the delegate
returns `item.action` directly. The router interprets the resulting `nav.PushMsg`.

### delegateKeyMap

```go
type delegateKeyMap struct {
    choose key.Binding  // "enter" — select item
}

func newDelegateKeyMap() delegateKeyMap
```

### MenuScreen struct

```go
type MenuScreen struct {
    list         list.Model
    delegateKeys delegateKeyMap
    keys         keys.GlobalKeyMap
    theme        styles.Theme
    isDark       bool
    width, height int
    helpExpanded  bool
}
```

### MenuScreen methods

```go
// NewMenuScreen creates the screen. Pass 0,0 for width/height —
// WindowSizeMsg will call updateListSize() before first render.
func NewMenuScreen(title string, items []list.Item, isDark bool) *MenuScreen

// Critical initialization sequence in NewMenuScreen:
//   l := list.New(items, delegate, 0, 0)
//   l.Title = title
//   l.Styles = list.DefaultStyles(isDark)
//   l.Styles.Title = theme.Title          // override with branded title
//   l.DisableQuitKeybindings()             // REQUIRED: prevents list from eating "ctrl+c"/"q"
//   l.SetShowHelp(false)                   // disable list's own help — we render it in Model

func (s *MenuScreen) Init() tea.Cmd                              // returns nil
func (s *MenuScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd)
func (s *MenuScreen) View() string                               // s.theme.App.Render(s.list.View())

// Implements nav.Themeable:
func (s *MenuScreen) SetTheme(isDark bool)

// Implements nav.ScreenKeyMap:
func (s *MenuScreen) HelpKeys() []key.Binding  // [delegateKeys.choose]

// Dynamic list items:
func (s *MenuScreen) SetItems(items []list.Item) tea.Cmd     // list.SetItems returns Cmd
func (s *MenuScreen) InsertItem(index int, item list.Item) tea.Cmd
func (s *MenuScreen) RemoveItem(index int)

// private:
func (s *MenuScreen) updateListSize()
```

### MenuScreen.Update — ESC/Back/Filter conflict resolution

```go
func (s *MenuScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.WindowSizeMsg:
        s.width, s.height = msg.Width, msg.Height
        s.updateListSize()

    case tea.BackgroundColorMsg:
        s.isDark = msg.IsDark()
        s.theme = styles.New(s.isDark)
        s.list.Styles = list.DefaultStyles(s.isDark)
        s.list.Styles.Title = s.theme.Title
        // Recreate delegate with updated styles:
        d := list.NewDefaultDelegate()
        d.Styles = list.NewDefaultItemStyles(s.isDark)
        d.UpdateFunc = ...   // same delegate logic, new styles
        s.list.SetDelegate(d)
        return s, nil

    case tea.KeyPressMsg:
        // ESC: pop stack only when NOT filtering
        if msg.String() == "esc" {
            switch s.list.FilterState() {
            case list.Unfiltered:
                return s, nav.Pop()
            // list.Filtering, list.FilterApplied: fall through to list.Update
            }
        }
    }

    var cmd tea.Cmd
    s.list, cmd = s.list.Update(msg)
    return s, cmd
}
```

### updateListSize — dynamic list height

```go
func (s *MenuScreen) updateListSize() {
    frameH, frameV := s.theme.App.GetFrameSize()
    helpH := styles.HelpBarHeight(s.helpExpanded)
    s.list.SetSize(s.width-frameH, s.height-frameV-helpH)
}
```

**This is the core of dynamic height support.** Called on every `WindowSizeMsg` and whenever
`helpExpanded` toggles. The help bar's vertical footprint is subtracted from the list height.

### Delegate

```go
// newMenuDelegate creates a list.DefaultDelegate with:
// - UpdateFunc: on "enter", returns item.action (the navigation Cmd)
// - ShortHelpFunc: [choose binding]
// - FullHelpFunc: [[choose binding]]
func newMenuDelegate(dKeys delegateKeyMap, isDark bool) list.DefaultDelegate {
    d := list.NewDefaultDelegate()
    d.Styles = list.NewDefaultItemStyles(isDark)
    d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
        if kp, ok := msg.(tea.KeyPressMsg); ok {
            if key.Matches(kp, dKeys.choose) {
                if item, ok := m.SelectedItem().(MenuItem); ok && item.action != nil {
                    return item.action
                }
            }
        }
        return nil
    }
    d.ShortHelpFunc = func() []key.Binding { return []key.Binding{dKeys.choose} }
    d.FullHelpFunc  = func() [][]key.Binding { return [][]key.Binding{{dKeys.choose}} }
    return d
}
```

---

## File 5: `internal/ui/screens/detail.go`

**Package:** `screens`

```go
type DetailScreen struct {
    title, content string
    keys           keys.GlobalKeyMap
    theme          styles.Theme
    isDark         bool
    width, height  int
    vp             viewport.Model
    ready          bool   // false until first WindowSizeMsg
}

func NewDetailScreen(title, content string, isDark bool) *DetailScreen

func (s *DetailScreen) Init() tea.Cmd  // nil
func (s *DetailScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd)
func (s *DetailScreen) View() string

// Implements nav.Themeable:
func (s *DetailScreen) SetTheme(isDark bool)

func (s *DetailScreen) SetContent(content string)  // update viewport content
func (s *DetailScreen) headerView() string          // private: renders title bar
```

### DetailScreen.Update

```go
case tea.WindowSizeMsg:
    s.width, s.height = msg.Width, msg.Height
    headerH := lipgloss.Height(s.headerView())
    frameH, frameV := s.theme.App.GetFrameSize()
    helpH  := styles.HelpBarHeight(false)
    s.vp.SetWidth(s.width - frameH)
    s.vp.SetHeight(s.height - frameV - headerH - helpH)
    if !s.ready {
        s.vp.SetContent(s.content)
        s.ready = true
    }

case tea.KeyPressMsg:
    if key.Matches(msg, s.keys.Back) {
        return s, nav.Pop()   // always pop — no filter state to check
    }
```

Viewport receives remaining keys for scroll (j/k, PageUp/Down, etc.).

```go
// Use viewport.New() with options in v2:
s.vp = viewport.New()
s.vp.MouseWheelEnabled = true
```

---

## File 6: `internal/ui/model.go` (Rewrite)

**Package:** `ui`

### Model struct

```go
type Model struct {
    screens      []nav.Screen
    width, height int
    isDark        bool
    quitting      bool
    help          help.Model
    keys          keys.GlobalKeyMap
    helpExpanded  bool
}
```

### New()

```go
// New accepts the loaded config so UI behaviour respects config.UI.AltScreen,
// config.UI.MouseEnabled, and config.App.Title.
func New(cfg *config.Config) Model {
    globalKeys := keys.New()

    items := []list.Item{
        screens.NewMenuItem("Details", "View a detail screen",
            nav.Push(screens.NewDetailScreen("Details", detailContent, false))),
        screens.NewMenuItem("About", "About this application",
            nav.Push(screens.NewDetailScreen("About", aboutContent, false))),
    }
    root := screens.NewMenuScreen(cfg.App.Title, items, false)

    h := help.New()
    h.Styles = help.DefaultStyles(false)  // corrected on BackgroundColorMsg

    return Model{
        screens:      []nav.Screen{root},
        help:         h,
        keys:         globalKeys,
        altScreen:    cfg.UI.AltScreen,
        mouseEnabled: cfg.UI.MouseEnabled,
        windowTitle:  cfg.App.Title,
    }
}
```

### Init()

```go
func (m Model) Init() tea.Cmd {
    cmds := []tea.Cmd{tea.RequestBackgroundColor}
    if len(m.screens) > 0 {
        cmds = append(cmds, m.screens[len(m.screens)-1].Init())
    }
    return tea.Batch(cmds...)
}
```

### Update() — routing logic

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {

    case tea.KeyPressMsg:
        if key.Matches(msg, m.keys.Quit) {
            m.quitting = true
            return m, tea.Quit
        }
        if key.Matches(msg, m.keys.Help) {
            m.helpExpanded = !m.helpExpanded
            m.help.ShowAll = m.helpExpanded
            m.notifyActiveScreenResize()   // re-sends WindowSizeMsg to active screen
            return m, nil
        }

    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        m.help.Width = msg.Width   // Note: Width field, not a method (bubbles v2 help)
        // fall through to delegate to active screen

    case tea.BackgroundColorMsg:
        m.isDark = msg.IsDark()
        m.help.Styles = help.DefaultStyles(m.isDark)
        // Propagate theme to ALL screens in stack (not just top)
        for i := range m.screens {
            if t, ok := m.screens[i].(nav.Themeable); ok {
                t.SetTheme(m.isDark)
            }
        }
        // fall through to deliver msg to active screen too

    case nav.PushMsg:
        s := msg.Screen
        if cmd := s.Init(); cmd != nil {
            cmds = append(cmds, cmd)
        }
        // Inject current theme and dimensions immediately
        if t, ok := s.(nav.Themeable); ok {
            t.SetTheme(m.isDark)
        }
        s, cmd := s.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
        cmds = append(cmds, cmd)
        m.screens = append(m.screens, s)
        return m, tea.Batch(cmds...)

    case nav.PopMsg:
        if len(m.screens) > 1 {
            m.screens = m.screens[:len(m.screens)-1]
        }
        return m, nil

    case nav.ReplaceMsg:
        if len(m.screens) > 0 {
            s := msg.Screen
            if cmd := s.Init(); cmd != nil { cmds = append(cmds, cmd) }
            if t, ok := s.(nav.Themeable); ok { t.SetTheme(m.isDark) }
            m.screens[len(m.screens)-1] = s
        }
        return m, tea.Batch(cmds...)
    }

    // Delegate to active screen
    if len(m.screens) > 0 {
        top := m.screens[len(m.screens)-1]
        updated, cmd := top.Update(msg)
        m.screens[len(m.screens)-1] = updated
        cmds = append(cmds, cmd)
    }

    return m, tea.Batch(cmds...)
}
```

### View() — final assembly

```go
func (m Model) View() tea.View {
    if m.quitting {
        return tea.NewView("")
    }

    var screenContent string
    if len(m.screens) > 0 {
        screenContent = m.screens[len(m.screens)-1].View()
    }

    helpView := m.help.View(m.buildKeyMap())

    content := lipgloss.JoinVertical(lipgloss.Left, screenContent, helpView)

    v := tea.NewView(content)
    v.AltScreen   = m.altScreen    // from config.UI.AltScreen (default: false)
    v.WindowTitle = m.windowTitle  // from config.App.Title
    if m.mouseEnabled {            // from config.UI.MouseEnabled (default: true)
        v.MouseMode = tea.MouseModeCellMotion
    }
    return v
}

// buildKeyMap merges global keys with active screen's keys for the help bar.
func (m Model) buildKeyMap() help.KeyMap {
    var extra []key.Binding
    if len(m.screens) > 0 {
        if sk, ok := m.screens[len(m.screens)-1].(nav.ScreenKeyMap); ok {
            extra = sk.HelpKeys()
        }
    }
    return combinedKeyMap{global: m.keys, extra: extra}
}

// notifyActiveScreenResize re-sends WindowSizeMsg to the active screen
// so it can recalculate its component height after help bar changes.
func (m *Model) notifyActiveScreenResize() {
    if len(m.screens) == 0 { return }
    top := m.screens[len(m.screens)-1]
    updated, _ := top.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
    m.screens[len(m.screens)-1] = updated
}

// combinedKeyMap merges GlobalKeyMap + active screen's extra bindings.
// Unexported type, local to model.go.
type combinedKeyMap struct {
    global keys.GlobalKeyMap
    extra  []key.Binding
}
func (c combinedKeyMap) ShortHelp() []key.Binding { return append(c.global.ShortHelp(), c.extra...) }
func (c combinedKeyMap) FullHelp() [][]key.Binding {
    full := c.global.FullHelp()
    if len(c.extra) > 0 { full = append(full, c.extra) }
    return full
}
```

---

## Critical Implementation Notes

| Issue | Resolution |
|---|---|
| `list.DisableQuitKeybindings()` | MUST be called in `NewMenuScreen` or list will intercept `ctrl+c`/`q` |
| `l.SetShowHelp(false)` | MUST be set so the list doesn't render its own help bar underneath ours |
| `list.New(items, delegate, 0, 0)` | Always 0,0 — `WindowSizeMsg` drives sizing via `updateListSize()` |
| `list.SetItems` returns `tea.Cmd` | Never ignore this return — it contains a status message animation |
| `help.Width` is a field, not method | `m.help.Width = msg.Width`, not `m.help.SetWidth()` |
| `list.DefaultStyles(isDark bool)` | Must be called with correct `isDark`; zero value `false` OK as initial |
| `list.NewDefaultItemStyles(isDark)` | Same — update in `SetTheme()` by recreating delegate |
| ESC in menu | Check `list.FilterState() == list.Unfiltered` before calling `nav.Pop()` |
| Theme on new screens | Inject via `nav.Themeable.SetTheme()` + send `WindowSizeMsg` immediately on push |
| Viewport in v2 | `viewport.New()` with no positional args (changed from v1) |

---

## Verification

```sh
# 1. Add dependencies
cd template-v2-enhanced
go get charm.land/bubbles/v2 charm.land/lipgloss/v2

# 2. Build all packages
go build ./...

# 3. Run the app
go run .

# Manual tests:
# a) Arrow keys / j/k navigate the list
# b) Enter on "Details" pushes DetailScreen
# c) ESC from DetailScreen pops back to menu
# d) "/" enters filter mode; ESC exits filter (does NOT pop)
# e) "?" toggles full help — list height adjusts
# f) Resize terminal — list height adjusts dynamically
# g) ctrl+c quits from any screen

# 4. Run with debug logging
go run . --debug
tail -f debug.log
```

## Config Integration (AltScreen + Mouse + Window Title)

`config.UIConfig` has `AltScreen bool`, `MouseEnabled bool`, and `ThemeName string` fields.
`config.AppConfig` has `Title string` (window title). These must flow into `Model`.

### Change: `main.go` loads config and passes it to `ui.New()`

```go
// In main(), after initLogger() and before ui.Run():
cfg := loadConfig()   // new helper in main.go
if err := ui.Run(ui.New(cfg)); err != nil { ... }

// loadConfig tries cfgFile path, falls back to defaults (no error on missing file)
func loadConfig() *config.Config {
    if path := cmd.GetConfigFile(); path != "" {
        if c, err := config.Load(path); err == nil {
            return c
        }
    }
    return config.DefaultConfig()
}
```

### Change: `ui.New()` signature

```go
// New creates a Model configured from cfg.
func New(cfg *config.Config) Model {
    // ...
    return Model{
        screens:      []nav.Screen{root},
        help:         h,
        keys:         globalKeys,
        altScreen:    cfg.UI.AltScreen,
        mouseEnabled: cfg.UI.MouseEnabled,
        windowTitle:  cfg.App.Title,
    }
}
```

### Change: `Model` struct adds UI config fields

```go
type Model struct {
    screens      []nav.Screen
    width, height int
    isDark        bool
    quitting      bool
    help          help.Model
    keys          keys.GlobalKeyMap
    helpExpanded  bool
    // from config:
    altScreen    bool   // config.UI.AltScreen
    mouseEnabled bool   // config.UI.MouseEnabled
    windowTitle  string // config.App.Title
}
```

### Change: `Model.View()` uses config values

```go
func (m Model) View() tea.View {
    if m.quitting { return tea.NewView("") }
    // ... assemble content ...
    v := tea.NewView(content)
    v.AltScreen   = m.altScreen    // NOT hardcoded true
    v.WindowTitle = m.windowTitle
    if m.mouseEnabled {
        v.MouseMode = tea.MouseModeCellMotion
    }
    return v
}
```

---

## Files to Modify vs Create

| File | Action | Description |
|------|--------|-------------|
| `go.mod` | MODIFY | Add bubbles/v2 and lipgloss/v2 |
| `go.sum` | MODIFY | Updated by `go get` |
| `main.go` | MODIFY | Add `loadConfig()`, pass cfg to `ui.New(cfg)` |
| `internal/ui/model.go` | REWRITE | Stack router replaces skeleton |
| `internal/ui/nav/nav.go` | CREATE | Screen interface + nav messages |
| `internal/ui/styles/styles.go` | CREATE | Adaptive lipgloss theme |
| `internal/ui/keys/keys.go` | CREATE | Global key bindings (help.KeyMap) |
| `internal/ui/screens/menu.go` | CREATE | List-based navigation menu |
| `internal/ui/screens/detail.go` | CREATE | Scrollable text detail screen |

`cmd/`, `config/`, `internal/logger/` are **unchanged**.
