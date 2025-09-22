// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package iac

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/kemadev/kemutil/internal/gomodtool"
	"github.com/kemadev/kemutil/pkg/str"
	"github.com/kemadev/kemutil/pkg/wgo"
	"github.com/spf13/cobra"
)

type templatedFile struct {
	Name    string
	Content string
}

const (
	StackNameDev  = "dev"
	StackNameNext = "next"
	StackNameProd = "main"
)

var (
	// Debug is a flag to enable debug output for Pulumi commands.
	//nolint:gochecknoglobals // Cobra flags are global
	Debug bool
	// Refresh is a flag to refresh the Pulumi stack before updating.
	//nolint:gochecknoglobals // Cobra flags are global
	Refresh bool
)

const MainGoContent = `package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		return nil
	})
}
`

func renderTemplates() ([]templatedFile, error) {
	moduleName, err := gomodtool.GetGoModExpectedName()
	if err != nil {
		return nil, fmt.Errorf("error getting expected Go module name: %w", err)
	}

	moduleNameDash := str.KebabCase(moduleName)

	pulumiYaml := templatedFile{
		Name: "Pulumi.yaml",
		Content: `name: ` + moduleNameDash + `
description: IaC for ` + moduleName + `
runtime: go
config:
  pulumi:disable-default-providers:
    description: Disable default providers to enforce using configured ones
    value:
      - '*'
`,
	}
	pulumiDevYaml := templatedFile{
		Name: "Pulumi." + StackNameDev + ".yaml",
		Content: `config: {}
`,
	}
	pulumiNextYaml := templatedFile{
		Name: "Pulumi." + StackNameNext + ".yaml",
		Content: `config: {}
`,
	}
	pulumiProdYaml := templatedFile{
		Name: "Pulumi." + StackNameProd + ".yaml",
		Content: `config: {}
`,
	}
	mainGo := templatedFile{
		Name:    "main.go",
		Content: MainGoContent,
	}
	templatedInitFiles := []templatedFile{
		pulumiYaml,
		pulumiDevYaml,
		pulumiNextYaml,
		pulumiProdYaml,
		mainGo,
	}

	return templatedInitFiles, nil
}

// Init initializes a IaC module in the current directory.
func Init(_ *cobra.Command, _ []string) error {
	slog.Info("Initializing IaC module")

	err := wgo.Init(nil, nil)
	if err != nil {
		return fmt.Errorf("error initializing Go module: %w", err)
	}

	templatedInitFiles, err := renderTemplates()
	if err != nil {
		return fmt.Errorf("error rendering templates: %w", err)
	}

	for _, file := range templatedInitFiles {
		filePath := path.Join(".", file.Name)

		FilePermReadWriteCurrentUser := 0o600

		err := os.WriteFile(
			filePath,
			[]byte(file.Content),
			os.FileMode(FilePermReadWriteCurrentUser),
		)
		if err != nil {
			return fmt.Errorf("error writing templated file %s: %w", filePath, err)
		}

		slog.Debug("Created templated file", slog.String("filePath", filePath))
	}

	err = wgo.Update(nil, nil)
	if err != nil {
		return fmt.Errorf("error updating Go module after initialization: %w", err)
	}

	return nil
}
