// Package help provides a help component wrapping bubbles v2 help.
package help

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"template-v2-enhanced/internal/ui/keys"
)

// Model wraps bubbles/help with app-specific configuration.
// It provides automatic theme support and common key bindings.
type Model struct {
	help  help.Model
	width int
}

// New creates a new help component with default configuration.
// The help component will use shared key bindings by default.
func New() Model {
	h := help.New()
	return Model{help: h}
}

// NewWithWidth creates a new help component with a specific width.
func NewWithWidth(width int) Model {
	m := New()
	m.width = width
	m.help.SetWidth(width)
	return m
}

// Update handles messages for the help component.
// It forwards messages to the underlying bubbles help component.
func (m Model) Update(msg tea.Msg) Model {
	var cmd tea.Cmd
	m.help, cmd = m.help.Update(msg)
	_ = cmd // Ignore commands from help component
	return m
}

// SetWidth sets the help width.
func (m *Model) SetWidth(width int) {
	m.width = width
	m.help.SetWidth(width)
}

// SetStyles updates help styles for the current theme.
// This should be called when the terminal background color changes.
func (m *Model) SetStyles(isDark bool) {
	m.help.Styles = help.DefaultStyles(isDark)
}

// View renders the help with common app bindings.
// It displays the help text using the shared key bindings.
func (m Model) View() string {
	common := keys.CommonBindings()
	return m.help.View(common)
}

// ViewWithBindings renders the help with custom bindings.
// Use this when a screen needs custom key bindings beyond the common ones.
func (m Model) ViewWithBindings(bindings help.KeyMap) string {
	return m.help.View(bindings)
}

// ViewMinimal renders the help with only the specified bindings.
// Use this to highlight specific key bindings in the help display.
func (m Model) ViewMinimal(bindings []key.Binding) string {
	return m.help.View(MinimalKeyMap{Bindings: bindings})
}

// MinimalKeyMap implements help.KeyMap for a slice of bindings.
type MinimalKeyMap struct {
	Bindings []key.Binding
}

// ShortHelp returns the bindings for short help display.
func (m MinimalKeyMap) ShortHelp() []key.Binding {
	return m.Bindings
}

// FullHelp returns the bindings for full help display.
func (m MinimalKeyMap) FullHelp() [][]key.Binding {
	if len(m.Bindings) == 0 {
		return nil
	}
	// Split into groups of 4 for display
	var groups [][]key.Binding
	for i := 0; i < len(m.Bindings); i += 4 {
		end := i + 4
		if end > len(m.Bindings) {
			end = len(m.Bindings)
		}
		groups = append(groups, m.Bindings[i:end])
	}
	return groups
}
