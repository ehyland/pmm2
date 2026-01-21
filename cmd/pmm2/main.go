package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/ehyland/pmm2/internal/config"
	"github.com/ehyland/pmm2/internal/defaults"
	"github.com/ehyland/pmm2/internal/executor"
	"github.com/ehyland/pmm2/internal/inspector"
	"github.com/ehyland/pmm2/internal/installer"
	"github.com/ehyland/pmm2/internal/registry"
	"github.com/ehyland/pmm2/internal/shim"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	exeName := filepath.Base(os.Args[0])
	conf := config.LoadConfig()

	switch exeName {
	case "npm", "pnpm", "yarn", "bun":
		if err := executor.RunPackageManager(conf, exeName, exeName, os.Args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	case "pnpx":
		if err := executor.RunPackageManager(conf, "pnpm", "pnpx", os.Args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	case "npx":
		if err := executor.RunPackageManager(conf, "npm", "npx", os.Args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	case "bunx":
		if err := executor.RunPackageManager(conf, "bun", "bunx", os.Args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	rootCmd := &cobra.Command{
		Use:     "pmm2",
		Short:   "Package Manager Manager v2",
		Version: fmt.Sprintf("%s (commit: %s, date: %s)", version, commit, date),
	}

	updateLocalCmd := &cobra.Command{
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

	updateDefaultCmd := &cobra.Command{
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

	updateSelfCmd := &cobra.Command{
		Use:   "update-self",
		Short: "Update pmm2 itself",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdateSelf()
		},
	}

	pinCmd := &cobra.Command{
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

	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "Ensure all shims (npm, pnpm, yarn, etc.) are correctly linked",
		RunE: func(cmd *cobra.Command, args []string) error {
			return shim.EnsureShims()
		},
	}

	rootCmd.AddCommand(updateLocalCmd, updateDefaultCmd, updateSelfCmd, pinCmd, setupCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runUpdateSelf() error {
	updater, err := selfupdate.NewUpdater(selfupdate.Config{})
	if err != nil {
		return err
	}
	latest, found, err := updater.DetectLatest(context.Background(), selfupdate.ParseSlug("ehyland/pmm2"))
	if err != nil {
		return err
	}
	if !found {
		fmt.Println("No updates found.")
		return nil
	}

	if latest.LessOrEqual(version) {
		fmt.Printf("pmm2 %s is already up to date\n", version)
		return nil
	}

	fmt.Printf("Updating to %s...\n", latest.Version())
	if err := updater.UpdateTo(context.Background(), latest, os.Args[0]); err != nil {
		return err
	}

	fmt.Println("Self-update successful. Synchronizing shims...")

	// Use syscall.Exec to run the NEW binary with the 'setup' command
	// This ensures the new config.Shims list from the updated binary is used.
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	return syscall.Exec(exePath, []string{filepath.Base(exePath), "setup"}, os.Environ())
}
