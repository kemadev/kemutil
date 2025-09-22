// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"github.com/kemadev/kemutil/pkg/workflow"
	"github.com/spf13/cobra"
)

func init() {
	workflowCmd := &cobra.Command{
		Use:    "workflow",
		Short:  "Run workflows",
		Long:   `Run workflows that usually run in CI/CD pipelines, locally`,
		Args:   cobra.MinimumNArgs(1),
		PreRun: setLogLevel,
	}

	workflowCiCmd := &cobra.Command{
		Use:    "ci",
		Short:  "Run CI workflows",
		Long:   `Run all CI pipelines`,
		RunE:   workflow.Ci,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}

	workflowCustomCmd := &cobra.Command{
		Use:    "custom",
		Short:  "Run custom commands",
		Long:   `Run custom commands using the CI/CD runner`,
		RunE:   workflow.Custom,
		Args:   cobra.MinimumNArgs(1),
		PreRun: setLogLevel,
	}

	rootCmd.AddCommand(workflowCmd)
	workflowCmd.PersistentFlags().
		BoolVar(&workflow.Hot, "hot", false, "Enable hot reload mode")
	workflowCmd.PersistentFlags().
		BoolVar(&workflow.RunnerDebug, "runner-debug", false, "Enable debug mode for the CI/CD runner")
	workflowCmd.PersistentFlags().
		BoolVar(&workflow.ExportNetrc, "netrc", false, "Export netrc")
	workflowCmd.AddCommand(workflowCiCmd)
	workflowCmd.AddCommand(workflowCustomCmd)
	workflowCmd.PersistentFlags().BoolVar(&workflow.Fix, "fix", false, "Enable fix mode")
}
