// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func setLogLevel(cmd *cobra.Command, _ []string) {
	if cmd.Flag("debug").Value.String() == "true" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Debug mode is enabled, setting log level to debug")
	} else if cmd.Flag("silent").Value.String() == "true" {
		slog.SetLogLoggerLevel(slog.LevelError)
		slog.Debug("Silent mode is enabled, setting log level to error")
	}
}
