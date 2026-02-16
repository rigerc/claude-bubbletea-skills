// Package ui provides the BubbleTea UI model for the application.
// This file contains the minimal model skeleton following the Elm Architecture.
package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"

	"template-v2-enhanced/config"
	applogger "template-v2-enhanced/internal/logger"
)

// Model represents the application state.
type Model struct {
	// cfg drives UI behaviour (AltScreen, MouseEnabled, etc.).
	cfg config.Config

	// width and height store the current terminal dimensions.
	width  int
	height int

	// isDark indicates if the terminal has a dark background.
	isDark bool

	// quitting is set to true when the app is about to exit.
	quitting bool
}

// New creates a new Model with the provided configuration.
func New(cfg config.Config) Model {
	return Model{cfg: cfg}
}

// Init returns the initial command. It requests the terminal background color
// so that styles can be adapted to light or dark themes.
func (m Model) Init() tea.Cmd {
	applogger.Debug().Msg("Initializing UI model")
	return tea.RequestBackgroundColor
}

// Update handles incoming messages and returns an updated model and command.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			applogger.Debug().Msg("Quit key pressed")
			m.quitting = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		applogger.Debug().Msgf("Window resized: %dx%d", m.width, m.height)

	case tea.BackgroundColorMsg:
		m.isDark = msg.IsDark()
		applogger.Debug().Msgf("Background color detected: isDark=%v", m.isDark)
	}

	return m, nil
}

// View renders the current model state as a tea.View.
func (m Model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	content := fmt.Sprintf(
		"BubbleTea v2 skeleton\n\nTerminal: %dx%d  dark=%v\n\nPress q to quit.",
		m.width, m.height, m.isDark,
	)
	v := tea.NewView(content)
	v.AltScreen = m.cfg.UI.AltScreen
	if m.cfg.UI.MouseEnabled {
		v.MouseMode = tea.MouseModeCellMotion
	}
	return v
}

// Run starts the BubbleTea program with the given model.
func Run(m Model) error {
	applogger.Info().Msg("Starting BubbleTea program")

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("running program: %w", err)
	}

	applogger.Info().Msg("Program exited successfully")
	return nil
}
