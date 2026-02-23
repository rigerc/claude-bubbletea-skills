package screens

import (
	"slices"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	huhadapter "ralphio/internal/ui/huh"
	"ralphio/internal/ui/nav"
	"ralphio/internal/ui/theme"
	"ralphio/internal/adapter"
	"ralphio/internal/orchestrator"
)

// AdapterChangedMsg is sent when the user submits the adapter form.
// The root model forwards it to the orchestrator as a ChangeAdapterCmd.
type AdapterChangedMsg struct {
	Agent adapter.AgentType
	Model string
}

// AdapterScreen wraps a huh form that lets the user change the active adapter
// and, for adapters that support it, the active model.
type AdapterScreen struct {
	ScreenBase
	form          *huh.Form
	selectedAgent string
	selectedModel string
	formBuilder   func() *huh.Form
	needsReset    bool
}

// NewAdapterScreen creates an AdapterScreen pre-populated with the current
// agent and model.
func NewAdapterScreen(current adapter.AgentType, currentModel string, isDark bool, appName string) *AdapterScreen {
	s := &AdapterScreen{
		ScreenBase:    NewBase(isDark, appName),
		selectedAgent: string(current),
		selectedModel: currentModel,
	}

	s.formBuilder = func() *huh.Form {
		return s.buildForm()
	}
	s.form = s.formBuilder()
	s.form.WithTheme(theme.HuhThemeFunc())
	s.form.WithKeyMap(huhadapter.KeyMap(s.Keys))
	return s
}

// buildForm constructs the huh form with the current binding pointers.
// The form always shows the agent selector. A second group with the model
// input is conditionally shown based on whether the selected agent supports
// model selection, using group-level WithHideFunc.
func (s *AdapterScreen) buildForm() *huh.Form {
	agentOptions := make([]huh.Option[string], len(adapter.ValidAgents))
	for i, a := range adapter.ValidAgents {
		agentOptions[i] = huh.NewOption(string(a), string(a))
	}

	agentField := huh.NewSelect[string]().
		Title("Agent").
		Description("Select the AI adapter to use").
		Options(agentOptions...).
		Value(&s.selectedAgent)

	modelField := huh.NewInput().
		Title("Model").
		Description("Optional model override (leave blank for default)").
		Placeholder("e.g. anthropic/claude-sonnet-4-20250514").
		Value(&s.selectedModel)

	modelGroup := huh.NewGroup(modelField).
		WithHideFunc(func() bool {
			return !supportsModelSelection(adapter.AgentType(s.selectedAgent))
		})

	return huh.NewForm(
		huh.NewGroup(agentField),
		modelGroup,
	).WithShowHelp(true)
}

// supportsModelSelection reports whether the given agent type supports model selection.
func supportsModelSelection(a adapter.AgentType) bool {
	return slices.Contains(adapter.AgentsSupportingModel, a)
}

// Init returns the form's initial command.
func (s *AdapterScreen) Init() tea.Cmd {
	return s.form.Init()
}

// Update handles incoming messages.
func (s *AdapterScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.Width, s.Height = msg.Width, msg.Height

	case tea.KeyPressMsg:
		switch {
		case msg.String() == "esc":
			return s, nav.Pop()
		case msg.String() == "ctrl+c":
			return s, tea.Quit
		}
	}

	form, cmd := s.form.Update(msg)
	s.form = form.(*huh.Form)

	switch s.form.State {
	case huh.StateCompleted:
		changed := AdapterChangedMsg{
			Agent: adapter.AgentType(s.selectedAgent),
			Model: s.selectedModel,
		}
		// Rebuild form so the screen stays usable if we navigate back.
		s.form = s.formBuilder()
		s.form.WithTheme(theme.HuhThemeFunc())
		s.form.WithKeyMap(huhadapter.KeyMap(s.Keys))
		return s, tea.Batch(
			func() tea.Msg { return changed },
			nav.Pop(),
		)
	case huh.StateAborted:
		s.form = s.formBuilder()
		s.form.WithTheme(theme.HuhThemeFunc())
		s.form.WithKeyMap(huhadapter.KeyMap(s.Keys))
		return s, tea.Batch(cmd, s.form.Init(), nav.Pop())
	}

	return s, cmd
}

// View renders the adapter management form screen.
func (s *AdapterScreen) View() string {
	return s.Theme.App.Render(
		s.HeaderView() + "\n" + s.form.View(),
	)
}

// SetTheme updates the theme. Implements nav.Themeable.
func (s *AdapterScreen) SetTheme(isDark bool) {
	s.ApplyTheme(isDark)
	s.form.WithTheme(theme.HuhThemeFunc())
}

// Ensure AdapterChangedMsg satisfies tea.Msg.
var _ tea.Msg = AdapterChangedMsg{}

// Ensure ChangeAdapterCmd from orchestrator is referenced so we don't have
// an import cycle — the root model translates AdapterChangedMsg into
// orchestrator.ChangeAdapterCmd before sending to the cmdCh.
var _ = orchestrator.ChangeAdapterCmd{}
