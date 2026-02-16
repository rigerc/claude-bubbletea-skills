// Package screens provides the individual screen implementations for the application.
package screens

import (
	lipglossv2 "charm.land/lipgloss/v2"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"template-v2-enhanced/internal/ui/nav"
	"template-v2-enhanced/internal/ui/styles"
)

// DetailScreen displays scrollable text content with a header bar.
// It implements nav.Screen and nav.Themeable.
type DetailScreen struct {
	title, content string
	theme          styles.Theme
	isDark         bool
	width, height  int
	vp             viewport.Model
	ready          bool // false until first WindowSizeMsg
}

// NewDetailScreen creates a new DetailScreen with the given title and content.
// The isDark parameter should be false initially; the correct value will be
// set via SetTheme when the screen is pushed onto the stack.
func NewDetailScreen(title, content string, isDark bool) *DetailScreen {
	vp := viewport.New()
	vp.MouseWheelEnabled = true

	return &DetailScreen{
		title:  title,
		content: content,
		theme:  styles.New(isDark),
		isDark: isDark,
		vp:     vp,
	}
}

// Init returns nil (no initial commands needed).
func (s *DetailScreen) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and returns an updated screen and command.
func (s *DetailScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width, s.height = msg.Width, msg.Height
		headerH := lipglossv2.Height(s.headerView())
		frameH, frameV := s.theme.App.GetFrameSize()
		s.vp.SetWidth(s.width - frameH)
		s.vp.SetHeight(s.height - frameV - headerH)
		if !s.ready {
			s.vp.SetContent(s.content)
			s.ready = true
		}

	case tea.KeyPressMsg:
		if msg.String() == "esc" {
			return s, nav.Pop()
		}
	}

	// Always pass messages to the viewport for scroll handling
	var cmd tea.Cmd
	s.vp, cmd = s.vp.Update(msg)
	return s, cmd
}

// View renders the detail screen with a header and scrollable content.
func (s *DetailScreen) View() string {
	if !s.ready {
		return "Loading..."
	}
	return s.theme.App.Render(
		lipglossv2.JoinVertical(lipglossv2.Left, s.headerView(), s.vp.View()),
	)
}

// SetTheme updates the screen's theme based on the terminal background.
// Implements nav.Themeable.
func (s *DetailScreen) SetTheme(isDark bool) {
	s.isDark = isDark
	s.theme = styles.New(isDark)
}

// SetContent updates the viewport content.
func (s *DetailScreen) SetContent(content string) {
	s.content = content
	if s.ready {
		s.vp.SetContent(content)
	}
}

// headerView renders the title bar at the top of the detail screen.
func (s *DetailScreen) headerView() string {
	return s.theme.Title.Render(s.title)
}

// Height returns the height of the header.
func (s *DetailScreen) Height() int {
	return lipglossv2.Height(s.headerView())
}
