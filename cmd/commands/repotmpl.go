// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"github.com/kemadev/kemutil/pkg/repotmpl"
	"github.com/spf13/cobra"
)

func init() {
	repotmplCmd := &cobra.Command{
		Use:    "repotmpl",
		Short:  "Wrapper for Repotmpl tasks",
		Long:   `Run everyday Repotmpl tasks like initializing a Repotmpl stack, deploying it, ...`,
		Args:   cobra.ExactArgs(1),
		PreRun: setLogLevel,
	}
	repotmplInit := &cobra.Command{
		Use:    "init",
		Short:  "Initialize a new repository",
		Long:   `Initialize a new repository in the current directory, using a template`,
		RunE:   repotmpl.Init,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}

	rootCmd.AddCommand(repotmplCmd)
	repotmplCmd.AddCommand(repotmplInit)
	repotmplInit.PersistentFlags().
		StringVar(&repotmpl.TargetTemplate, "template", "go-app", "Template to use")
}
