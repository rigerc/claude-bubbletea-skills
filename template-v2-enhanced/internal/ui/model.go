// Package ui provides the BubbleTea UI model for the application.
// It implements a stack-based navigation router with theme support.
package ui

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"template-v2-enhanced/config"
	applogger "template-v2-enhanced/internal/logger"
	"template-v2-enhanced/internal/ui/nav"
	"template-v2-enhanced/internal/ui/screens"
)

// Model represents the application state with a navigation stack.
type Model struct {
	// screens holds the navigation stack. The last element is the active screen.
	screens []nav.Screen

	// width and height store the current terminal dimensions.
	width, height int

	// isDark indicates if the terminal has a dark background.
	isDark bool

	// quitting is set to true when the app is about to exit.
	quitting bool

	// Config-derived fields (extracted from config.Config at construction).
	altScreen    bool
	mouseEnabled bool
	windowTitle  string
}

// New creates a new Model with the provided configuration.
// It accepts config.Config as a value type (main.go passes *cfg dereferenced).
func New(cfg config.Config) Model {
	// Define sample content for detail screens.
	detailContent := `This is a detail screen.

You can scroll this content using:
  • j/k or up/down arrows — line by line
  • page up/down       — page by page
  • home/end           — jump to top/bottom

Press ESC to return to the menu.`

	aboutContent := `template-v2-enhanced

A BubbleTea v2 skeleton application with:
  • Stack-based navigation
  • Adaptive light/dark theming
  • List-based menu navigation
  • Scrollable detail screens

Built with:
  • charm.land/bubbletea/v2
  • charm.land/bubbles/v2
  • charm.land/lipgloss/v2
  • github.com/spf13/cobra
  • github.com/knadh/koanf/v2
  • github.com/rs/zerolog

Press ESC to return to the menu.`

	// Create menu items.
	items := []list.Item{
		screens.NewMenuItem("Details", "View a detail screen",
			nav.Push(screens.NewDetailScreen("Details", detailContent, false))),
		screens.NewMenuItem("About", "About this application",
			nav.Push(screens.NewDetailScreen("About", aboutContent, false))),
	}

	root := screens.NewMenuScreen(cfg.App.Title, items, false)

	return Model{
		screens:      []nav.Screen{root},
		altScreen:    cfg.UI.AltScreen,
		mouseEnabled: cfg.UI.MouseEnabled,
		windowTitle:  cfg.App.Title,
	}
}

// Init returns the initial command. It requests the terminal background color
// and initializes the root screen.
func (m Model) Init() tea.Cmd {
	applogger.Debug().Msg("Initializing UI model")
	cmds := []tea.Cmd{tea.RequestBackgroundColor}
	if len(m.screens) > 0 {
		cmds = append(cmds, m.screens[len(m.screens)-1].Init())
	}
	return tea.Batch(cmds...)
}

// Update handles incoming messages and returns an updated model and command.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			applogger.Debug().Msg("Quit key pressed")
			m.quitting = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		applogger.Debug().Msgf("Window resized: %dx%d", m.width, m.height)
		// fall through to delegate to active screen

	case tea.BackgroundColorMsg:
		m.isDark = msg.IsDark()
		applogger.Debug().Msgf("Background color detected: isDark=%v", m.isDark)
		// Propagate theme to ALL screens in stack
		for i := range m.screens {
			if t, ok := m.screens[i].(nav.Themeable); ok {
				t.SetTheme(m.isDark)
			}
		}
		// fall through to deliver msg to active screen

	case nav.PushMsg:
		s := msg.Screen
		if cmd := s.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
		if t, ok := s.(nav.Themeable); ok {
			t.SetTheme(m.isDark)
		}
		s, cmd := s.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		cmds = append(cmds, cmd)
		m.screens = append(m.screens, s)
		return m, tea.Batch(cmds...)

	case nav.PopMsg:
		if len(m.screens) > 1 {
			m.screens = m.screens[:len(m.screens)-1]
		}
		return m, nil

	case nav.ReplaceMsg:
		if len(m.screens) > 0 {
			s := msg.Screen
			if cmd := s.Init(); cmd != nil {
				cmds = append(cmds, cmd)
			}
			if t, ok := s.(nav.Themeable); ok {
				t.SetTheme(m.isDark)
			}
			s, cmd := s.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
			cmds = append(cmds, cmd)
			m.screens[len(m.screens)-1] = s
		}
		return m, tea.Batch(cmds...)
	}

	// Delegate to active screen
	if len(m.screens) > 0 {
		top := m.screens[len(m.screens)-1]
		updated, cmd := top.Update(msg)
		m.screens[len(m.screens)-1] = updated
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the current model state as a tea.View.
func (m Model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	var content string
	if len(m.screens) > 0 {
		content = m.screens[len(m.screens)-1].View()
	}

	v := tea.NewView(content)
	v.AltScreen = m.altScreen   // from cfg.UI.AltScreen
	v.WindowTitle = m.windowTitle // from cfg.App.Title
	if m.mouseEnabled {          // from cfg.UI.MouseEnabled
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
