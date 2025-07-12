package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:    "kemutil",
	Short:  "Little helpers for easy development",
	Long:   `kemutil is a collection of small utilities to help with development tasks`,
	Args:   cobra.MinimumNArgs(1),
	PreRun: setLogLevel,
}

// Execute runs the root command, and thus its subcommands.
// It is the entry point for the CLI application.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		slog.Error("Error executing root command", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

var (
	// debug is a flag to enable debug output, actually unused
	debug bool
	// silent is a flag to enable silent mode, actually unused
	silent bool
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")
	rootCmd.PersistentFlags().BoolVar(&silent, "silent", false, "Enable silent mode (less logs)")
}
