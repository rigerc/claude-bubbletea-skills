package theme

// ThemeSwitchMsg requests a theme switch by name.
// Sent by screens (e.g., settings). The root model handles it.
type ThemeSwitchMsg struct {
	Name string
}

// ThemeChangedMsg is broadcast when the active theme or dark mode changes.
// All components should rebuild their styles upon receiving this.
type ThemeChangedMsg struct {
	Colors ColorSet
	Name   string
	IsDark bool
}
