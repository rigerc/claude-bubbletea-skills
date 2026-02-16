// Package screens provides layout components for TUI screens.
package screens

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// Layout provides a three-part layout with sticky header, content area, and footer.
// The header stays at the top, footer stays at the bottom, and content fills the middle.
type Layout struct {
	// Header is the sticky header text.
	Header string

	// HeaderStyle is the style applied to the header.
	HeaderStyle lipgloss.Style

	// Footer is the footer text (typically help/keys).
	Footer string

	// FooterStyle is the style applied to the footer.
	FooterStyle lipgloss.Style

	// Content is the main content area.
	Content string

	// BorderStyle is an optional border around the content area.
	// If this is an empty style string, no border is applied.
	BorderStyle lipgloss.Style

	// UseBorder indicates whether to apply the border style.
	UseBorder bool

	// width and height are the terminal dimensions.
	width  int
	height int
}

// NewLayout creates a new layout with the given dimensions.
func NewLayout(width, height int) Layout {
	return Layout{
		width:  width,
		height: height,
	}
}

// SetHeader sets the header text and style.
func (l *Layout) SetHeader(text string, style lipgloss.Style) {
	l.Header = text
	l.HeaderStyle = style
}

// SetFooter sets the footer text and style.
func (l *Layout) SetFooter(text string, style lipgloss.Style) {
	l.Footer = text
	l.FooterStyle = style
}

// SetContent sets the content text.
func (l *Layout) SetContent(content string) {
	l.Content = content
}

// SetBorder sets the border style for the content area.
func (l *Layout) SetBorder(style lipgloss.Style) {
	l.BorderStyle = style
	l.UseBorder = true
}

// ClearBorder removes the border from the content area.
func (l *Layout) ClearBorder() {
	l.UseBorder = false
}

// SetDimensions updates the terminal dimensions.
func (l *Layout) SetDimensions(width, height int) {
	l.width = width
	l.height = height
}

// Render renders the full layout with header, content, and footer positioned correctly.
// Returns the complete rendered string ready for display.
func (l Layout) Render() string {
	// Render styled components
	var header, footer, content string

	if l.Header != "" {
		header = l.HeaderStyle.Render(l.Header)
	}

	if l.Footer != "" {
		footer = l.FooterStyle.Render(l.Footer)
	}

	if l.Content != "" {
		content = l.Content
		if l.UseBorder {
			content = l.BorderStyle.Render(content)
		}
	}

	// Build the full layout
	var result strings.Builder

	// Add header at top
	if header != "" {
		result.WriteString(header)
		result.WriteByte('\n')
	}

	// Add content in middle
	if content != "" {
		result.WriteString(content)
	}

	// Add footer at bottom
	if footer != "" {
		if content != "" {
			result.WriteByte('\n')
		}
		result.WriteString(footer)
	}

	return result.String()
}

// RenderWithPlace uses lipgloss.Place to position header, content, and footer.
// This provides better positioning control for alt screen mode.
func (l Layout) RenderWithPlace() string {
	// Build the content section first
	content := l.Content
	if l.UseBorder {
		content = l.BorderStyle.Render(content)
	}

	// Calculate heights
	headerHeight := 1
	footerHeight := 1
	contentHeight := l.height - headerHeight - footerHeight
	if contentHeight < 1 {
		contentHeight = 1
	}

	// Position header at top
	header := ""
	if l.Header != "" {
		renderedHeader := l.HeaderStyle.Render(l.Header)
		header = lipgloss.Place(l.width, headerHeight,
			lipgloss.Left, lipgloss.Top, renderedHeader)
	}

	// Position footer at bottom
	footer := ""
	if l.Footer != "" {
		renderedFooter := l.FooterStyle.Render(l.Footer)
		footer = lipgloss.Place(l.width, footerHeight,
			lipgloss.Left, lipgloss.Bottom, renderedFooter)
	}

	// Position content in middle
	middleContent := lipgloss.Place(l.width, contentHeight,
		lipgloss.Left, lipgloss.Top, content)

	// Stack them vertically
	parts := []string{header, middleContent, footer}
	return strings.Join(parts, "\n")
}

// SimpleLayout is a simpler version that doesn't use complex positioning.
// It's useful for non-alt-screen mode.
func SimpleLayout(header, content, footer string, headerStyle, footerStyle lipgloss.Style) string {
	var result strings.Builder

	if header != "" {
		result.WriteString(headerStyle.Render(header))
		result.WriteString("\n\n")
	}

	if content != "" {
		result.WriteString(content)
		result.WriteString("\n")
	}

	if footer != "" {
		result.WriteString("\n")
		result.WriteString(footerStyle.Render(footer))
	}

	return result.String()
}
