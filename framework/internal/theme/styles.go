package theme

import "charm.land/lipgloss/v2"

// Styles holds pre-built lipgloss styles derived from a ColorSet.
// These are value types — rebuild entirely on theme change.
type Styles struct {
	Header    lipgloss.Style
	Title     lipgloss.Style
	NavBar    lipgloss.Style
	NavItem   lipgloss.Style
	NavActive lipgloss.Style
	Body      lipgloss.Style
	Footer    lipgloss.Style
	Error     lipgloss.Style
	Warning   lipgloss.Style
	Success   lipgloss.Style
}

// NewStyles builds Styles from a ColorSet and terminal width.
func NewStyles(c ColorSet, width int) Styles {
	return Styles{
		Header: lipgloss.NewStyle().
			Foreground(c.Text).
			Background(c.Surface).
			Width(width).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(c.Primary),
		NavBar: lipgloss.NewStyle().
			Background(c.Surface).
			Width(width).
			Padding(0, 1),
		NavItem: lipgloss.NewStyle().
			Foreground(c.Subtle).
			Padding(0, 1),
		NavActive: lipgloss.NewStyle().
			Foreground(c.Primary).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.Border{Bottom: "─"}, false, false, true, false).
			BorderForeground(c.Primary),
		Body: lipgloss.NewStyle().
			Foreground(c.Text),
		Footer: lipgloss.NewStyle().
			Foreground(c.Subtle).
			Width(width).
			Padding(0, 1),
		Error: lipgloss.NewStyle().
			Foreground(c.Error).
			Bold(true),
		Warning: lipgloss.NewStyle().
			Foreground(c.Warning),
		Success: lipgloss.NewStyle().
			Foreground(c.Success),
	}
}
