---
name: lipgloss-v2
description: >-
  Generates, explains, and debugs terminal styling and layout code using
  Charm's Lip Gloss v2 library for Go. Use when building styled terminal
  output, TUI layouts, color gradients, bordered boxes, tables, lists, or
  trees. Covers Style API, adaptive colors, LightDark detection, compositing
  layers, and all subpackages (list, table, tree). Also handles migration
  from Lip Gloss v1 to v2.
allowed-tools: Read, Write, Edit, Bash(go *)
metadata:
  category: go-cli
  version: "2.0"
---

# Lip Gloss v2

Lip Gloss v2 (`charm.land/lipgloss/v2`) provides style definitions for
terminal layouts. It is a pure value type library — no global renderer, no
hidden state. Colors are downsampled at the output layer.

```bash
go get charm.land/lipgloss/v2
```

## Quick Reference

| v1 | v2 |
|---|---|
| `github.com/charmbracelet/lipgloss` | `charm.land/lipgloss/v2` |
| `type Color string` | `func Color(string) color.Color` |
| `lipgloss.AdaptiveColor{...}` | `compat.AdaptiveColor{...}` or `LightDark(isDark)(light, dark)` |
| `HasDarkBackground()` | `HasDarkBackground(os.Stdin, os.Stdout)` |
| `fmt.Println(s.Render("hi"))` | `lipgloss.Println(s.Render("hi"))` |
| `renderer.NewStyle()` | `lipgloss.NewStyle()` |

See `references/MIGRATION.md` for the full migration guide.

---

## Style API

Styles are immutable value types. Chain setters, end with `.Render(str)`.

```go
s := lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FF5F87")).
    Background(lipgloss.Color("#1a1a2e")).
    Padding(1, 2).
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#874BFD")).
    Width(40).
    Align(lipgloss.Center)

lipgloss.Println(s.Render("Hello!"))
```

### Text properties

| Method | Description |
|---|---|
| `Bold(bool)` | Bold text |
| `Italic(bool)` | Italic text |
| `Underline(bool)` | Underline (single) |
| `UnderlineStyle(Underline)` | `UnderlineSingle/Double/Curly/Dotted/Dashed/None` |
| `UnderlineColor(color.Color)` | Color the underline separately |
| `Strikethrough(bool)` | Strikethrough |
| `Faint(bool)` | Faint/dim text |
| `Blink(bool)` | Blinking text |
| `Reverse(bool)` | Reverse foreground/background |
| `Foreground(color.Color)` | Text color |
| `Background(color.Color)` | Background color |
| `Hyperlink(url, params...)` | Clickable hyperlink |

### Layout properties

| Method | Description |
|---|---|
| `Width(int)` | Set width (pads to fill) |
| `Height(int)` | Set height (pads to fill) |
| `MaxWidth(int)` | Truncate at max width |
| `MaxHeight(int)` | Truncate at max height |
| `Align(Position)` | Horizontal text alignment |
| `AlignVertical(Position)` | Vertical text alignment |
| `Inline(bool)` | Render inline (no newlines) |
| `TabWidth(int)` | Tab-to-space width (`NoTabConversion` to disable) |
| `Transform(func(string) string)` | Post-render transform |

Position constants: `lipgloss.Left`, `lipgloss.Center`, `lipgloss.Right`,
`lipgloss.Top`, `lipgloss.Bottom`.

### Padding & margin

```go
s.Padding(1, 2)          // top/bottom=1, left/right=2
s.Padding(1, 2, 3, 4)    // top, right, bottom, left
s.PaddingTop(1).PaddingLeft(2)
s.PaddingChar('·')       // fill padding with this character

s.Margin(1, 2)
s.MarginChar('·')
s.MarginBackground(color.Color)
```

### Borders

```go
s.Border(lipgloss.RoundedBorder())
s.Border(lipgloss.NormalBorder(), true)       // all sides
s.Border(lipgloss.ThickBorder(), true, false) // top+bottom only
s.BorderForeground(lipgloss.Color("#874BFD"))
s.BorderForegroundBlend(from, to)             // gradient border
s.BorderStyle(lipgloss.NormalBorder()).BorderBottom(true)
```

