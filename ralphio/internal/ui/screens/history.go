package screens

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	appkeys "ralphio/internal/ui/keys"
	"ralphio/internal/ui/nav"
	"ralphio/internal/orchestrator"
)

// historyHelpKeys implements help.KeyMap for the history screen.
type historyHelpKeys struct {
	vp  viewport.KeyMap
	app appkeys.GlobalKeyMap
}

func (k historyHelpKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.vp.Up, k.vp.Down, k.app.Back}
}

func (k historyHelpKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.vp.Up, k.vp.Down, k.vp.HalfPageUp, k.vp.HalfPageDown},
		{k.vp.PageUp, k.vp.PageDown, k.app.Back, k.app.Help},
	}
}

// HistoryScreen shows the iteration history in a scrollable viewport.
type HistoryScreen struct {
	ScreenBase
	history []orchestrator.IterationCompleteMsg
	vp      viewport.Model
	ready   bool
}

// NewHistoryScreen creates a HistoryScreen.
// Call AppendHistory to add entries before or after creation.
func NewHistoryScreen(appName string, isDark bool) *HistoryScreen {
	vp := viewport.New()
	vp.MouseWheelEnabled = true
	vp.SoftWrap = true

	return &HistoryScreen{
		ScreenBase: NewBase(isDark, appName),
		vp:         vp,
	}
}

// AppendHistory adds an iteration result to the history list and refreshes the
// viewport content. This is called by the root model whenever a new
// IterationCompleteMsg arrives.
func (s *HistoryScreen) AppendHistory(entry orchestrator.IterationCompleteMsg) {
	s.history = append(s.history, entry)
	s.vp.SetContent(s.buildContent())
	s.vp.GotoBottom()
}

// SetHistory replaces the full history slice and refreshes the viewport.
func (s *HistoryScreen) SetHistory(history []orchestrator.IterationCompleteMsg) {
	s.history = history
	s.vp.SetContent(s.buildContent())
	s.vp.GotoBottom()
}

// Init returns nil (no startup commands needed).
func (s *HistoryScreen) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages.
func (s *HistoryScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.Width, s.Height = msg.Width, msg.Height
		s.updateViewportSize()
		if !s.ready {
			s.vp.SetContent(s.buildContent())
			s.ready = true
		}

	case orchestrator.IterationCompleteMsg:
		s.AppendHistory(msg)

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, s.Keys.Help):
			s.Help.ShowAll = !s.Help.ShowAll
			s.updateViewportSize()
			return s, nil
		case key.Matches(msg, s.Keys.Back):
			return s, nav.Pop()
		}
	}

	var cmd tea.Cmd
	s.vp, cmd = s.vp.Update(msg)
	return s, cmd
}

// View renders the history screen.
func (s *HistoryScreen) View() string {
	if !s.ready {
		return "Loading..."
	}
	helpKeys := historyHelpKeys{vp: s.vp.KeyMap, app: s.Keys}
	return s.Theme.App.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			s.HeaderView(),
			s.vp.View(),
			s.footerView(),
			s.RenderHelp(helpKeys),
		),
	)
}

// SetTheme updates the theme. Implements nav.Themeable.
func (s *HistoryScreen) SetTheme(isDark bool) {
	s.ApplyTheme(isDark)
	s.vp.SetContent(s.buildContent())
}

// buildContent formats the history slice as a table-style string.
func (s *HistoryScreen) buildContent() string {
	t := s.Theme

	if len(s.history) == 0 {
		return t.Subtle.Render("No iterations recorded yet.")
	}

	var sb strings.Builder
	sb.WriteString(t.Title.Render("Iteration History") + "\n\n")

	header := fmt.Sprintf("%-6s %-10s %-6s %s", "#", "Task", "Result", "Duration")
	sb.WriteString(t.Subtle.Render(header) + "\n")
	sb.WriteString(t.Subtle.Render(strings.Repeat("─", 36)) + "\n")

	for _, entry := range s.history {
		result := "PASS"
		style := t.StatusComplete
		if !entry.Passed {
			result = "FAIL"
			style = t.StatusFailed
		}

		taskID := entry.TaskID
		if len(taskID) > 10 {
			taskID = taskID[:9] + "…"
		}

		dur := entry.Duration.Round(1000000000) // round to seconds
		line := fmt.Sprintf("%-6d %-10s %-6s %s",
			entry.Iteration,
			taskID,
			style.Render(result),
			dur,
		)
		sb.WriteString(line + "\n")
	}

	return sb.String()
}

func (s *HistoryScreen) footerView() string {
	b := lipgloss.RoundedBorder()
	b.Left = "┤"
	info := lipgloss.NewStyle().
		BorderStyle(b).
		BorderForeground(s.Theme.Palette.Primary).
		Padding(0, 1).
		Render(fmt.Sprintf("%3.f%%", s.vp.ScrollPercent()*100))

	lineW := max(0, s.ContentWidth()-lipgloss.Width(info))
	line := s.Theme.Subtle.Render(strings.Repeat("─", lineW))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (s *HistoryScreen) updateViewportSize() {
	if !s.IsSized() {
		return
	}
	s.Help.SetWidth(s.ContentWidth())
	helpKeys := historyHelpKeys{vp: s.vp.KeyMap, app: s.Keys}
	headerH := lipgloss.Height(s.HeaderView())
	footerH := lipgloss.Height(s.footerView())
	helpH := lipgloss.Height(s.RenderHelp(helpKeys))

	vpH := s.CalculateContentHeight(headerH+footerH, helpH)
	s.vp.SetWidth(s.ContentWidth())
	s.vp.SetHeight(vpH)
}
