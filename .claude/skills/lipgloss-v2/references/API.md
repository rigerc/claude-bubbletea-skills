# Lip Gloss v2 — API Cheat Sheet

Import: `charm.land/lipgloss/v2`

---

## Package-level Functions

### Color

| Function | Signature | Description |
|---|---|---|
| `Color` | `func Color(s string) color.Color` | Hex (`#rrggbb`) or ANSI256 (`"21"`) |
| `LightDark` | `func LightDark(isDark bool) LightDarkFunc` | Returns helper that picks light or dark value |
| `Complete` | `func Complete(p colorprofile.Profile) CompleteFunc` | Returns helper that picks by color profile |
| `HasDarkBackground` | `func HasDarkBackground(in, out term.File) bool` | Query terminal background (standalone only) |
| `BackgroundColor` | `func BackgroundColor(in, out term.File) (color.Color, error)` | Query raw background color |
| `Alpha` | `func Alpha(c color.Color, alpha float64) color.Color` | Adjust alpha 0.0–1.0 |
| `Darken` | `func Darken(c color.Color, pct float64) color.Color` | Darken by 0.0–1.0 |
| `Lighten` | `func Lighten(c color.Color, pct float64) color.Color` | Lighten by 0.0–1.0 |
| `Complementary` | `func Complementary(c color.Color) color.Color` | 180° hue shift |
| `Blend1D` | `func Blend1D(steps int, stops ...color.Color) []color.Color` | Linear gradient, N steps |
| `Blend2D` | `func Blend2D(w, h int, angle float64, stops ...color.Color) []color.Color` | 2D gradient, row-major |

### Layout

| Function | Signature | Description |
|---|---|---|
| `JoinHorizontal` | `func JoinHorizontal(pos Position, strs ...string) string` | Side-by-side, aligned Top/Center/Bottom |
| `JoinVertical` | `func JoinVertical(pos Position, strs ...string) string` | Stacked, aligned Left/Center/Right |
| `Place` | `func Place(w, h int, hPos, vPos Position, str string, opts ...WhitespaceOption) string` | Center content in box |
| `PlaceHorizontal` | `func PlaceHorizontal(w int, pos Position, str string, opts ...WhitespaceOption) string` | Horizontal placement |
| `PlaceVertical` | `func PlaceVertical(h int, pos Position, str string, opts ...WhitespaceOption) string` | Vertical placement |

Whitespace options:
- `WithWhitespaceChars(str string)` — fill character(s)
- `WithWhitespaceStyle(s Style)` — style for fill area

### Text utilities

| Function | Description |
|---|---|
| `Width(str string) int` | Cell width (ANSI-aware, Unicode-aware) |
| `Height(str string) int` | Number of lines |
| `Size(str string) (int, int)` | width, height |
| `Wrap(s string, width int, breakpoints string) string` | Word-wrap preserving ANSI |
| `StyleRunes(str string, indices []int, matched, unmatched Style) string` | Style by rune index |
| `StyleRanges(s string, ranges ...Range) string` | Style byte ranges |

### Printing (always use these, not `fmt`)

| Function | Description |
|---|---|
| `Print(v ...any)` | Print to stdout with downsampling |
| `Printf(format string, v ...any)` | Formatted print to stdout |
| `Println(v ...any)` | Print + newline to stdout |
| `Fprint(w io.Writer, v ...any)` | Print to writer with downsampling |
| `Fprintf(w io.Writer, format string, v ...any)` | Formatted print to writer |
| `Fprintln(w io.Writer, v ...any)` | Print + newline to writer |
| `Sprint(v ...any) string` | String for stdout profile |
| `Sprintf(format string, v ...any) string` | Formatted string for stdout profile |
| `Sprintln(v ...any) string` | String + newline for stdout profile |

Default writer: `lipgloss.Writer` (targets `os.Stdout`).
Override: `lipgloss.Writer = colorprofile.NewWriter(os.Stderr, os.Environ())`

---

## Style Methods

All methods return `Style` (value type — immutable chainable API).

### Text decoration

