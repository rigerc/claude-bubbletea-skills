package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "myapp",
		Short: "Demonstrates shell completions",
	}

	// Static completions for positional arguments
	var getCmd = &cobra.Command{
		Use:       "get [resource]",
		Short:     "Get a resource",
		ValidArgs: []string{"pod", "service", "deployment", "configmap", "secret"},
		Args:      cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Getting %s\n", args[0])
		},
	}

	// Dynamic completions
	var statusCmd = &cobra.Command{
		Use:   "status [release]",
		Short: "Show release status",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			// Simulate fetching releases from cluster
			releases := []string{"production", "staging", "dev", "test"}
			var completions []cobra.Completion
			for _, r := range releases {
				if strings.HasPrefix(r, toComplete) {
					completions = append(completions, cobra.CompletionWithDesc(r, "Release: "+r))
				}
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				fmt.Printf("Status of release: %s\n", args[0])
			}
		},
	}

	// Flag completions
	var outputFormat string
	statusCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "output format")
	statusCmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		return []cobra.Completion{
			cobra.CompletionWithDesc("json", "JSON output format"),
			cobra.CompletionWithDesc("yaml", "YAML output format"),
			cobra.CompletionWithDesc("text", "Plain text output"),
		}, cobra.ShellCompDirectiveNoFileComp
	})

	// File extension filtering
	var templateCmd = &cobra.Command{
		Use:   "template [file]",
		Short: "Process a template file",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				fmt.Printf("Processing template: %s\n", args[0])
			}
		},
	}
	templateCmd.Flags().String("input", "", "input file")
	templateCmd.MarkFlagFilename("input", "yaml", "yml", "json")

	// Directory filtering
	var cdCmd = &cobra.Command{
		Use:   "cd [directory]",
		Short: "Change directory context",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				fmt.Printf("Changing to: %s\n", args[0])
			}
		},
	}
	cdCmd.Flags().String("path", "", "target path")
	cdCmd.MarkFlagDirname("path")

	rootCmd.AddCommand(getCmd, statusCmd, templateCmd, cdCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
