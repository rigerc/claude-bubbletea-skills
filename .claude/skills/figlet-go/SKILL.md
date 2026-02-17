---
name: figlet-go
description: >-
  Guide for using the figlet-go library to generate ASCII art text in Go applications.
  Use when working with the figlet-go package, creating ASCII art output, or integrating
  FIGlet functionality into Go CLI tools. Covers rendering, fonts, colors, and configuration.
---

# figlet-go Skill

Guide for using the figlet-go library (`github.com/lsferreira42/figlet-go/figlet`) to generate ASCII art text in Go applications.

## Installation

```bash
go get github.com/lsferreira42/figlet-go/figlet
```

## Quick Start

### Simple Rendering

```go
import (
    "fmt"
    "log"
    "github.com/lsferreira42/figlet-go/figlet"
)

// Basic usage with default font
result, err := figlet.Render("Hello!")
if err != nil {
    log.Fatal(err)
}
fmt.Print(result)
```

### With Specific Font

```go
// Using a specific font (146 fonts embedded)
result, err := figlet.RenderWithFont("Go!", "slant")
if err != nil {
    log.Fatal(err)
}
fmt.Print(result)
```

## API Patterns

### Functional Options Pattern

Use option functions for configuration:

```go
result, err := figlet.Render("Text",
    figlet.WithFont("big"),                    // Font name
    figlet.WithWidth(60),                      // Output width
    figlet.WithJustification(1),               // 0=left, 1=center, 2=right
    figlet.WithParser("terminal-color"),       // Output format
    figlet.WithColors(figlet.ColorRed, figlet.ColorGreen),
)
```

### Direct Config Usage

For advanced control, use the Config struct directly:

```go
cfg := figlet.New()
cfg.Fontname = "banner"
cfg.Outputwidth = 100

if err := cfg.LoadFont(); err != nil {
    log.Fatal(err)
}

result := cfg.RenderString("Config")
```

## Available Options

| Option | Description |
|--------|-------------|
| `WithFont(name)` | Set font (146 available, see ListFonts()) |
| `WithWidth(width)` | Set output width in characters |
| `WithJustification(n)` | 0=left, 1=center, 2=right |
| `WithParser(type)` | `terminal`, `terminal-color`, `html` |
| `WithColors(colors...)` | ANSI colors or TrueColor |

## Colors

### ANSI Colors

```go
result, err := figlet.Render("Colors!",
    figlet.WithColors(figlet.ColorRed, figlet.ColorGreen, figlet.ColorBlue),
)
```

### TrueColor (24-bit RGB)

```go
tcRed, _ := figlet.NewTrueColorFromHexString("FF0000")
tcGreen, _ := figlet.NewTrueColorFromHexString("00FF00")
result, err := figlet.Render("TrueColor",
    figlet.WithColors(tcRed, tcGreen),
)
```

## Utility Functions

```go
// List all available fonts
fonts := figlet.ListFonts()

// Get version info
version := figlet.GetVersion()      // String version
versionInt := figlet.GetVersionInt() // Integer version
```

## Popular Fonts

- `standard` - Default font
- `banner` - Large banner style
- `big` - Large block letters
- `slant` - Slanted text
- `shadow` - Shadow effect
- `script` - Script style
- `doom` - DOOM game style
- `starwars` - Star Wars style

## Output Formats

- `terminal` - Plain text (default)
- `terminal-color` - ANSI colored output
- `html` - HTML formatted output

## Example

See complete working example:
- [references/_examples/example.go](references/_examples/example.go)

## Resources

- GitHub: https://github.com/lsferreira42/figlet-go
- Documentation: https://lsferreira42.github.io/figlet-go/
- Library Docs: https://github.com/lsferreira42/figlet-go/blob/main/lib.md
