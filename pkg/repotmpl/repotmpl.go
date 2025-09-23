// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package repotmpl

import (
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	g "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/storage/memory"
	"github.com/kemadev/go-framework/pkg/git"
	"github.com/spf13/cobra"
)

const (
	BaseTemplateSubDir        = "template"
	RepositoryNameTemplateKey = "REPONAMETMPL"
)

var (
	//nolint:gochecknoglobals // Used as a const
	RepoTemplateURL = url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   "kemadev/repo-tmpl",
	}
	//nolint:gochecknoglobals // Cobra flags are global
	TargetTemplate string
)

var (
	ErrNotInGitRepo = fmt.Errorf(
		"not in a git repository, please run this command inside a git repository",
	)
	ErrGitURLFormatInvalid = fmt.Errorf("git repository URL format is invalid")
)

// Init initializes the repository template.
func Init(_ *cobra.Command, _ []string) error {
	slog.Info("Initializing repository template")

	repo, err := git.NewGitService().GetGitBasePath()
	if err != nil {
		return fmt.Errorf("error getting git base path: %w", err)
	}

	if repo == "" {
		return ErrNotInGitRepo
	}

	repoParts := strings.Split(repo, "/")
	if len(repoParts) != 3 {
		return fmt.Errorf("repository %q: %w", repo, ErrGitURLFormatInvalid)
	}

	repoName := repoParts[2]

	slog.Debug("found current repository", slog.String("repository", repoName))

	r, err := g.Clone(memory.NewStorage(), nil, &g.CloneOptions{
		URL:           RepoTemplateURL.String(),
		ReferenceName: plumbing.HEAD,
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		return fmt.Errorf("error cloning template repository: %w", err)
	}

	ref, err := r.Head()
	if err != nil {
		return fmt.Errorf("error getting template repository HEAD: %w", err)
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return fmt.Errorf(
			"error getting template repository commit objects for ref %q: %w",
			ref.Hash(),
			err,
		)
	}

	tree, err := commit.Tree()
	if err != nil {
		return fmt.Errorf(
			"error getting template repository commit tree for commit %q: %w",
			commit.Hash,
			err,
		)
	}

	filesBasePath := BaseTemplateSubDir + "/" + TargetTemplate

	err = tree.Files().ForEach(func(file *object.File) error {
		if !strings.HasPrefix(file.Name, filesBasePath) {
			return nil
		}

		relativePath := strings.TrimPrefix(file.Name, filesBasePath)
		if relativePath == "" {
			return nil
		}

		relativePath = strings.ReplaceAll(relativePath, RepositoryNameTemplateKey, repoName)

		targetPath := filepath.Join(".", relativePath)

		dir := filepath.Dir(targetPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		reader, err := file.Reader()
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", file.Name, err)
		}
		defer reader.Close()

		targetFile, err := os.Create(targetPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", targetPath, err)
		}
		defer targetFile.Close()

		content, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("failed to read file contents from %s: %w", file.Name, err)
		}

		contentStr := string(content)
		contentStr = strings.ReplaceAll(contentStr, RepositoryNameTemplateKey, repoName)

		_, err = targetFile.WriteString(contentStr)
		if err != nil {
			return fmt.Errorf("failed to write file contents to %s: %w", targetPath, err)
		}

		slog.Debug("copied file", slog.String("source", file.Name), slog.String("dest", targetPath))

		return nil
	})
	if err != nil {
		return fmt.Errorf("error ooping over tree files: %w", err)
	}

	slog.Debug("repository initialized successfully", slog.String("repository", repoName))

	return nil
}
