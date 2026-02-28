package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the application version.
const Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version and exit",
	PreRun: func(cmd *cobra.Command, args []string) {
		runUI = false
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("framework v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
