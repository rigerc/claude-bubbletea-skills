package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"ralphio/internal/ui/banner"
	appkeys "ralphio/internal/ui/keys"
	"ralphio/internal/ui/nav"
	"ralphio/internal/orchestrator"
	"ralphio/internal/plan"
)

// --- User-action message types sent from dashboard to root model ---

// RetryUserMsg asks the root model to forward RetryCmd to the orchestrator.
type RetryUserMsg struct{}

// SkipUserMsg asks the root model to forward SkipCmd to the orchestrator.
type SkipUserMsg struct{}

// PauseUserMsg asks the root model to forward TogglePauseCmd to the orchestrator.
type PauseUserMsg struct{}

// ToggleModeUserMsg asks the root model to toggle the loop mode.
type ToggleModeUserMsg struct{}

// EditPlanUserMsg asks the root model to open tasks.json in $EDITOR.
type EditPlanUserMsg struct{}

// RegenPlanUserMsg asks the root model to switch to planning mode.
type RegenPlanUserMsg struct{}

// TurnKind identifies the visual style of an output turn.
type TurnKind int

const (
	// TurnAgent is numbered iteration output from the LLM.
	TurnAgent TurnKind = iota
	// TurnSystem is validation results, mode changes, or loop events.
	TurnSystem
)

// OutputTurn holds one discrete unit of output.
type OutputTurn struct {
	Kind      TurnKind
	Iteration int      // 0 for system turns
	Lines     []string // accumulated lines; last may be partial (streaming)
	Streaming bool     // true while this turn is still receiving content
}

// statusBarHeight is the fixed height (in rows) of the persistent status bar.
const statusBarHeight = 1

// DashboardScreen is the root screen of the ralphio TUI. It shows the loop
// status, a chat-like output panel for agent turns, and a persistent status bar.
type DashboardScreen struct {
	ScreenBase

	keys         appkeys.DashboardKeyMap
	state        orchestrator.LoopStateMsg
	turns        []OutputTurn
	chatViewport viewport.Model
	autoScroll   bool
	lastPassed   *bool  // nil until the first IterationCompleteMsg arrives
	projectDir   string // for file status checks
	bannerText   string // pre-rendered figlet text
	bannerHeight int    // lipgloss.Height(bannerText) + 1 for sub-line + 1 for divider
}

// NewDashboardScreen creates a new DashboardScreen.
func NewDashboardScreen(isDark bool, appName, projectDir string) *DashboardScreen {
	vp := viewport.New()
	vp.MouseWheelEnabled = true
	vp.SoftWrap = true

	s := &DashboardScreen{
		ScreenBase:   NewBase(isDark, appName),
		keys:         appkeys.NewDashboard(),
		chatViewport: vp,
		autoScroll:   true,
		projectDir:   projectDir,
	}
	s.initBanner()
	return s
}

// initBanner pre-renders the figlet banner and measures its height.
func (s *DashboardScreen) initBanner() {
	cfg := banner.BannerConfig{
		Text: "RALPHIO",
		Font: "standard",
	}
	rendered, err := banner.RenderBanner(cfg, 120)
	if err != nil {
		rendered = "RALPHIO"
	}
	s.bannerText = rendered
	// +1 for sub-line, +1 for divider line.
	s.bannerHeight = lipgloss.Height(rendered) + 2
}

