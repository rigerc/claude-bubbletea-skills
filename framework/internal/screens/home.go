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
	Settings key.Binding
}

func (k homeKeyMap) ShortHelp() []key.Binding  { return []key.Binding{k.Settings} }
func (k homeKeyMap) FullHelp() [][]key.Binding { return [][]key.Binding{{k.Settings}} }

// HomeScreen is the default landing screen.
type HomeScreen struct {
	keys      homeKeyMap
	colors    theme.ColorSet
	themeName string
}

// NewHome creates a new HomeScreen.
func NewHome() *HomeScreen {
	return &HomeScreen{
		keys: homeKeyMap{
			Settings: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "settings"),
			),
		},
	}
}

func (h *HomeScreen) Init() tea.Cmd { return nil }

func (h *HomeScreen) Update(msg tea.Msg) (router.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if key.Matches(msg, h.keys.Settings) {
			return h, func() tea.Msg {
				return router.NavigateMsg{Screen: NewSettings()}
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

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(h.colors.Primary)

	accentStyle := lipgloss.NewStyle().
		Foreground(h.colors.Secondary)

	subtleStyle := lipgloss.NewStyle().
		Foreground(h.colors.Subtle)

	b.WriteString(titleStyle.Render("Welcome to Framework"))
	b.WriteString("\n\n")
	b.WriteString(subtleStyle.Render("A BubbleTea v2 application framework with"))
	b.WriteString("\n")
	b.WriteString(accentStyle.Render("theming, routing, and Elm architecture."))
	b.WriteString("\n\n")

	if h.themeName != "" {
		b.WriteString(subtleStyle.Render(fmt.Sprintf("Current theme: %s", h.themeName)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(subtleStyle.Render("Press 's' to open settings."))

	return b.String()
}

func (h *HomeScreen) KeyMap() help.KeyMap { return h.keys }
func (h *HomeScreen) Title() string       { return "Home" }
