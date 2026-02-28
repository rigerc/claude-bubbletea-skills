package theme

import (
	"sort"

	"charm.land/lipgloss/v2"
)

// Registry holds named themes and tracks the active theme + dark mode state.
type Registry struct {
	themes  map[string]Theme
	order   []string // insertion order for cycling
	current string
	isDark  bool
}

// NewRegistry creates a registry pre-populated with built-in themes.
func NewRegistry() *Registry {
	r := &Registry{
		themes:  make(map[string]Theme),
		current: "default",
	}
	r.Register(defaultTheme())
	r.Register(draculaTheme())
	r.Register(catppuccinTheme())
	return r
}

// Register adds a theme to the registry.
func (r *Registry) Register(t Theme) {
	if _, exists := r.themes[t.Name]; !exists {
		r.order = append(r.order, t.Name)
	}
	r.themes[t.Name] = t
}

// Names returns sorted theme names.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.themes))
	for name := range r.themes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Current returns the active theme.
func (r *Registry) Current() Theme { return r.themes[r.current] }

// CurrentName returns the active theme name.
func (r *Registry) CurrentName() string { return r.current }

// SetCurrent switches the active theme by name. Returns false if not found.
func (r *Registry) SetCurrent(name string) bool {
	if _, ok := r.themes[name]; !ok {
		return false
	}
	r.current = name
	return true
}

// CycleNext advances to the next theme in insertion order.
func (r *Registry) CycleNext() string {
	for i, name := range r.order {
		if name == r.current {
			next := r.order[(i+1)%len(r.order)]
			r.current = next
			return next
		}
	}
	return r.current
}

// Colors returns the active ColorSet based on current theme + isDark.
func (r *Registry) Colors() ColorSet {
	return r.Current().Colors(r.isDark)
}

// SetDark updates the dark mode state.
func (r *Registry) SetDark(isDark bool) { r.isDark = isDark }

// IsDark returns the current dark mode state.
func (r *Registry) IsDark() bool { return r.isDark }

// --- Built-in themes ---

func defaultTheme() Theme {
	return Theme{
		Name: "default",
		Light: ColorSet{
			Primary:   lipgloss.Color("#6C63FF"),
			Secondary: lipgloss.Color("#2BC4C4"),
			Text:      lipgloss.Color("#1A1A2E"),
			Subtle:    lipgloss.Color("#888899"),
			Surface:   lipgloss.Color("#F0F0F5"),
			Overlay:   lipgloss.Color("#E0E0E8"),
			Error:     lipgloss.Color("#E53E3E"),
			Warning:   lipgloss.Color("#DD8800"),
			Success:   lipgloss.Color("#38A169"),
			Border:    lipgloss.Color("#C0C0CC"),
		},
		Dark: ColorSet{
			Primary:   lipgloss.Color("#7C73FF"),
			Secondary: lipgloss.Color("#3BD4D4"),
			Text:      lipgloss.Color("#E1E1EF"),
			Subtle:    lipgloss.Color("#606072"),
			Surface:   lipgloss.Color("#1E1E2E"),
			Overlay:   lipgloss.Color("#313145"),
			Error:     lipgloss.Color("#FC8181"),
			Warning:   lipgloss.Color("#F6AD55"),
			Success:   lipgloss.Color("#68D391"),
			Border:    lipgloss.Color("#44445A"),
		},
	}
}

func draculaTheme() Theme {
	return Theme{
		Name: "dracula",
		Light: ColorSet{
			Primary:   lipgloss.Color("#7C3AED"),
			Secondary: lipgloss.Color("#EC4899"),
			Text:      lipgloss.Color("#282A36"),
			Subtle:    lipgloss.Color("#6272A4"),
			Surface:   lipgloss.Color("#F8F8F2"),
			Overlay:   lipgloss.Color("#E8E8EF"),
			Error:     lipgloss.Color("#FF5555"),
			Warning:   lipgloss.Color("#FFB86C"),
			Success:   lipgloss.Color("#50FA7B"),
			Border:    lipgloss.Color("#BD93F9"),
		},
		Dark: ColorSet{
			Primary:   lipgloss.Color("#BD93F9"),
			Secondary: lipgloss.Color("#FF79C6"),
			Text:      lipgloss.Color("#F8F8F2"),
			Subtle:    lipgloss.Color("#6272A4"),
			Surface:   lipgloss.Color("#282A36"),
			Overlay:   lipgloss.Color("#44475A"),
			Error:     lipgloss.Color("#FF5555"),
			Warning:   lipgloss.Color("#FFB86C"),
			Success:   lipgloss.Color("#50FA7B"),
			Border:    lipgloss.Color("#6272A4"),
		},
	}
}

func catppuccinTheme() Theme {
	return Theme{
		Name: "catppuccin",
		Light: ColorSet{
			Primary:   lipgloss.Color("#8839EF"), // Mauve (Latte)
			Secondary: lipgloss.Color("#1E66F5"), // Blue (Latte)
			Text:      lipgloss.Color("#4C4F69"), // Text (Latte)
			Subtle:    lipgloss.Color("#9CA0B0"), // Overlay1 (Latte)
			Surface:   lipgloss.Color("#EFF1F5"), // Base (Latte)
			Overlay:   lipgloss.Color("#E6E9EF"), // Mantle (Latte)
			Error:     lipgloss.Color("#D20F39"), // Red (Latte)
			Warning:   lipgloss.Color("#DF8E1D"), // Yellow (Latte)
			Success:   lipgloss.Color("#40A02B"), // Green (Latte)
			Border:    lipgloss.Color("#ACB0BE"), // Overlay0 (Latte)
		},
		Dark: ColorSet{
			Primary:   lipgloss.Color("#CBA6F7"), // Mauve (Mocha)
			Secondary: lipgloss.Color("#89B4FA"), // Blue (Mocha)
			Text:      lipgloss.Color("#CDD6F4"), // Text (Mocha)
			Subtle:    lipgloss.Color("#7F849C"), // Overlay1 (Mocha)
			Surface:   lipgloss.Color("#1E1E2E"), // Base (Mocha)
			Overlay:   lipgloss.Color("#313244"), // Surface0 (Mocha)
			Error:     lipgloss.Color("#F38BA8"), // Red (Mocha)
			Warning:   lipgloss.Color("#F9E2AF"), // Yellow (Mocha)
			Success:   lipgloss.Color("#A6E3A1"), // Green (Mocha)
			Border:    lipgloss.Color("#6C7086"), // Overlay0 (Mocha)
		},
	}
}
