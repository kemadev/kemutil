// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	//nolint:gochecknoglobals // Used as a const
	ciImageProdURL = url.URL{
		Host: "ghcr.io",
		Path: "kemadev/ci-cd:latest",
	}
	//nolint:gochecknoglobals // Used as a const
	ciImageDevURL = url.URL{
		Host: "ghcr.io",
		Path: "kemadev/ci-cd:hot-latest",
	}
)

var (
	// Hot is a flag to enable hot reload mode.
	//nolint:gochecknoglobals // Cobra flags are global
	Hot bool
	// Fix is a flag to enable fix mode.
	//nolint:gochecknoglobals // Cobra flags are global
	Fix bool
	// RunnerDebug is a flag to enable debug mode for the CI/CD runner.
	//nolint:gochecknoglobals // Cobra flags are global
	RunnerDebug bool
	//nolint:gochecknoglobals // Used as a const
	dockerArgs = []string{
		"run",
		"--rm",
		"--interactive",
		"--tty",
		"-v",
		".:/src:Z",
	}
)

const GitTokenEnvVarKey string = "GIT_TOKEN"

func getImageURL() url.URL {
	if Hot {
		slog.Debug("Hot reload mode enabled", slog.String("imageUrl", ciImageDevURL.String()))

		return ciImageDevURL
	}

	slog.Debug("Hot reload mode not enabled", slog.String("imageUrl", ciImageDevURL.String()))

	return ciImageProdURL
}

// Ci runs the CI workflows.
func Ci(cmd *cobra.Command, _ []string) error {
	slog.Debug("Running workflow CI")

	imageURL := getImageURL()

	binary, err := exec.LookPath("docker")
	if err != nil {
		return fmt.Errorf("docker binary not found: %w", err)
	}

	baseArgs := dockerArgs

	if RunnerDebug {
		slog.Debug("Debug mode is enabled, adding debug flag to base arguments")

		baseArgs = append(baseArgs, "-e", "RUNNER_DEBUG=1")
	}

	if cmd.Flag("silent").Value.String() == "true" {
		slog.Debug("Silent mode is enabled, adding silent flag to base arguments")

		baseArgs = append(baseArgs, "-e", "RUNNER_SILENT=1")
	}

	baseArgs = append(baseArgs, strings.TrimPrefix(imageURL.String(), "//"))

	baseArgs = append(baseArgs, "ci")
	if Fix {
		baseArgs = append(baseArgs, "--fix")
	}

	slog.Debug("Running command", slog.Any("binary", binary), slog.Any("baseArgs", baseArgs))

	// nosemgrep: go.lang.security.audit.dangerous-syscall-exec.dangerous-syscall-exec // exec.LookPath() is used to locate the binary via $PATH, however we run on trusted developer machines
	err = syscall.Exec(binary, append([]string{binary}, baseArgs...), os.Environ())
	if err != nil {
		return fmt.Errorf("error running workflow ci command: %w", err)
	}

	return nil
}

// Custom runs custom commands using the CI/CD runner.
func Custom(_ *cobra.Command, args []string) error {
	slog.Debug("Running workflow custom")

	imageURL := getImageURL()

	binary, err := exec.LookPath("docker")
	if err != nil {
		return fmt.Errorf("docker binary not found: %w", err)
	}

	baseArgs := dockerArgs

	if RunnerDebug {
		slog.Debug("Debug mode is enabled, adding debug flag to base arguments")

		baseArgs = append(baseArgs, "-e", "RUNNER_DEBUG=1")
	}

	baseArgs = append(baseArgs, strings.TrimPrefix(imageURL.String(), "//"))

	baseArgs = append(baseArgs, args...)

	if Fix {
		slog.Debug("Fix mode is enabled, adding fix flag to base arguments")

		baseArgs = append(baseArgs, "--fix")
	}

	slog.Debug("Running command", slog.Any("binary", binary), slog.Any("baseArgs", baseArgs))

	// nosemgrep: go.lang.security.audit.dangerous-syscall-exec.dangerous-syscall-exec // The purpose of the command is to run a custom command
	err = syscall.Exec(binary, append([]string{binary}, baseArgs...), os.Environ())
	if err != nil {
		return fmt.Errorf("error running workflow custom command: %w", err)
	}

	return nil
}
