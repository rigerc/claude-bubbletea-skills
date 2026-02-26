package screens

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"scaffold/internal/ui/theme"
)

// Detail is a detail screen that shows information about a selected menu item.
type Detail struct {
	theme.ThemeAware

	title       string
	description string
	screenID    string
	width       int
}

// NewDetail creates a new Detail screen.
func NewDetail(title, description, screenID string) Detail {
	return Detail{
		title:       title,
		description: description,
		screenID:    screenID,
	}
}

// SetWidth sets the screen width.
func (d Detail) SetWidth(w int) Screen {
	d.width = w
	return d
}

// ApplyTheme implements theme.Themeable.
func (d *Detail) ApplyTheme(state theme.State) {
	d.ApplyThemeState(state)
	// Styles built in View() from Palette()
}

// Init initializes the detail screen.
func (d Detail) Init() tea.Cmd {
	return nil
}

// Update handles messages for the detail screen.
func (d Detail) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch keyMsg.String() {
		case "esc":
			return d, func() tea.Msg { return BackMsg{} }
		}
	}
	return d, nil
}

// View renders the detail screen.
func (d Detail) View() tea.View {
	return tea.NewView(d.Body())
}

// Body returns the body content for layout composition.
func (d Detail) Body() string {
	p := d.Palette()

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(p.Primary).MarginBottom(1)
	descStyle := lipgloss.NewStyle().Foreground(p.TextMuted).MarginBottom(2)
	contentStyle := lipgloss.NewStyle().Foreground(p.TextPrimary)
	infoStyle := lipgloss.NewStyle().Foreground(p.TextSecondary).Italic(true)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(d.title),
		descStyle.Render(d.description),
		contentStyle.Render(fmt.Sprintf("Screen ID: %s", d.screenID)),
		"",
		infoStyle.Render("Press Esc to go back to the menu"),
	)

	return content
}
