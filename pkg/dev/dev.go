// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package dev

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/kemadev/ci-cd/pkg/auth"
	"github.com/kemadev/go-framework/pkg/git"
	"github.com/spf13/cobra"
)

var ErrRepoURLInvalid = errors.New("repository URL is invalid")

var (
	// Debug is a flag to enable debug profile
	//nolint:gochecknoglobals // Cobra flags are global
	Debug bool
	// Live is a flag to enable hot reload.
	//nolint:gochecknoglobals // Cobra flags are global
	Live bool
	// ExportNetrc is a flag to export netrc environment variable
	//nolint:gochecknoglobals // Cobra flags are global
	ExportNetrc bool
)

// StartLocal starts the live development server.
func StartLocal(_ *cobra.Command, _ []string) error {
	slog.Info("Starting local development server")

	binary, err := exec.LookPath("docker")
	if err != nil {
		return fmt.Errorf("docker binary not found: %w", err)
	}

	profile := "dev"
	if Debug {
		profile = "debug"
	}

	baseArgs := []string{
		"compose",
		"--profile",
		profile,
		"--file",
		"./tool/dev/docker-compose.yaml",
		"up",
		"--build",
	}

	if ExportNetrc {
		machine, err := git.NewGitService().GetGitBasePath()
		if err != nil {
			return fmt.Errorf("error getting git repository: %w", err)
		}

		machineParts := strings.Split(machine, "/")
		if len(machineParts) < 1 {
			return fmt.Errorf("error parsing git repository URL: %w", ErrRepoURLInvalid)
		}

		machine = machineParts[0]

		ghBinary, err := exec.LookPath("gh")
		if err != nil {
			return fmt.Errorf("error finding gh binary: %w", err)
		}

		ghArgs := []string{
			"auth",
			"token",
		}

		// nosemgrep: gitlab.gosec.G204-1 // exec.LookPath() is used to locate the binary via $PATH, however we run on trusted developer machines
		com := exec.Command(ghBinary, ghArgs...)

		token, err := com.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error getting git token command output: %w", err)
		}

		netrc := `machine ` + machine + `
login git
password ` + string(token) + `
`
		os.Setenv(auth.NetrcEnvVarKey, netrc)
	}

	if Live {
		baseArgs = append(baseArgs, "--watch")
	}

	slog.Debug("Running command", slog.Any("binary", binary), slog.Any("baseArgs", baseArgs))

	// nosemgrep: go.lang.security.audit.dangerous-syscall-exec.dangerous-syscall-exec // exec.LookPath() is used to locate the binary via $PATH, however we run on trusted developer machines
	err = syscall.Exec(binary, append([]string{binary}, baseArgs...), os.Environ())
	if err != nil {
		return fmt.Errorf("error running workflow ci command: %w", err)
	}

	return nil
}

// StopLocal stops the live development server.
func StopLocal(_ *cobra.Command, _ []string) error {
	slog.Info("Shutting down local development server")

	binary, err := exec.LookPath("docker")
	if err != nil {
		return fmt.Errorf("docker binary not found: %w", err)
	}

	profile := "dev"
	if Debug {
		profile = "debug"
	}

	baseArgs := []string{
		"compose",
		"--profile",
		profile,
		"--file",
		"./tool/dev/docker-compose.yaml",
		"down",
	}

	slog.Debug("Running command", slog.Any("binary", binary), slog.Any("baseArgs", baseArgs))

	// nosemgrep: go.lang.security.audit.dangerous-syscall-exec.dangerous-syscall-exec // exec.LookPath() is used to locate the binary via $PATH, however we run on trusted developer machines
	err = syscall.Exec(binary, append([]string{binary}, baseArgs...), os.Environ())
	if err != nil {
		return fmt.Errorf("error running workflow ci command: %w", err)
	}

	return nil
}