```
Bold(bool)
Italic(bool)
Underline(bool)
UnderlineStyle(Underline)   // UnderlineNone/Single/Double/Curly/Dotted/Dashed
UnderlineColor(color.Color)
Strikethrough(bool)
StrikethroughSpaces(bool)
Faint(bool)
Blink(bool)
Reverse(bool)
```

### Color

```
Foreground(color.Color)
Background(color.Color)
```

### Size & alignment

```
Width(int)
Height(int)
MaxWidth(int)
MaxHeight(int)
Align(Position)             // Left/Center/Right
AlignVertical(Position)     // Top/Center/Bottom
Inline(bool)
TabWidth(int)               // NoTabConversion = -1 to disable
Transform(func(string) string)
```

### Padding

```
Padding(values ...int)      // top [right [bottom [left]]]
PaddingTop(int)
PaddingRight(int)
PaddingBottom(int)
PaddingLeft(int)
PaddingChar(rune)           // fill character
UnsetPadding()
```

### Margin

```
Margin(values ...int)
MarginTop(int)
MarginRight(int)
MarginBottom(int)
MarginLeft(int)
MarginBackground(color.Color)
MarginChar(rune)
UnsetMargin()
```

### Border

```
Border(border Border, sides ...bool)
BorderStyle(Border)
BorderTop(bool) / BorderRight(bool) / BorderBottom(bool) / BorderLeft(bool)
BorderForeground(colors ...color.Color)        // per-side or single
BorderBackground(colors ...color.Color)
BorderTopForeground(color.Color)
BorderRightForeground(color.Color)
BorderBottomForeground(color.Color)
BorderLeftForeground(color.Color)
BorderForegroundBlend(from, to color.Color)    // gradient
BorderForegroundBlendOffset(int)
UnsetBorder()
```

### Inheritance & utilities

```
Inherit(Style)              // fills unset props from parent
SetString(str string) Style // bind content
String() string             // render with SetString content
Render(str ...string) string
Value() string              // raw underlying string
```

### Unset methods

All properties have corresponding unset methods:
```
UnsetBold(), UnsetItalic(), UnsetUnderline(), UnsetStrikethrough(),
UnsetForeground(), UnsetBackground(), UnsetWidth(), UnsetHeight(),
UnsetAlign(), UnsetPadding(), UnsetMargin(), UnsetBorder(), etc.
```

---

## Constants

### Position

```go
lipgloss.Left    = 0.0
lipgloss.Right   = 1.0
lipgloss.Center  = 0.5
lipgloss.Top     = 0.0
lipgloss.Bottom  = 1.0
```

### ANSI 4-bit colors

```go
lipgloss.Black, Red, Green, Yellow, Blue, Magenta, Cyan, White
lipgloss.BrightBlack, BrightRed, BrightGreen, BrightYellow
lipgloss.BrightBlue, BrightMagenta, BrightCyan, BrightWhite
```

### Underline styles

```go
lipgloss.UnderlineNone
lipgloss.UnderlineSingle   // default when Underline(true)
lipgloss.UnderlineDouble
lipgloss.UnderlineCurly
lipgloss.UnderlineDotted
lipgloss.UnderlineDashed
```

### Other

```go
lipgloss.NBSP              // non-breaking space rune '\u00A0'
lipgloss.NoTabConversion   // -1 — pass to TabWidth to disable replacement
```

---

## Border Constructors

```go
lipgloss.NormalBorder()           // ─ │ ┌ ┐ └ ┘
lipgloss.RoundedBorder()          // ─ │ ╭ ╮ ╰ ╯
lipgloss.ThickBorder()            // ━ ┃ ┏ ┓ ┗ ┛
lipgloss.DoubleBorder()           // ═ ║ ╔ ╗ ╚ ╝
lipgloss.BlockBorder()            // █ (solid blocks)
lipgloss.HiddenBorder()           // invisible
lipgloss.ASCIIBorder()            // - | + + + +
lipgloss.MarkdownBorder()         // markdown table style
lipgloss.InnerHalfBlockBorder()   // half block inner
lipgloss.OuterHalfBlockBorder()   // half block outer
```

