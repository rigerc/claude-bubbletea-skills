package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "myapp",
		Short: "Demonstrates PreRun/PostRun hooks",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("[Root PersistentPreRun] Global setup (runs for all commands)")
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("[Root PreRun] Root-specific setup")
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("[Root Run] Main root command logic")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("[Root PostRun] Root-specific cleanup")
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("[Root PersistentPostRun] Global cleanup (runs for all commands)")
		},
	}

	var subCmd = &cobra.Command{
		Use:   "sub",
		Short: "A subcommand",
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("[Sub PreRun] Sub-specific setup")
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("[Sub Run] Sub command logic")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("[Sub PostRun] Sub-specific cleanup")
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("[Sub PersistentPostRun] Sub's persistent cleanup")
		},
	}

	rootCmd.AddCommand(subCmd)

	// Run with sub command to see hook execution order:
	// go run hooks.go sub
	// Output:
	// [Root PersistentPreRun] Global setup (runs for all commands)
	// [Sub PreRun] Sub-specific setup
	// [Sub Run] Sub command logic
	// [Sub PostRun] Sub-specific cleanup
	// [Sub PersistentPostRun] Sub's persistent cleanup

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
