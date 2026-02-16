# Bubbles v2 — Full Component API Reference

Import: `charm.land/bubbles/v2`

---

## spinner

```go
import "charm.land/bubbles/v2/spinner"
```

### Constants

Pre-defined spinner animations:

```go
spinner.Line       // | / - \  (10 FPS)
spinner.Dot        // braille dots ⣾ ⣽ ⣻ ⢿ ⡿ ⣟ ⣯ ⣷ (10 FPS)
spinner.MiniDot    // small braille ⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏ (12 FPS)
spinner.Jump       // ⢄ ⢂ ⢁ ⡁ ⡈ ⡐ ⡠ (10 FPS)
spinner.Pulse      // █ ▓ ▒ ░ (8 FPS)
spinner.Points     // ∙∙∙ ●∙∙ ∙●∙ ∙∙● (7 FPS)
spinner.Globe      // 🌍 🌎 🌏 (4 FPS)
spinner.Moon       // 🌑 🌒 🌓 🌔 🌕 🌖 🌗 🌘 (8 FPS)
spinner.Monkey     // 🙈 🙉 🙊 (3 FPS)
spinner.Meter      // ▱▱▱ ▰▱▱ ▰▰▱ ▰▰▰ (7 FPS)
spinner.Hamburger  // ☱ ☲ ☴ ☲ (3 FPS)
spinner.Ellipsis   // "" . .. ... (3 FPS)
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
    // Has unexported fields.
}

type TickMsg struct {
    Time time.Time
    ID   int
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
    EchoNormal   EchoMode = iota  // displays text as is (default)
    EchoPassword                  // shows EchoCharacter mask
    EchoNone                      // displays nothing
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
    // Has unexported fields.
}

type ValidateFunc func(string) error

type KeyMap struct {
    CharacterForward        key.Binding
    CharacterBackward       key.Binding
    WordForward             key.Binding
    WordBackward            key.Binding
    DeleteWordBackward      key.Binding
    DeleteWordForward       key.Binding
    DeleteAfterCursor       key.Binding
    DeleteBeforeCursor      key.Binding
    DeleteCharacterBackward key.Binding
    DeleteCharacterForward  key.Binding
    LineStart               key.Binding
    LineEnd                 key.Binding
    Paste                   key.Binding
    AcceptSuggestion        key.Binding
    NextSuggestion          key.Binding
    PrevSuggestion          key.Binding
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
func Blink() tea.Msg
func Paste() tea.Msg
func DefaultKeyMap() KeyMap       // function in v2, was variable in v1
func DefaultStyles(isDark bool) Styles
func DefaultDarkStyles() Styles
func DefaultLightStyles() Styles

// Lifecycle
func (m *Model) Focus() tea.Cmd
func (m *Model) Blur()
func (m Model) Focused() bool

// Dimensions
func (m *Model) SetWidth(w int)
func (m Model) Width() int

// Value
func (m *Model) SetValue(s string)
func (m Model) Value() string
func (m *Model) Reset()

// Cursor
func (m *Model) SetCursor(pos int)
func (m Model) Position() int
func (m *Model) CursorStart()
func (m *Model) CursorEnd()
func (m *Model) SetVirtualCursor(v bool)
func (m Model) VirtualCursor() bool
func (m Model) Cursor() *tea.Cursor

// Suggestions
func (m *Model) SetSuggestions(s []string)
func (m Model) AvailableSuggestions() []string
func (m Model) MatchedSuggestions() []string
func (m Model) CurrentSuggestion() string
func (m Model) CurrentSuggestionIndex() int

// Styling
func (m *Model) SetStyles(s Styles)
func (m Model) Styles() Styles

// Update loop
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
    // Has unexported fields.
}

type KeyMap struct {
    CharacterBackward       key.Binding
    CharacterForward        key.Binding
    DeleteAfterCursor       key.Binding
    DeleteBeforeCursor      key.Binding
    DeleteCharacterBackward key.Binding
    DeleteCharacterForward  key.Binding
    DeleteWordBackward      key.Binding
    DeleteWordForward       key.Binding
    InsertNewline           key.Binding
    LineEnd                 key.Binding
    LineNext                key.Binding
    LinePrevious            key.Binding
    LineStart               key.Binding
    PageUp                  key.Binding
    PageDown                key.Binding
    Paste                   key.Binding
    WordBackward            key.Binding
    WordForward             key.Binding
    InputBegin              key.Binding
    InputEnd                key.Binding
    UppercaseWordForward    key.Binding
    LowercaseWordForward    key.Binding
    CapitalizeWordForward   key.Binding
    TransposeCharacterBackward key.Binding
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

type CursorStyle struct {
    Color      color.Color
    Shape      tea.CursorShape
    Blink      bool
    BlinkSpeed time.Duration
}

type LineInfo struct {
    Width        int  // columns in line
    CharWidth    int  // characters (accounting for double-width runes)
    Height       int  // rows in line
    StartColumn  int
    ColumnOffset int
    RowOffset    int
    CharOffset   int
}

type PromptInfo struct {
    LineNumber int
    Focused    bool
}
```

