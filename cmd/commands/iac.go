// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"github.com/kemadev/kemutil/pkg/iac"
	"github.com/spf13/cobra"
)

func init() {
	iacCmd := &cobra.Command{
		Use:    "iac",
		Short:  "Wrapper for IaC tasks",
		Long:   `Run everyday IaC tasks like initializing a Iac stack, deploying it, ...`,
		Args:   cobra.ExactArgs(1),
		PreRun: setLogLevel,
	}
	iacInit := &cobra.Command{
		Use:    "init",
		Short:  "Initialize a Pulumi stack",
		Long:   `Initialize a Pulumi stack in the current directory, using a template`,
		RunE:   iac.Init,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}

	rootCmd.AddCommand(iacCmd)
	iacCmd.AddCommand(iacInit)
	iacCmd.PersistentFlags().
		BoolVar(&iac.Debug, "iac-debug", false, "Enable debug output for IaC commands")
	iacCmd.PersistentFlags().
		BoolVar(&iac.Refresh, "refresh", false, "Refresh the IaC stack before updating")
}
