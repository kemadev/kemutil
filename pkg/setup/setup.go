// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package setup

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	maasExport "github.com/kemadev/infrastructure-components/deploy/infra/10-vars/export"
	rootCAExport "github.com/kemadev/infrastructure-components/deploy/pki/30-root-ca/export"
	"github.com/kemadev/infrastructure-components/pkg/hardware/datacenter"
	"github.com/kemadev/infrastructure-components/pkg/hardware/router"
	"github.com/kemadev/infrastructure-components/pkg/private/constant/contact"
	"github.com/kemadev/infrastructure-components/pkg/private/constant/host"
	"github.com/kemadev/infrastructure-components/pkg/private/hardware/datacenter/datacenters"
	"github.com/spf13/cobra"
)

// Region to setup.
//
//nolint:gochecknoglobals // Cobra flags are global
var Region string

var RootCA1CertFilePath = ConfigPathBase + RootCASubPath + contact.OrganizationName + "-root-ca-1.crt"

const (
	RootCASubPath = "pki/"
	SSHSubPath    = "ssh/"
)

var ConfigPathBase = os.Getenv("HOME") + "/." + contact.OrganizationName + "/"

// RootCA downloads kema's Root Certificate Authority certificate, so it can be trusted by browsers and alike
func RootCA(_ *cobra.Command, _ []string) error {
	slog.Info("Setting up Root CA")

	bin, err := exec.LookPath("pulumi")
	if err != nil {
		return fmt.Errorf("error finding binary: %w", err)
	}

	cmd := exec.Command(
		bin,
		[]string{
			"stack",
			"output",
			"--stack",
			"bpthdt4i/github-com-kemadev-infrastructure-components-deploy-pki-30-root-ca/main",
			"--show-secrets=true",
			"--json",
		}...)

	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	res := map[string]any{}
	err = json.Unmarshal(out, &res)
	if err != nil {
		return fmt.Errorf("error unmarshalling command output: %w", err)
	}

	content, ok := res[rootCAExport.PulumiStackReferenceRootCA1Certificate].(string)
	if !ok {
		return fmt.Errorf(
			"error extracting %q: %w",
			rootCAExport.PulumiStackReferenceRootCA1Certificate,
			err,
		)
	}

	dir := filepath.Dir(RootCA1CertFilePath)
	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	err = os.WriteFile(
		RootCA1CertFilePath,
		[]byte(content),
		os.FileMode(0o644),
	)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

var ErrRegionUnspecified = errors.New("region flag ")

// MAASMachinesSSHKeys downloads MAAS-deployed machines SSH keys
func MAASMachinesSSHKeys(_ *cobra.Command, _ []string) error {
	slog.Info("Setting up MAAS machines SSH keys for region " + Region)

	bin, err := exec.LookPath("pulumi")
	if err != nil {
		return fmt.Errorf("error finding binary: %w", err)
	}

	cmd := exec.Command(
		bin,
		[]string{
			"stack",
			"output",
			"--stack",
			"bpthdt4i/github-com-kemadev-infrastructure-components-deploy-infra-10-vars/" + Region,
			"--show-secrets=true",
			"--json",
		}...)

	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	res := map[string]any{}
	err = json.Unmarshal(out, &res)
	if err != nil {
		return fmt.Errorf("error unmarshalling command output: %w", err)
	}

	for _, refID := range []string{maasExport.PulumiStackReferenceMAASMachinesPrivateKey, maasExport.PulumiStackReferenceMAASMachinesPublicKey} {
		content, ok := res[refID].(string)
		if !ok {
			return fmt.Errorf(
				"error extracting %q: %w",
				refID,
				err,
			)
		}

		filePath := ConfigPathBase + SSHSubPath + "maas-machines-" + Region + ".key"
		if strings.Contains(refID, "Public") {
			filePath += ".pub"
		}

		dir := filepath.Dir(filePath)
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}

		err = os.WriteFile(
			filePath,
			[]byte(content),
			os.FileMode(0o600),
		)
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	return nil
}

// MAASControllersSSHKeys downloads MAAS-deployed controllers SSH keys
func MAASControllersSSHKeys(_ *cobra.Command, _ []string) error {
	slog.Info("Setting up MAAS controllers SSH keys for region " + Region)

	bin, err := exec.LookPath("pulumi")
	if err != nil {
		return fmt.Errorf("error finding binary: %w", err)
	}

	cmd := exec.Command(
		bin,
		[]string{
			"stack",
			"output",
			"--stack",
			"bpthdt4i/github-com-kemadev-infrastructure-components-deploy-infra-10-vars/" + Region,
			"--show-secrets=true",
			"--json",
		}...)

	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	res := map[string]any{}
	err = json.Unmarshal(out, &res)
	if err != nil {
		return fmt.Errorf("error unmarshalling command output: %w", err)
	}

	for _, refID := range []string{maasExport.PulumiStackReferenceMAASControllerPrivateKey, maasExport.PulumiStackReferenceMAASControllerPublicKey} {
		content, ok := res[refID].(string)
		if !ok {
			return fmt.Errorf(
				"error extracting %q: %w",
				refID,
				err,
			)
		}

		filePath := ConfigPathBase + SSHSubPath + "maas-controllers-" + Region + ".key"
		if strings.Contains(refID, "Public") {
			filePath += ".pub"
		}

		dir := filepath.Dir(filePath)
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}

		err = os.WriteFile(
			filePath,
			[]byte(content),
			os.FileMode(0o600),
		)
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	return nil
}

var ErrRouterNameInvalid = errors.New("router name is invalid")

// SSHConfig creates an SSH config file from all delared nodes for given region
func SSHConfig(_ *cobra.Command, _ []string) error {
	slog.Info("Setting up MAAS machines SSH config for region " + Region)

	edgeHostsHostnames := []string{}
	for name := range datacenters.RoutersInRegion(datacenter.Region(Region)) {
		routerPlacementInfo, err := datacenter.ParseRouterHostname(name)
		if err != nil {
			return fmt.Errorf(
				"error parsing hostname for router %q: %w",
				name,
				err,
			)
		}

		if routerPlacementInfo.Type == string(router.TypeEdge) {
			edgeHostsHostnames = append(edgeHostsHostnames, name)
		}
	}

	confs := []string{}
	for _, dc := range datacenters.InRegion(datacenter.Region(Region)) {
		for routerType, routers := range dc.AllRoutersByType() {
			for name, rtr := range routers {
				conf := fmt.Sprintf(`Host %s
	HostName %s
	Port %d
	IdentityFile %s
	User %s
	RequestTTY yes`, rtr.Name(dc.String(), routerType, name), host.HostNetworkRouter(dc, routerType, name).Host, 22, ConfigPathBase+SSHSubPath+"maas-controllers-"+Region+".key", "ubuntu")
				if !(routerType == router.TypeEdge) {
					conf += "\n    ProxyJump " + edgeHostsHostnames[rand.Intn(len(edgeHostsHostnames))]
				}
				confs = append(confs, conf)
			}
		}
	}

	content := strings.Join(confs, "\n\n")

	filePath := ConfigPathBase + SSHSubPath + "ssh-" + Region + ".conf"

	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	err = os.WriteFile(
		filePath,
		[]byte(content),
		os.FileMode(0o644),
	)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
