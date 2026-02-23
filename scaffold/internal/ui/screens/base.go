package screens

import (
	"strings"

	"charm.land/bubbles/v2/help"
	lipgloss "charm.land/lipgloss/v2"

	appkeys "scaffold/internal/ui/keys"
	"scaffold/internal/ui/layout"
	"scaffold/internal/ui/theme"
)

// Sizing constants for consistent view dimensions across all screens.
const (
	// MaxContentHeight is the maximum height for any screen's main content area.
	// Pass this to Column.BodyMaxHeight() on screens that should not grow
	// beyond a fixed size (e.g. detail / viewport screens on large terminals).
	MaxContentHeight = 25

	// MinContentHeight is the minimum height for any screen's main content area.
	// This ensures UI elements remain usable even on small terminals.
	MinContentHeight = 10
)

// ScreenBase holds state shared by every screen: adaptive theme, terminal
// dimensions, global key bindings, and a help bar component.
// Embed it in your Screen struct and call its helpers to avoid repeating
// layout and theming boilerplate.
type ScreenBase struct {
	Theme   theme.Theme
	IsDark  bool
	Width   int
	Height  int
	Keys    appkeys.GlobalKeyMap
	Help    help.Model
	AppName string // application name shown in every screen's header badge
}

// NewBase initialises a ScreenBase for the given terminal background.
func NewBase(isDark bool, appName string) ScreenBase {
	b := ScreenBase{Keys: appkeys.New(), Help: help.New(), AppName: appName}
	b.ApplyTheme(isDark)
	return b
}

// ApplyTheme regenerates theme and help styles for the given background.
// Call at the start of SetTheme() before any component-specific updates.
func (b *ScreenBase) ApplyTheme(isDark bool) {
	b.IsDark = isDark
	b.Theme = theme.New(isDark)
	b.Help.Styles = help.DefaultStyles(isDark)
}

// Layout returns a Column pre-configured with the screen's current dimensions
// and the theme's App container style. This is the primary way screens should
// build their view:
//
//	func (s *MyScreen) View() string {
//	    return s.Layout().
//	        Header(s.HeaderView()).
//	        Body(myContent).
//	        Help(s.RenderHelp(myKeys)).
//	        Render()
//	}
//
// To calculate body height for sizing viewports or forms during Update():
//
//	bodyH := s.Layout().
//	    Header(s.HeaderView()).
//	    Footer(s.footerView()).
//	    Help(s.RenderHelp(myKeys)).
//	    BodyHeight()
func (b *ScreenBase) Layout() *layout.Column {
	return layout.NewColumn(b.Width, b.Height).
		Container(b.Theme.App)
}

// ContentWidth returns the inner width after the App container's horizontal
// frame (margin + padding). Use this when sizing child components (viewports,
// gutters, decorators) that must fit inside the App container.
func (b *ScreenBase) ContentWidth() int {
	frameH, _ := b.Theme.App.GetFrameSize()
	return b.Width - frameH
}

// IsSized returns true once the screen has received a non-zero WindowSizeMsg.
func (b *ScreenBase) IsSized() bool {
	return b.Width > 0 && b.Height > 0
}

// HeaderView renders the app name badge followed by a horizontal rule that
// fills the remaining content width. Visible on every screen.
// The bottom margin creates consistent spacing between the header and the
// screen content below; all sizing calculations use lipgloss.Height() on
// this output so the space is automatically accounted for.
func (b *ScreenBase) HeaderView() string {
	t := b.Theme.Title.Padding(1, 2).Render(b.AppName)
	lineW := max(0, b.ContentWidth()-lipgloss.Width(t))
	line := b.Theme.Subtle.Render(strings.Repeat("─", lineW))
	return lipgloss.NewStyle().MarginBottom(1).Render(
		lipgloss.JoinHorizontal(lipgloss.Center, t, line),
	)
}

// RenderHelp renders the help bar from any help.KeyMap, with a top margin.
func (b *ScreenBase) RenderHelp(km help.KeyMap) string {
	return lipgloss.NewStyle().MarginTop(1).Render(b.Help.View(km))
}
