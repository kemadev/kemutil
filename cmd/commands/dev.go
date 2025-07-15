// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"github.com/kemadev/kemutil/pkg/dev"
	"github.com/spf13/cobra"
)

func init() {
	devCmd := &cobra.Command{
		Use:    "dev",
		Short:  "Run developement tasks",
		Long:   `Run everyday development tasks like starting a local live dev server, ...`,
		Args:   cobra.MinimumNArgs(1),
		PreRun: setLogLevel,
	}

	devLive := &cobra.Command{
		Use:    "live",
		Short:  "Run live development server",
		Long:   `Run a live development server that watches for changes and reloads the application`,
		RunE:   dev.StartLive,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}

	rootCmd.AddCommand(devCmd)
	devCmd.AddCommand(devLive)
	devLive.PersistentFlags().
		BoolVar(&dev.DebugEnabled, "debugger", false, "Enable debugger start for the live development server")
}