Custom border:
```go
b := lipgloss.Border{
    Top: "─", Bottom: "─", Left: "│", Right: "│",
    TopLeft: "╭", TopRight: "╮",
    BottomLeft: "╰", BottomRight: "╯",
}
```

---

## Compositing

### Layer

```go
// Create layers
base := lipgloss.NewLayer(content)
modal := lipgloss.NewLayer(content, childLayer...).
    X(xOffset).Y(yOffset).Z(zIndex).ID("modal")

// Layer methods
layer.GetContent() string
layer.GetID() string
layer.GetX() / GetY() / GetZ() int
layer.Width() / Height() int
layer.AddLayers(layers ...*Layer) *Layer
layer.GetLayer(id string) *Layer
layer.MaxZ() int
```

### Compositor

```go
// Compose
comp := lipgloss.NewCompositor(base, modal, ...)

// Methods
comp.AddLayers(layers ...*Layer) *Compositor
comp.Bounds() image.Rectangle
comp.Hit(x, y int) LayerHit
comp.GetLayer(id string) *Layer
comp.Refresh()                    // re-flatten after changes
comp.Render() string

// LayerHit methods
hit.Empty() bool
hit.ID() string
hit.Layer() *Layer
hit.Bounds() image.Rectangle
```

### Canvas

```go
canvas := lipgloss.NewCanvas(width, height)
canvas.Resize(width, height)
canvas.Clear()
canvas.Compose(drawer) *Canvas
canvas.Render() string
canvas.Width() / Height() int
```

---

## Subpackage: list

Import: `charm.land/lipgloss/v2/list`

### Constructors

```go
list.New(items...)
```

### Methods

```go
l.Enumerator(enumerator) *List                    // Bullet, Dash, Roman, Arabic, Alphabet, Asterisk
l.EnumeratorStyle(style) *List
l.EnumeratorStyleFunc(f StyleFunc) *List
l.Item(item) *List                                // add single item
l.Items(items...) *List                           // add multiple items
l.ItemStyle(style) *List
l.ItemStyleFunc(f StyleFunc) *List
l.Indenter(indenter) *List
l.Hide(hide bool) *List
l.Hidden() bool
l.Offset(start, end int) *List
l.String() string
l.Value() string
```

### Types

```go
list.Enumerator func(items Items, index int) string
list.Indenter func(items Items, index int) string
list.StyleFunc func(items Items, index int) lipgloss.Style
list.Items // alias for tree.Children
```

### Built-in enumerators

```go
list.Bullet(items, i)      // •
list.Dash(items, i)        // -
list.Roman(items, i)       // I. II. III.
list.Arabic(items, i)      // 1. 2. 3.
list.Alphabet(items, i)    // a. b. c.
list.Asterisk(items, i)    // *
```

---

## Subpackage: table

Import: `charm.land/lipgloss/v2/table`

### Constants

```go
table.HeaderRow = -1  // use in StyleFunc for header row
```

### Constructors

```go
table.New() *Table
table.NewStringData(rows...[]string) *StringData
table.NewFilter(data Data) *Filter
```

### Table methods

```go
t.Border(border lipgloss.Border) *Table
t.BorderTop/Right/Bottom/Left(bool) *Table
t.BorderColumn(bool) *Table        // column separators
t.BorderHeader(bool) *Table        // header separator
t.BorderRow(bool) *Table           // row separators
t.BorderStyle(style) *Table
t.Data(data Data) *Table
t.Headers(cols ...string) *Table
t.Row(cols ...string) *Table
t.Rows(rows ...[]string) *Table
t.ClearRows() *Table
t.StyleFunc(fn StyleFunc) *Table
t.Width(int) *Table
t.Height(int) *Table
t.Offset(n int) *Table             // scroll Y offset
t.Wrap(bool) *Table
t.Render() string
t.String() string
```

### Data interface

```go
type Data interface {
    At(row, cell int) string
    Rows() int
    Columns() int
}
```

### Types

```go
table.StyleFunc func(row, col int) lipgloss.Style
table.StringData // implements Data
table.Filter     // wraps Data with filtering
```

