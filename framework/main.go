package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"framework/cmd"
	"framework/internal/ui"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if !cmd.ShouldRunUI() {
		return
	}

	m := ui.New()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
}
