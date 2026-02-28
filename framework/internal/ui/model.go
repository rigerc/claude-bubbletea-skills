package ui

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"

	"framework/internal/router"
	"framework/internal/theme"
)

// NavItem defines a top-level navigation destination.
type NavItem struct {
	Label  string
	Screen func() router.Screen // factory — creates a fresh screen on navigation
}

type rootModel struct {
	router   *router.Router
	registry *theme.Registry
	styles   theme.Styles
	help     help.Model
	keys     globalKeyMap

	navItems  []NavItem
	navCursor int

	width  int
	height int
	ready  bool
}

func newRootModel(navItems []NavItem) rootModel {
	if len(navItems) == 0 {
		panic("at least one nav item required")
	}

	reg := theme.NewRegistry()
	colors := reg.Colors()

	return rootModel{
		router:   router.New(navItems[0].Screen()),
		registry: reg,
		styles:   theme.NewStyles(colors, 0),
		keys:     defaultGlobalKeyMap(),
		help:     help.New(),
		navItems: navItems,
	}
}

func (m rootModel) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
		m.router.Current().Init(),
	)
}
