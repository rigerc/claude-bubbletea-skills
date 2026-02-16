// Package ui provides the BubbleTea UI model for the application.
package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	appkeys "template-v2-enhanced/internal/ui/keys"
)

// keyMap wraps the shared key bindings for the main UI model.
// It implements help.KeyMap for automatic help generation.
type keyMap struct {
	common appkeys.Common
}

// ShortHelp returns the key bindings to show in the short help view.
func (k keyMap) ShortHelp() []key.Binding {
	return k.common.ShortHelp()
}

// FullHelp returns the key bindings to show in the full help view.
func (k keyMap) FullHelp() [][]key.Binding {
	return k.common.FullHelp()
}

// defaultKeyMap returns the default key bindings for the main UI.
func defaultKeyMap() keyMap {
	return keyMap{
		common: appkeys.CommonBindings(),
	}
}

// newHelpModel creates a new help component with shared bindings.
func newHelpModel() help.Model {
	h := help.New()
	h.ShowAll = true
	return h
}

// handleKeyPress handles a key press message and returns any commands to execute.
// This is a helper function that can be used in the Update method.
func handleKeyPress(msg tea.KeyPressMsg, keys keyMap) tea.Cmd {
	switch {
	case key.Matches(msg, keys.common.Quit):
		return tea.Quit
	case key.Matches(msg, keys.common.Help):
		// Toggle help visibility
		return nil
	default:
		return nil
	}
}
