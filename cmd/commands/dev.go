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

	devLocal := &cobra.Command{
		Use:    "local",
		Short:  "Run local development server",
		Long:   `Run a local development server`,
		RunE:   dev.StartLocal,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}

	rootCmd.AddCommand(devCmd)
	devCmd.AddCommand(devLocal)
	devLocal.PersistentFlags().
		BoolVar(&dev.Debug, "debugger", false, "Enable debugger startup")
	devLocal.PersistentFlags().
		BoolVar(&dev.Live, "live", false, "Enable hot reload")
}
