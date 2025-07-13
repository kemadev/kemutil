// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package wgo

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/kemadev/ci-cd/pkg/filesfind"
	"github.com/kemadev/kemutil/internal/app/util"
	"github.com/spf13/cobra"
)

// Init initializes a Go module in the current directory.
func Init(_ *cobra.Command, _ []string) error {
	slog.Info("Initializing Go module")

	modName, err := util.GetGoModExpectedName()
	if err != nil {
		return fmt.Errorf("error getting expected Go module name: %w", err)
	}

	binary, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("go binary not found: %w", err)
	}

	baseArgs := []string{"mod", "init", modName}

	slog.Debug("Running command", slog.Any("binary", binary), slog.Any("baseArgs", baseArgs))

	// nosemgrep: gitlab.gosec.G204-1 // exec.LookPath() is used to locate the binary via $PATH and git repo is variable too, however we run on trusted developer machines
	command := exec.Command(binary, baseArgs...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err = command.Run()
	if err != nil {
		return fmt.Errorf("error running go mod init: %w", err)
	}

	slog.Info("Initialized Go module", slog.String("modName", modName))

	return nil
}

// Update updates all Go modules dependencies found in the current directory and subdirectories.
func Update(_ *cobra.Command, _ []string) error {
	slog.Info("Updating Go modules")

	mods, err := filesfind.FindFilesByExtension(filesfind.FilesFindingArgs{
		Extension: "go.mod",
		Recursive: true,
	})
	if err != nil {
		return fmt.Errorf("error finding go.mod files: %w", err)
	}

	if len(mods) == 0 {
		return nil
	}

	slog.Debug("Found go.mod files", slog.Any("mods", mods))

	binary, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("go binary not found: %w", err)
	}

	baseArgs := []string{"get", "-u", "./..."}

	for _, mod := range mods {
		slog.Debug("Updating Go module", slog.String("mod", mod))

		// nosemgrep: go.lang.security.audit.dangerous-syscall-exec.dangerous-syscall-exec // exec.LookPath() is used to locate the binary via $PATH, however we run on trusted developer machines
		// nosemgrep: gitlab.gosec.G204-1 // Same
		command := exec.Command(binary, baseArgs...)
		command.Dir = path.Dir(mod)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err = command.Run()
		if err != nil {
			return fmt.Errorf("error updating Go module %s: %w", mod, err)
		}

		slog.Info("Updated Go module", slog.String("mod", mod))
	}

	return nil
}

// Tidy tidies all Go modules dependencies found in the current directory and subdirectories.
func Tidy(_ *cobra.Command, _ []string) error {
	slog.Info("Tidying Go modules")

	mods, err := filesfind.FindFilesByExtension(filesfind.FilesFindingArgs{
		Extension: "go.mod",
		Recursive: true,
	})
	if err != nil {
		return fmt.Errorf("error finding go.mod files: %w", err)
	}

	if len(mods) == 0 {
		return nil
	}

	slog.Debug("Found go.mod files", slog.Any("mods", mods))

	binary, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("go binary not found: %w", err)
	}

	baseArgs := []string{"mod", "tidy"}

	for _, mod := range mods {
		slog.Debug("Tidying Go module", slog.String("mod", mod))

		// nosemgrep: go.lang.security.audit.dangerous-syscall-exec.dangerous-syscall-exec // exec.LookPath() is used to locate the binary via $PATH, however we run on trusted developer machines
		// nosemgrep: gitlab.gosec.G204-1 // Same
		command := exec.Command(binary, baseArgs...)
		command.Dir = path.Dir(mod)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err = command.Run()
		if err != nil {
			return fmt.Errorf("error tidying Go module %s: %w", mod, err)
		}

		slog.Info("Tidied Go module", slog.String("mod", mod))
	}

	return nil
}

// UpdateGoVersion updates the Go version in all `go.mod` files found in the current directory and subdirectories to the latest version.
func UpdateGoVersion(_ *cobra.Command, _ []string) error {
	slog.Info("Updating Go version in go.mod files")

	mods, err := filesfind.FindFilesByExtension(filesfind.FilesFindingArgs{
		Extension: "go.mod",
		Recursive: true,
	})
	if err != nil {
		return fmt.Errorf("error finding go.mod files: %w", err)
	}

	if len(mods) == 0 {
		return nil
	}

	slog.Debug("Found go.mod files", slog.Any("mods", mods))

	binary, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("go binary not found: %w", err)
	}

	curlBinary, err := exec.LookPath("curl")
	if err != nil {
		return fmt.Errorf("curl binary not found: %w", err)
	}

	// nosemgrep: go.lang.security.audit.dangerous-syscall-exec.dangerous-syscall-exec // exec.LookPath() is used to locate the binary via $PATH, however we run on trusted developer machines
	// nosemgrep: gitlab.gosec.G204-1 // Same
	latestGoVersionCmd := exec.Command(curlBinary, "-fsSL", "https://go.dev/VERSION?m=text")

	latestGoVersionOutput, err := latestGoVersionCmd.Output()
	if err != nil {
		return fmt.Errorf("error getting latest Go version: %w", err)
	}

	// Version number, time and final newline
	expectedPartsNum := 3
	latestGoVersionParts := strings.Split(string(latestGoVersionOutput), "\n")

	if len(latestGoVersionParts) != expectedPartsNum {
		return fmt.Errorf("unexpected output from go version command: %s", latestGoVersionOutput)
	}

	latestGoVersion := strings.TrimPrefix(latestGoVersionParts[0], "go")
	slog.Debug("Latest Go version", slog.String("version", latestGoVersion))

	baseArgs := []string{"mod", "edit", "-go=" + string(latestGoVersion)}

	for _, mod := range mods {
		slog.Debug("Updating Go version", slog.String("mod", mod))

		// nosemgrep: go.lang.security.audit.dangerous-syscall-exec.dangerous-syscall-exec // exec.LookPath() is used to locate the binary via $PATH, however we run on trusted developer machines
		// nosemgrep: gitlab.gosec.G204-1 // Same
		command := exec.Command(binary, baseArgs...)
		command.Dir = path.Dir(mod)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err = command.Run()
		if err != nil {
			return fmt.Errorf("error updating Go version in Go module %s: %w", mod, err)
		}

		slog.Info("Updated Go version in Go module", slog.String("mod", mod))
	}

	return nil
}
