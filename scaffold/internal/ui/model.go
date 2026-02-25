package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"scaffold/config"
	"scaffold/internal/ui/banner"
	"scaffold/internal/ui/keys"
	"scaffold/internal/ui/screens"
	"scaffold/internal/ui/theme"
)

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
		if h, ok := m.current.(interface{ SetWidth(int) screens.Home }); ok {
			m.current = h.SetWidth(m.width)
		}
		return m, nil

	case tea.BackgroundColorMsg:
		m.isDark = msg.IsDark()
		m.styles = theme.New(m.isDark, m.width)
		m.help.Styles = help.DefaultStyles(m.isDark)
		return m, nil

	case tea.KeyPressMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}
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

// helpView renders the persistent help box showing global keybindings.
func (m rootModel) helpView() string {
	return m.styles.Help.Render(m.help.View(m.keys))
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
