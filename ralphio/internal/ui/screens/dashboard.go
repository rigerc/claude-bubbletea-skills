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
	"ralphio/internal/plan"
	"ralphio/internal/validator"
)

// --- User-action message types sent from dashboard to root model ---

// RetryUserMsg asks the root model to forward RetryCmd to the orchestrator.
type RetryUserMsg struct{}

// SkipUserMsg asks the root model to forward SkipCmd to the orchestrator.
type SkipUserMsg struct{}

// PauseUserMsg asks the root model to forward TogglePauseCmd to the orchestrator.
type PauseUserMsg struct{}

// DashboardScreen is the root screen of the ralphio TUI. It shows the loop
// status, task list, validation results, and a streaming agent output log.
type DashboardScreen struct {
	ScreenBase

	keys           appkeys.DashboardKeyMap
	state          orchestrator.LoopStateMsg
	logLines       []string
	logViewport    viewport.Model
	logReady       bool
	lastValidation *validator.Result
}

// NewDashboardScreen creates a new DashboardScreen.
func NewDashboardScreen(isDark bool, appName string) *DashboardScreen {
	vp := viewport.New()
	vp.MouseWheelEnabled = true
	vp.SoftWrap = true

	return &DashboardScreen{
		ScreenBase:  NewBase(isDark, appName),
		keys:        appkeys.NewDashboard(),
		logViewport: vp,
	}
}

// Init returns nil; the orchestrator drives updates via channel messages.
func (s *DashboardScreen) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and returns an updated screen and command.
func (s *DashboardScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.Width, s.Height = msg.Width, msg.Height
		s.syncLogViewport()
		if !s.logReady {
			s.logReady = true
		}

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, s.keys.Retry):
			return s, func() tea.Msg { return RetryUserMsg{} }
		case key.Matches(msg, s.keys.Skip):
			return s, func() tea.Msg { return SkipUserMsg{} }
		case key.Matches(msg, s.keys.Pause):
			return s, func() tea.Msg { return PauseUserMsg{} }
		case key.Matches(msg, s.keys.Detail):
			if s.state.CurrentTask != nil {
				return s, nav.Push(NewTaskDetailScreen(*s.state.CurrentTask, s.IsDark, s.AppName))
			}
		case key.Matches(msg, s.keys.History):
			return s, nav.Push(NewHistoryScreen(s.AppName, s.IsDark))
		case key.Matches(msg, s.keys.Client):
			return s, nav.Push(NewAdapterScreen(s.state.ActiveAgent, s.state.ActiveModel, s.IsDark, s.AppName))
		case key.Matches(msg, s.keys.Quit):
			return s, tea.Quit
		}

	case orchestrator.LoopStateMsg:
		s.state = msg
		s.syncLogViewport()

	case orchestrator.AgentOutputMsg:
		if msg.Text != "" {
			s.logLines = append(s.logLines, msg.Text)
			// Keep the log bounded to avoid unbounded memory growth.
			if len(s.logLines) > 500 {
				s.logLines = s.logLines[len(s.logLines)-500:]
			}
			s.syncLogContent()
			s.logViewport.GotoBottom()
		}

	case orchestrator.IterationCompleteMsg:
		if msg.ValidationResult != nil {
			vr := *msg.ValidationResult
			s.lastValidation = &vr
		}
	}

	var cmd tea.Cmd
	s.logViewport, cmd = s.logViewport.Update(msg)
	return s, cmd
}

// View renders the dashboard: header, two-column task/validation panes,
// agent log pane, and help bar.
func (s *DashboardScreen) View() string {
	if !s.IsSized() {
		return "Loading..."
	}

	header := s.HeaderView()
	helpBar := s.RenderHelp(s.keys)
	headerH := lipgloss.Height(header)
	helpH := lipgloss.Height(helpBar)

	// Status line below header
	statusLine := s.renderStatusLine()
	statusH := lipgloss.Height(statusLine)

	// Calculate remaining height for panes.
	_, frameV := s.Theme.App.GetFrameSize()
	remaining := s.Height - frameV - headerH - statusH - helpH
	remaining = max(remaining, 4)

	// Split: upper panes take ~40% of remaining, log takes the rest.
	upperH := remaining * 2 / 5
	upperH = max(upperH, 4)
	logH := remaining - upperH
	logH = max(logH, 3)

	cw := s.ContentWidth()

	var upperSection string
	if cw < 80 {
		// Narrow: single-column layout.
		taskPane := s.renderTaskList(cw-2, upperH)
		upperSection = taskPane
	} else {
		// Wide: two-column layout.
		leftW := cw * 2 / 3
		rightW := cw - leftW - 1
		taskPane := s.renderTaskList(leftW-2, upperH)
		validPane := s.renderValidation(rightW-2, upperH)
		upperSection = lipgloss.JoinHorizontal(lipgloss.Top, taskPane, " ", validPane)
	}

	// Sync log viewport dimensions.
	s.logViewport.SetWidth(cw)
	s.logViewport.SetHeight(logH)
	logPane := s.renderLogPane(cw, logH)

	body := lipgloss.JoinVertical(lipgloss.Left,
		header,
		statusLine,
		upperSection,
		logPane,
		helpBar,
	)

	return s.Theme.App.Render(body)
}

