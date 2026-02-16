# Bubbles v2 — Full Component API Reference

Import: `charm.land/bubbles/v2`

---

## spinner

```go
import "charm.land/bubbles/v2/spinner"
```

### Available spinners

```go
spinner.Line       // | / - \
spinner.Dot        // braille dots
spinner.MiniDot    // small braille
spinner.Jump       // jump
spinner.Pulse      // █ ▓ ▒ ░
spinner.Points     // ∙∙∙ ●∙∙ ...
spinner.Globe      // 🌍 🌎 🌏
spinner.Moon       // moon phases
spinner.Monkey     // 🙈 🙉 🙊
spinner.Meter      // ▱▱▱ → ▰▰▰
spinner.Hamburger  // ☱ ☲ ☴ ☲
spinner.Ellipsis   // "" . .. ...
```

### Types

```go
type Spinner struct {
    Frames []string
    FPS    time.Duration
}

type Model struct {
    Spinner Spinner
    Style   lipgloss.Style
}
```

### Functions

```go
func New(opts ...Option) Model
func WithSpinner(s Spinner) Option
func WithStyle(style lipgloss.Style) Option

func (m Model) ID() int
func (m Model) Tick() tea.Msg        // method in v2 (was package func in v1)
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
```

---

## textinput

```go
import "charm.land/bubbles/v2/textinput"
```

### Types

```go
type EchoMode int
const (
    EchoNormal   EchoMode = iota
    EchoPassword           // shows EchoCharacter mask
    EchoNone               // shows nothing
)

type Model struct {
    Err             error
    Prompt          string
    Placeholder     string
    EchoMode        EchoMode
    EchoCharacter   rune
    CharLimit       int
    KeyMap          KeyMap
    Validate        ValidateFunc
    ShowSuggestions bool
}

type Styles struct {
    Focused StyleState
    Blurred StyleState
    Cursor  CursorStyle
}

type StyleState struct {
    Text        lipgloss.Style
    Placeholder lipgloss.Style
    Suggestion  lipgloss.Style
    Prompt      lipgloss.Style
}

type CursorStyle struct {
    Color      color.Color
    Shape      tea.CursorShape  // CursorBlock | CursorUnderline | CursorBar
    Blink      bool
    BlinkSpeed time.Duration
}
```

### Functions

```go
func New() Model
func DefaultKeyMap() KeyMap       // function in v2, was variable in v1
func DefaultStyles(isDark bool) Styles
func DefaultDarkStyles() Styles
func DefaultLightStyles() Styles
func Blink() tea.Msg
func Paste() tea.Msg

// Methods
func (m *Model) Focus() tea.Cmd
func (m *Model) Blur()
func (m Model) Focused() bool
func (m *Model) SetWidth(w int)
func (m Model) Width() int
func (m *Model) SetValue(s string)
func (m Model) Value() string
func (m *Model) Reset()
func (m *Model) SetCursor(pos int)
func (m Model) Position() int
func (m *Model) CursorStart()
func (m *Model) CursorEnd()
func (m *Model) SetStyles(s Styles)
func (m Model) Styles() Styles
func (m *Model) SetVirtualCursor(v bool)
func (m Model) VirtualCursor() bool
func (m Model) Cursor() *tea.Cursor
func (m *Model) SetSuggestions(s []string)
func (m *Model) AvailableSuggestions() []string
func (m *Model) MatchedSuggestions() []string
func (m *Model) CurrentSuggestion() string
func (m *Model) CurrentSuggestionIndex() int
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
```

---

## textarea

```go
import "charm.land/bubbles/v2/textarea"
```

### Types

```go
type Model struct {
    Err                  error
    Prompt               string
    Placeholder          string
    ShowLineNumbers      bool
    EndOfBufferCharacter rune
    KeyMap               KeyMap
    CharLimit            int
    MaxHeight            int
    MaxWidth             int
}

type Styles struct {
    Focused StyleState
    Blurred StyleState
    Cursor  CursorStyle
}

type StyleState struct {
    Base             lipgloss.Style
    Text             lipgloss.Style
    LineNumber       lipgloss.Style
    CursorLineNumber lipgloss.Style
    CursorLine       lipgloss.Style
    EndOfBuffer      lipgloss.Style
    Placeholder      lipgloss.Style
    Prompt           lipgloss.Style
}
```

