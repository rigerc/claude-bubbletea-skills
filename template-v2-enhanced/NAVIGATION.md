# Navigation Framework for BubbleTea v2

This document describes the navigation framework integrated into the template-v2-enhanced project. The framework provides a stack-based navigation system for managing multiple screens in BubbleTea v2 applications.

## Overview

The navigation framework is located in `internal/ui/nav/` and provides:

- **Stack-based navigation**: Push, pop, and replace screens
- **Lifecycle events**: Screens can be notified when they appear or disappear
- **Type-safe screen management**: The `Screen` interface ensures all screens are compatible
- **BubbleTea v2 compatibility**: Fully updated for the v2 API with `tea.View` return types

## Architecture

### Screen Interface

All navigable screens must implement the `Screen` interface:

```go
type Screen interface {
    Init() tea.Cmd
    Update(tea.Msg) (Screen, tea.Cmd)
    View() tea.View
}
```

### LifecycleScreen Interface (Optional)

For screens that need to react to visibility changes:

```go
type LifecycleScreen interface {
    Screen
    Appeared() tea.Cmd  // Called when screen becomes active
    Disappeared()       // Called when screen loses active status
}
```

### Navigation Stack

The `Stack` type manages an ordered collection of screens:

```go
type Stack struct {
    screens     []Screen
    pendingOps  []tea.Msg
    inLifecycle bool
}
```

## Usage

### Basic Navigation

```go
import (
    "template-v2-enhanced/internal/ui/nav"
    "template-v2-enhanced/internal/ui/screens"
)

// Create a navigation stack with a root screen
rootScreen := screens.NewHomeScreen()
stack := nav.NewStack(rootScreen)

// Push a new screen onto the stack
return nav.Push(screens.NewDetailsScreen())

// Pop the top screen from the stack
return nav.Pop()

// Replace the top screen with a new one
return nav.Replace(screens.NewSettingsScreen())
```

### Creating a Custom Screen

```go
package screens

import (
    tea "charm.land/bubbletea/v2"
    "charm.land/lipgloss/v2"
    "template-v2-enhanced/internal/ui/nav"
)

// MyScreen is a custom screen implementation
type MyScreen struct {
    // Screen state
    isDark bool
    width  int
    height int
}

// NewMyScreen creates a new instance
func NewMyScreen() *MyScreen {
    return &MyScreen{}
}

// Init initializes the screen
func (m *MyScreen) Init() tea.Cmd {
    return nil
}

// Update handles incoming messages
func (m *MyScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.BackgroundColorMsg:
        m.isDark = msg.IsDark()
        return m, nil

    case tea.KeyPressMsg:
        switch msg.String() {
        case "q":
            return m, tea.Quit
        case "b":
            return m, nav.Pop()
        }
    }
    return m, nil
}

// View renders the screen
func (m *MyScreen) View() tea.View {
    v := tea.NewView("My Screen Content")
    v.AltScreen = true
    return v
}
```

### Screen with Lifecycle Events

```go
// LifecycleScreen implementation
func (m *MyScreen) Appeared() tea.Cmd {
    // Called when screen becomes active
    // Return a command for async initialization
    return nil
}

func (m *MyScreen) Disappeared() {
    // Called when screen loses active status
    // Perform cleanup here
}
```

## Key Concepts

### Navigation Messages

The framework uses three navigation message types:

- **PushMsg**: Pushes a new screen onto the stack
- **PopMsg**: Pops the top screen from the stack
- **ReplaceMsg**: Replaces the top screen with a new one

### Screen Transitions

When navigation occurs:

1. **Push**: Old top screen disappears, new screen appears
2. **Pop**: Popped screen disappears, revealed screen appears
3. **Replace**: Old screen disappears, new screen appears

### Protected Operations

Navigation operations during lifecycle callbacks are queued and processed after the current transition completes. This prevents re-entrancy issues.

## Example Screens

The template includes four example screens:

1. **HomeScreen** (`screens/home.go`): Main menu with navigation options
2. **DetailsScreen** (`screens/details.go`): Detail view with back button
3. **SettingsScreen** (`screens/settings.go`): Settings with toggle options
4. **AboutScreen** (`screens/about.go`): Information screen

## Integration with Main Model

The main `Model` in `internal/ui/model.go` wraps the navigation stack:

```go
type Model struct {
    navStack nav.Stack
    // ... other fields
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Forward messages to navigation stack
    var updatedModel tea.Model
    var cmd tea.Cmd
    updatedModel, cmd = m.navStack.Update(msg)
    m.navStack = updatedModel.(nav.Stack)
    return m, cmd
}

func (m Model) View() tea.View {
    // Get view from active screen
    return m.navStack.View()
}
```

## Best Practices

1. **Always return `nav.Screen`** from `Update()`, not `tea.Model`
2. **Use `tea.View`** return type in `View()` for BubbleTea v2
3. **Handle background color** in each screen for adaptive theming
4. **Implement `LifecycleScreen`** for screens with async initialization
5. **Keep screens focused** on their specific responsibility
6. **Use `nav.Pop()`** to go back, not `tea.Quit` (except for root screen)

## Key Bindings

Common key bindings used across screens:

- **q**: Quit application
- **b** / **esc**: Go back to previous screen
- **↑/k**: Move up in lists
- **↓/j**: Move down in lists
- **enter**: Select/toggle option

## File Structure

```
internal/ui/
├── nav/
│   ├── nav.go          # Core navigation stack implementation
│   └── lifecycle.go    # LifecycleScreen interface
├── screens/
│   ├── home.go         # Home screen
│   ├── details.go      # Details screen
│   ├── settings.go     # Settings screen
│   └── about.go        # About screen
├── model.go            # Main model wrapping navigation stack
└── help.go             # Key bindings and help
```

## API Reference

### Stack Methods

```go
// Create a new stack with root screen
func NewStack(root Screen) Stack

// Initialize the stack
func (s Stack) Init() tea.Cmd

// Update handles messages
func (s Stack) Update(msg tea.Msg) (tea.Model, tea.Cmd)

// View renders active screen
func (s Stack) View() tea.View

// Depth returns number of screens
func (s Stack) Depth() int
```

### Navigation Commands

```go
// Push a screen onto the stack
func Push(screen Screen) tea.Cmd

// Pop the top screen
func Pop() tea.Cmd

// Replace the top screen
func Replace(screen Screen) tea.Cmd
```

## Migration from v1

If migrating from BubbleTea v1:

1. Change import from `github.com/charmbracelet/bubbletea` to `charm.land/bubbletea/v2`
2. Change `View() string` to `View() tea.View`
3. Use `tea.NewView(content)` to wrap rendered content
4. Use `view.AltScreen = true` instead of `tea.WithAltScreen()` option
5. Handle `tea.KeyPressMsg` instead of `tea.KeyMsg`
6. Use `msg.String()` for key comparison (e.g., "q", "ctrl+c")

## Testing

To test the navigation framework:

```bash
# Build the project
go build -o template-v2-enhanced .

# Run with debug logging
./template-v2-enhanced --debug --log-level trace

# Check debug.log for lifecycle events
tail -f debug.log
```

## License

This navigation framework is part of the template-v2-enhanced project.