### Functions

```go
func New() Model
func Blink() tea.Msg
func Paste() tea.Msg
func DefaultKeyMap() KeyMap
func DefaultStyles(isDark bool) Styles
func DefaultDarkStyles() Styles
func DefaultLightStyles() Styles

// Lifecycle
func (m *Model) Focus() tea.Cmd
func (m *Model) Blur()
func (m Model) Focused() bool

// Dimensions
func (m *Model) SetWidth(w int)
func (m Model) Width() int
func (m *Model) SetHeight(h int)
func (m Model) Height() int

// Value
func (m *Model) SetValue(s string)
func (m Model) Value() string
func (m *Model) Reset()
func (m Model) Length() int
func (m *Model) InsertRune(r rune)
func (m *Model) InsertString(s string)

// Cursor / Navigation
func (m Model) Line() int              // 0-indexed row
func (m Model) Column() int            // 0-indexed column
func (m Model) LineCount() int
func (m Model) LineInfo() LineInfo
func (m *Model) CursorUp()
func (m *Model) CursorDown()
func (m *Model) CursorStart()
func (m *Model) CursorEnd()
func (m *Model) MoveToBegin()          // renamed from MoveToBeginning in v2
func (m *Model) MoveToEnd()
func (m *Model) PageUp()               // v2 addition
func (m *Model) PageDown()             // v2 addition
func (m *Model) SetCursorColumn(col int)

// Word
func (m *Model) Word() string          // word at cursor position

// Scroll
func (m Model) ScrollPercent() float64
func (m Model) ScrollYOffset() int     // v2 addition

// Prompt
func (m *Model) SetPromptFunc(promptWidth int, fn func(PromptInfo) string)

// Styling
func (m *Model) SetStyles(s Styles)
func (m Model) Styles() Styles
func (m *Model) SetVirtualCursor(v bool)
func (m Model) VirtualCursor() bool
func (m Model) Cursor() *tea.Cursor

// Update loop
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
    // Has unexported fields.
}

type FilterState int
const (
    Unfiltered    FilterState = iota
    Filtering
    FilterApplied
)

type KeyMap struct {
    CursorUp    key.Binding
    CursorDown  key.Binding
    NextPage    key.Binding
    PrevPage    key.Binding
    GoToStart   key.Binding
    GoToEnd     key.Binding
    Filter      key.Binding
    ClearFilter key.Binding
    CancelWhileFiltering key.Binding
    AcceptWhileFiltering key.Binding
    ShowFullHelp  key.Binding
    CloseFullHelp key.Binding
    Quit key.Binding
    ForceQuit key.Binding
}

type Styles struct {
    TitleBar lipgloss.Style
    Title    lipgloss.Style
    Spinner  lipgloss.Style
    Filter   textinput.Styles
    DefaultFilterCharacterMatch lipgloss.Style
    StatusBar             lipgloss.Style
    StatusEmpty           lipgloss.Style
    StatusBarActiveFilter lipgloss.Style
    StatusBarFilterCount  lipgloss.Style
    NoItems lipgloss.Style
    PaginationStyle lipgloss.Style
    HelpStyle       lipgloss.Style
    ActivePaginationDot   lipgloss.Style
    InactivePaginationDot lipgloss.Style
    ArabicPagination      lipgloss.Style
    DividerDot            lipgloss.Style
}

type DefaultDelegate struct {
    ShowDescription bool
    Styles          DefaultItemStyles
    UpdateFunc      func(tea.Msg, *Model) tea.Cmd
    ShortHelpFunc   func() []key.Binding
    FullHelpFunc    func() [][]key.Binding
}

type DefaultItemStyles struct {
    NormalTitle lipgloss.Style
    NormalDesc  lipgloss.Style
    SelectedTitle lipgloss.Style
    SelectedDesc  lipgloss.Style
    DimmedTitle lipgloss.Style
    DimmedDesc  lipgloss.Style
    FilterMatch lipgloss.Style
}

type Rank struct {
    Index          int
    MatchedIndexes []int
}

type FilterMatchesMsg []filteredItem
```