### Functions

```go
func New() Model
func DefaultKeyMap() KeyMap       // function in v2, was variable in v1
func DefaultStyles(isDark bool) Styles
func DefaultDarkStyles() Styles
func DefaultLightStyles() Styles
func Blink() tea.Msg
func Paste() tea.Msg

// Methods
func (m *Model) Focus() tea.Cmd
func (m *Model) Blur()
func (m Model) Focused() bool
func (m *Model) SetWidth(w int)
func (m Model) Width() int
func (m *Model) SetHeight(h int)
func (m Model) Height() int
func (m *Model) SetValue(s string)
func (m Model) Value() string
func (m *Model) Reset()
func (m Model) Line() int           // 0-indexed row
func (m Model) Column() int         // 0-indexed column (v2 addition)
func (m Model) LineCount() int
func (m *Model) Length() int
func (m *Model) InsertRune(r rune)
func (m *Model) InsertString(s string)
func (m Model) LineInfo() LineInfo
func (m *Model) CursorUp()
func (m *Model) CursorDown()
func (m *Model) CursorStart()
func (m *Model) CursorEnd()
func (m *Model) MoveToBegin()       // renamed from MoveToBeginning in v2
func (m *Model) MoveToEnd()         // renamed from MoveToEnd in v2
func (m *Model) PageUp()            // v2 addition
func (m *Model) PageDown()          // v2 addition
func (m *Model) SetCursorColumn(col int)  // renamed from SetCursor in v2
func (m *Model) SetStyles(s Styles)
func (m Model) Styles() Styles
func (m *Model) SetVirtualCursor(v bool)
func (m Model) VirtualCursor() bool
func (m Model) Cursor() *tea.Cursor
func (m Model) ScrollPercent() float64
func (m Model) ScrollYOffset() int       // v2 addition
func (m *Model) SetPromptFunc(promptWidth int, fn func(PromptInfo) string)
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
```

---

## list

```go
import "charm.land/bubbles/v2/list"
```

### Interfaces

```go
type Item interface {
    FilterValue() string
}

// For DefaultDelegate:
type DefaultItem interface {
    Item
    Title() string
    Description() string
}

type ItemDelegate interface {
    Render(w io.Writer, m Model, index int, item Item)
    Height() int
    Spacing() int
    Update(msg tea.Msg, m *Model) tea.Cmd
}
```

### Types

```go
type Model struct {
    Title             string
    Styles            Styles
    InfiniteScrolling bool
    KeyMap            KeyMap
    Filter            FilterFunc
    AdditionalShortHelpKeys func() []key.Binding
    AdditionalFullHelpKeys  func() []key.Binding
    Paginator         paginator.Model
    Help              help.Model
    FilterInput       textinput.Model
    StatusMessageLifetime time.Duration
}

type FilterState int
const (
    Unfiltered    FilterState = iota
    Filtering
    FilterApplied
)
```

### Functions

```go
func New(items []Item, delegate ItemDelegate, width, height int) Model
func DefaultKeyMap() KeyMap
func DefaultStyles(isDark bool) Styles          // isDark required (v2 change)
func NewDefaultDelegate() DefaultDelegate
func NewDefaultItemStyles(isDark bool) DefaultItemStyles  // isDark required (v2 change)
func DefaultFilter(term string, targets []string) []Rank
func UnsortedFilter(term string, targets []string) []Rank

// Methods
func (m *Model) SetSize(width, height int)
func (m *Model) SetWidth(v int)
func (m *Model) SetHeight(v int)
func (m Model) Width() int
func (m Model) Height() int
func (m *Model) SetItems(i []Item) tea.Cmd
func (m Model) Items() []Item
func (m Model) VisibleItems() []Item
func (m Model) SelectedItem() Item
func (m Model) Index() int
func (m Model) GlobalIndex() int
func (m Model) Cursor() int
func (m *Model) Select(index int)
func (m *Model) InsertItem(index int, item Item) tea.Cmd
func (m *Model) SetItem(index int, item Item) tea.Cmd
func (m *Model) RemoveItem(index int)
func (m *Model) SetDelegate(d ItemDelegate)
func (m *Model) CursorUp()
func (m *Model) CursorDown()
func (m *Model) GoToStart()
func (m *Model) GoToEnd()
func (m *Model) NextPage()
func (m *Model) PrevPage()
func (m *Model) ResetSelected()
func (m *Model) ResetFilter()
func (m Model) FilterState() FilterState
func (m Model) FilterValue() string
func (m Model) IsFiltered() bool
func (m *Model) SetFilterState(state FilterState)
func (m *Model) SetFilterText(filter string)
func (m *Model) SetFilteringEnabled(v bool)
func (m Model) FilteringEnabled() bool
func (m Model) SettingFilter() bool
func (m Model) MatchesForItem(index int) []int
func (m *Model) SetShowTitle(v bool)
func (m Model) ShowTitle() bool
func (m *Model) SetShowFilter(v bool)
func (m Model) ShowFilter() bool
func (m *Model) SetShowStatusBar(v bool)
func (m Model) ShowStatusBar() bool
func (m *Model) SetStatusBarItemName(singular, plural string)
func (m Model) StatusBarItemName() (string, string)
func (m *Model) SetShowPagination(v bool)
func (m *Model) ShowPagination() bool
func (m *Model) SetShowHelp(v bool)
func (m Model) ShowHelp() bool
func (m *Model) DisableQuitKeybindings()
func (m *Model) SetSpinner(s spinner.Spinner)
func (m *Model) StartSpinner() tea.Cmd
func (m *Model) StopSpinner()
func (m *Model) ToggleSpinner() tea.Cmd
func (m *Model) NewStatusMessage(s string) tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
```

