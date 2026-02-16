// Package screens provides common screen functionality for the UI.
package screens

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	applogger "template-v2-enhanced/internal/logger"
)

// BaseScreen provides common functionality for all screens.
// It handles standard initialization, theme detection, and view wrapping.
// Screens should embed this struct to inherit its behavior.
type BaseScreen struct {
	// AltScreen indicates whether to use alternate screen mode.
	AltScreen bool

	// IsDark indicates if the terminal has a dark background.
	IsDark bool

	// Width and Height store the terminal dimensions.
	Width  int
	Height int

	// LoggerName is the name prefix for log messages (e.g., "HomeScreen").
	// If empty, logs use the default applogger.
	LoggerName string

	// Header is the sticky header text displayed at the top.
	Header string

	// Footer is the footer text displayed at the bottom (typically help/keys).
	Footer string

	// AppTitle is the application title shown in the header.
	AppTitle string
}

// Init provides default initialization behavior for screens.
// It requests background color detection and window size.
// Screens can override this if they need additional initialization.
func (b *BaseScreen) Init() tea.Cmd {
	return tea.Batch(tea.RequestBackgroundColor, tea.RequestWindowSize)
}

// UpdateBackgroundColor handles background color messages.
// It updates IsDark and returns true if the message was handled.
// Screens should call this from their Update method when handling tea.BackgroundColorMsg.
func (b *BaseScreen) UpdateBackgroundColor(msg tea.BackgroundColorMsg) bool {
	b.IsDark = msg.IsDark()
	b.LogDebugf("Background color detected: isDark=%v", b.IsDark)
	return true
}

// UpdateWindowSize handles window resize messages.
// It updates Width and Height and returns true if the message was handled.
// Screens should call this from their Update method when handling tea.WindowSizeMsg.
func (b *BaseScreen) UpdateWindowSize(msg tea.WindowSizeMsg) bool {
	b.Width = msg.Width
	b.Height = msg.Height
	b.LogDebugf("Window resized: width=%d, height=%d", b.Width, b.Height)
	return true
}

// View wraps content with alt screen support.
// It returns a tea.View with AltScreen set based on the BaseScreen configuration.
// Screens should call this from their View method.
func (b *BaseScreen) View(content string) tea.View {
	v := tea.NewView(content)
	v.AltScreen = b.AltScreen
	return v
}

// CenterContent places content in the center if alt screen is active.
// If AltScreen is false or dimensions are not set, content is returned as-is.
// This is useful for fullscreen screens that want centered content.
func (b *BaseScreen) CenterContent(content string) string {
	if b.AltScreen && b.Width > 0 && b.Height > 0 {
		return lipgloss.Place(b.Width, b.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

// LogDebug logs a debug message with the screen's name prefix.
// If LoggerName is set, it formats as "LoggerName: message".
// Otherwise, it logs the message directly.
func (b *BaseScreen) LogDebug(msg string) {
	if b.LoggerName != "" {
		applogger.Debug().Msgf("%s: %s", b.LoggerName, msg)
	} else {
		applogger.Debug().Msg(msg)
	}
}

// LogDebugf logs a formatted debug message with the screen's name prefix.
// If LoggerName is set, it formats as "LoggerName: formatted message".
// Otherwise, it logs the formatted message directly.
func (b *BaseScreen) LogDebugf(format string, args ...interface{}) {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	b.LogDebug(msg)
}

// LogInfo logs an info message with the screen's name prefix.
func (b *BaseScreen) LogInfo(msg string) {
	if b.LoggerName != "" {
		applogger.Info().Msgf("%s: %s", b.LoggerName, msg)
	} else {
		applogger.Info().Msg(msg)
	}
}

// LogError logs an error message with the screen's name prefix.
func (b *BaseScreen) LogError(msg string) {
	if b.LoggerName != "" {
		applogger.Error().Msgf("%s: %s", b.LoggerName, msg)
	} else {
		applogger.Error().Msg(msg)
	}
}

// RenderLayout renders a three-part layout with header, content, and footer.
// This is a convenience method for screens that want the standard layout.
func (b *BaseScreen) RenderLayout(header, content, footer string, headerStyle, footerStyle lipgloss.Style) string {
	return SimpleLayout(header, content, footer, headerStyle, footerStyle)
}

// RenderLayoutWithBorder renders a three-part layout with header, bordered content, and footer.
func (b *BaseScreen) RenderLayoutWithBorder(header, content, footer string, headerStyle, footerStyle, borderStyle lipgloss.Style) string {
	var result strings.Builder

	// Header - combine app title with screen title
	fullHeader := b.buildHeader(header)
	if fullHeader != "" {
		result.WriteString(headerStyle.Render(fullHeader))
		result.WriteByte('\n')
	}

	// Content with border
	if content != "" {
		result.WriteString(borderStyle.Render(content))
	}

	// Footer
	if footer != "" {
		result.WriteByte('\n')
		result.WriteString(footerStyle.Render(footer))
	}

	return result.String()
}

// buildHeader combines the app title with the screen header.
func (b *BaseScreen) buildHeader(screenHeader string) string {
	if b.AppTitle == "" {
		b.AppTitle = "Template-v2"
	}
	if screenHeader == "" {
		return b.AppTitle
	}
	return fmt.Sprintf("%s | %s", b.AppTitle, screenHeader)
}

// GetContentWidth returns the width to use for content based on terminal size.
// For alt screen, uses a larger portion of the terminal width.
func (b *BaseScreen) GetContentWidth() int {
	if b.Width <= 0 {
		return 80 // Default
	}

	if b.AltScreen {
		// Use 80% of terminal width for alt screen, max 120
		width := int(float64(b.Width) * 0.8)
		if width > 120 {
			return 120
		}
		if width < 60 {
			return 60
		}
		return width
	}

	// For regular mode, use smaller width
	width := b.Width - 10
	if width > 100 {
		return 100
	}
	if width < 50 {
		return 50
	}
	return width
}
