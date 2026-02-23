// Package keys defines application-wide key bindings.
package keys

import "charm.land/bubbles/v2/key"

// GlobalKeyMap defines the application-wide key bindings.
// It implements help.KeyMap so it can be passed directly to help.View().
type GlobalKeyMap struct {
	Back key.Binding // "esc"    — go to previous screen
	Quit key.Binding // "ctrl+c" — always quit (no conflict with list filter "q")
	Help key.Binding // "?"      — toggle help expansion
}

// New returns a GlobalKeyMap with default bindings.
func New() GlobalKeyMap {
	return GlobalKeyMap{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

// ShortHelp returns the short-form key bindings for the help bar.
// Implements help.KeyMap.
func (k GlobalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

// FullHelp returns the full-form key bindings for the expanded help view.
// Implements help.KeyMap.
func (k GlobalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Back, k.Help}, {k.Quit}}
}

// DashboardKeyMap defines key bindings specific to the dashboard screen.
// It extends GlobalKeyMap with ralphio workflow control keys.
type DashboardKeyMap struct {
	Retry   key.Binding // "r" — retry current task
	Skip    key.Binding // "s" — skip current task
	Detail  key.Binding // "v" — view task detail
	Client  key.Binding // "c" — manage adapter/client
	History key.Binding // "h" — view iteration history
	Pause   key.Binding // "p" — pause/resume loop
	Quit    key.Binding // "q" or "ctrl+c" — quit
}

// NewDashboard returns a DashboardKeyMap with default bindings.
func NewDashboard() DashboardKeyMap {
	return DashboardKeyMap{
		Retry: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "retry"),
		),
		Skip: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "skip"),
		),
		Detail: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "detail"),
		),
		Client: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "client"),
		),
		History: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "history"),
		),
		Pause: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "pause"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns the short-form key bindings for the dashboard help bar.
// Implements help.KeyMap.
func (k DashboardKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Retry, k.Skip, k.Detail, k.Client, k.History, k.Pause, k.Quit}
}

// FullHelp returns the full-form key bindings for the expanded dashboard help view.
// Implements help.KeyMap.
func (k DashboardKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Retry, k.Skip, k.Pause},
		{k.Detail, k.Client, k.History},
		{k.Quit},
	}
}