Border constructors: `NormalBorder()`, `RoundedBorder()`, `ThickBorder()`,
`DoubleBorder()`, `BlockBorder()`, `HiddenBorder()`, `ASCIIBorder()`,
`MarkdownBorder()`, `InnerHalfBlockBorder()`, `OuterHalfBlockBorder()`.

Custom border:
```go
b := lipgloss.Border{
    Top: "─", Bottom: "─", Left: "│", Right: "│",
    TopLeft: "╭", TopRight: "╮",
    BottomLeft: "╰", BottomRight: "╯",
}
```

### Style inheritance

```go
base := lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("5"))

// Copy and override
active := base.Bold(true).Background(lipgloss.Color("#FF5F87"))

// Inherit (fills only unset props)
child := lipgloss.NewStyle().Inherit(base)
```

### Render helpers

```go
s.Render("text")        // returns styled string
s.String()              // alias for Render with SetString content
s.SetString("content")  // bind content to style
```

---

## Color System

`color.Color` is the standard `image/color.Color` interface.

```go
import "charm.land/lipgloss/v2"

// Hex or ANSI256
c := lipgloss.Color("#FF5F87")
c := lipgloss.Color("200")

// Named 4-bit ANSI constants
c := lipgloss.Red      // lipgloss.Black/Red/Green/.../BrightWhite

// ANSI256
c := lipgloss.ANSIColor(134)
```

### Adaptive colors (light/dark terminals)

**Standalone:**
```go
hasDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
lightDark := lipgloss.LightDark(hasDark)

fg := lightDark(lipgloss.Color("#333333"), lipgloss.Color("#f1f1f1"))
s := lipgloss.NewStyle().Foreground(fg)
```

**Bubble Tea:**
```go
func (m model) Init() tea.Cmd {
    return tea.RequestBackgroundColor
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.BackgroundColorMsg:
        m.styles = newStyles(msg.IsDark())
    }
    return m, nil
}

func newStyles(isDark bool) styles {
    ld := lipgloss.LightDark(isDark)
    return styles{
        title: lipgloss.NewStyle().Foreground(ld(
            lipgloss.Color("#333333"),
            lipgloss.Color("#f1f1f1"),
        )),
    }
}
```

### Color profile selection (`Complete`)

```go
p := colorprofile.Detect(os.Stdout, os.Environ())
complete := lipgloss.Complete(p)
c := complete(
    lipgloss.Color("5"),       // ANSI fallback
    lipgloss.Color("200"),     // ANSI256 fallback
    lipgloss.Color("#ff34ac"), // TrueColor
)
```

### Color manipulation

```go
lipgloss.Darken(c, 0.2)        // 20% darker
lipgloss.Lighten(c, 0.2)       // 20% lighter
lipgloss.Alpha(c, 0.5)         // 50% alpha
lipgloss.Complementary(c)      // 180° on color wheel
```

### Color gradients

```go
// 1D linear blend across N steps
colors := lipgloss.Blend1D(40, from, to, extra...)

// 2D gradient (angle 0=left→right, 180=right→left)
colors := lipgloss.Blend2D(width, height, 45.0, c1, c2, c3)
// Access: colors[y*width + x]
```

---

## Layout Functions

```go
// Join side-by-side, aligned at Top/Center/Bottom
row := lipgloss.JoinHorizontal(lipgloss.Top, blockA, blockB, blockC)

// Stack vertically, aligned at Left/Center/Right
col := lipgloss.JoinVertical(lipgloss.Left, blockA, blockB)

// Place content in a box (fills with whitespace)
out := lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content,
    lipgloss.WithWhitespaceChars("·"),
    lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Foreground(subtle)),
)

// One-axis variants
lipgloss.PlaceHorizontal(width, lipgloss.Center, content)
lipgloss.PlaceVertical(height, lipgloss.Center, content)
```

---

## Printing (Important!)

