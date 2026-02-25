// Package theme provides styling for the TUI.
package theme

import (
	"image/color"

	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

// Palette defines semantic colors for the application theme.
type Palette struct {
	// Brand
	Accent      color.Color // Charple - primary brand
	AccentHover color.Color // Hazy - hover state

	// Foreground (adaptive)
	Foreground color.Color // Primary text
	Muted      color.Color // Secondary text
	Subtle     color.Color // Borders, subtle elements

	// Status (always visible)
	Success color.Color // Julep
	Error   color.Color // Sriracha
	Warning color.Color // Tang
	Info    color.Color // Thunder

	// Special
	Inverse color.Color // Text on accent backgrounds (Butter)
}

// NewPalette creates a semantic color palette based on the background.
func NewPalette(isDark bool) Palette {
	ld := lipgloss.LightDark(isDark)

	return Palette{
		Accent:      charmtone.Charple,
		AccentHover: charmtone.Hazy,
		Foreground:  ld(charmtone.Pepper, charmtone.Salt),
		Muted:       ld(charmtone.Charcoal, charmtone.Ash),
		Subtle:      ld(charmtone.Squid, charmtone.Oyster),
		Success:     charmtone.Julep,
		Error:       charmtone.Sriracha,
		Warning:     charmtone.Tang,
		Info:        charmtone.Thunder,
		Inverse:     charmtone.Butter,
	}
}

// Styles holds all styled components for the UI.
type Styles struct {
	App         lipgloss.Style
	Header      lipgloss.Style
	Body        lipgloss.Style
	Help        lipgloss.Style
	Footer      lipgloss.Style
	StatusLeft  lipgloss.Style
	StatusRight lipgloss.Style
	MaxWidth    int
}

// newStylesFromPalette creates Styles from a Palette.
func newStylesFromPalette(p Palette, width int) Styles {
	maxWidth := width * 50 / 100
	if maxWidth < 40 {
		maxWidth = width - 4
	}

	return Styles{
		MaxWidth: maxWidth,
		App:      lipgloss.NewStyle().Width(maxWidth).Padding(0, 0),
		Header:   lipgloss.NewStyle().Padding(5).PaddingBottom(1),
		Body:     lipgloss.NewStyle().Padding(0, 3).Foreground(p.Foreground),
		Help:     lipgloss.NewStyle().MarginTop(0).Padding(0, 3),
		Footer: lipgloss.NewStyle().
			MarginTop(1).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(p.Subtle).
			PaddingLeft(1),
		StatusLeft: lipgloss.NewStyle().
			Background(p.Accent).
			Foreground(p.Inverse).
			Bold(true),
		StatusRight: lipgloss.NewStyle().Foreground(p.Subtle),
	}
}

// New creates Styles with adaptive colors. Backward compatible.
func New(isDark bool, width int) Styles {
	return newStylesFromPalette(NewPalette(isDark), width)
}

// DetailStyles holds styles for the detail screen.
type DetailStyles struct {
	Title   lipgloss.Style
	Desc    lipgloss.Style
	Content lipgloss.Style
	Info    lipgloss.Style
}

// newDetailStylesFromPalette creates DetailStyles from a Palette.
func newDetailStylesFromPalette(p Palette) DetailStyles {
	return DetailStyles{
		Title:   lipgloss.NewStyle().Bold(true).Foreground(p.Accent).MarginBottom(1),
		Desc:    lipgloss.NewStyle().Foreground(p.Subtle).MarginBottom(2),
		Content: lipgloss.NewStyle().Foreground(p.Foreground),
		Info:    lipgloss.NewStyle().Foreground(p.Muted).Italic(true),
	}
}

// NewDetailStyles creates detail styles with adaptive colors. Backward compatible.
func NewDetailStyles(isDark bool) DetailStyles {
	return newDetailStylesFromPalette(NewPalette(isDark))
}

// StatusStyles provides pre-built styles for status messages.
type StatusStyles struct {
	Success lipgloss.Style
	Error   lipgloss.Style
	Warning lipgloss.Style
	Info    lipgloss.Style
}

// NewStatusStyles creates status styles from a Palette.
func NewStatusStyles(isDark bool) StatusStyles {
	p := NewPalette(isDark)
	return StatusStyles{
		Success: lipgloss.NewStyle().Foreground(p.Success).Bold(true),
		Error:   lipgloss.NewStyle().Foreground(p.Error).Bold(true),
		Warning: lipgloss.NewStyle().Foreground(p.Warning),
		Info:    lipgloss.NewStyle().Foreground(p.Info),
	}
}

// ListStyles creates list.Styles from a Palette.
func ListStyles(p Palette) list.Styles {
	s := list.DefaultStyles(false)

	s.Title = lipgloss.NewStyle().
		Foreground(p.Accent).
		Bold(true).
		Padding(0, 1)
	s.TitleBar = lipgloss.NewStyle()
	s.Spinner = lipgloss.NewStyle().Foreground(p.Accent)
	s.PaginationStyle = lipgloss.NewStyle().Foreground(p.Subtle)
	s.HelpStyle = lipgloss.NewStyle().Foreground(p.Muted)
	s.StatusBar = lipgloss.NewStyle().Foreground(p.Muted)
	s.StatusEmpty = lipgloss.NewStyle().Foreground(p.Muted)
	s.NoItems = lipgloss.NewStyle().Foreground(p.Muted)
	s.ActivePaginationDot = lipgloss.NewStyle().Foreground(p.Accent)
	s.InactivePaginationDot = lipgloss.NewStyle().Foreground(p.Subtle)
	s.DividerDot = lipgloss.NewStyle().Foreground(p.Subtle)

	return s
}

// ListItemStyles creates list.DefaultItemStyles from a Palette.
func ListItemStyles(p Palette) list.DefaultItemStyles {
	s := list.NewDefaultItemStyles(false)

	// Normal state
	s.NormalTitle = lipgloss.NewStyle().Foreground(p.Foreground)
	s.NormalDesc = lipgloss.NewStyle().Foreground(p.Muted)

	// Selected state
	s.SelectedTitle = lipgloss.NewStyle().
		Foreground(p.Accent).
		Bold(true)
	s.SelectedDesc = lipgloss.NewStyle().Foreground(p.Subtle)

	// Dimmed state
	s.DimmedTitle = lipgloss.NewStyle().Foreground(p.Muted)
	s.DimmedDesc = lipgloss.NewStyle().Foreground(p.Subtle)

	// Filter match
	s.FilterMatch = lipgloss.NewStyle().Foreground(p.Accent)

	return s
}
