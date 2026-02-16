package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "myapp",
		Short: "Application with nested subcommands",
	}

	var echoTimes int

	var cmdEcho = &cobra.Command{
		Use:   "echo [string to echo]",
		Short: "Echo anything to the screen",
		Long: `Echo works like print, except it has a child command
for repeating the output.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Echo: " + args[0])
		},
	}

	var cmdTimes = &cobra.Command{
		Use:   "times [string to echo]",
		Short: "Echo anything multiple times",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < echoTimes; i++ {
				fmt.Println("Echo: " + args[0])
			}
		},
	}
	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo")

	var cmdPrint = &cobra.Command{
		Use:   "print [string to print]",
		Short: "Print anything to the screen",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print: " + args[0])
		},
	}

	rootCmd.AddCommand(cmdEcho, cmdPrint)
	cmdEcho.AddCommand(cmdTimes)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
