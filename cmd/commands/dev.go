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

	devLocalUp := &cobra.Command{
		Use:    "local",
		Short:  "Run local development server",
		Long:   `Run a local development server`,
		RunE:   dev.StartLocal,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}

	devLocalDown := &cobra.Command{
		Use:    "down",
		Short:  "Stop local development server",
		Long:   `Stop the local development server`,
		RunE:   dev.StopLocal,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}

	rootCmd.AddCommand(devCmd)
	devCmd.AddCommand(devLocalUp)
	devLocalUp.PersistentFlags().
		BoolVar(&dev.Debug, "debugger", false, "Enable debugger startup")
	devLocalUp.PersistentFlags().
		BoolVar(&dev.Live, "live", false, "Enable hot reload")
	devCmd.AddCommand(devLocalDown)
	devLocalDown.PersistentFlags().
		BoolVar(&dev.Debug, "debugger", false, "Enable debugger startup")
}
