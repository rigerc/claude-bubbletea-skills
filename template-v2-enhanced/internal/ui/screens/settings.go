// Package screens provides example screen implementations for the navigation demo.
package screens

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"template-v2-enhanced/internal/ui/nav"
	appkeys "template-v2-enhanced/internal/ui/keys"
	"template-v2-enhanced/internal/ui/styles"
)

// SettingsScreen demonstrates settings options with toggle switches.
type SettingsScreen struct {
	BaseScreen

	// selectedOption tracks which setting is selected.
	selectedOption int

	// settings tracks the current values of settings.
	settings map[string]bool

	// Styles contains the theme-aware styles.
	Styles styles.MenuStyles

	// Keys contains the key bindings.
	Keys appkeys.Common
}

// settingOptions defines the available settings.
var settingOptions = []struct {
	key   string
	text  string
	field string
}{
	{"1", "Dark Mode", "darkMode"},
	{"2", "Notifications", "notifications"},
	{"3", "Auto Save", "autoSave"},
}

// NewSettingsScreen creates a new settings screen.
func NewSettingsScreen(altScreen bool) *SettingsScreen {
	return &SettingsScreen{
		BaseScreen: BaseScreen{
			AltScreen:  altScreen,
			LoggerName: "SettingsScreen",
			Header:     "Settings",
			AppTitle:   "Template-v2",
		},
		selectedOption: 0,
		settings: map[string]bool{
			"darkMode":      true,
			"notifications": false,
			"autoSave":      true,
		},
		Keys:   appkeys.CommonBindings(),
		Styles: styles.NewMenuStyles(false), // Default to light, will update in Appeared
	}
}

// Init initializes the settings screen.
func (s *SettingsScreen) Init() tea.Cmd {
	s.LogDebug("Initialized")
	return s.BaseScreen.Init()
}

// Update handles incoming messages.
func (s *SettingsScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		if s.BaseScreen.UpdateBackgroundColor(msg) {
			s.Styles = styles.NewMenuStyles(s.IsDark)
			return s, nil
		}

	case tea.KeyPressMsg:
		return s.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		s.BaseScreen.UpdateWindowSize(msg)
		return s, nil
	}

	return s, nil
}

// handleKeyPress processes keyboard input.
func (s *SettingsScreen) handleKeyPress(msg tea.KeyPressMsg) (nav.Screen, tea.Cmd) {
	switch {
	case key.Matches(msg, s.Keys.Quit):
		s.LogDebug("Quit key pressed")
		return s, tea.Quit

	case key.Matches(msg, s.Keys.Back):
		s.LogDebug("Back key pressed")
		return s, nav.Pop()

	case key.Matches(msg, s.Keys.Up):
		if s.selectedOption > 0 {
			s.selectedOption--
		}
		return s, nil

	case key.Matches(msg, s.Keys.Down):
		if s.selectedOption < len(settingOptions)-1 {
			s.selectedOption++
		}
		return s, nil

	case key.Matches(msg, s.Keys.Enter), key.Matches(msg, s.Keys.Space):
		field := settingOptions[s.selectedOption].field
		s.settings[field] = !s.settings[field]
		s.LogDebugf("Toggled %s to %v", field, s.settings[field])
		return s, nil
	}

	return s, nil
}

// View renders the settings screen.
func (s *SettingsScreen) View() tea.View {
	return s.BaseScreen.View(s.render())
}

// render builds the visual representation of the settings screen.
func (s *SettingsScreen) render() string {
	// Build header style
	headerStyle := s.Styles.CommonStyles.Header.Width(s.GetContentWidth())

	// Build settings options content
	var content strings.Builder
	for i, opt := range settingOptions {
		var style lipgloss.Style
		if i == s.selectedOption {
			style = s.Styles.Selected
		} else {
			style = s.Styles.Option
		}

		// Get current value
		value := "OFF"
		if s.settings[opt.field] {
			value = "ON"
		}

		content.WriteString(style.Render(fmt.Sprintf("%s. %s: %s", opt.key, opt.text, value)))
		content.WriteByte('\n')
	}

	// Build help text
	helpText := "Navigate: ↑/k | ↓/j | Toggle: Enter/Space | Back: b | Quit: q"

	// Build border style with dynamic width
	borderStyle := s.Styles.Border.Width(s.GetContentWidth())

	// Use the layout with header, content, and footer
	layoutContent := s.BaseScreen.RenderLayoutWithBorder(
		s.Header,
		content.String(),
		helpText,
		headerStyle,
		s.Styles.Help,
		borderStyle,
	)

	// Center the whole layout if in alt screen
	return s.BaseScreen.CenterContent(layoutContent)
}

// Appeared is called when the screen becomes active.
func (s *SettingsScreen) Appeared() tea.Cmd {
	s.LogDebug("Appeared")
	s.Styles = styles.NewMenuStyles(s.IsDark)
	return nil
}

// Disappeared is called when the screen loses active status.
func (s *SettingsScreen) Disappeared() {
	s.LogDebug("Disappeared")
}
