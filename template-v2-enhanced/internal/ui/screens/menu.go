// Package screens provides the individual screen implementations for the application.
package screens

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"template-v2-enhanced/internal/ui/nav"
)

// MenuItem represents an item in the navigation menu.
// It implements list.Item so it can be used with the list component.
type MenuItem struct {
	title, description string
	action             tea.Cmd // typically nav.Push(NewDetailScreen(...))
}

// Title returns the display title of the menu item.
// Implements list.Item.
func (i MenuItem) Title() string {
	return i.title
}

// Description returns the description text shown next to the title.
// Implements list.Item.
func (i MenuItem) Description() string {
	return i.description
}

// FilterValue returns the text used for filtering.
// Implements list.Item.
func (i MenuItem) FilterValue() string {
	return i.title
}

// NewMenuItem creates a new MenuItem with the given title, description, and action.
func NewMenuItem(title, description string, action tea.Cmd) MenuItem {
	return MenuItem{
		title:       title,
		description: description,
		action:      action,
	}
}

// delegateKeyMap defines key bindings specific to the list delegate.
type delegateKeyMap struct {
	choose key.Binding
}

// newDelegateKeyMap creates a new delegateKeyMap with initialized bindings.
func newDelegateKeyMap() delegateKeyMap {
	return delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
	}
}

// newMenuDelegate creates a styled list delegate with the given key bindings.
// The delegate handles the "enter" key to trigger the selected item's action.
func newMenuDelegate(dKeys delegateKeyMap, isDark bool) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles = list.NewDefaultItemStyles(isDark)
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		if kp, ok := msg.(tea.KeyPressMsg); ok {
			if key.Matches(kp, dKeys.choose) {
				if item, ok := m.SelectedItem().(MenuItem); ok && item.action != nil {
					return item.action
				}
			}
		}
		return nil
	}
	d.ShortHelpFunc = func() []key.Binding { return []key.Binding{dKeys.choose} }
	d.FullHelpFunc = func() [][]key.Binding { return [][]key.Binding{{dKeys.choose}} }
	return d
}

// MenuScreen displays a navigable list of menu items.
// It implements nav.Screen and nav.Themeable.
type MenuScreen struct {
	ScreenBase
	list         list.Model
	delegateKeys delegateKeyMap
}

// NewMenuScreen creates a new MenuScreen with the given title and items.
// The isDark parameter should be false initially; the correct value will be
// set via SetTheme when the screen is pushed onto the stack.
func NewMenuScreen(title string, items []list.Item, isDark bool) *MenuScreen {
	s := &MenuScreen{
		ScreenBase:   NewBase(isDark),
		delegateKeys: newDelegateKeyMap(),
	}
	// Use true as the initial isDark to match the library's own default (which
	// hardcodes isDark=true). The correct value is applied via SetTheme /
	// BackgroundColorMsg before the first meaningful render.
	d := newMenuDelegate(s.delegateKeys, true)

	l := list.New(items, d, 0, 0) // 0,0: WindowSizeMsg drives size
	l.Title = title
	l.Styles = list.DefaultStyles(isDark)
	l.Styles.Title = s.Theme.Title // branded title override
	l.DisableQuitKeybindings()     // REQUIRED: prevents list eating ctrl+c/q
	s.list = l

	return s
}

// Init returns nil (no initial commands needed).
func (s *MenuScreen) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and returns an updated screen and command.
func (s *MenuScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.Width, s.Height = msg.Width, msg.Height
		s.updateListSize()

	case tea.KeyPressMsg:
		// ESC pops stack ONLY when not filtering
		if msg.String() == "esc" && s.list.FilterState() == list.Unfiltered {
			return s, nav.Pop()
		}
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

// View renders the menu list wrapped with the theme's margin style.
func (s *MenuScreen) View() string {
	return s.Theme.App.Render(s.list.View())
}

// SetTheme updates the screen's theme based on the terminal background.
// Implements nav.Themeable.
func (s *MenuScreen) SetTheme(isDark bool) {
	s.ApplyTheme(isDark)
	s.list.Styles = list.DefaultStyles(s.IsDark)
	s.list.Styles.Title = s.Theme.Title
	s.list.SetDelegate(newMenuDelegate(s.delegateKeys, s.IsDark))
}

// updateListSize recalculates the list dimensions based on window size
// and theme frame. Height accounts for title, items (with descriptions), and help.
func (s *MenuScreen) updateListSize() {
	if !s.IsSized() {
		return
	}
	frameH, frameV := s.Theme.App.GetFrameSize()

	// Calculate available height after frame
	availH := s.Height - frameV

	// Get actual item count (not filtered)
	itemCount := len(s.list.Items())
	if itemCount == 0 {
		itemCount = 1
	}

	// Each item is 2 lines when description is shown (title + description)
	// Add space for title bar (1 line) + help section (4 lines)
	itemLines := itemCount * 4
	titleLines := 1
	helpLines := 4
	targetH := itemLines + titleLines + helpLines

	// Clamp to available height
	if targetH > availH {
		targetH = availH
	}
	if targetH < 10 {
		targetH = 10
	}

	s.list.SetSize(s.Width-frameH, targetH)
}
