// Package theme provides centralized theming for the application UI.
// It defines semantic color palettes and Lip Gloss styles that adapt
// to light and dark terminal backgrounds.
package theme

import (
	"image/color"
	lipgloss "charm.land/lipgloss/v2"
)

// ThemePalette defines semantic colors used throughout the application.
// Colors are adaptive and automatically adjust for light/dark backgrounds.
type ThemePalette struct {
	// Brand colors
	Primary   color.Color // Main brand color (green)
	PrimaryFg color.Color // Text color on primary background (#FFFDF5)

	// Accent colors
	Secondary color.Color // Secondary accent color (purple)

	// Semantic colors
	Success color.Color
	Warning color.Color
	Alert   color.Color

	// Text colors
	Text   color.Color // Primary text
	Muted  color.Color // Secondary text
	Subtle color.Color // De-emphasized text

	// UI elements
	Border color.Color

	// Status colors for task and loop state badges.
	StatusRunning  color.Color
	StatusPaused   color.Color
	StatusFailed   color.Color
	StatusComplete color.Color
	StatusSkipped  color.Color
	StatusPending  color.Color
}

// NewPalette creates a palette with adaptive colors using lipgloss.LightDark().
// The isDark parameter should come from tea.BackgroundColorMsg.IsDark().
func NewPalette(isDark bool) ThemePalette {
	ld := lipgloss.LightDark(isDark)
	return ThemePalette{
		// Brand colors - green theme
		Primary:   ld(lipgloss.Color("#04B575"), lipgloss.Color("#10CC85")),
		PrimaryFg: lipgloss.Color("#FFFDF5"), // Constant - white/cream text on green

		// Accent
		Secondary: ld(lipgloss.Color("#7D56F4"), lipgloss.Color("#9B7CFF")),

		// Semantic
		Success: ld(lipgloss.Color("#00CC66"), lipgloss.Color("#00FF9F")),
		Warning: ld(lipgloss.Color("#FFCC00"), lipgloss.Color("#FFD700")),
		Alert:   ld(lipgloss.Color("#FF5F87"), lipgloss.Color("#FF7AA3")),

		// Text
		Text:   ld(lipgloss.Color("#1A1A1A"), lipgloss.Color("#F0F0F0")),
		Muted:  ld(lipgloss.Color("#626262"), lipgloss.Color("#9B9B9B")),
		Subtle: ld(lipgloss.Color("#9B9B9B"), lipgloss.Color("#626262")),

		// UI
		Border: ld(lipgloss.Color("#D0D0D0"), lipgloss.Color("#3A3A3A")),

		// Status colors
		StatusRunning:  ld(lipgloss.Color("#0080FF"), lipgloss.Color("#4DA6FF")),
		StatusPaused:   ld(lipgloss.Color("#FFAA00"), lipgloss.Color("#FFD060")),
		StatusFailed:   ld(lipgloss.Color("#FF3030"), lipgloss.Color("#FF6060")),
		StatusComplete: ld(lipgloss.Color("#00AA44"), lipgloss.Color("#00DD66")),
		StatusSkipped:  ld(lipgloss.Color("#888888"), lipgloss.Color("#AAAAAA")),
		StatusPending:  ld(lipgloss.Color("#555555"), lipgloss.Color("#CCCCCC")),
	}
}
