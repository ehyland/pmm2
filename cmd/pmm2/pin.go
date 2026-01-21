package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ehyland/pmm2/internal/config"
	"github.com/ehyland/pmm2/internal/inspector"
	"github.com/ehyland/pmm2/internal/registry"
	"github.com/spf13/cobra"
)

func newPinCmd(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "pin <package-manager> <path-to-package>",
		Short: "Write packageManager field to package.json",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			path := args[1]

			if !config.IsSupported(name) {
				return fmt.Errorf("unsupported package manager: %s", name)
			}

			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			pkgJSONPath := absPath
			if !strings.HasSuffix(absPath, "package.json") {
				pkgJSONPath = filepath.Join(absPath, "package.json")
			}

			if _, err := os.Stat(pkgJSONPath); err != nil {
				return fmt.Errorf("package.json not found at %s", pkgJSONPath)
			}

			latest, err := registry.GetLatestVersion(conf, name)
			if err != nil {
				return err
			}

			return inspector.UpdateSpecInPackageJSON(pkgJSONPath, *latest)
		},
	}
}
