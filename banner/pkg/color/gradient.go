package color

// Gradient represents a color gradient with multiple color stops.
type Gradient struct {
	Colors []RGB
}

// NewGradient creates a gradient from a start and end color.
func NewGradient(start, end RGB) Gradient {
	return Gradient{Colors: []RGB{start, end}}
}

// NewMultiGradient creates a gradient from multiple hex color strings.
func NewMultiGradient(hexColors ...string) Gradient {
	colors := make([]RGB, len(hexColors))
	for i, hex := range hexColors {
		colors[i] = FromHex(hex)
	}
	return Gradient{Colors: colors}
}

// NewGradientFromRGBs creates a gradient from multiple RGB colors.
func NewGradientFromRGBs(colors ...RGB) Gradient {
	return Gradient{Colors: colors}
}

// At returns the color at position t (0.0 to 1.0) in the gradient.
func (g Gradient) At(t float64) RGB {
	if len(g.Colors) == 0 {
		return White
	}
	if len(g.Colors) == 1 {
		return g.Colors[0]
	}
	if t <= 0 {
		return g.Colors[0]
	}
	if t >= 1 {
		return g.Colors[len(g.Colors)-1]
	}

	segments := float64(len(g.Colors) - 1)
	segment := int(t * segments)
	if segment >= len(g.Colors)-1 {
		segment = len(g.Colors) - 2
	}

	localT := (t*segments - float64(segment))
	return g.Colors[segment].Blend(g.Colors[segment+1], localT)
}

// Apply applies the gradient horizontally across a string.
func (g Gradient) Apply(text string) string {
	if len(text) == 0 {
		return ""
	}

	runes := []rune(text)
	result := ""

	for i, r := range runes {
		var t float64
		if len(runes) > 1 {
			t = float64(i) / float64(len(runes)-1)
		} else {
			t = 0.5
		}
		c := g.At(t)
		result += c.ANSI() + string(r)
	}

	return result + Reset
}

// ApplyLines applies the gradient horizontally across multiple lines.
// The gradient flows continuously from the first character of the first line
// to the last character of the last line.
func (g Gradient) ApplyLines(lines []string) []string {
	result := make([]string, len(lines))

	totalChars := 0
	for _, line := range lines {
		totalChars += len([]rune(line))
	}

	charIndex := 0
	for i, line := range lines {
		runes := []rune(line)
		coloredLine := ""

		for _, r := range runes {
			var t float64
			if totalChars > 1 {
				t = float64(charIndex) / float64(totalChars-1)
			} else {
				t = 0.5
			}
			c := g.At(t)
			coloredLine += c.ANSI() + string(r)
			charIndex++
		}

		result[i] = coloredLine + Reset
	}

	return result
}

// ApplyVertical applies the gradient vertically across lines.
// Each line gets a single color from the gradient.
func (g Gradient) ApplyVertical(lines []string) []string {
	result := make([]string, len(lines))

	for i, line := range lines {
		var t float64
		if len(lines) > 1 {
			t = float64(i) / float64(len(lines)-1)
		} else {
			t = 0.5
		}
		c := g.At(t)
		result[i] = c.ANSI() + line + Reset
	}

	return result
}

// ApplyDiagonal applies the gradient diagonally across lines.
// The gradient flows from top-left to bottom-right.
func (g Gradient) ApplyDiagonal(lines []string) []string {
	result := make([]string, len(lines))

	maxLen := 0
	for _, line := range lines {
		if len([]rune(line)) > maxLen {
			maxLen = len([]rune(line))
		}
	}

	for i, line := range lines {
		runes := []rune(line)
		coloredLine := ""

		for j, r := range runes {
			var diagonal float64
			if len(lines)+maxLen > 2 {
				diagonal = float64(i+j) / float64(len(lines)+maxLen-2)
			} else {
				diagonal = 0.5
			}
			c := g.At(diagonal)
			coloredLine += c.ANSI() + string(r)
		}

		result[i] = coloredLine + Reset
	}

	return result
}