---

## Subpackage: tree

Import: `charm.land/lipgloss/v2/tree`

### Constructors

```go
tree.New() *Tree
tree.Root(root any) *Tree           // shorthand for tree.New().Root(root)
tree.NewStringData(data ...string) Children
tree.NewFilter(data Children) *Filter
tree.NewLeaf(value any, hidden bool) *Leaf
```

### Tree methods

```go
t.Root(root any) *Tree
t.Child(children...any) *Tree       // strings or *Tree
t.Enumerator(enum Enumerator) *Tree // DefaultEnumerator, RoundedEnumerator
t.EnumeratorStyle(style) *Tree
t.EnumeratorStyleFunc(fn StyleFunc) *Tree
t.ItemStyle(style) *Tree
t.ItemStyleFunc(fn StyleFunc) *Tree
t.Indenter(indenter) *Tree
t.RootStyle(style) *Tree
t.Hide(hide bool) *Tree
t.Hidden() bool
t.Offset(start, end int) *Tree
t.Children() Children
t.SetHidden(bool)
t.SetValue(value any)
t.String() string
t.Value() string
```

### Types

```go
tree.Enumerator func(children Children, index int) string
tree.Indenter func(children Children, index int) string
tree.StyleFunc func(children Children, i int) lipgloss.Style
tree.Children interface {
    At(index int) Node
    Length() int
}
tree.Node interface {
    fmt.Stringer
    Value() string
    Children() Children
    Hidden() bool
    SetHidden(bool)
    SetValue(any)
}
tree.NodeChildren []Node           // implements Children
tree.Leaf                          // implements Node
```

### Built-in enumerators

```go
tree.DefaultEnumerator   // ├── └──
tree.RoundedEnumerator   // ├── ╰──
```

---

## Integration with huh v2

huh v2 forms use Lip Gloss v2 for theming. Import both packages:

```go
import (
    "charm.land/huh/v2"
    "charm.land/lipgloss/v2"
)
```

### Theme structure

```go
type Styles struct {
    Form           FormStyles
    Group          GroupStyles
    FieldSeparator lipgloss.Style
    Blurred        FieldStyles
    Focused        FieldStyles
    Help           help.Styles
}
```

### Creating a custom theme

```go
theme := huh.ThemeFunc(func(isDark bool) *huh.Styles {
    s := huh.ThemeCharm(isDark) // Start with base

    // Override focused styles
    s.Focused.Title = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FF5F87"))

    s.Focused.Base = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#874BFD"))

    return s
})

form.WithTheme(theme)
```

### Key FieldStyles properties

```go
type FieldStyles struct {
    Base           lipgloss.Style
    Title          lipgloss.Style
    Description    lipgloss.Style
    ErrorIndicator lipgloss.Style
    ErrorMessage   lipgloss.Style
    SelectSelector lipgloss.Style
    Option         lipgloss.Style
    NextIndicator  lipgloss.Style
    PrevIndicator  lipgloss.Style
    Directory      lipgloss.Style
    File           lipgloss.Style
    MultiSelectSelector lipgloss.Style
    SelectedOption      lipgloss.Style
    SelectedPrefix      lipgloss.Style
    UnselectedOption    lipgloss.Style
    UnselectedPrefix    lipgloss.Style
    TextInput           TextInputStyles
    FocusedButton       lipgloss.Style
    BlurredButton       lipgloss.Style
    Card                lipgloss.Style
    NoteTitle           lipgloss.Style
    Next                lipgloss.Style
}
```

---

## compat package

For drop-in v1 replacement of `AdaptiveColor`, `CompleteColor`, etc.:

```go
import "charm.land/lipgloss/v2/compat"

color := compat.AdaptiveColor{
    Light: lipgloss.Color("#f1f1f1"),
    Dark:  lipgloss.Color("#cccccc"),
}
// compat reads os.Stdin/os.Stdout globally, just like v1
```

Available types:
- `compat.AdaptiveColor`
- `compat.CompleteColor`
- `compat.CompleteAdaptiveColor`