Always use Lip Gloss writers to ensure correct color downsampling:

```go
lipgloss.Println(s)           // stdout + newline
lipgloss.Print(s)             // stdout
lipgloss.Printf("%s\n", s)    // stdout, formatted
lipgloss.Fprintln(os.Stderr, s)
lipgloss.Sprint(s)            // returns string for stdout profile
lipgloss.Sprintf("...", s)

// Customize default writer
lipgloss.Writer = colorprofile.NewWriter(os.Stderr, os.Environ())
```

In Bubble Tea, downsampling is automatic — just return strings normally.

---

## Text Utilities

```go
w := lipgloss.Width(str)           // cell width (Unicode-aware)
h := lipgloss.Height(str)          // line count
w, h := lipgloss.Size(str)         // both

wrapped := lipgloss.Wrap(str, 80, "")  // wrap preserving ANSI

// Style specific rune indices
out := lipgloss.StyleRunes(str, []int{0, 1, 2}, matched, unmatched)

// Style character ranges
out := lipgloss.StyleRanges(str,
    lipgloss.NewRange(0, 5, boldStyle),
    lipgloss.NewRange(6, 10, italicStyle),
)
```

---

## Compositing (Layers)

Compose styled strings at arbitrary positions:

```go
base := lipgloss.NewLayer(document)
modal := lipgloss.NewLayer(floatingBox).X(10).Y(5)
overlay := lipgloss.NewLayer(badge).X(80).Y(2).Z(10)

comp := lipgloss.NewCompositor(base, modal, overlay)
lipgloss.Println(comp.Render())
```

---

## Subpackages

### list (`charm.land/lipgloss/v2/list`)

```go
import "charm.land/lipgloss/v2/list"

l := list.New("Item A", "Item B", "Item C").
    Enumerator(list.Roman).
    EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)).
    ItemStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("255")))

// Nested lists
l.Item(
    list.New("Sub A", "Sub B"),
)

lipgloss.Println(l)
```

Built-in enumerators: `list.Bullet`, `list.Dash`, `list.Roman`, `list.Arabic`.
Custom: `func(items list.Items, i int) string`.

### table (`charm.land/lipgloss/v2/table`)

```go
import "charm.land/lipgloss/v2/table"

t := table.New().
    Border(lipgloss.ThickBorder()).
    BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
    StyleFunc(func(row, col int) lipgloss.Style {
        if row == table.HeaderRow {
            return lipgloss.NewStyle().Bold(true).Align(lipgloss.Center)
        }
        if row%2 == 0 {
            return lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
        }
        return lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
    }).
    Headers("NAME", "LANG", "STARS").
    Rows([][]string{
        {"BubbleTea", "Go", "★★★★★"},
        {"Lipgloss", "Go", "★★★★★"},
    }...)

lipgloss.Println(t)
```

### tree (`charm.land/lipgloss/v2/tree`)

```go
import "charm.land/lipgloss/v2/tree"

t := tree.New().
    Root("project/").
    Child(
        ".git",
        tree.Root("src/").
            Child("main.go", "util.go"),
        tree.Root("docs/").
            Child("README.md"),
        "go.mod",
    ).
    EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
    ItemStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("255")))

lipgloss.Println(t)
```

---

## Examples Index

See `references/_examples/` for runnable examples:

| File | Demonstrates |
|---|---|
| `color-standalone.go` | `HasDarkBackground` + `LightDark` in standalone mode |
| `color-bubbletea.go` | `tea.BackgroundColorMsg` + adaptive styles in BubbleTea |
| `layout.go` | Full layout: tabs, title, dialog, lists, status bar, compositing |
| `table-languages.go` | Styled table with headers, alternating rows, custom column widths |
| `list-grocery.go` | Custom enumerator + per-item style functions |
| `blending-1d.go` | `Blend1D` color gradients rendered as color bars |

---

## References

- `references/API.md` — complete API cheat sheet
- `references/MIGRATION.md` — v1 → v2 migration guide (also suitable for LLM-driven migration)