### Functions

```go
func New(items []Item, delegate ItemDelegate, width, height int) Model
func DefaultKeyMap() KeyMap
func DefaultStyles(isDark bool) Styles
func NewDefaultDelegate() DefaultDelegate
func NewDefaultItemStyles(isDark bool) DefaultItemStyles
func DefaultFilter(term string, targets []string) []Rank
func UnsortedFilter(term string, targets []string) []Rank

// Dimensions
func (m *Model) SetSize(width, height int)
func (m *Model) SetWidth(v int)
func (m *Model) SetHeight(v int)
func (m Model) Width() int
func (m Model) Height() int

// Items
func (m *Model) SetItems(i []Item) tea.Cmd
func (m Model) Items() []Item
func (m Model) VisibleItems() []Item
func (m *Model) InsertItem(index int, item Item) tea.Cmd
func (m *Model) SetItem(index int, item Item) tea.Cmd
func (m *Model) RemoveItem(index int)

// Selection
func (m Model) SelectedItem() Item
func (m Model) Index() int
func (m Model) GlobalIndex() int
func (m Model) Cursor() int
func (m *Model) Select(index int)

// Navigation
func (m *Model) CursorUp()
func (m *Model) CursorDown()
func (m *Model) GoToStart()
func (m *Model) GoToEnd()
func (m *Model) NextPage()
func (m *Model) PrevPage()
func (m *Model) ResetSelected()

// Filtering
func (m Model) FilterState() FilterState
func (m Model) FilterValue() string
func (m Model) IsFiltered() bool
func (m *Model) SetFilterState(state FilterState)
func (m *Model) SetFilterText(filter string)
func (m *Model) SetFilteringEnabled(v bool)
func (m Model) FilteringEnabled() bool
func (m Model) SettingFilter() bool
func (m *Model) ResetFilter()
func (m Model) MatchesForItem(index int) []int

// Visibility
func (m *Model) SetShowTitle(v bool)
func (m Model) ShowTitle() bool
func (m *Model) SetShowFilter(v bool)
func (m Model) ShowFilter() bool
func (m *Model) SetShowStatusBar(v bool)
func (m Model) ShowStatusBar() bool
func (m *Model) SetStatusBarItemName(singular, plural string)
func (m Model) StatusBarItemName() (string, string)
func (m *Model) SetShowPagination(v bool)
func (m Model) ShowPagination() bool
func (m *Model) SetShowHelp(v bool)
func (m Model) ShowHelp() bool
func (m *Model) DisableQuitKeybindings()

// Delegate
func (m *Model) SetDelegate(d ItemDelegate)

// Spinner
func (m *Model) SetSpinner(s spinner.Spinner)
func (m *Model) StartSpinner() tea.Cmd
func (m *Model) StopSpinner()
func (m *Model) ToggleSpinner() tea.Cmd

// Status
func (m *Model) NewStatusMessage(s string) tea.Cmd

// Help
func (m Model) ShortHelp() []key.Binding
func (m Model) FullHelp() [][]key.Binding

// Update loop
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

type Model struct {
    KeyMap KeyMap
    Help   help.Model
    // Has unexported fields.
}

type KeyMap struct {
    LineUp       key.Binding
    LineDown     key.Binding
    PageUp       key.Binding
    PageDown     key.Binding
    HalfPageUp   key.Binding
    HalfPageDown key.Binding
    GotoTop      key.Binding
    GotoBottom   key.Binding
}

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
func DefaultStyles() Styles

// Options
func WithColumns(cols []Column) Option
func WithRows(rows []Row) Option
func WithFocused(f bool) Option
func WithHeight(h int) Option
func WithWidth(w int) Option
func WithStyles(s Styles) Option
func WithKeyMap(km KeyMap) Option

// Focus
func (m *Model) Focus()
func (m *Model) Blur()
func (m Model) Focused() bool

// Dimensions
func (m Model) Width() int
func (m *Model) SetWidth(w int)
func (m Model) Height() int
func (m *Model) SetHeight(h int)

// Columns & Rows
func (m Model) Columns() []Column
func (m *Model) SetColumns(c []Column)
func (m Model) Rows() []Row
func (m *Model) SetRows(r []Row)
func (m *Model) FromValues(value, separator string)

// Selection
func (m Model) Cursor() int
func (m *Model) SetCursor(n int)
func (m Model) SelectedRow() Row

// Navigation
func (m *Model) MoveUp(n int)
func (m *Model) MoveDown(n int)
func (m *Model) GotoTop()
func (m *Model) GotoBottom()

// Styling
func (m *Model) SetStyles(s Styles)

// Help
func (m Model) HelpView() string

// Update loop
func (m *Model) UpdateViewport()
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
    DefaultFullCharHalfBlock = '▌'  // allows more granular color blending
    DefaultFullCharFullBlock = '█'  // disables higher resolution blending
    DefaultEmptyCharBlock    = '░'
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
    // Has unexported fields.
}

type FrameMsg struct {
    // Has unexported fields.
}

type ColorFunc func(total, current float64) color.Color
```

