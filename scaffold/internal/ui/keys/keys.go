// Package keys provides global key bindings for the TUI.
package keys

import "charm.land/bubbles/v2/key"

// GlobalKeyMap holds global key bindings.
type GlobalKeyMap struct {
	Quit key.Binding
}

// DefaultGlobalKeyMap returns the default global key bindings.
func DefaultGlobalKeyMap() GlobalKeyMap {
	return GlobalKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
	}
}

// ShortHelp returns a slice of bindings for short help view.
func (k GlobalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

// FullHelp returns grouped bindings for full help view.
func (k GlobalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Quit}}
}
