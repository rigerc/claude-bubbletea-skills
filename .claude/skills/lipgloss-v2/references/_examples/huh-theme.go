// This example demonstrates how to create a custom huh theme using Lip Gloss styles.
//
// It shows how to customize all aspects of a form's appearance including
// colors, borders, typography, and field states.
package main

import (
	"fmt"
	"os"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

func main() {
	// Create a custom theme using Lip Gloss styles
	theme := createCustomTheme()

	// Variables to store form data
	var (
		name     string
		email    string
		plan     string
		features []string
		confirm  bool
	)

	// Create and run the form with our custom theme
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Welcome!").
				Description("Let's set up your account with a custom themed form."),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("What's your name?").
				Placeholder("e.g., Jane Doe").
				Value(&name),
			huh.NewInput().
				Title("Email address").
				Placeholder("you@example.com").
				Value(&email),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a plan").
				Options(
					huh.NewOption("Starter (Free)", "starter"),
					huh.NewOption("Pro ($20/month)", "pro"),
					huh.NewOption("Enterprise (Custom)", "enterprise"),
				).
				Value(&plan),
			huh.NewMultiSelect[string]().
				Title("Select features").
				Options(
					huh.NewOption("Analytics", "analytics"),
					huh.NewOption("API Access", "api"),
					huh.NewOption("Priority Support", "support"),
					huh.NewOption("Custom Domain", "domain"),
				).
				Value(&features),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Ready to continue?").
				Affirmative("Yes, let's go!").
				Negative("No, cancel").
				Value(&confirm),
		),
	).WithTheme(theme)

	if err := form.Run(); err != nil {
		if err == huh.ErrUserAborted {
			fmt.Println("Cancelled.")
			os.Exit(0)
		}
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Success output
	successStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF00")).
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	fmt.Println(successStyle.Render(fmt.Sprintf(
		"Welcome aboard, %s!\nPlan: %s\nFeatures: %v",
		name, plan, features,
	)))
}

// createCustomTheme creates a fully customized huh theme using Lip Gloss styles.
func createCustomTheme() huh.Theme {
	return huh.ThemeFunc(func(isDark bool) *huh.Styles {
		// Base styles
		base := lipgloss.NewStyle().Padding(0, 1)

		// Color palette
		var (
			primary   = lipgloss.Color("#7D56F4") // Purple
			secondary = lipgloss.Color("#FF6AD2") // Pink
			success   = lipgloss.Color("#00FF9F") // Green
			warning   = lipgloss.Color("#FFCC00") // Yellow
			error     = lipgloss.Color("#FF5F87") // Red
			text      = lipgloss.Color("#FFFFFF")
			subtle    = lipgloss.Color("#666666")
		)

		// Create the styles
		s := &huh.Styles{}

		// Form styles
		s.Form = huh.FormStyles{
			Base: base.Copy().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primary).
				Padding(2, 4),
		}

		s.Group = huh.GroupStyles{
			Base: base.Copy().Margin(1, 0),
		}

		// Focused (active) field styles
		s.Focused = huh.FieldStyles{
			Base: base.Copy().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primary).
				Padding(1, 2),
			Title: lipgloss.NewStyle().
				Bold(true).
				Foreground(primary).
				MarginBottom(1),
			Description: lipgloss.NewStyle().
				Foreground(subtle).
				MarginBottom(1),
			ErrorIndicator: lipgloss.NewStyle().
				Foreground(error).
				SetString("✗ "),
			ErrorMessage: lipgloss.NewStyle().
				Foreground(error).
				PaddingLeft(2),
			// Select styles
			SelectSelector: lipgloss.NewStyle().
				Foreground(secondary).
				SetString("▸ "),
			Option: lipgloss.NewStyle().
				PaddingLeft(2),
			NextIndicator: lipgloss.NewStyle().
				Foreground(subtle).
				SetString(" →"),
			PrevIndicator: lipgloss.NewStyle().
				Foreground(subtle).
				SetString("← "),
			// Multi-select styles
			MultiSelectSelector: lipgloss.NewStyle().
				Foreground(secondary).
				SetString("▸ "),
			SelectedOption: lipgloss.NewStyle().
				Foreground(success),
			SelectedPrefix: lipgloss.NewStyle().
				Foreground(success).
				SetString("[✓] "),
			UnselectedOption: lipgloss.NewStyle().
				Foreground(text),
			UnselectedPrefix: lipgloss.NewStyle().
				Foreground(subtle).
				SetString("[ ] "),
			// File picker styles
			Directory: lipgloss.NewStyle().
				Foreground(primary).
				Bold(true),
			File: lipgloss.NewStyle().
				Foreground(text),
			// Confirm button styles
			FocusedButton: lipgloss.NewStyle().
				Background(secondary).
				Foreground(text).
				Padding(0, 3).
				MarginRight(2),
			BlurredButton: lipgloss.NewStyle().
				Background(subtle).
				Foreground(text).
				Padding(0, 3).
				MarginRight(2),
			// Note/Card styles
			Card: lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(primary).
				Padding(1, 2).
				MarginBottom(1),
			NoteTitle: lipgloss.NewStyle().
				Bold(true).
				Foreground(primary),
			Next: lipgloss.NewStyle().
				Foreground(success).
				SetString("Next →"),
			// Text input styles
			TextInput: huh.TextInputStyles{
				Prompt: lipgloss.NewStyle().
					Foreground(primary),
				Placeholder: lipgloss.NewStyle().
					Foreground(subtle),
			},
		}

		// Blurred (inactive) field styles
		s.Blurred = s.Focused
		s.Blurred.Base = s.Blurred.Base.BorderForeground(subtle)
		s.Blurred.Title = s.Blurred.Title.Foreground(subtle)
		s.Blurred.SelectSelector = lipgloss.NewStyle().SetString("  ")
		s.Blurred.MultiSelectSelector = lipgloss.NewStyle().SetString("  ")

		// Help styles
		s.Help = huh.HelpStyles{
			ShortKey: lipgloss.NewStyle().
				Foreground(subtle),
			ShortDesc: lipgloss.NewStyle().
				Foreground(subtle),
			ShortSeparator: lipgloss.NewStyle().
				Foreground(subtle).
				SetString(" • "),
			FullKey: lipgloss.NewStyle().
				Foreground(subtle),
			FullDesc: lipgloss.NewStyle().
				Foreground(subtle),
			FullSeparator: lipgloss.NewStyle().
				Foreground(subtle).
				SetString(" • "),
		}

		return s
	})
}
