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

type settingsKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
}

func (k settingsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select}
}

func (k settingsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Select}}
}

// SettingsScreen allows the user to pick a theme.
type SettingsScreen struct {
	themes  []string
	cursor  int
	current string
	keys    settingsKeyMap
	colors  theme.ColorSet
}

// NewSettings creates a new SettingsScreen.
func NewSettings() *SettingsScreen {
	reg := theme.NewRegistry()
	return &SettingsScreen{
		themes:  reg.Names(),
		current: reg.CurrentName(),
		keys: settingsKeyMap{
			Up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("up/k", "up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("dn/j", "down"),
			),
			Select: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select"),
			),
		},
	}
}

func (s *SettingsScreen) Init() tea.Cmd { return nil }

func (s *SettingsScreen) Update(msg tea.Msg) (router.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, s.keys.Up):
			if s.cursor > 0 {
				s.cursor--
			}
		case key.Matches(msg, s.keys.Down):
			if s.cursor < len(s.themes)-1 {
				s.cursor++
			}
		case key.Matches(msg, s.keys.Select):
			selected := s.themes[s.cursor]
			return s, func() tea.Msg {
				return theme.ThemeSwitchMsg{Name: selected}
			}
		}
	case theme.ThemeChangedMsg:
		s.colors = msg.Colors
		s.current = msg.Name
	}
	return s, nil
}

func (s *SettingsScreen) View(width, height int) string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(s.colors.Primary)

	selectedStyle := lipgloss.NewStyle().
		Foreground(s.colors.Secondary).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(s.colors.Text)

	activeStyle := lipgloss.NewStyle().
		Foreground(s.colors.Subtle)

	b.WriteString(titleStyle.Render("Theme Settings"))
	b.WriteString("\n\n")

	for i, name := range s.themes {
		cursor := "  "
		style := normalStyle
		if i == s.cursor {
			cursor = "> "
			style = selectedStyle
		}

		line := fmt.Sprintf("%s%s", cursor, style.Render(name))
		if name == s.current {
			line += activeStyle.Render(" (active)")
		}
		b.WriteString(line + "\n")
	}

	return b.String()
}

func (s *SettingsScreen) KeyMap() help.KeyMap { return s.keys }
func (s *SettingsScreen) Title() string       { return "Settings" }
