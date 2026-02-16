package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	verbose   int
	user      string
	password  string
	outputFmt string
	files     []string
	jsonOut   bool
	yamlOut   bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "myapp",
		Short: "Demonstrates various flag patterns",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Config: %s\n", cfgFile)
			fmt.Printf("Verbose level: %d\n", verbose)
			fmt.Printf("Files: %v\n", files)
		},
	}

	// Persistent flags - available to this command and all subcommands
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.myapp.yaml)")
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "verbosity level (-v, -vv, -vvv)")

	// Local flags - only for this command
	rootCmd.Flags().StringArrayVarP(&files, "file", "f", []string{}, "input files (can be repeated)")

	// Required flag
	rootCmd.Flags().StringVarP(&user, "user", "u", "", "username (required)")
	rootCmd.MarkFlagRequired("user")

	// Flag groups - require together
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "password (required if --user is set)")
	rootCmd.MarkFlagsRequiredTogether("user", "password")

	// Flag groups - mutually exclusive
	rootCmd.Flags().BoolVar(&jsonOut, "json", false, "output in JSON format")
	rootCmd.Flags().BoolVar(&yamlOut, "yaml", false, "output in YAML format")
	rootCmd.MarkFlagsMutuallyExclusive("json", "yaml")

	// Add a subcommand to demonstrate persistent flags
	var subCmd = &cobra.Command{
		Use:   "process",
		Short: "Process something",
		Run: func(cmd *cobra.Command, args []string) {
			// --config and --verbose are available here due to persistent flags
			fmt.Printf("Processing with config: %s, verbose: %d\n", cfgFile, verbose)
		},
	}
	subCmd.Flags().StringVarP(&outputFmt, "output", "o", "text", "output format")

	rootCmd.AddCommand(subCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
