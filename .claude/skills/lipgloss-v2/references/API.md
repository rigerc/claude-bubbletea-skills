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

## Compositing

```go
// Create layers
base := lipgloss.NewLayer(content)
modal := lipgloss.NewLayer(content, childLayer...).
    X(xOffset).Y(yOffset).Z(zIndex).ID("modal")

// Compose
comp := lipgloss.NewCompositor(base, modal, ...)
output := comp.Render()

// Hit testing
hit := comp.Hit(x, y) // returns LayerHit with ID

// Low-level canvas
canvas := lipgloss.NewCanvas(width, height)
canvas.Compose(layer)
canvas.Render()
```

---

## Subpackage: list

Import: `charm.land/lipgloss/v2/list`

```go
l := list.New(items...).
    Enumerator(list.Bullet).           // or Dash, Roman, Arabic, custom func
    EnumeratorStyle(style).
    EnumeratorStyleFunc(func(items list.Items, i int) lipgloss.Style).
    ItemStyle(style).
    ItemStyleFunc(func(items list.Items, i int) lipgloss.Style).
    Item(nestedList)
```

Custom enumerator: `func(items list.Items, i int) string`

---

## Subpackage: table

Import: `charm.land/lipgloss/v2/table`

```go
t := table.New().
    Border(lipgloss.Border).
    BorderStyle(lipgloss.Style).
    BorderTop/Right/Bottom/Left(bool).
    Headers(cols ...string).
    Rows(rows ...[]string).
    Row(cols ...string).
    StyleFunc(func(row, col int) lipgloss.Style).
    Width(int).
    Height(int).
    Offset(n int).               // scroll Y offset
    String() string              // implements fmt.Stringer
```

`row == table.HeaderRow` in StyleFunc to style the header.

---

## Subpackage: tree

Import: `charm.land/lipgloss/v2/tree`

```go
t := tree.New().
    Root("name").
    Child(items...).             // strings or *Tree
    EnumeratorStyle(style).
    EnumeratorStyleFunc(func(children tree.Children, i int) lipgloss.Style).
    ItemStyle(style).
    ItemStyleFunc(func(children tree.Children, i int) lipgloss.Style).
    IndenterStyle(style).
    IndenterStyleFunc(func(children tree.Children, i int) lipgloss.Style).
    Width(int)

// Subtrees
tree.Root("dir/").Child("file.go", "file_test.go")
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
