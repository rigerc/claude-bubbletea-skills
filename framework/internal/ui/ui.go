package ui

import (
	tea "charm.land/bubbletea/v2"

	"framework/internal/screens"
)

// New creates the root model with the home screen as initial view.
func New() tea.Model {
	return newRootModel(screens.NewHome())
}
