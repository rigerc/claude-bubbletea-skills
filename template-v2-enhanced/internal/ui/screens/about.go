// Package screens provides example screen implementations for the navigation demo.
package screens

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"template-v2-enhanced/internal/ui/nav"
	appkeys "template-v2-enhanced/internal/ui/keys"
	"template-v2-enhanced/internal/ui/styles"
)

// AboutScreen shows information about the application.
type AboutScreen struct {
	BaseScreen
	Styles styles.ContentStyles
	Keys   appkeys.Common
}

// NewAboutScreen creates a new about screen.
func NewAboutScreen(altScreen bool) *AboutScreen {
	return &AboutScreen{
		BaseScreen: BaseScreen{
			AltScreen:  altScreen,
			LoggerName: "AboutScreen",
			Header:     "About",
			AppTitle:   "Template-v2",
		},
		Keys:   appkeys.CommonBindings(),
		Styles: styles.NewContentStyles(false), // Default to light, will update in Appeared
	}
}

// Init initializes the about screen.
func (a *AboutScreen) Init() tea.Cmd {
	a.LogDebug("Initialized")
	return a.BaseScreen.Init()
}

// Update handles incoming messages.
func (a *AboutScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		if a.BaseScreen.UpdateBackgroundColor(msg) {
			a.Styles = styles.NewContentStyles(a.IsDark)
			return a, nil
		}

	case tea.KeyPressMsg:
		return a.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		a.BaseScreen.UpdateWindowSize(msg)
		return a, nil
	}

	return a, nil
}

// handleKeyPress processes keyboard input.
func (a *AboutScreen) handleKeyPress(msg tea.KeyPressMsg) (nav.Screen, tea.Cmd) {
	switch {
	case key.Matches(msg, a.Keys.Quit):
		a.LogDebug("Quit key pressed")
		return a, tea.Quit

	case key.Matches(msg, a.Keys.Back):
		a.LogDebug("Back key pressed")
		return a, nav.Pop()
	}

	return a, nil
}

// View renders the about screen.
func (a *AboutScreen) View() tea.View {
	return a.BaseScreen.View(a.render())
}

// render builds the visual representation of the about screen.
func (a *AboutScreen) render() string {
	// Build header style
	headerStyle := a.Styles.Header.Width(a.GetContentWidth())

	// Build content sections
	content := fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s",
		a.Styles.Title.Render("Template V2 Enhanced"),
		a.Styles.Label.Render("Version: 1.0.0"),
		a.Styles.Content.Render("A production-ready scaffold for building"),
		a.Styles.Content.Render("terminal user interface applications"),
		a.Styles.Content.Render("using BubbleTea v2, Bubbles v2, and Lip Gloss v2."),
		a.Styles.Content.Render("Built with navigation support for complex UIs."),
	)

	// Build help text
	helpText := "Back: b/esc | Quit: q"

	// Build border style with dynamic width
	borderStyle := a.Styles.Border.Width(a.GetContentWidth())

	// Use the layout with header, content, and footer
	layoutContent := a.BaseScreen.RenderLayoutWithBorder(
		a.Header,
		content,
		helpText,
		headerStyle,
		a.Styles.Help,
		borderStyle,
	)

	// Center the whole layout if in alt screen
	return a.BaseScreen.CenterContent(layoutContent)
}

// Appeared is called when the screen becomes active.
func (a *AboutScreen) Appeared() tea.Cmd {
	a.LogDebug("Appeared")
	a.Styles = styles.NewContentStyles(a.IsDark)
	return nil
}

// Disappeared is called when the screen loses active status.
func (a *AboutScreen) Disappeared() {
	a.LogDebug("Disappeared")
}