// SetTheme updates the screen's theme. Implements nav.Themeable.
func (s *DashboardScreen) SetTheme(isDark bool) {
	s.ApplyTheme(isDark)
}

// syncLogViewport resizes the log viewport to match current dimensions.
func (s *DashboardScreen) syncLogViewport() {
	if !s.IsSized() {
		return
	}
	s.syncLogContent()
}

// syncLogContent rebuilds the viewport content from logLines.
func (s *DashboardScreen) syncLogContent() {
	content := strings.Join(s.logLines, "\n")
	s.logViewport.SetContent(content)
}

// renderStatusLine renders the one-line status summary below the header.
func (s *DashboardScreen) renderStatusLine() string {
	t := s.Theme

	status := s.state.Status
	var statusStyle lipgloss.Style
	switch status {
	case "running":
		statusStyle = t.StatusRunning
		status = "RUNNING"
	case "paused":
		statusStyle = t.StatusPaused
		status = "PAUSED"
	case "error":
		statusStyle = t.StatusFailed
		status = "ERROR"
	case "stopped":
		statusStyle = t.StatusSkipped
		status = "STOPPED"
	default:
		statusStyle = t.StatusPending
		status = strings.ToUpper(status)
	}

	agent := string(s.state.ActiveAgent)
	if agent == "" {
		agent = "—"
	}
	modelPart := ""
	if s.state.ActiveModel != "" {
		modelPart = " / " + s.state.ActiveModel
	}

	parts := []string{
		fmt.Sprintf("Iteration: %d", s.state.Iteration),
		statusStyle.Render(status),
		fmt.Sprintf("Task %d/%d", s.state.CompletedTasks, s.state.TotalTasks),
		fmt.Sprintf("[%s%s]", agent, modelPart),
	}

	return t.Subtle.Render(strings.Join(parts, "   ")) + "\n"
}

// renderTaskList renders the bordered task list pane.
func (s *DashboardScreen) renderTaskList(width, height int) string {
	t := s.Theme

	var sb strings.Builder
	for _, task := range s.state.Tasks {
		icon, style := s.taskIconAndStyle(task)
		title := task.Title
		maxTitleW := max(width-6, 1)
		if len(title) > maxTitleW {
			title = title[:maxTitleW-1] + "…"
		}
		line := fmt.Sprintf("[%s] %s", icon, title)
		sb.WriteString(style.Render(line) + "\n")
	}
	if s.state.TotalTasks == 0 {
		sb.WriteString(t.Subtle.Render("No tasks loaded"))
	}

	content := lipgloss.NewStyle().
		Width(width).
		Height(height - 2).
		Render(sb.String())

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Palette.Border).
		Width(width + 2).
		Render("Tasks\n" + content)

	return panel
}

// taskIconAndStyle returns the status icon and display style for a task.
func (s *DashboardScreen) taskIconAndStyle(task plan.Task) (string, lipgloss.Style) {
	t := s.Theme
	switch task.Status {
	case plan.StatusCompleted:
		return "✓", t.StatusComplete
	case plan.StatusInProgress:
		return ">", t.StatusRunning
	case plan.StatusFailed:
		return "!", t.StatusFailed
	case plan.StatusSkipped:
		return "-", t.StatusSkipped
	default:
		return " ", t.StatusPending
	}
}

// renderValidation renders the validation results pane.
func (s *DashboardScreen) renderValidation(width, height int) string {
	t := s.Theme

	var sb strings.Builder
	if s.lastValidation != nil {
		v := s.lastValidation
		cmd := v.Command
		maxCmdW := max(width-8, 1)
		if len(cmd) > maxCmdW {
			cmd = cmd[:maxCmdW-1] + "…"
		}
		result := "PASS"
		style := t.StatusComplete
		if !v.Passed {
			result = "FAIL"
			style = t.StatusFailed
		}
		fmt.Fprintf(&sb, "%s: %s\n", cmd, style.Render(result))
	} else {
		sb.WriteString(t.Subtle.Render("No results yet"))
	}

	content := lipgloss.NewStyle().
		Width(width).
		Height(height - 2).
		Render(sb.String())

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Palette.Border).
		Width(width + 2).
		Render("Validation\n" + content)

	return panel
}

// renderLogPane renders the agent output log viewport.
func (s *DashboardScreen) renderLogPane(width, height int) string {
	t := s.Theme

	label := t.Subtle.Render("Agent Output")
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Palette.Border).
		Width(width).
		Height(height)

	s.logViewport.SetWidth(width - 2)
	s.logViewport.SetHeight(height - 2)

	return border.Render(label + "\n" + s.logViewport.View())
}