// subLine returns the one-line status summary rendered beneath the banner.
func (s *DashboardScreen) subLine() string {
	modeStr := "BUILDING"
	if s.state.LoopMode == plan.ModePlanning {
		modeStr = "PLANNING"
	}
	agent := string(s.state.ActiveAgent)
	if agent == "" {
		agent = "—"
	}
	return fmt.Sprintf("Project: %s  |  Mode: %s  |  Adapter: %s  |  Iteration: %d",
		s.projectDir,
		modeStr,
		agent,
		s.state.Iteration,
	)
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
		s.rebuildViewport()

	case tea.KeyPressMsg:
		// Scroll key handling takes priority so users can control auto-scroll.
		switch msg.String() {
		case "up", "pgup":
			s.autoScroll = false
		case "end":
			s.autoScroll = true
			s.chatViewport.GotoBottom()
		}

		switch {
		case key.Matches(msg, s.keys.Retry):
			return s, func() tea.Msg { return RetryUserMsg{} }
		case key.Matches(msg, s.keys.Skip):
			return s, func() tea.Msg { return SkipUserMsg{} }
		case key.Matches(msg, s.keys.Pause):
			return s, func() tea.Msg { return PauseUserMsg{} }
		case key.Matches(msg, s.keys.ToggleMode):
			return s, func() tea.Msg { return ToggleModeUserMsg{} }
		case key.Matches(msg, s.keys.EditPlan):
			return s, func() tea.Msg { return EditPlanUserMsg{} }
		case key.Matches(msg, s.keys.RegenPlan):
			return s, func() tea.Msg { return RegenPlanUserMsg{} }
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

	case orchestrator.IterationStartMsg:
		// Open a new agent turn for this iteration.
		s.turns = append(s.turns, OutputTurn{
			Kind:      TurnAgent,
			Iteration: msg.Iteration,
			Streaming: true,
		})
		s.autoScroll = true
		s.rebuildViewport()

	case orchestrator.AgentOutputMsg:
		if msg.Text != "" {
			// Append to the last agent turn, or create one if none exists.
			if len(s.turns) == 0 || s.turns[len(s.turns)-1].Kind != TurnAgent {
				s.turns = append(s.turns, OutputTurn{
					Kind:      TurnAgent,
					Iteration: s.state.Iteration,
					Streaming: true,
				})
			}
			last := &s.turns[len(s.turns)-1]
			last.Lines = append(last.Lines, msg.Text)

			// Bound total lines across all turns to avoid memory growth.
			totalLines := 0
			for _, t := range s.turns {
				totalLines += len(t.Lines)
			}
			if totalLines > 1000 {
				// Drop lines from the oldest turn.
				for i := range s.turns {
					if len(s.turns[i].Lines) > 0 {
						s.turns[i].Lines = s.turns[i].Lines[1:]
						break
					}
				}
			}
			s.autoScroll = true
			s.rebuildViewport()
		}

	case orchestrator.IterationCompleteMsg:
		// Close the streaming turn and append a system result.
		if len(s.turns) > 0 {
			s.turns[len(s.turns)-1].Streaming = false
		}
		passed := msg.Passed
		s.lastPassed = &passed
		result := "PASS"
		if !passed {
			result = "FAIL"
		}
		s.turns = append(s.turns, OutputTurn{
			Kind:  TurnSystem,
			Lines: []string{"Iteration " + fmt.Sprintf("%d", msg.Iteration) + ": " + result},
		})
		s.autoScroll = true
		s.rebuildViewport()

	case orchestrator.LoopStateMsg:
		s.state = msg
		// No viewport rebuild needed here; layout() handles dimensions.

	case orchestrator.LoopDoneMsg:
		s.turns = append(s.turns, OutputTurn{
			Kind:  TurnSystem,
			Lines: []string{"All tasks completed"},
		})
		s.rebuildViewport()

	case orchestrator.LoopErrorMsg:
		s.turns = append(s.turns, OutputTurn{
			Kind:  TurnSystem,
			Lines: []string{"Error: " + msg.Err.Error()},
		})
		s.rebuildViewport()

	case orchestrator.LoopPausedMsg:
		s.turns = append(s.turns, OutputTurn{Kind: TurnSystem, Lines: []string{"Loop paused"}})
		s.rebuildViewport()

	case orchestrator.LoopResumedMsg:
		s.turns = append(s.turns, OutputTurn{Kind: TurnSystem, Lines: []string{"Loop resumed"}})
		s.rebuildViewport()
	}

	var cmd tea.Cmd
	s.chatViewport, cmd = s.chatViewport.Update(msg)
	return s, cmd
}