### Functions

```go
func New(opts ...Option) Model

// Options
func WithColors(colors ...color.Color) Option   // 0: clear, 1: solid, 2+: blend
func WithDefaultBlend() Option                  // purple to pink gradient
func WithScaled(enabled bool) Option            // scale gradient to filled portion
func WithColorFunc(fn ColorFunc) Option         // dynamic color based on percent
func WithFillCharacters(full, empty rune) Option
func WithoutPercentage() Option
func WithSpringOptions(frequency, damping float64) Option
func WithWidth(w int) Option

// Dimensions
func (m *Model) SetWidth(w int)
func (m Model) Width() int

// Percent
func (m *Model) SetPercent(p float64) tea.Cmd
func (m Model) Percent() float64
func (m *Model) IncrPercent(v float64) tea.Cmd
func (m *Model) DecrPercent(v float64) tea.Cmd

// Animation
func (m Model) Init() tea.Cmd
func (m *Model) IsAnimating() bool
func (m *Model) SetSpringOptions(frequency, damping float64)

// Update loop
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
func (m Model) ViewAs(percent float64) string  // static render
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
    LeftGutterFunc         GutterFunc
    HighlightStyle         lipgloss.Style
    SelectedHighlightStyle lipgloss.Style
    StyleLineFunc          func(int) lipgloss.Style
    // Has unexported fields.
}

type KeyMap struct {
    PageDown     key.Binding
    PageUp       key.Binding
    HalfPageUp   key.Binding
    HalfPageDown key.Binding
    Down         key.Binding
    Up           key.Binding
    Left         key.Binding
    Right        key.Binding
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

// Dimensions
func (m *Model) SetWidth(w int)
func (m Model) Width() int
func (m *Model) SetHeight(h int)
func (m Model) Height() int

// Content
func (m *Model) SetContent(s string)
func (m *Model) SetContentLines(lines []string)  // v2 addition
func (m Model) GetContent() string               // v2 addition
func (m Model) TotalLineCount() int
func (m Model) VisibleLineCount() int

// Offset
func (m *Model) SetYOffset(n int)
func (m Model) YOffset() int
func (m *Model) SetXOffset(n int)
func (m Model) XOffset() int

// Position
func (m Model) AtTop() bool
func (m Model) AtBottom() bool
func (m Model) PastBottom() bool
func (m Model) ScrollPercent() float64
func (m Model) HorizontalScrollPercent() float64  // v2 addition

// Navigation
func (m *Model) GotoTop() []string
func (m *Model) GotoBottom() []string
func (m *Model) ScrollUp(n int)
func (m *Model) ScrollDown(n int)
func (m *Model) ScrollLeft(n int)   // v2 addition
func (m *Model) ScrollRight(n int)  // v2 addition
func (m *Model) PageUp()
func (m *Model) PageDown()
func (m *Model) HalfPageUp()
func (m *Model) HalfPageDown()

// Highlights (v2 additions)
func (m *Model) SetHighlights(matches [][]int)
func (m *Model) HighlightNext()
func (m *Model) HighlightPrevious()
func (m *Model) ClearHighlights()

// Other
func (m *Model) EnsureVisible(line, colstart, colend int)  // v2 addition
func (m *Model) SetHorizontalStep(n int)                   // v2 addition

// Update loop
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
    // Has unexported fields.
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
func New() Model
func DefaultStyles(isDark bool) Styles
func DefaultDarkStyles() Styles
func DefaultLightStyles() Styles

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
type Binding struct {
    // Has unexported fields.
}

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

type KeyMap struct {
    PrevPage key.Binding
    NextPage key.Binding
}
```

