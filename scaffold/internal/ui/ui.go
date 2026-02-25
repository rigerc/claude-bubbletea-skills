// Package ui provides the TUI entry point for scaffold.
package ui

import (
	tea "charm.land/bubbletea/v2"

	"scaffold/config"
)

// New creates a new root model from the config.
// configPath is the path to persist settings; empty means no file save.
func New(cfg config.Config, configPath string) rootModel {
	return newRootModel(cfg, configPath)
}

// Run starts the TUI program.
func Run(m rootModel) error {
	_, err := tea.NewProgram(m).Run()
	return err
}
