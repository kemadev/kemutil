// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package wgo

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"

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
