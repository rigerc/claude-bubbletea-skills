package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"framework/internal/router"
)

func (m rootModel) View() tea.View {
	if !m.ready {
		return tea.NewView("")
	}

	header := m.renderHeader()
	footer := m.renderFooter()

	headerH := lipgloss.Height(header)
	footerH := lipgloss.Height(footer)
	bodyH := max(1, m.height-headerH-footerH)

	body := m.router.Current().View(m.width, bodyH)
	bodyStyled := m.styles.Body.
		Width(m.width).
		Height(bodyH).
		MaxHeight(bodyH).
		Padding(1, 2).
		Render(body)

	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		bodyStyled,
		footer,
	)

	v := tea.NewView(content)
	v.AltScreen = false
	return v
}

func (m rootModel) renderHeader() string {
	title := m.styles.Title.Render("Framework")
	route := m.styles.Route.Render(" > " + m.router.Current().Title())
	return m.styles.Header.Render(title + route)
}

func (m rootModel) renderFooter() string {
	km := combinedKeyMap{
		screen: m.router.Current().KeyMap(),
		global: m.keys,
	}
	return m.styles.Footer.Render(m.help.View(km))
}

// combinedKeyMap merges screen-specific and global keybindings for the help bar.
type combinedKeyMap struct {
	screen help.KeyMap
	global globalKeyMap
}

func (c combinedKeyMap) ShortHelp() []key.Binding {
	bindings := c.screen.ShortHelp()
	return append(bindings, c.global.ShortHelp()...)
}

func (c combinedKeyMap) FullHelp() [][]key.Binding {
	groups := c.screen.FullHelp()
	return append(groups, c.global.FullHelp()...)
}

// HeaderStyle returns the header style (for screens that need it).
func (m rootModel) HeaderStyle() lipgloss.Style { return m.styles.Header }

// CurrentScreen returns the current screen title for breadcrumb.
func (m rootModel) CurrentScreen() router.Screen { return m.router.Current() }
