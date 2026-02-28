package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"framework/internal/router"
	"framework/internal/theme"
)

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleResize(msg)
	case tea.BackgroundColorMsg:
		return m.handleBgColor(msg)
	case theme.ThemeSwitchMsg:
		return m.handleThemeSwitch(msg)
	case theme.ThemeChangedMsg:
		return m.handleThemeChanged(msg)
	case router.NavigateMsg:
		cmd := m.router.Navigate(msg.Screen)
		return m, cmd
	case router.BackMsg:
		m.router.Back()
		return m, nil
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	default:
		return m.routeToScreen(msg)
	}
}

func (m rootModel) handleResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	m.help.SetWidth(msg.Width - 2) // account for footer padding
	m.styles = theme.NewStyles(m.registry.Colors(), m.width)
	m.ready = true
	return m.routeToScreen(msg)
}

func (m rootModel) handleBgColor(msg tea.BackgroundColorMsg) (tea.Model, tea.Cmd) {
	isDark := msg.IsDark()
	m.registry.SetDark(isDark)
	m.help.Styles = help.DefaultStyles(isDark)
	colors := m.registry.Colors()
	m.styles = theme.NewStyles(colors, m.width)
	return m.routeToScreen(theme.ThemeChangedMsg{
		Colors: colors,
		Name:   m.registry.CurrentName(),
		IsDark: isDark,
	})
}

func (m rootModel) handleThemeSwitch(msg theme.ThemeSwitchMsg) (tea.Model, tea.Cmd) {
	if !m.registry.SetCurrent(msg.Name) {
		return m, nil
	}
	colors := m.registry.Colors()
	m.styles = theme.NewStyles(colors, m.width)
	return m.routeToScreen(theme.ThemeChangedMsg{
		Colors: colors,
		Name:   msg.Name,
		IsDark: m.registry.IsDark(),
	})
}

func (m rootModel) handleThemeChanged(msg theme.ThemeChangedMsg) (tea.Model, tea.Cmd) {
	m.styles = theme.NewStyles(msg.Colors, m.width)
	return m.routeToScreen(msg)
}

func (m rootModel) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Back):
		if m.router.Depth() > 0 {
			m.router.Back()
			return m, nil
		}
		return m, tea.Quit
	case key.Matches(msg, m.keys.ThemeCycle):
		return m.cycleTheme()
	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll
		return m, nil
	default:
		return m.routeToScreen(msg)
	}
}

func (m rootModel) cycleTheme() (tea.Model, tea.Cmd) {
	m.registry.CycleNext()
	colors := m.registry.Colors()
	m.styles = theme.NewStyles(colors, m.width)
	return m.routeToScreen(theme.ThemeChangedMsg{
		Colors: colors,
		Name:   m.registry.CurrentName(),
		IsDark: m.registry.IsDark(),
	})
}

func (m rootModel) routeToScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	screen, cmd := m.router.Current().Update(msg)
	m.router.SetCurrent(screen)
	return m, cmd
}
