// Package nav provides a navigation stack framework for BubbleTea v2.
//
// # Overview
//
// The nav package implements a stack-based navigation system for terminal
// user interfaces built with BubbleTea v2. It allows applications to manage
// multiple screens with push, pop, and replace operations.
//
// # Basic Usage
//
//	import "template-v2-enhanced/internal/ui/nav"
//
//	// Create a stack with a root screen
//	stack := nav.NewStack(rootScreen)
//
//	// Push a new screen
//	return nav.Push(newScreen)
//
//	// Pop back to previous screen
//	return nav.Pop()
//
//	// Replace current screen
//	return nav.Push(replacementScreen)
//
// # Screen Interface
//
// All screens must implement the Screen interface:
//
//	type Screen interface {
//	    Init() tea.Cmd
//	    Update(tea.Msg) (Screen, tea.Cmd)
//	    View() tea.View
//	}
//
// # Lifecycle Events
//
// Screens can optionally implement LifecycleScreen to receive visibility
// notifications:
//
//	type LifecycleScreen interface {
//	    Screen
//	    Appeared() tea.Cmd
//	    Disappeared()
//	}
//
// # Navigation Flow
//
//	Stack State            Operation            Result
//	───────────────────────────────────────────────────────────────
//	[A]                   Push B               [A, B]
//	[A, B]                Push C               [A, B, C]
//	[A, B, C]             Pop                  [A, B]
//	[A, B]                Replace D            [A, D]
//	[A, D]                Pop                  [A]
//
// # Message Flow
//
//	1. User input → Model.Update(msg)
//	2. Model forwards to Stack.Update(msg)
//	3. Stack intercepts navigation messages (PushMsg, PopMsg, ReplaceMsg)
//	4. Stack forwards other messages to active screen
//	5. Screen returns (Screen, tea.Cmd)
//	6. Navigation commands modify the stack
//	7. Stack returns (Stack, tea.Cmd)
//
// # Error Handling
//
// The navigation stack does not handle errors internally. Screens should
// handle their own errors and return appropriate commands.
package nav
