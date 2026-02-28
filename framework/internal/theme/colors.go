package theme

import "image/color"

// ColorSet holds all semantic color slots for one variant (light or dark).
type ColorSet struct {
	Primary   color.Color
	Secondary color.Color
	Text      color.Color
	Subtle    color.Color
	Surface   color.Color
	Overlay   color.Color
	Error     color.Color
	Warning   color.Color
	Success   color.Color
	Border    color.Color
}

// Theme holds a named theme with light and dark color sets.
type Theme struct {
	Name  string
	Light ColorSet
	Dark  ColorSet
}

// Colors selects the correct ColorSet based on the isDark flag.
func (t Theme) Colors(isDark bool) ColorSet {
	if isDark {
		return t.Dark
	}
	return t.Light
}