### Functions

```go
func New(opts ...Option) Model
func DefaultKeyMap() KeyMap
func WithPerPage(perPage int) Option
func WithTotalPages(totalPages int) Option

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
type Model struct {
    Timeout  time.Duration
    Interval time.Duration
    // Has unexported fields.
}

type TickMsg struct {
    ID       int
    Timeout  bool
}

type TimeoutMsg struct {
    ID int
}

type StartStopMsg struct {
    ID int
}

type Option func(*Model)
```

### Functions

```go
// Constructor changed in v2: was NewWithInterval(timeout, interval)
func New(timeout time.Duration, opts ...Option) Model
func WithInterval(d time.Duration) Option

func (m Model) ID() int
func (m Model) Init() tea.Cmd
func (m Model) Running() bool
func (m Model) Timedout() bool
func (m *Model) Start() tea.Cmd
func (m *Model) Stop() tea.Cmd
func (m *Model) Toggle() tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
```

---

## stopwatch

```go
import "charm.land/bubbles/v2/stopwatch"
```

### Types

```go
type Model struct {
    Interval time.Duration
    // Has unexported fields.
}

type TickMsg struct {
    ID int
}

type StartStopMsg struct {
    ID int
}

type ResetMsg struct {
    ID int
}

type Option func(*Model)
```

### Functions

```go
// Constructor changed in v2: was NewWithInterval(d)
func New(opts ...Option) Model
func WithInterval(d time.Duration) Option

func (m Model) ID() int
func (m Model) Init() tea.Cmd
func (m Model) Running() bool
func (m Model) Elapsed() time.Duration
func (m Model) Start() tea.Cmd
func (m Model) Stop() tea.Cmd
func (m Model) Toggle() tea.Cmd
func (m Model) Reset() tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
```

---

## filepicker

```go
import "charm.land/bubbles/v2/filepicker"
```

### Types

```go
type Model struct {
    Path             string
    CurrentDirectory string
    AllowedTypes     []string
    KeyMap           KeyMap
    ShowPermissions  bool
    ShowSize         bool
    ShowHidden       bool
    DirAllowed       bool
    FileAllowed      bool
    FileSelected     string
    AutoHeight       bool
    Cursor           string
    Styles           Styles
    // Has unexported fields.
}

type KeyMap struct {
    GoToTop  key.Binding
    GoToLast key.Binding
    Down     key.Binding
    Up       key.Binding
    PageUp   key.Binding
    PageDown key.Binding
    Back     key.Binding
    Open     key.Binding
    Select   key.Binding
}

type Styles struct {
    DisabledCursor   lipgloss.Style
    Cursor           lipgloss.Style
    Symlink          lipgloss.Style
    Directory        lipgloss.Style
    File             lipgloss.Style
    DisabledFile     lipgloss.Style
    Permission       lipgloss.Style
    Selected         lipgloss.Style
    DisabledSelected lipgloss.Style
    FileSize         lipgloss.Style
    EmptyDirectory   lipgloss.Style
}
```

### Functions

```go
func New() Model
func IsHidden(file string) (bool, error)
func DefaultStyles() Styles

func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string

func (m Model) Height() int
func (m *Model) SetHeight(h int)

func (m Model) DidSelectFile(msg tea.Msg) (bool, string)
func (m Model) DidSelectDisabledFile(msg tea.Msg) (bool, string)
func (m Model) HighlightedPath() string
```

---

## cursor

```go
import "charm.land/bubbles/v2/cursor"
```

> Usually embedded within textinput/textarea rather than used directly.

### Types

```go
type Mode int
const (
    CursorBlink Mode = iota
    CursorStatic
    CursorHide
)

type Model struct {
    Style      lipgloss.Style
    TextStyle  lipgloss.Style
    BlinkSpeed time.Duration
    IsBlinked  bool
    // Has unexported fields.
}

type BlinkMsg struct {
    // Has unexported fields.
}
```

### Functions

```go
func New() Model
func Blink() tea.Msg

func (m Model) Mode() Mode
func (m *Model) SetMode(mode Mode) tea.Cmd
func (m *Model) SetChar(char string)
func (m *Model) Focus() tea.Cmd
func (m *Model) Blur()
func (m *Model) Blink() tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
```

### Changes from v1

| v1 | v2 |
|---|---|
| `model.Blink` | `model.IsBlinked` |
| `model.BlinkCmd()` | `model.Blink()` |
