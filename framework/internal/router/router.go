package router

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
)

// Screen is the interface all navigable screens implement.
type Screen interface {
	Init() tea.Cmd
	Update(tea.Msg) (Screen, tea.Cmd)
	View(width, height int) string
	KeyMap() help.KeyMap
	Title() string
}

// NavigateMsg tells the root model to push a new screen.
type NavigateMsg struct {
	Screen Screen
}

// BackMsg tells the root model to pop the current screen.
type BackMsg struct{}

// Router manages the current screen and a navigation stack.
type Router struct {
	current Screen
	stack   []Screen
}

// New creates a Router with the given initial screen.
func New(initial Screen) *Router {
	return &Router{current: initial}
}

// Current returns the active screen.
func (r *Router) Current() Screen { return r.current }

// Navigate pushes the current screen onto the stack and sets a new one.
func (r *Router) Navigate(s Screen) tea.Cmd {
	r.stack = append(r.stack, r.current)
	r.current = s
	return s.Init()
}

// Back pops the stack. Returns false if already at root.
func (r *Router) Back() bool {
	if len(r.stack) == 0 {
		return false
	}
	idx := len(r.stack) - 1
	r.current = r.stack[idx]
	r.stack = r.stack[:idx]
	return true
}

// SetCurrent replaces the active screen without modifying the stack.
func (r *Router) SetCurrent(s Screen) { r.current = s }

// Replace switches the current screen and clears the navigation stack.
// Use this for top-level navigation (e.g., nav bar) where you want to
// switch between peer screens rather than drill down.
func (r *Router) Replace(s Screen) tea.Cmd {
	r.current = s
	r.stack = r.stack[:0]
	return s.Init()
}

// Depth returns the stack depth (0 = at root screen).
func (r *Router) Depth() int { return len(r.stack) }
