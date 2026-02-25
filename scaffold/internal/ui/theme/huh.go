package theme

import (
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

// HuhTheme returns a huh.Theme that matches the application palette.
// Uses huh.ThemeFunc so huh drives isDark on every View() call.
func HuhTheme() huh.Theme {
	return huh.ThemeFunc(func(isDark bool) *huh.Styles {
		p := NewPalette(isDark)
		s := huh.ThemeCharm(isDark)

		s.Group.Title = s.Group.Title.Foreground(p.Accent).Bold(true)
		s.Group.Description = s.Group.Description.Foreground(p.Muted)

		s.Focused.Base = s.Focused.Base.BorderForeground(p.Accent)
		s.Focused.Title = s.Focused.Title.Foreground(p.Accent)
		s.Focused.Description = s.Focused.Description.Foreground(p.Muted)
		s.Focused.SelectSelector = s.Focused.SelectSelector.Foreground(p.AccentHover)
		s.Focused.NextIndicator = s.Focused.NextIndicator.Foreground(p.AccentHover)
		s.Focused.PrevIndicator = s.Focused.PrevIndicator.Foreground(p.AccentHover)
		s.Focused.FocusedButton = s.Focused.FocusedButton.Background(p.Accent).Foreground(p.Inverse)
		s.Focused.TextInput.Cursor = s.Focused.TextInput.Cursor.Foreground(p.AccentHover)
		s.Focused.TextInput.Prompt = s.Focused.TextInput.Prompt.Foreground(p.Accent)
		s.Focused.ErrorMessage = s.Focused.ErrorMessage.Foreground(p.Error)
		s.Focused.ErrorIndicator = s.Focused.ErrorIndicator.Foreground(p.Error)

		s.Blurred = s.Focused
		s.Blurred.Base = s.Focused.Base.BorderStyle(lipgloss.HiddenBorder())

		return s
	})
}
