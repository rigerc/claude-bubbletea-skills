package screens

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"framework/internal/router"
	"framework/internal/theme"
)

type homeKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
}

func (k homeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select}
}

func (k homeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Select}}
}

// menuItem represents a navigation destination.
type menuItem struct {
	label       string
	description string
	screen      func() router.Screen
}

// HomeScreen is the main menu screen.
type HomeScreen struct {
	keys       homeKeyMap
	colors     theme.ColorSet
	themeName  string
	menuItems  []menuItem
	cursor     int
}

// NewHome creates a new HomeScreen with navigation menu.
func NewHome() *HomeScreen {
	return &HomeScreen{
		keys: homeKeyMap{
			Up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "down"),
			),
			Select: key.NewBinding(
				key.WithKeys("enter", "l", "right"),
				key.WithHelp("enter", "select"),
			),
		},
		menuItems: []menuItem{
			{label: "Settings", description: "Change theme and preferences", screen: func() router.Screen { return NewSettings() }},
		},
	}
}

func (h *HomeScreen) Init() tea.Cmd { return nil }

func (h *HomeScreen) Update(msg tea.Msg) (router.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, h.keys.Up):
			if h.cursor > 0 {
				h.cursor--
			}
		case key.Matches(msg, h.keys.Down):
			if h.cursor < len(h.menuItems)-1 {
				h.cursor++
			}
		case key.Matches(msg, h.keys.Select):
			if h.cursor < len(h.menuItems) {
				return h, func() tea.Msg {
					return router.NavigateMsg{Screen: h.menuItems[h.cursor].screen()}
				}
			}
		}
	case theme.ThemeChangedMsg:
		h.colors = msg.Colors
		h.themeName = msg.Name
	}
	return h, nil
}

func (h *HomeScreen) View(width, height int) string {
	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(h.colors.Primary).
		MarginBottom(1)

	b.WriteString(titleStyle.Render("Framework"))
	b.WriteString("\n\n")

	// Menu items
	for i, item := range h.menuItems {
		cursor := "  "
		labelStyle := lipgloss.NewStyle().Foreground(h.colors.Text)
		descStyle := lipgloss.NewStyle().Foreground(h.colors.Subtle)

		if i == h.cursor {
			cursor = "> "
			labelStyle = lipgloss.NewStyle().Foreground(h.colors.Primary).Bold(true)
			descStyle = lipgloss.NewStyle().Foreground(h.colors.Secondary)
		}

		fmt.Fprintf(&b, "%s%s\n", cursor, labelStyle.Render(item.label))
		fmt.Fprintf(&b, "  %s\n", descStyle.Render(item.description))
		b.WriteString("\n")
	}

	// Theme indicator
	if h.themeName != "" {
		subtleStyle := lipgloss.NewStyle().Foreground(h.colors.Subtle)
		b.WriteString(subtleStyle.Render(fmt.Sprintf("Theme: %s", h.themeName)))
	}

	return b.String()
}

func (h *HomeScreen) KeyMap() help.KeyMap { return h.keys }
func (h *HomeScreen) Title() string       { return "Home" }
