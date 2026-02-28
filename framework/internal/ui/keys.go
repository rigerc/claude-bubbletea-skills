package ui

import "charm.land/bubbles/v2/key"

// globalKeyMap defines application-wide keybindings.
type globalKeyMap struct {
	Quit       key.Binding
	Back       key.Binding
	ThemeCycle key.Binding
	Help       key.Binding
}

func defaultGlobalKeyMap() globalKeyMap {
	return globalKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		ThemeCycle: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "theme"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}

func (k globalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.ThemeCycle, k.Quit}
}

func (k globalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Back, k.ThemeCycle, k.Help, k.Quit}}
}