---

## table

```go
import "charm.land/bubbles/v2/table"
```

### Types

```go
type Column struct {
    Title string
    Width int
}

type Row []string

type Styles struct {
    Header   lipgloss.Style
    Cell     lipgloss.Style
    Selected lipgloss.Style
}
```

### Functions

```go
func New(opts ...Option) Model
func DefaultKeyMap() KeyMap
func DefaultStyles() Styles     // no isDark needed (table unchanged in v2)
func WithColumns(cols []Column) Option
func WithRows(rows []Row) Option
func WithFocused(f bool) Option
func WithHeight(h int) Option
func WithWidth(w int) Option
func WithStyles(s Styles) Option
func WithKeyMap(km KeyMap) Option

// Methods
func (m *Model) Focus()
func (m *Model) Blur()
func (m Model) Focused() bool
func (m Model) Columns() []Column
func (m *Model) SetColumns(c []Column)
func (m Model) Rows() []Row
func (m *Model) SetRows(r []Row)
func (m Model) Cursor() int
func (m *Model) SetCursor(n int)
func (m Model) SelectedRow() Row
func (m *Model) MoveUp(n int)
func (m *Model) MoveDown(n int)
func (m *Model) GotoTop()
func (m *Model) GotoBottom()
func (m Model) Width() int
func (m *Model) SetWidth(w int)
func (m Model) Height() int
func (m *Model) SetHeight(h int)
func (m *Model) SetStyles(s Styles)
func (m *Model) UpdateViewport()
func (m Model) HelpView() string
func (m *Model) FromValues(value, separator string)
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
```

---

## progress

```go
import "charm.land/bubbles/v2/progress"
```

### Constants

```go
const (
    DefaultFullCharHalfBlock = '▌'
    DefaultFullCharFullBlock  = '█'
    DefaultEmptyCharBlock     = '░'
)
```

### Types

```go
type Model struct {
    Full            rune
    FullColor       color.Color
    Empty           rune
    EmptyColor      color.Color
    ShowPercentage  bool
    PercentFormat   string
    PercentageStyle lipgloss.Style
}

type ColorFunc func(total, current float64) color.Color
```

### Functions

```go
func New(opts ...Option) Model
func WithColors(colors ...color.Color) Option   // replaces WithGradient/WithSolidFill
func WithDefaultBlend() Option                  // replaces WithDefaultGradient
func WithScaled(enabled bool) Option            // replaces WithScaledGradient/WithDefaultScaledGradient
func WithColorFunc(fn ColorFunc) Option         // v2 addition: dynamic color
func WithFillCharacters(full, empty rune) Option
func WithoutPercentage() Option
func WithSpringOptions(frequency, damping float64) Option
func WithWidth(w int) Option

// Methods
func (m *Model) SetWidth(w int)
func (m Model) Width() int
func (m *Model) SetPercent(p float64) tea.Cmd
func (m Model) Percent() float64
func (m *Model) IncrPercent(v float64) tea.Cmd
func (m *Model) DecrPercent(v float64) tea.Cmd
func (m *Model) IsAnimating() bool
func (m *Model) SetSpringOptions(frequency, damping float64)
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
func (m Model) ViewAs(percent float64) string   // static render without animation
```

