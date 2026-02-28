package ui

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"

	"framework/internal/router"
	"framework/internal/theme"
)

type rootModel struct {
	router   *router.Router
	registry *theme.Registry
	styles   theme.Styles
	help     help.Model
	keys     globalKeyMap

	width  int
	height int
	ready  bool
}

func newRootModel(initial router.Screen) rootModel {
	reg := theme.NewRegistry()
	colors := reg.Colors()

	return rootModel{
		router:   router.New(initial),
		registry: reg,
		styles:   theme.NewStyles(colors, 0),
		keys:     defaultGlobalKeyMap(),
		help:     help.New(),
	}
}

func (m rootModel) Init() tea.Cmd {
	// Send initial theme state immediately so screens have colors on first render
	initialTheme := func() tea.Msg {
		return theme.ThemeChangedMsg{
			Colors: m.registry.Colors(),
			Name:   m.registry.CurrentName(),
			IsDark: m.registry.IsDark(),
		}
	}
	return tea.Batch(
		tea.RequestBackgroundColor,
		initialTheme,
		m.router.Current().Init(),
	)
}
