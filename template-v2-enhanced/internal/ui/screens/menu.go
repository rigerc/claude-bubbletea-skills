// Package screens provides a generic menu screen for selection-based UIs.
package screens

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"template-v2-enhanced/internal/ui/nav"
	appkeys "template-v2-enhanced/internal/ui/keys"
)

// MenuOption represents a single menu option.
type MenuOption struct {
	// Key is the shortcut key displayed (e.g., "1", "2", "a").
	Key string

	// Text is the display text for the option.
	Text string

	// Action is the command to execute when this option is selected.
	Action tea.Cmd
}

// menuItem adapts MenuOption to the list.Item interface.
// It implements list.DefaultItem so the built-in delegate can render it.
type menuItem struct {
	opt MenuOption
}

func (i menuItem) Title() string       { return i.opt.Key + ". " + i.opt.Text }
func (i menuItem) Description() string { return "" }
func (i menuItem) FilterValue() string { return i.opt.Text }

// newMenuDelegate creates a single-line list delegate.
func newMenuDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.SetHeight(1)
	d.SetSpacing(1)
	return d
}

// listHeight returns an appropriate list height for the given item count.
// Items are 1 row with no spacing; chrome (title, status bar, help) ~8 rows.
func listHeight(itemCount, windowHeight int) int {
	h := itemCount + 10
	if h > windowHeight {
		h = windowHeight
	}
	return h
}

// buildListItems converts MenuOptions to list.Items.
func buildListItems(opts []MenuOption) []list.Item {
	items := make([]list.Item, len(opts))
	for i, o := range opts {
		items[i] = menuItem{opt: o}
	}
	return items
}

// MenuScreen is a generic menu selection screen.
// It displays a list of options and handles navigation and selection.
// Options can be set at construction time or updated programmatically
// via SetOptions.
type MenuScreen struct {
	BaseScreen

	// Title is the main title for the menu.
	Title string

	// Subtitle is an optional subtitle.
	Subtitle string

	// options stores the original MenuOption values so actions survive
	// list re-ordering / filtering.
	options []MenuOption

	// list is the Bubbles list component used for rendering and navigation.
	list list.Model

	// Keys contains the key bindings.
	Keys appkeys.Common
}

// NewMenuScreen creates a new menu screen.
func NewMenuScreen(title, subtitle string, options []MenuOption, altScreen bool) *MenuScreen {
	delegate := newMenuDelegate()
	l := list.New(buildListItems(options), delegate, 0, 0)
	l.Title = title
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	return &MenuScreen{
		BaseScreen: BaseScreen{
			AltScreen:  altScreen,
			LoggerName: "MenuScreen",
		},
		Title:    title,
		Subtitle: subtitle,
		options:  options,
		list:     l,
		Keys:     appkeys.CommonBindings(),
	}
}

// SetOptions replaces the list items programmatically.
// This satisfies the requirement that items be dynamically buildable.
func (m *MenuScreen) SetOptions(options []MenuOption) tea.Cmd {
	m.options = options
	return m.list.SetItems(buildListItems(options))
}

// Init initializes the menu screen.
func (m *MenuScreen) Init() tea.Cmd {
	m.LogDebug("Initialized")
	return m.BaseScreen.Init()
}

// Update handles messages for the menu screen.
func (m *MenuScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		if m.BaseScreen.UpdateBackgroundColor(msg) {
			m.applyTheme()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.BaseScreen.UpdateWindowSize(msg)
		m.list.SetSize(msg.Width, listHeight(len(m.options), msg.Height))
		return m, nil

	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// handleKeyPress processes keyboard input.
// Application-level keys (quit, enter) are handled here; all other
// key events are forwarded to the list for built-in navigation.
func (m *MenuScreen) handleKeyPress(msg tea.KeyPressMsg) (nav.Screen, tea.Cmd) {
	switch {
	case key.Matches(msg, m.Keys.Quit):
		m.LogDebug("Quit key pressed")
		return m, tea.Quit

	case key.Matches(msg, m.Keys.Enter):
		selected, ok := m.list.SelectedItem().(menuItem)
		if !ok {
			return m, nil
		}
		m.LogDebugf("Selected option: %s", selected.opt.Text)
		return m, selected.opt.Action
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the menu screen.
func (m *MenuScreen) View() tea.View {
	return m.BaseScreen.View(m.list.View())
}

// Appeared is called when the screen becomes active.
func (m *MenuScreen) Appeared() tea.Cmd {
	m.LogDebug("Appeared")
	m.applyTheme()
	m.list.SetSize(m.Width, listHeight(len(m.options), m.Height))
	return nil
}

// Disappeared is called when the screen loses active status.
func (m *MenuScreen) Disappeared() {
	m.LogDebug("Disappeared")
}

// applyTheme updates the list and delegate styles to match the current theme.
func (m *MenuScreen) applyTheme() {
	m.list.Styles = list.DefaultStyles(m.IsDark)
	delegate := newMenuDelegate()
	delegate.Styles = list.NewDefaultItemStyles(m.IsDark)
	m.list.SetDelegate(delegate)
}
