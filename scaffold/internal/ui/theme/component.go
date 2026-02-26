package theme

// ThemeAware provides reusable theme state holder for components.
type ThemeAware struct {
	themeState State
}

// ApplyThemeState updates the stored theme state.
func (t *ThemeAware) ApplyThemeState(state State) {
	t.themeState = state
}

// ThemeState returns current theme state.
func (t *ThemeAware) ThemeState() State {
	return t.themeState
}

// Palette is a convenience accessor.
func (t *ThemeAware) Palette() Palette {
	return t.themeState.Palette
}

// IsDark is a convenience accessor.
func (t *ThemeAware) IsDark() bool {
	return t.themeState.IsDark
}

// ThemeName is a convenience accessor.
func (t *ThemeAware) ThemeName() string {
	return t.themeState.Name
}

// ThemeWidth is a convenience accessor.
func (t *ThemeAware) ThemeWidth() int {
	return t.themeState.Width
}
