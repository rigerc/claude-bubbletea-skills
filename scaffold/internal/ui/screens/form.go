// Package screens provides the individual screen implementations for the application.
package screens

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	lipgloss "charm.land/lipgloss/v2"

	huhadapter "scaffold/internal/ui/huh"
	"scaffold/internal/ui/nav"
	"scaffold/internal/ui/theme"
)

// FormScreen wraps a huh.Form to implement nav.Screen.
// It provides integration between Huh forms and the app's navigation system,
// handling global keys (ESC, Ctrl+C, ?) alongside form-specific keys.
type FormScreen struct {
	ScreenBase
	form        *huh.Form
	onSubmit    func() tea.Cmd
	onAbort     func() tea.Cmd
	formBuilder func() *huh.Form // function to rebuild the form when needed
	bodyMaxH    int              // optional override for max body height (0 = default)
}

// NewFormScreen creates a form screen with theme and callbacks.
// The isDark parameter determines the initial theme; it will be updated
// via SetTheme when the screen is pushed onto the stack.
//
// onSubmit is called when the form is completed successfully.
// onAbort is called when the form is aborted (ESC or form abort).
func NewFormScreen(
	form *huh.Form,
	isDark bool,
	appName string,
	onSubmit func() tea.Cmd,
	onAbort func() tea.Cmd,
) *FormScreen {
	fs := &FormScreen{
		ScreenBase: NewBase(isDark, appName),
		form:       form,
		onSubmit:   onSubmit,
		onAbort:    onAbort,
	}

	// Apply theme and keymap
	form.WithTheme(theme.HuhThemeFunc())
	form.WithKeyMap(huhadapter.KeyMap(fs.Keys))

	return fs
}

// newFormScreenWithBuilder is like NewFormScreen but takes a form builder
// function that can rebuild the form when it needs to be reset (e.g., after
// navigation returns to this screen).
//
// The bodyMaxH parameter allows overriding the default body max height for
// forms that need more vertical space (e.g. menus). Pass 0 to use the default.
func newFormScreenWithBuilder(
	formBuilder func() *huh.Form,
	isDark bool,
	appName string,
	onSubmit func() tea.Cmd,
	onAbort func() tea.Cmd,
	bodyMaxH int,
) *FormScreen {
	form := formBuilder()
	fs := &FormScreen{
		ScreenBase:  NewBase(isDark, appName),
		form:        form,
		onSubmit:    onSubmit,
		onAbort:     onAbort,
		formBuilder: formBuilder,
		bodyMaxH:    bodyMaxH,
	}

	// Apply theme and keymap
	form.WithTheme(theme.HuhThemeFunc())
	form.WithKeyMap(huhadapter.KeyMap(fs.Keys))

	return fs
}

// Init returns the form's initial command.
func (s *FormScreen) Init() tea.Cmd {
	return s.form.Init()
}

// Update handles incoming messages and returns an updated screen and command.
// Global keys take precedence over form-specific keys:
//   - ? toggles help expansion
//   - ESC triggers abort callback (navigates back)
//   - Ctrl+C quits the application
func (s *FormScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case resetFormMsg:
		return s.handleResetMsg()

	case tea.WindowSizeMsg:
		s.Width, s.Height = msg.Width, msg.Height

	case tea.KeyPressMsg:
		// Global keys take precedence
		switch {
		case key.Matches(msg, s.Keys.Help):
			s.Help.ShowAll = !s.Help.ShowAll
			return s, nil
		case key.Matches(msg, s.Keys.Back):
			return s, s.onAbort()
		case key.Matches(msg, s.Keys.Quit):
			return s, tea.Quit
		}
	}

	// Delegate to form
	form, cmd := s.form.Update(msg)
	s.form = form.(*huh.Form)

	// Check form state and trigger callbacks.
	// After triggering callbacks, reset form state so it remains usable
	// when we navigate back to this screen (e.g. after ESC from a pushed screen).
	switch s.form.State {
	case huh.StateCompleted:
		cmd = tea.Batch(cmd, s.onSubmit(), s.resetFormCmd())
	case huh.StateAborted:
		cmd = tea.Batch(cmd, s.onAbort(), s.resetFormCmd())
	}

	return s, cmd
}

// View renders the form screen using the 1-column layout: header → form body.
// Huh forms render their own help bar internally, so we reserve space for it
// rather than using the layout's Help section.
func (s *FormScreen) View() string {
	headerView := s.HeaderView()

	// Huh renders its own help bar at the bottom of the form (~4 lines).
	// Reserve that space so the form is not clipped.
	const formInternalHelpH = 4

	maxH := s.bodyMaxH
	if maxH == 0 {
		// Derive a sensible default: available space minus the form's internal
		// help, clamped to [MinContentHeight, MaxContentHeight].
		_, frameV := s.Theme.App.GetFrameSize()
		headerH := lipgloss.Height(headerView)
		avail := s.Height - frameV - headerH - formInternalHelpH
		maxH = avail
		if maxH > MaxContentHeight {
			maxH = MaxContentHeight
		}
		if maxH < MinContentHeight {
			maxH = MinContentHeight
		}
	}

	return s.Layout().
		BodyMaxHeight(maxH).
		Header(headerView).
		Body(s.form.View()).
		Render()
}

// SetTheme updates the screen's theme based on the terminal background.
// Implements nav.Themeable.
func (s *FormScreen) SetTheme(isDark bool) {
	s.ApplyTheme(isDark)
	s.form.WithTheme(theme.HuhThemeFunc())
}

// resetFormCmd returns a command that triggers a form reset on the next Update.
func (s *FormScreen) resetFormCmd() tea.Cmd {
	return func() tea.Msg {
		return resetFormMsg{}
	}
}

// resetFormMsg is dispatched by resetFormCmd to trigger a form rebuild/re-init.
type resetFormMsg struct{}

// handleResetMsg rebuilds the form from formBuilder if one is set, then
// re-initialises it so it is ready for reuse after back-navigation.
func (s *FormScreen) handleResetMsg() (nav.Screen, tea.Cmd) {
	if s.formBuilder != nil {
		s.form = s.formBuilder()
		s.form.WithTheme(theme.HuhThemeFunc())
		s.form.WithKeyMap(huhadapter.KeyMap(s.Keys))
	}
	return s, s.form.Init()
}
