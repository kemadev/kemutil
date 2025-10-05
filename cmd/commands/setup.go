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
	setupMAASMachinesSSHKey := &cobra.Command{
		Use:    "maas-machines-ssh",
		Short:  "Download MAAS machines SSH private key",
		Long:   `Download MAAS machines SSH private key and place it in ` + setup.ConfigPathBase,
		RunE:   setup.MAASMachinesSSHKeys,
		Args:   cobra.NoArgs,
		PreRun: setLogLevel,
	}
	setupMAASControllersSSHKey := &cobra.Command{
		Use:    "maas-controllers-ssh",
		Short:  "Download MAAS controllers SSH private key",
		Long:   `Download MAAS controllers SSH private key and place it in ` + setup.ConfigPathBase,
		RunE:   setup.MAASControllersSSHKeys,
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
	setupCmd.AddCommand(setupMAASMachinesSSHKey)
	setupMAASMachinesSSHKey.PersistentFlags().
		StringVar(&setup.Region, "region", "", "Region to setup")
	setupMAASMachinesSSHKey.MarkPersistentFlagRequired("region")
	setupCmd.AddCommand(setupSSHConfig)
	setupCmd.AddCommand(setupMAASControllersSSHKey)
	setupMAASControllersSSHKey.PersistentFlags().
		StringVar(&setup.Region, "region", "", "Region to setup")
	setupMAASControllersSSHKey.MarkPersistentFlagRequired("region")
	setupCmd.AddCommand(setupSSHConfig)
	setupSSHConfig.PersistentFlags().
		StringVar(&setup.Region, "region", "", "Region to setup")
	setupSSHConfig.MarkPersistentFlagRequired("region")
}
