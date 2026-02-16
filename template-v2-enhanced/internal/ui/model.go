// Package ui provides the BubbleTea UI model for the application.
// This file contains the main model implementation following the Elm Architecture.
package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	applogger "template-v2-enhanced/internal/logger"
	"template-v2-enhanced/config"
	"template-v2-enhanced/internal/ui/nav"
	"template-v2-enhanced/internal/ui/screens"
)

// errMsg wraps an error for use as a BubbleTea message.
type errMsg error

// Model represents the application state.
// It holds all the state for the terminal UI application.
type Model struct {
	// navStack is the navigation stack that manages screen transitions.
	navStack nav.Stack

	// quitting indicates the app is about to exit.
	quitting bool

	// err holds any error that occurs during execution.
	err error

	// isDark indicates if the terminal has a dark background.
	isDark bool

	// config holds the application configuration.
	config *config.Config
}

// New creates a new model with default values.
// It initializes the navigation stack with the home screen as the root.
func New(cfg *config.Config) Model {
	// Create the home screen as the root with config
	homeScreen := screens.NewHomeScreen(cfg.UI.AltScreen)

	// Create navigation stack with home screen as root
	stack := nav.NewStack(homeScreen)

	return Model{
		navStack: stack,
		config:   cfg,
	}
}

// Init initializes the model and returns the initial command.
// It requests the terminal background color.
func (m Model) Init() tea.Cmd {
	applogger.Debug().Msg("Initializing UI model")

	return tea.Batch(
		m.navStack.Init(),
		tea.RequestBackgroundColor,
	)
}

// Update handles incoming messages and updates the model state.
// It returns the updated model and any commands to execute.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	applogger.Trace().Msgf("Update called with message type: %T", msg)

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		// Handle terminal background color change
		m.isDark = msg.IsDark()
		applogger.Debug().Msgf("Background color detected: isDark=%v", m.isDark)

		// Forward to navigation stack
		var updatedModel tea.Model
		var cmd tea.Cmd
		updatedModel, cmd = m.navStack.Update(msg)
		m.navStack = updatedModel.(nav.Stack)
		return m, cmd

	case tea.KeyPressMsg:
		// Let the navigation stack handle key press messages first
		// This allows screens to handle their own key bindings
		var updatedModel tea.Model
		var cmd tea.Cmd
		updatedModel, cmd = m.navStack.Update(msg)
		m.navStack = updatedModel.(nav.Stack)

		// If the command is nil, check if we should handle it at the model level
		if cmd == nil {
			// Handle quit key if not handled by screen
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				applogger.Debug().Msg("Quit key pressed, exiting")
				m.quitting = true
				return m, tea.Quit
			}
		}

		return m, cmd

	case tea.WindowSizeMsg:
		// Forward window resize to navigation stack
		var updatedModel tea.Model
		var cmd tea.Cmd
		updatedModel, cmd = m.navStack.Update(msg)
		m.navStack = updatedModel.(nav.Stack)
		return m, cmd

	case errMsg:
		// Handle error message
		m.err = msg
		applogger.Error().Msgf("Error received: %v", msg)
		return m, nil

	default:
		// Forward all other messages to the navigation stack
		var updatedModel tea.Model
		var cmd tea.Cmd
		updatedModel, cmd = m.navStack.Update(msg)
		m.navStack = updatedModel.(nav.Stack)
		return m, cmd
	}
}

// View renders the model state as a tea.View.
// It returns the view from the active screen in the navigation stack.
func (m Model) View() tea.View {
	if m.err != nil {
		v := tea.NewView(m.renderError(m.err.Error()))
		v.AltScreen = m.config.UI.AltScreen
		return v
	}

	// Get view from navigation stack
	return m.navStack.View()
}

// renderError builds the visual representation of an error state.
func (m Model) renderError(errMsg string) string {
	ld := lipgloss.LightDark(m.isDark)

	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF0000")).
		MarginLeft(2).
		MarginTop(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(ld(
			lipgloss.Color("#777777"),
			lipgloss.Color("#BBBBBB"),
		)).
		MarginLeft(2).
		MarginTop(1)

	return fmt.Sprintf(
		"%s\n\n%s",
		errorStyle.Render("Error: "+errMsg),
		helpStyle.Render("Press q to quit"),
	)
}

// Run starts the BubbleTea program with the given model.
// It handles any errors that occur during execution.
func Run(m Model) error {
	applogger.Info().Msg("Starting BubbleTea program")

	// Create and run the program
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		applogger.Error().Err(err).Msg("Program failed")
		return fmt.Errorf("running program: %w", err)
	}

	applogger.Info().Msg("Program exited successfully")
	return nil
}

// RunWithConfig starts the BubbleTea program with additional configuration options.
// This allows for customizing the program behavior before starting.
func RunWithConfig(m Model, opts ...tea.ProgramOption) error {
	applogger.Info().Msg("Starting BubbleTea program with custom options")

	p := tea.NewProgram(m, opts...)
	if _, err := p.Run(); err != nil {
		applogger.Error().Err(err).Msg("Program failed")
		return fmt.Errorf("running program: %w", err)
	}

	return nil
}
