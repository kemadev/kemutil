// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals // Cobra root command is global
var rootCmd = &cobra.Command{
	Use:    "kemutil",
	Short:  "Little helpers for easy development",
	Long:   `kemutil is a collection of small utilities to help with development tasks`,
	Args:   cobra.MinimumNArgs(1),
	PreRun: setLogLevel,
}

var (
	// Debug is a flag to enable debug output, actually unused
	//nolint:gochecknoglobals // Cobra flags are global
	debug bool
	// Silent is a flag to enable silent mode, actually unused
	//nolint:gochecknoglobals // Cobra flags are global
	silent bool
)

func init() {
	rootCmd.PersistentFlags().
		BoolVar(&debug, "debug", false, "Enable debug output. Mutually exclusive with --silent, will take precedence if both are set")
	rootCmd.PersistentFlags().
		BoolVar(&silent, "silent", false, "Enable silent mode (less logs). Mutually exclusive with --debug, will be ignored if --debug is set")
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
