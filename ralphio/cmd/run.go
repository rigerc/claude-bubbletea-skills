// Package cmd provides the CLI commands for the application.
package cmd

import "github.com/spf13/cobra"

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the Ralph workflow loop",
	Long:  `Run the Ralph autonomous task execution loop with the configured adapter.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Execution happens in main.go after ShouldRunUI() returns true.
		// This subcommand exists to provide an explicit entry point with its
		// own help text and to distinguish "run" from other subcommands.
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
