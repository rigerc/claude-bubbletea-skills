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

// DetailsScreen shows detailed information with a back button.
type DetailsScreen struct {
	BaseScreen
	Styles styles.ContentStyles
	Keys   appkeys.Common
}

// NewDetailsScreen creates a new details screen.
func NewDetailsScreen(altScreen bool) *DetailsScreen {
	return &DetailsScreen{
		BaseScreen: BaseScreen{
			AltScreen:  altScreen,
			LoggerName: "DetailsScreen",
			Header:     "Details",
			AppTitle:   "Template-v2",
		},
		Keys:   appkeys.CommonBindings(),
		Styles: styles.NewContentStyles(false), // Default to light, will update in Appeared
	}
}

// Init initializes the details screen.
func (d *DetailsScreen) Init() tea.Cmd {
	d.LogDebug("Initialized")
	return d.BaseScreen.Init()
}

// Update handles incoming messages.
func (d *DetailsScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		if d.BaseScreen.UpdateBackgroundColor(msg) {
			d.Styles = styles.NewContentStyles(d.IsDark)
			return d, nil
		}

	case tea.KeyPressMsg:
		return d.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		d.BaseScreen.UpdateWindowSize(msg)
		return d, nil
	}

	return d, nil
}

// handleKeyPress processes keyboard input.
func (d *DetailsScreen) handleKeyPress(msg tea.KeyPressMsg) (nav.Screen, tea.Cmd) {
	switch {
	case key.Matches(msg, d.Keys.Quit):
		d.LogDebug("Quit key pressed")
		return d, tea.Quit

	case key.Matches(msg, d.Keys.Back):
		d.LogDebug("Back key pressed")
		return d, nav.Pop()
	}

	return d, nil
}

// View renders the details screen.
func (d *DetailsScreen) View() tea.View {
	return d.BaseScreen.View(d.render())
}

// render builds the visual representation of the details screen.
func (d *DetailsScreen) render() string {
	// Build header style
	headerStyle := d.Styles.Header.Width(d.GetContentWidth())

	// Build content
	content := fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s",
		d.Styles.Content.Render("This is the details screen."),
		d.Styles.Content.Render("You can navigate back to the home screen."),
		d.Styles.Content.Render("The navigation stack maintains your history."),
		d.Styles.Content.Render("Try pressing 'b' or 'esc' to go back."),
	)

	// Build help text
	helpText := "Back: b/esc | Quit: q"

	// Build border style with dynamic width
	borderStyle := d.Styles.Border.Width(d.GetContentWidth())

	// Use the layout with header, content, and footer
	layoutContent := d.BaseScreen.RenderLayoutWithBorder(
		d.Header,
		content,
		helpText,
		headerStyle,
		d.Styles.Help,
		borderStyle,
	)

	// Center the whole layout if in alt screen
	return d.BaseScreen.CenterContent(layoutContent)
}

// Appeared is called when the screen becomes active.
func (d *DetailsScreen) Appeared() tea.Cmd {
	d.LogDebug("Appeared")
	d.Styles = styles.NewContentStyles(d.IsDark)
	return nil
}

// Disappeared is called when the screen loses active status.
func (d *DetailsScreen) Disappeared() {
	d.LogDebug("Disappeared")
}
