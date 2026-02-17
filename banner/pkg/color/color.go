// Package color provides RGB color manipulation and gradient functionality
// for terminal applications.
package color

import (
	"fmt"
	"strconv"
	"strings"
)

// RGB represents a color with red, green, and blue components.
type RGB struct {
	R, G, B uint8
}

// ANSI returns the ANSI escape sequence for setting the foreground color.
func (c RGB) ANSI() string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B)
}

// ANSIBg returns the ANSI escape sequence for setting the background color.
func (c RGB) ANSIBg() string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B)
}

// Hex returns the hex representation of the color.
func (c RGB) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// Blend returns a new color that is a blend between this color and another.
// t is a value between 0 and 1, where 0 returns this color and 1 returns other.
func (c RGB) Blend(other RGB, t float64) RGB {
	if t <= 0 {
		return c
	}
	if t >= 1 {
		return other
	}
	return RGB{
		R: uint8(float64(c.R) + t*(float64(other.R)-float64(c.R))),
		G: uint8(float64(c.G) + t*(float64(other.G)-float64(c.G))),
		B: uint8(float64(c.B) + t*(float64(other.B)-float64(c.B))),
	}
}

// FromHex parses a hex color string and returns an RGB value.
// Supports both 3-digit (#RGB) and 6-digit (#RRGGBB) formats.
func FromHex(hex string) RGB {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 3 {
		hex = string(hex[0]) + string(hex[0]) + string(hex[1]) + string(hex[1]) + string(hex[2]) + string(hex[2])
	}
	if len(hex) != 6 {
		return RGB{255, 255, 255}
	}

	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)

	return RGB{uint8(r), uint8(g), uint8(b)}
}

// ANSI escape sequences for text styling.
const (
	Reset     = "\x1b[0m"
	Bold      = "\x1b[1m"
	Dim       = "\x1b[2m"
	Italic    = "\x1b[3m"
	Underline = "\x1b[4m"
	Blink     = "\x1b[5m"
	Reverse   = "\x1b[7m"
	Hidden    = "\x1b[8m"
	Strike    = "\x1b[9m"
)

// Colorize applies a color to text.
func Colorize(text string, c RGB) string {
	return c.ANSI() + text + Reset
}

// ColorizeBold applies a bold color to text.
func ColorizeBold(text string, c RGB) string {
	return Bold + c.ANSI() + text + Reset
}

// ColorizeWithBg applies foreground and background colors to text.
func ColorizeWithBg(text string, fg, bg RGB) string {
	return fg.ANSI() + bg.ANSIBg() + text + Reset
}

// ApplyStyle applies multiple ANSI styles to text.
func ApplyStyle(text string, styles ...string) string {
	prefix := strings.Join(styles, "")
	return prefix + text + Reset
}

// Standard colors.
var (
	White     = RGB{255, 255, 255}
	Black     = RGB{0, 0, 0}
	Red       = RGB{239, 68, 68}
	Green     = RGB{34, 197, 94}
	Blue      = RGB{59, 130, 246}
	Yellow    = RGB{234, 179, 8}
	Cyan      = RGB{6, 182, 212}
	Magenta   = RGB{168, 85, 247}
	Orange    = RGB{249, 115, 22}
	Pink      = RGB{236, 72, 153}
	Gray      = RGB{107, 114, 128}
	DimGray   = RGB{75, 85, 99}
	LightGray = RGB{156, 163, 175}
	DarkGray  = RGB{55, 65, 81}
)
