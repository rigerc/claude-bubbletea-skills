// Package styles provides style factories using lipgloss v2's LightDark function.
package styles

import (
	"charm.land/lipgloss/v2"
)

// CommonStyles provides shared styles used across all screens.
// These styles are pre-configured with the application's color scheme
// and adapt automatically to light/dark themes.
type CommonStyles struct {
	// Title is for main screen titles.
	Title lipgloss.Style

	// Subtitle is for secondary headings.
	Subtitle lipgloss.Style

	// Header is for the sticky header at the top of screens.
	Header lipgloss.Style

	// Content is for body text.
	Content lipgloss.Style

	// Label is for labels and key-value pairs.
	Label lipgloss.Style

	// Selected is for the currently selected item in lists.
	Selected lipgloss.Style

	// Help is for help and instruction text.
	Help lipgloss.Style

	// Border is for bordered containers.
	Border lipgloss.Style
}

// NewCommonStyles creates common styles adapted for the current theme.
// It uses lipgloss.LightDark to select appropriate colors based on isDark.
func NewCommonStyles(isDark bool) CommonStyles {
	ld := lipgloss.LightDark(isDark)

	return CommonStyles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(ld(lipgloss.Color("#644AED"), lipgloss.Color("#874BFD"))).
			MarginBottom(1),

		Subtitle: lipgloss.NewStyle().
			Foreground(ld(lipgloss.Color("#555555"), lipgloss.Color("#CCCCCC"))).
			MarginBottom(2),

		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(ld(lipgloss.Color("#FFFFFF"), lipgloss.Color("#FFFFFF"))).
			Background(ld(lipgloss.Color("#644AED"), lipgloss.Color("#874BFD"))).
			Padding(0, 1).
			Width(80), // Will be adjusted per screen

		Content: lipgloss.NewStyle().
			Foreground(ld(lipgloss.Color("#555555"), lipgloss.Color("#CCCCCC"))).
			MarginLeft(1),

		Label: lipgloss.NewStyle().
			Bold(true).
			Foreground(ld(lipgloss.Color("#205A9C"), lipgloss.Color("#3F92DF"))).
			MarginLeft(1),

		Selected: lipgloss.NewStyle().
			Bold(true).
			Foreground(ld(lipgloss.Color("#205A9C"), lipgloss.Color("#3F92DF"))).
			Background(ld(lipgloss.Color("#E8E8E8"), lipgloss.Color("#2A2A2A"))).
			Padding(0, 1).
			MarginLeft(1),

		Help: lipgloss.NewStyle().
			Foreground(ld(lipgloss.Color("#666666"), lipgloss.Color("#AAAAAA"))).
			MarginTop(1),

		Border: lipgloss.NewStyle().
			Padding(1, 2),
	}
}

// MenuStyles provides styles specific to menu screens.
type MenuStyles struct {
	CommonStyles
	Option lipgloss.Style
}

// NewMenuStyles creates menu styles adapted for the current theme.
func NewMenuStyles(isDark bool) MenuStyles {
	ld := lipgloss.LightDark(isDark)
	common := NewCommonStyles(isDark)

	return MenuStyles{
		CommonStyles: common,
		Option: lipgloss.NewStyle().
			Foreground(ld(lipgloss.Color("#555555"), lipgloss.Color("#CCCCCC"))).
			MarginLeft(2),
	}
}

// ContentStyles provides styles for content display screens.
type ContentStyles struct {
	CommonStyles
	Section lipgloss.Style
}

// NewContentStyles creates content styles adapted for the current theme.
func NewContentStyles(isDark bool) ContentStyles {
	return ContentStyles{
		CommonStyles: NewCommonStyles(isDark),
		Section: lipgloss.NewStyle().
			MarginBottom(1),
	}
}
