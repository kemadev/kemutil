// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package repotpl

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/kemadev/ci-cd/pkg/git"
	"github.com/spf13/cobra"
)

var (
	RepoTemplateURL = url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   "kemadev/repo-template",
	}
	SkipAnswered bool
)

const copierConfigPath = "config/copier/.copier-answers.yml"

// SkipAnswered indicates whether to skip questions that have already been answered during repository template update,.

// Init initializes the repository template.
func Init(_ *cobra.Command, _ []string) error {
	slog.Info("Initializing repository template")

	binary, err := exec.LookPath("copier")
	if err != nil {
		return fmt.Errorf("copier binary not found: %w", err)
	}

	baseArgs := []string{"copy", RepoTemplateURL.String(), "."}

	repo, err := git.GetGitBasePath()
	if err != nil {
		return fmt.Errorf("error getting git base path: %w", err)
	}

	if repo == "" {
		return fmt.Errorf(
			"not in a git repository, please run this command inside a git repository",
		)
	}

	repoFull := "https://" + repo

	repoURL, err := url.Parse(repoFull)
	if err != nil {
		return fmt.Errorf("error parsing git repository URL: %w", err)
	}

	baseArgs = append(baseArgs, "--data", "vcs_url_scheme"+"="+repoURL.Scheme+"://")
	baseArgs = append(baseArgs, "--data", "vcs_server_host"+"="+repoURL.Host)

	vcsParts := strings.Split(strings.TrimPrefix(repoURL.Path, "/"), "/")
	if len(vcsParts) != 2 {
		return fmt.Errorf("unexpected git repository URL format: %s", repoURL.Path)
	}

	baseArgs = append(baseArgs, "--data", "vcs_namespace"+"="+vcsParts[0])
	baseArgs = append(baseArgs, "--data", "vcs_repo"+"="+vcsParts[1])

	if repoURL.Host == "github.com" {
		baseArgs = append(baseArgs, "--data", "vcs_provider"+"="+"github")
	}

	slog.Debug("Running command", slog.Any("binary", binary), slog.Any("baseArgs", baseArgs))

	err = syscall.Exec(binary, append([]string{binary}, baseArgs...), os.Environ())
	if err != nil {
		return fmt.Errorf("error running copier init command: %w", err)
	}

	return nil
}

// Update updates the repository template with the latest changes from upstream.
func Update(_ *cobra.Command, _ []string) error {
	slog.Info("Updating repository template")

	binary, err := exec.LookPath("copier")
	if err != nil {
		return fmt.Errorf("copier binary not found: %w", err)
	}

	baseArgs := []string{"update", "--answers-file", copierConfigPath}

	if SkipAnswered {
		slog.Debug("Skip answered questions enabled", slog.Bool("skipAnswered", SkipAnswered))

		baseArgs = append(baseArgs, "--skip-answered")
	}

	slog.Debug("Running command", slog.Any("binary", binary), slog.Any("baseArgs", baseArgs))

	err = syscall.Exec(binary, append([]string{binary}, baseArgs...), os.Environ())
	if err != nil {
		return fmt.Errorf("error running copier update command: %w", err)
	}

	return nil
}
