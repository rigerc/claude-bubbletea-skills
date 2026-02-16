// Package keys provides shared key bindings for the application.
// It uses bubbles v2 key.Binding for consistent key definition.
package keys

import (
	"charm.land/bubbles/v2/key"
)

// Common holds shared key bindings for the entire application.
// These bindings are used across all screens for consistency.
type Common struct {
	Quit  key.Binding
	Back  key.Binding
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Space key.Binding
	Help  key.Binding
	Esc   key.Binding
	Left  key.Binding
	Right key.Binding
}

// CommonBindings returns the default shared key bindings.
// These are the standard bindings used throughout the application.
func CommonBindings() Common {
	return Common{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Back: key.NewBinding(
			key.WithKeys("b", "esc", "left"),
			key.WithHelp("b", "back"),
		),
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "move down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		Space: key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "toggle"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Left: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("→/l", "move right"),
		),
	}
}

// ShortHelp returns bindings for short help display.
// This is typically shown at the bottom of the screen.
func (c Common) ShortHelp() []key.Binding {
	return []key.Binding{c.Quit, c.Back, c.Help}
}

// FullHelp returns bindings for full help display.
// This is typically shown when the user presses '?'.
func (c Common) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{c.Up, c.Down, c.Left, c.Right},
		{c.Enter, c.Space, c.Esc},
		{c.Quit, c.Back, c.Help},
	}
}
