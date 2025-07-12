// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"github.com/kemadev/kemutil/pkg/repotpl"
	"github.com/spf13/cobra"
)

func init() {
	repotplCmd := &cobra.Command{
		Use:    "repotpl",
		Short:  "Repository template helpers",
		Long:   `Run tasks related to repository templates, such as initializing or updating`,
		Args:   cobra.ExactArgs(1),
		PreRun: setLogLevel,
	}

	repotplInit := &cobra.Command{
		Use:    "init",
		Short:  "Initialize repository template",
		Long:   `Initialize the repository template from upstream`,
		RunE:   repotpl.Init,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}

	repotplUpdate := &cobra.Command{
		Use:    "update",
		Short:  "Update repository template",
		Long:   `Update the repository template with the latest changes from upstream`,
		RunE:   repotpl.Update,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}

	rootCmd.AddCommand(repotplCmd)
	repotplCmd.AddCommand(repotplInit)
	repotplCmd.AddCommand(repotplUpdate)
	repotplUpdate.Flags().
		BoolVar(&repotpl.SkipAnswered, "skip-answered", false, "Skip answered questions update")
}