---

## viewport

```go
import "charm.land/bubbles/v2/viewport"
```

### Types

```go
type Model struct {
    KeyMap                 KeyMap
    SoftWrap               bool
    FillHeight             bool
    MouseWheelEnabled      bool
    MouseWheelDelta        int
    YPosition              int
    Style                  lipgloss.Style
    LeftGutterFunc         GutterFunc      // v2 addition
    HighlightStyle         lipgloss.Style  // v2 addition
    SelectedHighlightStyle lipgloss.Style  // v2 addition
    StyleLineFunc          func(int) lipgloss.Style  // v2 addition
}

type GutterContext struct {
    Index      int
    TotalLines int
    Soft       bool
}

type GutterFunc func(GutterContext) string

var NoGutter = func(GutterContext) string { return "" }
```

### Functions

```go
// Constructor changed in v2: was New(w, h int)
func New(opts ...Option) Model
func WithWidth(w int) Option
func WithHeight(h int) Option
func DefaultKeyMap() KeyMap

// Methods
func (m *Model) SetWidth(w int)
func (m Model) Width() int
func (m *Model) SetHeight(h int)
func (m Model) Height() int
func (m *Model) SetYOffset(n int)   // replaces direct field access
func (m *Model) YOffset() int
func (m *Model) SetXOffset(n int)
func (m *Model) XOffset() int
func (m *Model) SetContent(s string)
func (m *Model) SetContentLines(lines []string)  // v2 addition
func (m Model) GetContent() string               // v2 addition
func (m Model) TotalLineCount() int
func (m Model) VisibleLineCount() int
func (m Model) AtTop() bool
func (m Model) AtBottom() bool
func (m Model) PastBottom() bool
func (m Model) ScrollPercent() float64
func (m Model) HorizontalScrollPercent() float64  // v2 addition
func (m *Model) GotoTop() []string
func (m *Model) GotoBottom() []string
func (m *Model) ScrollUp(n int)
func (m *Model) ScrollDown(n int)
func (m *Model) ScrollLeft(n int)                 // v2 addition
func (m *Model) ScrollRight(n int)                // v2 addition
func (m *Model) PageUp()
func (m *Model) PageDown()
func (m *Model) HalfPageUp()
func (m *Model) HalfPageDown()
func (m *Model) SetHighlights(matches [][]int)    // v2 addition
func (m *Model) HighlightNext()                   // v2 addition
func (m *Model) HighlightPrevious()               // v2 addition
func (m *Model) ClearHighlights()                 // v2 addition
func (m *Model) EnsureVisible(line, colstart, colend int)  // v2 addition
func (m *Model) SetHorizontalStep(n int)          // v2 addition
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
```

---

## help

```go
import "charm.land/bubbles/v2/help"
```

### Interface

```go
type KeyMap interface {
    ShortHelp() []key.Binding
    FullHelp() [][]key.Binding
}
```

### Types

```go
type Model struct {
    ShowAll        bool
    ShortSeparator string
    FullSeparator  string
    Ellipsis       string
    Styles         Styles
}

type Styles struct {
    Ellipsis       lipgloss.Style
    ShortKey       lipgloss.Style
    ShortDesc      lipgloss.Style
    ShortSeparator lipgloss.Style
    FullKey        lipgloss.Style
    FullDesc       lipgloss.Style
    FullSeparator  lipgloss.Style
}
```

### Functions

```go
func New() Model                           // was NewModel() in v1
func DefaultStyles(isDark bool) Styles     // isDark required (v2 change)
func DefaultDarkStyles() Styles
func DefaultLightStyles() Styles

// Methods
func (m *Model) SetWidth(w int)
func (m Model) Width() int
func (m Model) View(k KeyMap) string
func (m Model) ShortHelpView(bindings []key.Binding) string
func (m Model) FullHelpView(groups [][]key.Binding) string
func (m Model) Update(_ tea.Msg) (Model, tea.Cmd)
```

---

## key

```go
import "charm.land/bubbles/v2/key"
```

### Types

