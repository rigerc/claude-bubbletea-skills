package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"scaffold/config"
	"scaffold/internal/ui/banner"
	"scaffold/internal/ui/keys"
	"scaffold/internal/ui/menu"
	"scaffold/internal/ui/screens"
	"scaffold/internal/ui/theme"
)

// NavigateMsg is a message to navigate to a new screen.
type NavigateMsg struct {
	Screen screens.Screen
}

// BackMsg is a message to navigate back to the previous screen.
type BackMsg struct{}

// screenStack holds the navigation history.
type screenStack struct {
	screens []screens.Screen
}

// Push adds a screen to the stack.
func (s *screenStack) Push(screen screens.Screen) {
	s.screens = append(s.screens, screen)
}

// Pop removes and returns the top screen.
func (s *screenStack) Pop() screens.Screen {
	if len(s.screens) == 0 {
		return nil
	}
	idx := len(s.screens) - 1
	screen := s.screens[idx]
	s.screens = s.screens[:idx]
	return screen
}

// Peek returns the top screen without removing it.
func (s *screenStack) Peek() screens.Screen {
	if len(s.screens) == 0 {
		return nil
	}
	return s.screens[len(s.screens)-1]
}

// Len returns the stack depth.
func (s *screenStack) Len() int {
	return len(s.screens)
}

// rootModel is the root tea.Model — owns routing, WindowSize, header/footer.
type rootModel struct {
	cfg     config.Config
	width   int
	banner  string
	isDark  bool
	ready   bool
	styles  theme.Styles
	keys    keys.GlobalKeyMap
	help    help.Model
	current screens.Screen
	stack   screenStack
}

// newRootModel creates a new root model.
func newRootModel(cfg config.Config) rootModel {
	return rootModel{
		cfg:     cfg,
		current: screens.NewHome(),
		keys:    keys.DefaultGlobalKeyMap(),
		help:    help.New(),
	}
}

// Init initializes the root model.
func (m rootModel) Init() tea.Cmd {
	return tea.RequestBackgroundColor
}

// Update handles messages for the root model.
func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.ready = true
		m.styles = theme.New(m.isDark, m.width)
		m.help.SetWidth(m.styles.MaxWidth)

		// Render banner once with the window width
		b, err := banner.Render(banner.Config{
			Text:          m.cfg.App.Name,
			Width:         m.width - 10, // Account for padding
			Justification: 0,            // Left aligned
		})
		if err != nil {
			b = m.cfg.App.Name
		}
		m.banner = b

		// Propagate width to current screen
		if setter, ok := m.current.(interface{ SetWidth(int) screens.Screen }); ok {
			m.current = setter.SetWidth(m.width)
		}
		return m, nil

	case tea.BackgroundColorMsg:
		m.isDark = msg.IsDark()
		m.styles = theme.New(m.isDark, m.width)
		m.help.Styles = help.DefaultStyles(m.isDark)
		// Propagate theme to current screen
		if setter, ok := m.current.(interface{ SetStyles(bool) screens.Screen }); ok {
			m.current = setter.SetStyles(m.isDark)
		}
		return m, nil

	case tea.KeyPressMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}
		// Handle Back key (Escape)
		if key.Matches(msg, m.keys.Back) {
			if m.stack.Len() > 0 {
				m.current = m.stack.Pop()
				return m, nil
			}
		}

	case NavigateMsg:
		// Push current screen to stack and navigate to new screen
		m.stack.Push(m.current)
		m.current = msg.Screen
		// Propagate width and theme to new screen
		if setter, ok := m.current.(interface{ SetWidth(int) screens.Screen }); ok {
			m.current = setter.SetWidth(m.width)
		}
		if setter, ok := m.current.(interface{ SetStyles(bool) screens.Screen }); ok {
			m.current = setter.SetStyles(m.isDark)
		}
		return m, nil

	case menu.SelectionMsg:
		// Convert menu selection to navigation
		item := msg.Item
		detail := screens.NewDetail(item.Title(), item.Description(), item.ScreenID())
		return m.Update(NavigateMsg{Screen: detail})
	}

	// Delegate to current screen
	var cmd tea.Cmd
	updated, cmd := m.current.Update(msg)
	if s, ok := updated.(screens.Screen); ok {
		m.current = s
	}
	return m, cmd
}

// View renders the root model.
func (m rootModel) View() tea.View {
	if !m.ready {
		return tea.NewView("")
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		m.headerView(),
		m.styles.Body.Render(m.current.Body()),
		m.helpView(),
		m.footerView(),
	)

	return tea.NewView(m.styles.App.Render(content))
}

// headerView renders the header with the banner.
func (m rootModel) headerView() string {
	return m.styles.Header.Render(m.banner)
}

// helpView renders the persistent help box showing global and screen-specific keybindings.
func (m rootModel) helpView() string {
	combined := m.combinedKeys()
	return m.styles.Help.Render(m.help.View(combined))
}

// combinedKeys returns a key map that combines global keys with screen-specific keys.
func (m rootModel) combinedKeys() combinedKeyMap {
	return combinedKeyMap{
		global: m.keys,
		screen: m.current,
	}
}

// combinedKeyMap combines global and screen-specific key bindings.
type combinedKeyMap struct {
	global keys.GlobalKeyMap
	screen screens.Screen
}

// ShortHelp returns combined short help bindings.
func (c combinedKeyMap) ShortHelp() []key.Binding {
	bindings := c.global.ShortHelp()
	if kb, ok := c.screen.(screens.KeyBinder); ok {
		bindings = append(bindings, kb.ShortHelp()...)
	}
	return bindings
}

// FullHelp returns combined full help bindings.
func (c combinedKeyMap) FullHelp() [][]key.Binding {
	groups := c.global.FullHelp()
	if kb, ok := c.screen.(screens.KeyBinder); ok {
		groups = append(groups, kb.FullHelp()...)
	}
	return groups
}

// footerView renders the status bar footer.
func (m rootModel) footerView() string {
	status := "Ready"
	left := m.styles.StatusLeft.Render(" " + status + " ")
	right := m.styles.StatusRight.Render(" v" + m.cfg.App.Version + " ")

	// Account for footer border (2) and padding (1)
	innerWidth := m.styles.MaxWidth - 3

	gap := lipgloss.NewStyle().
		Width(innerWidth - lipgloss.Width(left) - lipgloss.Width(right)).
		Render("")
	footerContent := lipgloss.JoinHorizontal(lipgloss.Top, left, gap, right)
	return m.styles.Footer.Render(footerContent)
}
