package iac

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	ut "github.com/kemadev/infrastructure-components/pkg/util"
	"github.com/kemadev/kemutil/internal/app/util"
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
	StackNameProd = "prod"
)

var (
	// DebugEnabled is a flag to enable debug output for Pulumi commands
	DebugEnabled bool
	// Refresh is a flag to refresh the Pulumi stack before updating
	Refresh bool
)

func renderTemplates() ([]templatedFile, error) {
	moduleName, err := util.GetGoModExpectedName()
	if err != nil {
		return nil, fmt.Errorf("error getting expected Go module name: %w", err)
	}

	moduleNameDash := ut.KebabCase(moduleName)

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
		Name: "main.go",
		Content: `package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/kemadev/infrastructure-components/pkg/k8s/basichttpapp"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		err := basichttpapp.DeployBasicHTTPApp(ctx, basichttpapp.AppParms{
			AppNamespace:        "changeme",
			AppComponent:        "changeme",
			BusinessUnitId:      "changeme",
			CustomerId:          "changeme",
			CostCenter:          "changeme",
			CostAllocationOwner: "changeme",
			OperationsOwner:     "changeme",
			Rpo:                 0 * time.Second,
			MonitoringUrl: url.URL{
				Scheme: "https",
				Host:   "changeme",
				Path:   "changeme",
			},
		})
		if err != nil {
			return fmt.Errorf("error deploying basic HTTP app: %w", err)
		}
		return nil
	})
}
`,
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
func Init(cmd *cobra.Command, args []string) error {
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
		err := os.WriteFile(filePath, []byte(file.Content), 0o644)
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
