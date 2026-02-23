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
	"ralphio/internal/plan"
)

// taskDetailHelpKeys implements help.KeyMap for the task detail screen.
type taskDetailHelpKeys struct {
	vp  viewport.KeyMap
	app appkeys.GlobalKeyMap
}

func (k taskDetailHelpKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.vp.Up, k.vp.Down, k.app.Back}
}

func (k taskDetailHelpKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.vp.Up, k.vp.Down, k.vp.HalfPageUp, k.vp.HalfPageDown},
		{k.vp.PageUp, k.vp.PageDown, k.app.Back, k.app.Help},
	}
}

// TaskDetailScreen shows the full detail of a single task in a scrollable viewport.
type TaskDetailScreen struct {
	ScreenBase
	task  plan.Task
	vp    viewport.Model
	ready bool
}

// NewTaskDetailScreen creates a TaskDetailScreen for the given task.
func NewTaskDetailScreen(task plan.Task, isDark bool, appName string) *TaskDetailScreen {
	vp := viewport.New()
	vp.MouseWheelEnabled = true
	vp.SoftWrap = true

	return &TaskDetailScreen{
		ScreenBase: NewBase(isDark, appName),
		task:       task,
		vp:         vp,
	}
}

// Init returns nil (no startup commands needed).
func (s *TaskDetailScreen) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages.
func (s *TaskDetailScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.Width, s.Height = msg.Width, msg.Height
		s.updateViewportSize()
		if !s.ready {
			s.vp.SetContent(s.buildContent())
			s.ready = true
		}

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

// View renders the task detail screen.
func (s *TaskDetailScreen) View() string {
	if !s.ready {
		return "Loading..."
	}
	helpKeys := taskDetailHelpKeys{vp: s.vp.KeyMap, app: s.Keys}
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
func (s *TaskDetailScreen) SetTheme(isDark bool) {
	s.ApplyTheme(isDark)
	s.vp.SetContent(s.buildContent())
}

// buildContent assembles the formatted task detail text.
func (s *TaskDetailScreen) buildContent() string {
	t := s.Theme
	task := s.task

	statusStyle := s.statusStyle(task.Status)

	var sb strings.Builder
	sb.WriteString(t.Title.Render("Task Detail") + "\n\n")
	fmt.Fprintf(&sb, "ID:          %s\n", task.ID)
	fmt.Fprintf(&sb, "Status:      %s\n", statusStyle.Render(task.Status))
	fmt.Fprintf(&sb, "Priority:    %d\n", task.Priority)
	fmt.Fprintf(&sb, "Retries:     %d / %d\n", task.RetryCount, task.MaxRetries)
	if task.ValidationCommand != "" {
		fmt.Fprintf(&sb, "Validation:  %s\n", task.ValidationCommand)
	}
	sb.WriteString("\n")
	sb.WriteString(t.Subtle.Render(strings.Repeat("─", 40)) + "\n\n")
	sb.WriteString(t.Title.Render("Title") + "\n\n")
	sb.WriteString(task.Title + "\n\n")
	sb.WriteString(t.Subtle.Render(strings.Repeat("─", 40)) + "\n\n")
	sb.WriteString(t.Title.Render("Description") + "\n\n")
	if task.Description != "" {
		sb.WriteString(task.Description + "\n")
	} else {
		sb.WriteString(t.Subtle.Render("(no description)") + "\n")
	}
	return sb.String()
}

// statusStyle returns the appropriate style for a task status string.
func (s *TaskDetailScreen) statusStyle(status string) lipgloss.Style {
	t := s.Theme
	switch status {
	case plan.StatusCompleted:
		return t.StatusComplete
	case plan.StatusInProgress:
		return t.StatusRunning
	case plan.StatusFailed:
		return t.StatusFailed
	case plan.StatusSkipped:
		return t.StatusSkipped
	default:
		return t.StatusPending
	}
}

func (s *TaskDetailScreen) footerView() string {
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

func (s *TaskDetailScreen) updateViewportSize() {
	if !s.IsSized() {
		return
	}
	s.Help.SetWidth(s.ContentWidth())
	helpKeys := taskDetailHelpKeys{vp: s.vp.KeyMap, app: s.Keys}
	headerH := lipgloss.Height(s.HeaderView())
	footerH := lipgloss.Height(s.footerView())
	helpH := lipgloss.Height(s.RenderHelp(helpKeys))

	vpH := s.CalculateContentHeight(headerH+footerH, helpH)
	s.vp.SetWidth(s.ContentWidth())
	s.vp.SetHeight(vpH)
}
