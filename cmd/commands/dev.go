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
		Short:  "Run development tasks",
		Long:   `Run everyday development tasks like starting a local live dev server, ...`,
		Args:   cobra.MinimumNArgs(1),
		PreRun: setLogLevel,
	}

	localUp := &cobra.Command{
		Use:    "up",
		Short:  "Run local development server",
		Long:   `Run a local development server`,
		RunE:   dev.StartLocal,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}

	localDown := &cobra.Command{
		Use:    "down",
		Short:  "Stop local development server",
		Long:   `Stop the local development server`,
		RunE:   dev.StopLocal,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}

	rootCmd.AddCommand(devCmd)
	devCmd.AddCommand(localUp)
	localUp.PersistentFlags().
		BoolVar(&dev.Debug, "debugger", false, "Enable debugger startup")
	localUp.PersistentFlags().
		BoolVar(&dev.Live, "live", false, "Enable hot reload")
	devCmd.AddCommand(localDown)
	localDown.PersistentFlags().
		BoolVar(&dev.Debug, "debugger", false, "Enable debugger startup")
}
