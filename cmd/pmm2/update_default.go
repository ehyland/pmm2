package main

import (
	"fmt"

	"github.com/ehyland/pmm2/internal/config"
	"github.com/ehyland/pmm2/internal/defaults"
	"github.com/ehyland/pmm2/internal/inspector"
	"github.com/ehyland/pmm2/internal/installer"
	"github.com/ehyland/pmm2/internal/registry"
	"github.com/spf13/cobra"
)

func newUpdateDefaultCmd(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "update-default [package-manager] [version]",
		Short: "Update the default package manager version",
		Args:  cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := "all"
			if len(args) > 0 {
				name = args[0]
			}

			var toUpdate []inspector.PackageManagerSpec
			if name == "all" {
				for _, pm := range config.GetSupportedPackageManagers() {
					latest, err := registry.GetLatestVersion(conf, pm)
					if err != nil {
						return err
					}
					toUpdate = append(toUpdate, *latest)
				}
			} else {
				if !config.IsSupported(name) {
					return fmt.Errorf("unsupported package manager: %s", name)
				}
				latest, err := registry.GetLatestVersion(conf, name)
				if err != nil {
					return err
				}
				toUpdate = append(toUpdate, *latest)
			}

			for _, spec := range toUpdate {
				if err := installer.Install(conf, spec); err != nil {
					return err
				}
				if err := defaults.UpdateDefault(conf, spec); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
