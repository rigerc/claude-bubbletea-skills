// Package theme provides styling for the TUI.
package theme

import "charm.land/lipgloss/v2"

// Styles holds all styled components for the UI.
type Styles struct {
	App         lipgloss.Style
	Header      lipgloss.Style
	Body        lipgloss.Style
	Help        lipgloss.Style
	Footer      lipgloss.Style
	StatusLeft  lipgloss.Style
	StatusRight lipgloss.Style
	MaxWidth    int
}

// New creates a Styles struct with adaptive colors based on the background.
func New(isDark bool, width int) Styles {
	ld := lipgloss.LightDark(isDark)
	subtle := ld(lipgloss.Color("#555555"), lipgloss.Color("#999999"))
	accent := lipgloss.Color("#7D56F4")
	fg := ld(lipgloss.Color("#1a1a1a"), lipgloss.Color("#f1f1f1"))

	// Max width is 50% of terminal width
	maxWidth := width * 50 / 100
	if maxWidth < 40 {
		maxWidth = width - 4 // Minimum usable width
	}

	return Styles{
		MaxWidth: maxWidth,
		App: lipgloss.NewStyle().
			Width(maxWidth).
			Padding(0, 0),
		Header: lipgloss.NewStyle().
			Padding(5).
			PaddingBottom(1),
		Body: lipgloss.NewStyle().
			Padding(0, 3).
			Foreground(fg),
		Help: lipgloss.NewStyle().
			MarginTop(0).
			Padding(0, 3),
		Footer: lipgloss.NewStyle().
			MarginTop(1).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(subtle).
			PaddingLeft(1),
		StatusLeft: lipgloss.NewStyle().
			Background(accent).
			Foreground(lipgloss.Color("#ffffff")).
			Bold(true),
		StatusRight: lipgloss.NewStyle().
			Foreground(subtle),
	}
}

// DetailStyles holds styles for the detail screen.
type DetailStyles struct {
	Title   lipgloss.Style
	Desc    lipgloss.Style
	Content lipgloss.Style
	Info    lipgloss.Style
}

// NewDetailStyles creates detail screen styles with adaptive colors.
func NewDetailStyles(isDark bool) DetailStyles {
	ld := lipgloss.LightDark(isDark)
	accent := lipgloss.Color("#7D56F4")

	return DetailStyles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(accent).
			MarginBottom(1),
		Desc: lipgloss.NewStyle().
			Foreground(ld(lipgloss.Color("#555555"), lipgloss.Color("#999999"))).
			MarginBottom(2),
		Content: lipgloss.NewStyle().
			Foreground(ld(lipgloss.Color("#333333"), lipgloss.Color("#f1f1f1"))),
		Info: lipgloss.NewStyle().
			Foreground(ld(lipgloss.Color("#888888"), lipgloss.Color("#666666"))).
			Italic(true),
	}
}
