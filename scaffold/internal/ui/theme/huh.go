package theme

import (
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

// HuhTheme returns a huh.Theme that matches the application palette for the given theme name.
// Uses huh.ThemeFunc so huh drives isDark on every View() call.
func HuhTheme(name string) huh.Theme {
	return huh.ThemeFunc(func(isDark bool) *huh.Styles {
		p := NewPalette(name, isDark)
		s := huh.ThemeCharm(isDark)

		s.Focused.Base = s.Focused.Base.
			Padding(0, 1, 0, 2).
			BorderStyle(lipgloss.ThickBorder()).
			BorderTop(false).BorderRight(false).BorderBottom(false).BorderLeft(true).
			BorderForeground(p.Focus)
		s.Focused.Card = s.Focused.Base
		s.Focused.Title = s.Focused.Title.Foreground(p.Primary).Bold(true).MarginRight(1)
		s.Focused.NoteTitle = s.Focused.NoteTitle.Foreground(p.Primary).Bold(true)
		s.Focused.Directory = s.Focused.Directory.Foreground(p.Primary)
		s.Focused.Description = s.Focused.Description.Foreground(p.ForegroundMuted)
		s.Focused.ErrorIndicator = s.Focused.ErrorIndicator.Foreground(p.Error)
		s.Focused.ErrorMessage = s.Focused.ErrorMessage.Foreground(p.Error)
		s.Focused.SelectSelector = s.Focused.SelectSelector.Foreground(p.Focus)
		s.Focused.NextIndicator = s.Focused.NextIndicator.Foreground(p.Primary)
		s.Focused.PrevIndicator = s.Focused.PrevIndicator.Foreground(p.Primary)
		s.Focused.Option = s.Focused.Option.Foreground(p.ForegroundSubtle)
		s.Focused.MultiSelectSelector = s.Focused.MultiSelectSelector.Foreground(p.Primary)
		s.Focused.SelectedOption = s.Focused.SelectedOption.Foreground(p.OnPrimary).Background(p.Primary).Padding(0, 1)
		s.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(p.SecondaryMuted).SetString("✓ ")
		s.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(p.SecondaryMuted).SetString("• ")
		s.Focused.UnselectedOption = s.Focused.UnselectedOption.Foreground(p.ForegroundSubtle)
		s.Focused.FocusedButton = s.Focused.FocusedButton.Foreground(p.OnPrimary).Background(p.Primary)
		s.Focused.Next = s.Focused.FocusedButton
		s.Focused.BlurredButton = s.Focused.BlurredButton.Foreground(p.ForegroundSubtle)

		s.Focused.TextInput.Cursor = s.Focused.TextInput.Cursor.Foreground(p.Success)
		s.Focused.TextInput.Placeholder = s.Focused.TextInput.Placeholder.Foreground(p.ForegroundSubtle)
		s.Focused.TextInput.Prompt = s.Focused.TextInput.Prompt.Foreground(p.Focus)

		s.Blurred = s.Focused
		s.Blurred.Base = s.Focused.Base.BorderStyle(lipgloss.HiddenBorder()).Padding(0, 1, 0, 2)
		s.Blurred.Title = s.Blurred.Title.Foreground(p.ForegroundSubtle)
		s.Blurred.Card = s.Blurred.Base
		s.Blurred.SelectedOption = s.Blurred.SelectedOption.Foreground(p.ForegroundMuted).Background(p.PrimaryMuted).Padding(0, 1)
		s.Blurred.NextIndicator = lipgloss.NewStyle().Foreground(p.ForegroundSubtle)
		s.Blurred.PrevIndicator = lipgloss.NewStyle().Foreground(p.ForegroundSubtle)

		s.Group.Title = s.Focused.Title
		s.Group.Description = s.Focused.Description

		return s
	})
}
