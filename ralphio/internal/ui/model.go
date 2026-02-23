// Package ui provides the BubbleTea UI model for the application.
// It implements a stack-based navigation router with theme support.
package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"

	"ralphio/config"
	applogger "ralphio/internal/logger"
	"ralphio/internal/orchestrator"
	"ralphio/internal/ui/nav"
	"ralphio/internal/ui/screens"
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

	// Orchestrator communication channels. Both may be nil when running without
	// an orchestrator (e.g. during unit tests or non-run subcommands).
	orchMsgCh <-chan tea.Msg
	orchCmdCh chan<- any

	// history accumulates IterationCompleteMsg entries for the history screen.
	history []orchestrator.IterationCompleteMsg
}

// New creates a new Model with the provided configuration.
// It accepts config.Config as a value type (main.go passes *cfg dereferenced).
func New(cfg config.Config) Model {
	root := screens.NewDashboardScreen(false, cfg.App.Name)

	return Model{
		screens:      []nav.Screen{root},
		altScreen:    cfg.UI.AltScreen,
		mouseEnabled: cfg.UI.MouseEnabled,
		windowTitle:  cfg.App.Title,
	}
}

// WithOrchestrator returns a Model option function that wires the bidirectional
// orchestrator channels into the model. Call before starting the program.
func WithOrchestrator(msgCh <-chan tea.Msg, cmdCh chan<- any) func(*Model) {
	return func(m *Model) {
		m.orchMsgCh = msgCh
		m.orchCmdCh = cmdCh
	}
}

// Init returns the initial command. It requests the terminal background color,
// initializes the root screen, and starts listening to the orchestrator channel
// (when one is wired in).
func (m Model) Init() tea.Cmd {
	applogger.Debug().Msg("Initializing UI model")
	cmds := []tea.Cmd{tea.RequestBackgroundColor}
	if len(m.screens) > 0 {
		cmds = append(cmds, m.screens[len(m.screens)-1].Init())
	}
	if m.orchMsgCh != nil {
		cmds = append(cmds, m.listenOrchestrator())
	}
	return tea.Batch(cmds...)
}

// listenOrchestrator returns a Cmd that blocks on the orchestrator message
// channel and delivers one message to the BubbleTea runtime. The Update
// handler re-subscribes after every delivery so the channel is continually
// drained without busy-looping.
func (m Model) listenOrchestrator() tea.Cmd {
	return func() tea.Msg {
		return <-m.orchMsgCh
	}
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
		// Propagate theme to ALL screens in stack.
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
			// Refresh the newly-exposed screen with current window size.
			top := m.screens[len(m.screens)-1]
			updated, cmd := top.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
			m.screens[len(m.screens)-1] = updated
			return m, cmd
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

	// --- Orchestrator messages ---

	case orchestrator.LoopStateMsg,
		orchestrator.AgentOutputMsg,
		orchestrator.IterationStartMsg,
		orchestrator.LoopDoneMsg,
		orchestrator.LoopErrorMsg,
		orchestrator.LoopPausedMsg,
		orchestrator.LoopResumedMsg:
		// Re-subscribe so the next orchestrator message is delivered.
		if m.orchMsgCh != nil {
			cmds = append(cmds, m.listenOrchestrator())
		}
		// Accumulate history for the history screen.
		if ic, ok := msg.(orchestrator.IterationCompleteMsg); ok {
			m.history = append(m.history, ic)
		}
		// fall through to delegate to active screen

	case orchestrator.IterationCompleteMsg:
		m.history = append(m.history, msg)
		if m.orchMsgCh != nil {
			cmds = append(cmds, m.listenOrchestrator())
		}
		// fall through to delegate to active screen

	// --- User-action messages from dashboard → orchestrator ---

	case screens.RetryUserMsg:
		m.sendOrch(orchestrator.RetryCmd{})
		return m, nil

	case screens.SkipUserMsg:
		m.sendOrch(orchestrator.SkipCmd{})
		return m, nil

	case screens.PauseUserMsg:
		m.sendOrch(orchestrator.TogglePauseCmd{})
		return m, nil

	case screens.AdapterChangedMsg:
		m.sendOrch(orchestrator.ChangeAdapterCmd{
			Agent: msg.Agent,
			Model: msg.Model,
		})
		return m, nil
	}

	// Delegate to active screen.
	if len(m.screens) > 0 {
		top := m.screens[len(m.screens)-1]
		updated, cmd := top.Update(msg)
		m.screens[len(m.screens)-1] = updated
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// sendOrch sends a command to the orchestrator channel without blocking.
// If the channel is full the command is dropped (the orchestrator drains
// commands at the top of every iteration, so the next slot will arrive soon).
func (m *Model) sendOrch(cmd any) {
	if m.orchCmdCh == nil {
		return
	}
	select {
	case m.orchCmdCh <- cmd:
	default:
	}
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
	v.AltScreen = m.altScreen
	v.WindowTitle = m.windowTitle
	if m.mouseEnabled {
		v.MouseMode = tea.MouseModeCellMotion
	}
	return v
}

// Run starts the BubbleTea program with the given model.
func Run(m Model, opts ...func(*Model)) error {
	for _, opt := range opts {
		opt(&m)
	}

	applogger.Info().Msg("Starting BubbleTea program")

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("running program: %w", err)
	}

	applogger.Info().Msg("Program exited successfully")
	return nil
}
