// Package styles provides adaptive lipgloss styling for the application UI.
package styles

import (
	lipglossv2 "charm.land/lipgloss/v2"
	"image/color"
)

// Theme contains all lipgloss styles used throughout the application.
// Each style is constructed using lipgloss.LightDark() for adaptive colors
// that work well on both light and dark terminal backgrounds.
type Theme struct {
	App           lipglossv2.Style // outer container: Margin(1, 2)
	Title         lipglossv2.Style // list title bar (overrides list default)
	StatusMessage lipglossv2.Style // list status bar message
	HelpBar       lipglossv2.Style // help bar at bottom
	Detail        lipglossv2.Style // body text in detail screens
	Subtle        lipglossv2.Style // de-emphasized text
}

// New constructs a Theme using lipgloss.LightDark() for adaptive colors.
// The isDark parameter should come from tea.BackgroundColorMsg.IsDark().
func New(isDark bool) Theme {
	ld := lipglossv2.LightDark(isDark)

	return Theme{
		App: lipglossv2.NewStyle().Margin(1, 2),
		Title: lipglossv2.NewStyle().
			Bold(true).
			Foreground(lipglossv2.Color("#FFFDF5")).
			Background(lipglossv2.Color("#25A065")).
			Padding(0, 1),
		StatusMessage: lipglossv2.NewStyle().
			Foreground(ld(
				lipglossv2.Color("#04B575"),
				lipglossv2.Color("#10CC85"),
			)),
		HelpBar: lipglossv2.NewStyle().
			Foreground(ld(
				lipglossv2.Color("#626262"),
				lipglossv2.Color("#9B9B9B"),
			)),
		Detail:  lipglossv2.NewStyle().Margin(0, 2),
		Subtle: lipglossv2.NewStyle().
			Foreground(ld(
				lipglossv2.Color("#9B9B9B"),
				lipglossv2.Color("#626262"),
			)),
	}
}

// HelpBarHeight returns the number of lines consumed by the help bar.
// It returns 1 for short help and 3 for full (expanded) help.
func HelpBarHeight(showAll bool) int {
	if showAll {
		return 3
	}
	return 1
}

// Color is a convenience function that creates an adaptive color.
// It returns the lightColor on light backgrounds and darkColor on dark backgrounds.
func Color(isDark bool, lightColor, darkColor string) color.Color {
	ld := lipglossv2.LightDark(isDark)
	return ld(lipglossv2.Color(lightColor), lipglossv2.Color(darkColor))
}

// LightDark returns a color selector function based on the isDark flag.
// The returned function takes two colors (light, dark) and returns the
// appropriate one for the current terminal background.
func LightDark(isDark bool) func(light, dark color.Color) color.Color {
	return lipglossv2.LightDark(isDark)
}
