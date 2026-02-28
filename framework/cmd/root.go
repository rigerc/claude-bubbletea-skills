package cmd

import "github.com/spf13/cobra"

var runUI = true

var rootCmd = &cobra.Command{
	Use:   "framework",
	Short: "A BubbleTea v2 framework application",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil // fall through to TUI in main()
	},
}

// Execute runs the root command.
func Execute() error { return rootCmd.Execute() }

// ShouldRunUI returns true if the TUI should start after Cobra exits.
func ShouldRunUI() bool { return runUI }
