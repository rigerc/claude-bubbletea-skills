# Lip Gloss v1 → v2 Migration Guide

This guide covers migrating from Lip Gloss v1 (`github.com/charmbracelet/lipgloss`)
to Lip Gloss v2 (`charm.land/lipgloss/v2`). Written for both humans and LLMs
performing automated migrations.

---

## Quick Start (Two Steps)

### 1. Module path

```
# Search-replace:
github.com/charmbracelet/lipgloss → charm.land/lipgloss/v2
```

All subpackages follow the same pattern:
```
github.com/charmbracelet/lipgloss/table → charm.land/lipgloss/v2/table
github.com/charmbracelet/lipgloss/tree  → charm.land/lipgloss/v2/tree
github.com/charmbracelet/lipgloss/list  → charm.land/lipgloss/v2/list
```

Install: `go get charm.land/lipgloss/v2@v2.0.0-beta.3.0.20251205162909-7869489d8971`

### 2. Use Lip Gloss print functions

```go
// v1
fmt.Println(s.Render("hello"))

// v2 — required for color downsampling
lipgloss.Println(s.Render("hello"))
```

(If using Bubble Tea v2, this is handled automatically.)

---

## Color System

### `Color` is now a function, not a type

```go
// v1
var c lipgloss.Color = "#ff00ff"
var c lipgloss.Color = "21"

// v2
var c color.Color = lipgloss.Color("#ff00ff")
var c color.Color = lipgloss.Color("21")
```

The return type is `image/color.Color`. Add `import "image/color"`.

### `TerminalColor` interface removed

Replace all `lipgloss.TerminalColor` with `color.Color`.

### `AdaptiveColor` moved to `compat` or use `LightDark`

```go
// v1
color := lipgloss.AdaptiveColor{Light: "#0000ff", Dark: "#000099"}

// v2 — using compat (drop-in)
import "charm.land/lipgloss/v2/compat"
color := compat.AdaptiveColor{
    Light: lipgloss.Color("#0000ff"),
    Dark:  lipgloss.Color("#000099"),
}

// v2 — recommended pattern
hasDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
ld := lipgloss.LightDark(hasDark)
color := ld(lipgloss.Color("#0000ff"), lipgloss.Color("#000099"))
```

### `CompleteColor` moved to `compat` or use `Complete`

```go
// v1
color := lipgloss.CompleteColor{TrueColor: "#ff00ff", ANSI256: "200", ANSI: "5"}

// v2 — using compat
import "charm.land/lipgloss/v2/compat"
color := compat.CompleteColor{
    TrueColor: lipgloss.Color("#ff00ff"),
    ANSI256:   lipgloss.Color("200"),
    ANSI:      lipgloss.Color("5"),
}

// v2 — recommended
import "github.com/charmbracelet/colorprofile"

profile := colorprofile.Detect(os.Stdout, os.Environ())
complete := lipgloss.Complete(profile)
color := complete(lipgloss.Color("5"), lipgloss.Color("200"), lipgloss.Color("#ff00ff"))
```

### `NoColor` type added

```go
// v2 — specify absence of color
style := lipgloss.NewStyle().Background(lipgloss.NoColor{})
```

---

## Renderer Removal

The `Renderer` type is removed entirely. `Style` is now a plain value type.

```go
// v1 — these no longer exist
lipgloss.DefaultRenderer()
lipgloss.SetDefaultRenderer(r)
lipgloss.NewRenderer(w, opts...)
lipgloss.ColorProfile()
lipgloss.SetColorProfile(p)
renderer.NewStyle()

// v2 replacements
lipgloss.NewStyle()                            // instead of renderer.NewStyle()
colorprofile.Detect(os.Stdout, os.Environ())   // instead of ColorProfile()
// SetColorProfile → set lipgloss.Writer.Profile
```

---

## Background Detection

```go
// v1
hasDark := lipgloss.HasDarkBackground()  // no args

// v2 — must pass I/O files
hasDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)

// v1
lipgloss.SetHasDarkBackground(b)  // removed

// v2 — just pass the bool to LightDark
ld := lipgloss.LightDark(hasDark)
```

### With Bubble Tea (v2 pattern)

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
```

---

## Printing & Color Downsampling

In v1, downsampling happened inside `Style.Render()`. In v2, `Render()` always
emits full-fidelity ANSI — downsampling happens at the output layer.

```go
// v1
fmt.Println(style.Render("text"))

// v2
lipgloss.Println(style.Render("text"))
lipgloss.Fprintln(os.Stderr, style.Render("text"))
str := lipgloss.Sprint(style.Render("text"))
```

---

## Whitespace Options

```go
// v1
lipgloss.Place(w, h, hPos, vPos, str,
    lipgloss.WithWhitespaceForeground(c),
    lipgloss.WithWhitespaceBackground(c),
)

// v2
lipgloss.Place(w, h, hPos, vPos, str,
    lipgloss.WithWhitespaceStyle(
        lipgloss.NewStyle().
            Foreground(c).
            Background(c),
    ),
)
```

---

## Underline

```go
// v1 + v2 — still works
s.Underline(true)

