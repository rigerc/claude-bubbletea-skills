// Package layout provides a composable single-column layout for BubbleTea v2
// screens. It handles vertical stacking and automatic body-height calculation
// so that screens only need to provide their rendered section content.
//
// # Layout structure
//
//	┌──────────────────────────────────┐
//	│  Header   (fixed, measured)      │
//	│  Body     (fills remaining)      │
//	│  Footer   (fixed, measured)      │
//	│  Help     (fixed, measured)      │
//	└──────────────────────────────────┘
//
// # Typical usage in View()
//
//	func (s *MyScreen) View() string {
//	    return s.Layout().
//	        Header(s.HeaderView()).
//	        Body(s.vp.View()).
//	        Footer(s.footerView()).
//	        Help(s.RenderHelp(helpKeys)).
//	        Render()
//	}
//
// # Sizing dynamic components in Update()
//
// Call BodyHeight() with the same sections you pass in View() but without
// Body content — the result is the number of rows your body component
// (viewport, form, list, …) should be sized to:
//
//	func (s *MyScreen) updateViewportSize() {
//	    bodyH := s.Layout().
//	        Header(s.HeaderView()).
//	        Footer(s.footerView()).
//	        Help(s.RenderHelp(helpKeys)).
//	        BodyHeight()
//	    s.vp.SetHeight(bodyH)
//	}
package layout

import (
	lipgloss "charm.land/lipgloss/v2"
)

// DefaultBodyMinH is the minimum body height used when no explicit minimum
// has been set via BodyMinHeight().
const DefaultBodyMinH = 5

// Column is a single-column layout builder. It stacks four optional sections
// vertically — Header, Body, Footer, Help — and wraps them in an optional
// container style. The Body section receives all remaining vertical space
// after the other sections and the container frame are measured.
//
// All setter methods return *Column so calls can be chained:
//
//	layout.NewColumn(w, h).Container(th.App).Header(hdr).Body(content).Render()
type Column struct {
	width, height int
	container     lipgloss.Style
	hasContainer  bool

	header string
	body   string
	footer string
	help   string

	bodyMinH int // minimum body height; defaults to DefaultBodyMinH
	bodyMaxH int // maximum body height; 0 means "fill all available"
}

// NewColumn creates a Column for the given terminal dimensions.
//
// width and height should be the inner content area, i.e. already accounting
// for any outer border applied by the root model (borderOverhead in model.go).
func NewColumn(width, height int) *Column {
	return &Column{
		width:    width,
		height:   height,
		bodyMinH: DefaultBodyMinH,
	}
}

// Container sets the outer lipgloss style that wraps the entire rendered
// output. Typically this is theme.App, which provides margin and padding.
// The container's frame dimensions are subtracted automatically when
// calculating the body height.
func (c *Column) Container(s lipgloss.Style) *Column {
	c.container = s
	c.hasContainer = true
	return c
}

// Header sets the header section content (an already-rendered string).
// The header is placed at the top of the column. Its height is measured
// automatically and subtracted from the body's available space.
func (c *Column) Header(s string) *Column {
	c.header = s
	return c
}

// Body sets the body section content (an already-rendered string).
// The body fills all remaining vertical space between the header and
// footer/help sections, subject to any BodyMinHeight / BodyMaxHeight
// constraints.
func (c *Column) Body(s string) *Column {
	c.body = s
	return c
}

// Footer sets the footer section content (an already-rendered string).
// The footer is placed directly below the body. Its height is measured
// automatically.
func (c *Column) Footer(s string) *Column {
	c.footer = s
	return c
}

// Help sets the help-bar section content (an already-rendered string).
// The help bar is placed at the very bottom, below the footer. Its height
// is measured automatically.
func (c *Column) Help(s string) *Column {
	c.help = s
	return c
}

// BodyMinHeight sets the minimum height for the body section in rows.
// Defaults to DefaultBodyMinH. The body will never be rendered shorter
// than this value even if the terminal does not have enough space.
func (c *Column) BodyMinHeight(h int) *Column {
	c.bodyMinH = h
	return c
}

// BodyMaxHeight sets the maximum height for the body section in rows.
// A value of 0 (the default) means no cap — the body fills all available
// space. Use this to prevent overly tall body sections on very large
// terminals (e.g. pass MaxContentHeight from the screens package).
func (c *Column) BodyMaxHeight(h int) *Column {
	c.bodyMaxH = h
	return c
}

// BodyHeight returns the number of rows available for the body section.
//
// It measures all fixed sections (header, footer, help) and the container
// frame, subtracts them from the total height, then clamps the result to
// the configured min/max range.
//
// Call this during Update() — before rendering the body — to size dynamic
// components such as viewports and forms. You do not need to call Body()
// before calling BodyHeight().
func (c *Column) BodyHeight() int {
	if c.height <= 0 {
		return c.bodyMinH
	}
	available := c.height - c.frameHeight() - c.fixedHeight()
	return c.clampBodyH(available)
}

// ContentWidth returns the usable inner width after subtracting the
// container's horizontal frame (margin + padding + border). This is the
// width available for content inside each section.
func (c *Column) ContentWidth() int {
	if !c.hasContainer {
		return c.width
	}
	frameH, _ := c.container.GetFrameSize()
	return max(0, c.width-frameH)
}

// Render stacks all non-empty sections vertically in the order
// Header → Body → Footer → Help and wraps the result in the container
// style. Empty sections are omitted so callers can freely leave out any
// section without producing blank lines.
//
// The body content is constrained to BodyHeight() rows using lipgloss
// Height / MaxHeight so it never overflows the terminal.
func (c *Column) Render() string {
	bodyH := c.BodyHeight()

	var sections []string

	if c.header != "" {
		sections = append(sections, c.header)
	}

	if c.body != "" {
		constrained := lipgloss.NewStyle().
			Height(bodyH).
			MaxHeight(bodyH).
			Render(c.body)
		sections = append(sections, constrained)
	}

	if c.footer != "" {
		sections = append(sections, c.footer)
	}

	if c.help != "" {
		sections = append(sections, c.help)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	if c.hasContainer {
		return c.container.Render(content)
	}
	return content
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// frameHeight returns the vertical space consumed by the container style
// (top + bottom margin + border + padding).
func (c *Column) frameHeight() int {
	if !c.hasContainer {
		return 0
	}
	_, frameV := c.container.GetFrameSize()
	return frameV
}

// fixedHeight returns the total measured height of all fixed sections
// (header + footer + help). The body section is excluded because its
// height is what we are computing.
func (c *Column) fixedHeight() int {
	h := 0
	if c.header != "" {
		h += lipgloss.Height(c.header)
	}
	if c.footer != "" {
		h += lipgloss.Height(c.footer)
	}
	if c.help != "" {
		h += lipgloss.Height(c.help)
	}
	return h
}

// clampBodyH applies bodyMinH / bodyMaxH constraints to the computed height.
func (c *Column) clampBodyH(h int) int {
	if c.bodyMaxH > 0 && h > c.bodyMaxH {
		h = c.bodyMaxH
	}
	if h < c.bodyMinH {
		h = c.bodyMinH
	}
	return h
}
