// This example demonstrates compositing layers with hit testing.
//
// It creates a layered UI with a base document, a modal dialog, and a
// floating badge. The compositor handles z-indexing and hit testing.
package main

import (
	"fmt"
	"os"

	"charm.land/lipgloss/v2"
)

func main() {
	// Detect background for adaptive colors
	hasDarkBG := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark := lipgloss.LightDark(hasDarkBG)

	// Define styles
	baseStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lightDark(lipgloss.Color("#874BFD"), lipgloss.Color("#5B3DBF"))).
		Padding(2).
		Width(60).
		Height(20)

	modalStyle := lipgloss.NewStyle().
		Background(lightDark(lipgloss.Color("#FFF"), lipgloss.Color("#333"))).
		Foreground(lightDark(lipgloss.Color("#333"), lipgloss.Color("#FFF"))).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF6AD2")).
		Padding(1, 3)

	badgeStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#FF6AD2")).
		Foreground(lipgloss.Color("#FFF")).
		Bold(true).
		Padding(0, 1)

	// Create content
	baseContent := baseStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Render("Main Document"),
			"",
			"This is the base layer with some content.",
			"It sits at z-index 0 by default.",
			"",
			"Try clicking positions to see hit testing!",
		),
	)

	modalContent := modalStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.NewStyle().Bold(true).Render("Modal Dialog"),
			"",
			"This is a modal on top.",
			"It has z-index 5.",
		),
	)

	badgeContent := badgeStyle.Render("NEW!")

	// Create layers with positioning
	base := lipgloss.NewLayer(baseContent)

	modal := lipgloss.NewLayer(modalContent).
		X(15).
		Y(5).
		Z(5).
		ID("modal")

	badge := lipgloss.NewLayer(badgeContent).
		X(50).
		Y(2).
		Z(10).
		ID("badge")

	// Create compositor
	comp := lipgloss.NewCompositor(base, modal, badge)

	// Demonstrate hit testing
	fmt.Println("Compositor bounds:", comp.Bounds())
	fmt.Println()

	// Test hits at various positions
	testPoints := []struct {
		x, y int
		desc string
	}{
		{5, 5, "base layer"},
		{25, 10, "modal"},
		{55, 4, "badge"},
		{100, 100, "outside all layers"},
	}

	fmt.Println("Hit testing results:")
	for _, pt := range testPoints {
		hit := comp.Hit(pt.x, pt.y)
		if hit.Empty() {
			fmt.Printf("  (%d, %d) - %s: no hit\n", pt.x, pt.y, pt.desc)
		} else {
			fmt.Printf("  (%d, %d) - %s: hit layer '%s' at %v\n",
				pt.x, pt.y, pt.desc, hit.ID(), hit.Bounds())
		}
	}

	fmt.Println()

	// Render and display
	lipgloss.Println(comp.Render())

	// You can also retrieve layers by ID
	if modalLayer := comp.GetLayer("modal"); modalLayer != nil {
		fmt.Printf("\nRetrieved layer 'modal' at position (%d, %d)\n",
			modalLayer.GetX(), modalLayer.GetY())
	}
}
