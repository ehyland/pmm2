package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ehyland/pmm2/internal/config"
	"github.com/ehyland/pmm2/internal/executor"
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

	rootCmd.AddCommand(
		newUpdateLocalCmd(conf),
		newUpdateDefaultCmd(conf),
		newUpdateSelfCmd(version),
		newPinCmd(conf),
		newSetupCmd(conf),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