// View renders the three-region dashboard: banner, chat panel, and status bar.
func (s *DashboardScreen) View() string {
	if !s.IsSized() {
		return "Loading..."
	}

	t := s.Theme

	// Banner region.
	bannerStyle := lipgloss.NewStyle().Foreground(t.Palette.Primary).Bold(true)
	subLineStyle := lipgloss.NewStyle().Foreground(t.Palette.Muted)
	dividerStyle := lipgloss.NewStyle().Foreground(t.Palette.Border)
	banner := lipgloss.JoinVertical(lipgloss.Left,
		bannerStyle.Render(s.bannerText),
		subLineStyle.Render(s.subLine()),
		dividerStyle.Render(strings.Repeat("─", s.Width)),
	)

	// File status line (one line, below banner).
	fileStatusLine := s.renderFileStatus(s.Width)
	fileStatusH := lipgloss.Height(fileStatusLine)

	// Chat panel fills the middle between banner + file status and status bar.
	chatH := s.Height - s.bannerHeight - fileStatusH - statusBarHeight
	if chatH < 1 {
		chatH = 1
	}
	s.chatViewport.SetWidth(s.Width)
	s.chatViewport.SetHeight(chatH)
	s.rebuildViewport()

	// Status bar.
	statusBar := s.renderStatusBar()

	return lipgloss.JoinVertical(lipgloss.Left,
		banner,
		fileStatusLine,
		s.chatViewport.View(),
		statusBar,
	)
}

// SetTheme updates the screen's theme. Implements nav.Themeable.
func (s *DashboardScreen) SetTheme(isDark bool) {
	s.ApplyTheme(isDark)
}

// rebuildViewport rebuilds the chat viewport content from the current turns.
func (s *DashboardScreen) rebuildViewport() {
	var sb strings.Builder
	for _, t := range s.turns {
		sb.WriteString(renderTurn(t, s.IsDark, s.chatViewport.Width()))
		sb.WriteString("\n\n")
	}
	s.chatViewport.SetContent(sb.String())
	if s.autoScroll {
		s.chatViewport.GotoBottom()
	}
}

// renderTurn formats a single OutputTurn for display in the chat viewport.
func renderTurn(t OutputTurn, isDark bool, width int) string {
	ld := lipgloss.LightDark(isDark)
	_ = width // reserved for future explicit wrapping

	switch t.Kind {
	case TurnAgent:
		label := fmt.Sprintf("[#%d]", t.Iteration)
		headerStyle := lipgloss.NewStyle().
			Foreground(ld(lipgloss.Color("#0080FF"), lipgloss.Color("#4DA6FF"))).
			Bold(true)
		header := headerStyle.Render(label)
		body := strings.Join(t.Lines, "\n")
		if t.Streaming {
			body += " \u258c" // block cursor
		}
		return lipgloss.JoinVertical(lipgloss.Left, header, body)

	case TurnSystem:
		style := lipgloss.NewStyle().
			Foreground(ld(lipgloss.Color("#888888"), lipgloss.Color("#AAAAAA"))).
			Italic(true)
		return style.Render("[sys]  " + strings.Join(t.Lines, "\n"))
	}
	return ""
}

// renderStatusBar renders the persistent bottom key-binding bar.
func (s *DashboardScreen) renderStatusBar() string {
	t := s.Theme
	helpView := s.Help.View(s.keys)
	return lipgloss.NewStyle().
		Background(t.Palette.Border).
		Foreground(t.Palette.Text).
		Width(s.Width).
		Render(helpView)
}

// renderFileStatus renders a one-line indicator showing which Ralph files exist.
func (s *DashboardScreen) renderFileStatus(width int) string {
	t := s.Theme

	check := func(path string) string {
		if _, err := os.Stat(path); err == nil {
			return t.StatusComplete.Render("✓")
		}
		return t.StatusFailed.Render("✗")
	}

	base := s.projectDir
	parts := []string{
		fmt.Sprintf("PROMPT_build.md %s", check(filepath.Join(base, "PROMPT_build.md"))),
		fmt.Sprintf("PROMPT_plan.md %s", check(filepath.Join(base, "PROMPT_plan.md"))),
		fmt.Sprintf("AGENTS.md %s", check(filepath.Join(base, "AGENTS.md"))),
		fmt.Sprintf("tasks.json %s", check(filepath.Join(base, "tasks.json"))),
	}

	line := strings.Join(parts, "  ")
	return lipgloss.NewStyle().
		Width(width).
		Render(t.Subtle.Render(line))
}

// renderTaskList renders the bordered task list pane (used by other screens
// or callers that embed the dashboard layout).
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