// v2 — new fine-grained control
s.UnderlineStyle(lipgloss.UnderlineCurly)
s.UnderlineColor(lipgloss.Color("#FF0000"))
```

---

## Strikethrough Spaces

```go
// v2 — new option
s.StrikethroughSpaces(true)  // apply strikethrough to spaces between words
```

---

## Tree Subpackage (new in v2)

New methods on `*tree.Tree`:
- `IndenterStyle(lipgloss.Style)` — static indentation style
- `IndenterStyleFunc(func(Children, int) lipgloss.Style)` — conditional
- `Width(int)` — set tree width for padding
- `Offset(start, end int)` — set children offsets
- `Hide(bool)` — hide the tree node

New root constructor:
```go
// v2 — shorthand
t := tree.Root("project/").Child("file.go")
// instead of
t := tree.New().Root("project/").Child("file.go")
```

---

## Compositing (new in v2)

Layer-based composition is entirely new in v2:

```go
// Create layers at arbitrary positions
base := lipgloss.NewLayer(content)
modal := lipgloss.NewLayer(floatingBox).X(10).Y(5).Z(10).ID("modal")

// Compose and render
comp := lipgloss.NewCompositor(base, modal)
lipgloss.Println(comp.Render())

// Hit testing
hit := comp.Hit(x, y)
```

---

## New Border Types

```go
// v2 — new border constructors
lipgloss.MarkdownBorder()           // markdown table style
lipgloss.InnerHalfBlockBorder()     // half block inner
lipgloss.OuterHalfBlockBorder()     // half block outer
```

---

## New Color Functions

```go
// v2 — color manipulation
lipgloss.Alpha(c, 0.5)              // adjust alpha
lipgloss.Darken(c, 0.2)             // darken by percentage
lipgloss.Lighten(c, 0.2)            // lighten by percentage
lipgloss.Complementary(c)           // 180° on color wheel
lipgloss.Blend1D(steps, c1, c2)     // 1D gradient
lipgloss.Blend2D(w, h, angle, c1, c2) // 2D gradient
```

---

## Quick Reference Table

| Task | v1 | v2 |
|---|---|---|
| Import | `"github.com/charmbracelet/lipgloss"` | `"charm.land/lipgloss/v2"` |
| Create style | `lipgloss.NewStyle()` | `lipgloss.NewStyle()` |
| Hex color | `lipgloss.Color("#ff00ff")` | `lipgloss.Color("#ff00ff")` |
| ANSI color | `lipgloss.Color("5")` | `lipgloss.Color("5")` or `lipgloss.Magenta` |
| Adaptive | `lipgloss.AdaptiveColor{Light: "#fff", Dark: "#000"}` | `compat.AdaptiveColor{...}` or `LightDark(isDark)(light, dark)` |
| Detect dark bg | `lipgloss.HasDarkBackground()` | `lipgloss.HasDarkBackground(os.Stdin, os.Stdout)` |
| Print | `fmt.Println(s.Render("hi"))` | `lipgloss.Println(s.Render("hi"))` |
| Renderer | `renderer.NewStyle()` | `lipgloss.NewStyle()` |
| Whitespace fg | `WithWhitespaceForeground(c)` | `WithWhitespaceStyle(s.Foreground(c))` |
| Whitespace bg | `WithWhitespaceBackground(c)` | `WithWhitespaceStyle(s.Background(c))` |
| Underline | `s.Underline(true)` | `s.Underline(true)` or `s.UnderlineStyle(lipgloss.UnderlineCurly)` |
| No color | N/A | `lipgloss.NoColor{}` |
| Blend colors | N/A | `lipgloss.Blend1D()`, `lipgloss.Blend2D()` |
| Compositing | N/A | `lipgloss.NewCompositor()`, `lipgloss.NewLayer()` |

## Removed Symbols

| v1 Symbol | v2 Replacement |
|---|---|
| `type Renderer` | Removed |
| `DefaultRenderer()` | Not needed |
| `SetDefaultRenderer(r)` | Not needed |
| `NewRenderer(w, opts...)` | Not needed |
| `ColorProfile()` | `colorprofile.Detect(w, env)` |
| `SetColorProfile(p)` | Set `lipgloss.Writer.Profile` |
| `HasDarkBackground()` (no args) | `lipgloss.HasDarkBackground(in, out)` |
| `SetHasDarkBackground(b)` | Not needed — pass bool to `LightDark` |
| `type TerminalColor` | `image/color.Color` |
| `type Color string` | `func Color(string) color.Color` |
| `type ANSIColor uint` | `type ANSIColor = ansi.IndexedColor` |
| `type AdaptiveColor` | `compat.AdaptiveColor` or `LightDark` |
| `type CompleteColor` | `compat.CompleteColor` or `Complete` |
| `type CompleteAdaptiveColor` | `compat.CompleteAdaptiveColor` |
| `WithWhitespaceForeground(c)` | `WithWhitespaceStyle(s)` |
| `WithWhitespaceBackground(c)` | `WithWhitespaceStyle(s)` |
| `renderer.NewStyle()` | `lipgloss.NewStyle()` |

## New Symbols in v2

| Symbol | Description |
|---|---|
| `lipgloss.NoColor{}` | Absence of color |
| `lipgloss.Alpha(c, a)` | Adjust color alpha |
| `lipgloss.Darken(c, pct)` | Darken color |
| `lipgloss.Lighten(c, pct)` | Lighten color |
| `lipgloss.Complementary(c)` | Complementary color |
| `lipgloss.Blend1D(n, colors...)` | 1D color gradient |
| `lipgloss.Blend2D(w, h, angle, colors...)` | 2D color gradient |
| `lipgloss.NewLayer(content)` | Create compositing layer |
| `lipgloss.NewCompositor(layers...)` | Compose layers |
| `lipgloss.NewCanvas(w, h)` | Low-level cell buffer |
| `lipgloss.StrikethroughSpaces(bool)` | Style method |
| `lipgloss.UnderlineSpaces(bool)` | Style method |
| `list.Alphabet` | Enumerator |
| `list.Asterisk` | Enumerator |
| `tree.Root(root)` | Constructor shorthand |
| `tree.RoundedEnumerator` | Tree branch style |
