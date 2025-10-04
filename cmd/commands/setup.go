// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"github.com/kemadev/kemutil/pkg/setup"
	"github.com/spf13/cobra"
)

func init() {
	setupCmd := &cobra.Command{
		Use:    "setup",
		Short:  "Wrapper for Setup tasks",
		Long:   `Setup various things, such as Root CA cert, ...`,
		Args:   cobra.ExactArgs(1),
		PreRun: setLogLevel,
	}
	setupRootCA := &cobra.Command{
		Use:    "root-ca",
		Short:  "Download Root CA certificate",
		Long:   `Download Root CA certificate and place it in ` + setup.ConfigPathBase,
		RunE:   setup.RootCA,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}
	setupMAASMachinesSSHPrivateKey := &cobra.Command{
		Use:    "maas-machines-ssh",
		Short:  "Download MAAS machines SSH private key",
		Long:   `Download MAAS machines SSH private key and place it in ` + setup.ConfigPathBase,
		RunE:   setup.MAASMachinesSSHPrivateKey,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}
	setupSSHConfig := &cobra.Command{
		Use:    "ssh-config",
		Short:  "Create SSH config",
		Long:   `Creates SSH configuration file for alldelared nodes, placing the file in ` + setup.ConfigPathBase,
		RunE:   setup.SSHConfig,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}

	rootCmd.AddCommand(setupCmd)
	setupCmd.AddCommand(setupRootCA)
	setupCmd.AddCommand(setupMAASMachinesSSHPrivateKey)
	setupMAASMachinesSSHPrivateKey.PersistentFlags().
		StringVar(&setup.Region, "region", "", "Region to setup")
	setupMAASMachinesSSHPrivateKey.MarkPersistentFlagRequired("region")
	setupCmd.AddCommand(setupSSHConfig)
	setupSSHConfig.PersistentFlags().
		StringVar(&setup.Region, "region", "", "Region to setup")
	setupSSHConfig.MarkPersistentFlagRequired("region")
}
