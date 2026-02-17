// Package huh provides integration adapters for Huh-v2 forms.
package huh

import (
	"charm.land/huh/v2"
	lipgloss "charm.land/lipgloss/v2"
	"template-v2-enhanced/internal/ui/styles"
)

// ThemeFunc converts app.Theme to huh.ThemeFunc.
// It adapts the app's color scheme to Huh's styling system while
// preserving the green-based theme and adaptive light/dark colors.
func ThemeFunc(theme styles.Theme) huh.ThemeFunc {
	return func(isDark bool) *huh.Styles {
		// Start with Charm's default theme as base
		huhTheme := huh.ThemeCharm(isDark)
		ld := lipgloss.LightDark(isDark)

		// Override title color to match app's green theme
		huhTheme.Group.Title = huhTheme.Group.Title.
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065"))

		// Match status message colors to app theme
		huhTheme.Focused.Description = huhTheme.Focused.Description.
			Foreground(ld(
				lipgloss.Color("#04B575"),
				lipgloss.Color("#10CC85"),
			))

		// Remove borders - explicitly disable all sides
		huhTheme.Form.Base = huhTheme.Form.Base.
			BorderTop(false).
			BorderRight(false).
			BorderBottom(false).
			BorderLeft(false)
		huhTheme.Group.Base = huhTheme.Group.Base.
			BorderTop(false).
			BorderRight(false).
			BorderBottom(false).
			BorderLeft(false)
		huhTheme.Focused.Base = huhTheme.Focused.Base.
			BorderTop(false).
			BorderRight(false).
			BorderBottom(false).
			BorderLeft(false)
		huhTheme.Blurred.Base = huhTheme.Blurred.Base.
			BorderTop(false).
			BorderRight(false).
			BorderBottom(false).
			BorderLeft(false)

		// Increase spacing between select options using margin
		huhTheme.Focused.Option = huhTheme.Focused.Option.Margin(2, 0)
		huhTheme.Blurred.Option = huhTheme.Blurred.Option.Margin(2, 0)

		return huhTheme
	}
}
