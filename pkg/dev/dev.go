package dev

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	// Debug is a flag to enable debug profile
	// nolint:gochecknoglobals // Cobra flags are global
	Debug bool
	// Live is a flag to enable hot reload
	Live bool
)

// StartLocal starts the live development server.
func StartLocal(cmd *cobra.Command, args []string) error {
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