```go
type Binding struct{ /* unexported */ }

type BindingOpt func(*Binding)

type Help struct {
    Key  string
    Desc string
}
```

### Functions

```go
func NewBinding(opts ...BindingOpt) Binding
func WithKeys(keys ...string) BindingOpt
func WithHelp(key, desc string) BindingOpt
func WithDisabled() BindingOpt

func Matches[Key fmt.Stringer](k Key, b ...Binding) bool

// Methods
func (b Binding) Enabled() bool
func (b *Binding) SetEnabled(v bool)
func (b Binding) Keys() []string
func (b *Binding) SetKeys(keys ...string)
func (b Binding) Help() Help
func (b *Binding) SetHelp(key, desc string)
func (b *Binding) Unbind()
```

---

## paginator

```go
import "charm.land/bubbles/v2/paginator"
```

### Types

```go
type Type int
const (
    Arabic Type = iota
    Dots
)

type Model struct {
    Type         Type
    Page         int
    PerPage      int
    TotalPages   int
    ActiveDot    string
    InactiveDot  string
    ArabicFormat string
    KeyMap       KeyMap
}
```

### Functions

```go
func New(opts ...Option) Model
func DefaultKeyMap() KeyMap    // function in v2, was variable in v1
func WithPerPage(perPage int) Option
func WithTotalPages(totalPages int) Option

// Methods
func (m *Model) SetTotalPages(items int) int
func (m Model) ItemsOnPage(totalItems int) int
func (m *Model) GetSliceBounds(length int) (start, end int)
func (m Model) OnFirstPage() bool
func (m Model) OnLastPage() bool
func (m *Model) PrevPage()
func (m *Model) NextPage()
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
```

---

## timer

```go
import "charm.land/bubbles/v2/timer"
```

### Types

```go
type TickMsg struct { ... }
type TimeoutMsg struct { ... }

type Option func(*Model)
```

### Functions

```go
// Constructor changed in v2: was NewWithInterval(timeout, interval)
func New(timeout time.Duration, opts ...Option) Model
func WithInterval(d time.Duration) Option

// Methods
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
func (m Model) Running() bool
func (m Model) Timedout() bool
func (m *Model) Start() tea.Cmd
func (m *Model) Stop() tea.Cmd
func (m *Model) Toggle() tea.Cmd
```

---

## stopwatch

```go
import "charm.land/bubbles/v2/stopwatch"
```

### Functions

```go
// Constructor changed in v2: was NewWithInterval(d)
func New(opts ...Option) Model
func WithInterval(d time.Duration) Option

// Methods
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
func (m Model) Running() bool
func (m Model) Elapsed() time.Duration
func (m *Model) Start() tea.Cmd
func (m *Model) Stop() tea.Cmd
func (m *Model) Reset() tea.Cmd
func (m *Model) Toggle() tea.Cmd
```

---

## filepicker

```go
import "charm.land/bubbles/v2/filepicker"
```

### Functions

```go
func New() Model
func DefaultStyles() Styles    // was DefaultStylesWithRenderer(r) in v1

// Methods
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
func (m Model) DidSelectFile(msg tea.Msg) (bool, string)
func (m Model) DidSelectDisabledFile(msg tea.Msg) (bool, string)
func (m *Model) SetHeight(h int)   // method in v2, was field in v1
func (m Model) Height() int
```

### Key Model Fields

```go
type Model struct {
    CurrentDirectory string
    AllowedTypes     []string   // e.g. []string{".go", ".md"}
    DirAllowed       bool
    FileAllowed      bool
    ShowHidden       bool
    AutoHeight       bool
    Styles           Styles
    KeyMap           KeyMap
}
```

---

## cursor

```go
import "charm.land/bubbles/v2/cursor"
```

> Usually embedded within textinput/textarea rather than used directly.

### Changes from v1

| v1 | v2 |
|---|---|
| `model.Blink` | `model.IsBlinked` |
| `model.BlinkCmd()` | `model.Blink()` |

### Functions

```go
func New(opts ...Option) Model
func NewWithMode(m Mode) Model

// Methods
func (m Model) IsBlinked() bool      // renamed from Blink in v1
func (m Model) Blink() tea.Msg       // renamed from BlinkCmd in v1
func (m Model) Focus() tea.Cmd
func (m Model) Blur()
func (m Model) Focused() bool
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View(displayValue string) string
```
