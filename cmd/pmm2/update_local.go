package main

import (
	"fmt"

	"github.com/ehyland/pmm2/internal/config"
	"github.com/ehyland/pmm2/internal/inspector"
	"github.com/ehyland/pmm2/internal/installer"
	"github.com/ehyland/pmm2/internal/registry"
	"github.com/spf13/cobra"
)

func newUpdateLocalCmd(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "update-local",
		Short: "Update package manager version in package.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			search, err := inspector.FindPackageManagerSpec()
			if err != nil {
				return err
			}
			if search == nil {
				return fmt.Errorf("unable to find package.json with \"packageManager\" field")
			}

			latest, err := registry.GetLatestVersion(conf, search.Spec.Name)
			if err != nil {
				return err
			}

			if latest.Version == search.Spec.Version {
				fmt.Printf("Already on latest version %s@%s\n", search.Spec.Name, latest.Version)
				return nil
			}

			if err := installer.Install(conf, *latest); err != nil {
				return err
			}

			fmt.Printf("Updating %s to %s@%s\n", search.PackageJSONPath, latest.Name, latest.Version)
			return inspector.UpdateSpecInPackageJSON(search.PackageJSONPath, *latest)
		},
	}
}
