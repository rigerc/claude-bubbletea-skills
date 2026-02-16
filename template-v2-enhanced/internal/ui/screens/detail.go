// Package screens provides the individual screen implementations for the application.
package screens

import (
	lipglossv2 "charm.land/lipgloss/v2"
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	appkeys "template-v2-enhanced/internal/ui/keys"
	"template-v2-enhanced/internal/ui/nav"
	"template-v2-enhanced/internal/ui/styles"
)

// DetailScreen displays scrollable text content with a header bar and help footer.
// It implements nav.Screen and nav.Themeable.
type DetailScreen struct {
	title, content string
	keys           appkeys.GlobalKeyMap
	help           help.Model
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

	h := help.New()
	h.Styles = help.DefaultStyles(isDark)

	return &DetailScreen{
		title:   title,
		content: content,
		keys:    appkeys.New(),
		help:    h,
		theme:   styles.New(isDark),
		isDark:  isDark,
		vp:      vp,
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
		s.updateViewportSize()
		if !s.ready {
			s.vp.SetContent(s.content)
			s.ready = true
		}

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, s.keys.Help):
			s.help.ShowAll = !s.help.ShowAll
			s.updateViewportSize()
			return s, nil
		case key.Matches(msg, s.keys.Back):
			return s, nav.Pop()
		}
	}

	// Always pass messages to the viewport for scroll handling.
	var cmd tea.Cmd
	s.vp, cmd = s.vp.Update(msg)
	return s, cmd
}

// View renders the detail screen: header, scrollable content, help footer.
func (s *DetailScreen) View() string {
	if !s.ready {
		return "Loading..."
	}
	return s.theme.App.Render(
		lipglossv2.JoinVertical(lipglossv2.Left,
			s.headerView(),
			s.vp.View(),
			s.help.View(s.keys),
		),
	)
}

// SetTheme updates the screen's theme and help styles based on the terminal background.
// Implements nav.Themeable.
func (s *DetailScreen) SetTheme(isDark bool) {
	s.isDark = isDark
	s.theme = styles.New(isDark)
	s.help.Styles = help.DefaultStyles(isDark)
}

// SetContent updates the viewport content.
func (s *DetailScreen) SetContent(content string) {
	s.content = content
	if s.ready {
		s.vp.SetContent(content)
	}
}

// updateViewportSize recalculates viewport dimensions from the window size,
// theme frame, header height, and actual rendered help bar height.
func (s *DetailScreen) updateViewportSize() {
	if s.width == 0 || s.height == 0 {
		return
	}
	headerH := lipglossv2.Height(s.headerView())
	frameH, frameV := s.theme.App.GetFrameSize()
	contentW := s.width - frameH
	s.help.SetWidth(contentW)
	helpH := lipglossv2.Height(s.help.View(s.keys))
	s.vp.SetWidth(contentW)
	s.vp.SetHeight(s.height - frameV - headerH - helpH)
}

// headerView renders the title bar at the top of the detail screen.
func (s *DetailScreen) headerView() string {
	return s.theme.Title.Render(s.title)
}
