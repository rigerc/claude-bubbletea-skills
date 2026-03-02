package screens

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

// inlineSelect wraps huh.Select to render label, description, and options
// on a single line for compact settings forms.
// Navigation is delegated to the underlying Select via Update().
type inlineSelect struct {
	*huh.Select[string]

	label       string
	description string
	width       int
	focused     bool
	theme       huh.Theme
	hasDarkBg   bool
}

// newInlineSelect creates an inline select field with label and description.
func newInlineSelect(label, description string, select_ *huh.Select[string]) *inlineSelect {
	return &inlineSelect{
		Select:      select_,
		label:       label,
		description: description,
	}
}

// Init initializes the field.
func (f *inlineSelect) Init() tea.Cmd {
	return f.Select.Init()
}

// Update handles messages - delegates to underlying Select for navigation.
func (f *inlineSelect) Update(msg tea.Msg) (huh.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		f.hasDarkBg = msg.IsDark()
	}

	m, cmd := f.Select.Update(msg)
	if s, ok := m.(*huh.Select[string]); ok {
		f.Select = s
	}
	return f, cmd
}

// View renders the field with label, description, and select options on one line.
func (f *inlineSelect) View() string {
	styles := f.activeStyles()

	// Render our own inline select view using the Select's current value
	selectView := f.renderInlineOptions(styles)

	// Build inline layout: "Label: description [‹ Option ›]"
	var parts []string

	if f.label != "" {
		parts = append(parts, styles.Title.Render(f.label+":"))
	}

	if f.description != "" {
		parts = append(parts, styles.Description.Render(f.description))
	}

	parts = append(parts, selectView)

	line := lipgloss.JoinHorizontal(lipgloss.Left, parts...)

	return styles.Base.Width(f.width).Render(line)
}

// renderInlineOptions renders the select options inline with prev/next indicators.
func (f *inlineSelect) renderInlineOptions(styles *huh.FieldStyles) string {
	// Get the current value from the select
	value := f.Select.GetValue()
	valueStr, _ := value.(string)

	// Get display name for the value - try to get the option key
	displayValue := valueStr
	// The Select stores options internally but doesn't expose them
	// We'll just use the value directly

	prevIndicator := "‹"
	nextIndicator := "›"

	// Style the indicators
	prevStyle := styles.PrevIndicator
	nextStyle := styles.NextIndicator

	if !f.focused {
		prevIndicator = ""
		nextIndicator = ""
	}

	valueStyle := styles.SelectedOption
	if !f.focused {
		valueStyle = styles.UnselectedOption
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		prevStyle.Render(prevIndicator),
		" ",
		valueStyle.Render(displayValue),
		" ",
		nextStyle.Render(nextIndicator),
	)
}

// Focus focuses the field.
func (f *inlineSelect) Focus() tea.Cmd {
	f.focused = true
	return f.Select.Focus()
}

// Blur blurs the field.
func (f *inlineSelect) Blur() tea.Cmd {
	f.focused = false
	return f.Select.Blur()
}

// KeyBinds returns key bindings.
func (f *inlineSelect) KeyBinds() []key.Binding {
	return f.Select.KeyBinds()
}

// Error returns any validation error.
func (f *inlineSelect) Error() error {
	return f.Select.Error()
}

// Skip returns false - this field should not be skipped.
func (f *inlineSelect) Skip() bool {
	return f.Select.Skip()
}

// Zoom returns false.
func (f *inlineSelect) Zoom() bool {
	return f.Select.Zoom()
}

// WithTheme sets the theme.
func (f *inlineSelect) WithTheme(theme huh.Theme) huh.Field {
	f.theme = theme
	f.Select = f.Select.WithTheme(theme).(*huh.Select[string])
	return f
}

// WithKeyMap sets the keymap.
func (f *inlineSelect) WithKeyMap(k *huh.KeyMap) huh.Field {
	f.Select = f.Select.WithKeyMap(k).(*huh.Select[string])
	return f
}

// WithWidth sets the width.
func (f *inlineSelect) WithWidth(width int) huh.Field {
	f.width = width
	f.Select = f.Select.WithWidth(width).(*huh.Select[string])
	return f
}

// WithHeight sets the height (no-op for inline).
func (f *inlineSelect) WithHeight(height int) huh.Field {
	return f
}

// WithPosition sets the field position.
func (f *inlineSelect) WithPosition(p huh.FieldPosition) huh.Field {
	f.Select = f.Select.WithPosition(p).(*huh.Select[string])
	return f
}

// GetKey returns the field key.
func (f *inlineSelect) GetKey() string {
	return f.Select.GetKey()
}

// GetValue returns the field value.
func (f *inlineSelect) GetValue() any {
	return f.Select.GetValue()
}

func (f *inlineSelect) activeStyles() *huh.FieldStyles {
	theme := f.theme
	if theme == nil {
		theme = huh.ThemeFunc(huh.ThemeCharm)
	}
	if f.focused {
		return &theme.Theme(f.hasDarkBg).Focused
	}
	return &theme.Theme(f.hasDarkBg).Blurred
}

// Ensure inlineSelect implements huh.Field
var _ huh.Field = (*inlineSelect)(nil)
