package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func setLogLevel(cmd *cobra.Command, args []string) {
	if cmd.Flag("debug").Value.String() == "true" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Debug mode is enabled, setting log level to debug")
	} else if cmd.Flag("silent").Value.String() == "true" {
		slog.SetLogLoggerLevel(slog.LevelError)
		slog.Debug("Silent mode is enabled, setting log level to error")
	}
}
