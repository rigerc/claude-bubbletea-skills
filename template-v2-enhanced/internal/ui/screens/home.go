// Package screens provides example screen implementations for the navigation demo.
package screens

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"template-v2-enhanced/internal/ui/nav"
	appkeys "template-v2-enhanced/internal/ui/keys"
)

// HomeScreen is the main home screen with navigation options.
type HomeScreen struct {
	BaseScreen

	// options stores the menu options so SetSize can use the item count.
	options []MenuOption

	// list is the Bubbles list component used for rendering and navigation.
	list list.Model

	// Keys contains the key bindings.
	Keys appkeys.Common
}

// NewHomeScreen creates a new home screen.
func NewHomeScreen(altScreen bool) *HomeScreen {
	// opts enumerates every item on the home menu. Add new entries here;
	// no other code needs to change.
	opts := []MenuOption{
		{Key: "1", Text: "View Details", Action: nav.Push(NewDetailsScreen(altScreen))},
		{Key: "2", Text: "Settings", Action: nav.Push(NewSettingsScreen(altScreen))},
		{Key: "3", Text: "About", Action: nav.Push(NewAboutScreen(altScreen))},
	}

	delegate := newMenuDelegate()
	l := list.New(buildListItems(opts), delegate, 0, 0)
	l.Title = "Main Menu"
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	return &HomeScreen{
		BaseScreen: BaseScreen{
			AltScreen:  altScreen,
			LoggerName: "HomeScreen",
			Header:     "Header",
			AppTitle:   "Template-v2",
		},
		options: opts,
		list:    l,
		Keys:    appkeys.CommonBindings(),
	}
}

// Init initializes the home screen.
func (h *HomeScreen) Init() tea.Cmd {
	h.LogDebug("Initialized")
	return h.BaseScreen.Init()
}

// Update handles incoming messages.
func (h *HomeScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		if h.BaseScreen.UpdateBackgroundColor(msg) {
			h.applyTheme()
			return h, nil
		}

	case tea.KeyPressMsg:
		return h.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		h.BaseScreen.UpdateWindowSize(msg)
		h.list.SetSize(msg.Width, listHeight(len(h.options), msg.Height))
		return h, nil
	}

	var cmd tea.Cmd
	h.list, cmd = h.list.Update(msg)
	return h, cmd
}

// handleKeyPress processes keyboard input.
// Application-level keys (quit, enter) are handled here; all other
// key events are forwarded to the list for built-in navigation.
func (h *HomeScreen) handleKeyPress(msg tea.KeyPressMsg) (nav.Screen, tea.Cmd) {
	switch {
	case key.Matches(msg, h.Keys.Quit):
		h.LogDebug("Quit key pressed")
		return h, tea.Quit

	case key.Matches(msg, h.Keys.Enter):
		selected, ok := h.list.SelectedItem().(menuItem)
		if !ok {
			return h, nil
		}
		h.LogDebugf("Selected option: %s", selected.opt.Text)
		return h, selected.opt.Action
	}

	var cmd tea.Cmd
	h.list, cmd = h.list.Update(msg)
	return h, cmd
}

// View renders the home screen.
func (h *HomeScreen) View() tea.View {
	return h.BaseScreen.View(h.list.View())
}

// Appeared is called when the screen becomes active.
func (h *HomeScreen) Appeared() tea.Cmd {
	h.LogDebug("Appeared")
	h.applyTheme()
	h.list.SetSize(h.Width, listHeight(len(h.options), h.Height))
	return nil
}

// Disappeared is called when the screen loses active status.
func (h *HomeScreen) Disappeared() {
	h.LogDebug("Disappeared")
}

// applyTheme updates the list and delegate styles to match the current theme.
func (h *HomeScreen) applyTheme() {
	h.list.Styles = list.DefaultStyles(h.IsDark)
	delegate := newMenuDelegate()
	delegate.Styles = list.NewDefaultItemStyles(h.IsDark)
	h.list.SetDelegate(delegate)
}

