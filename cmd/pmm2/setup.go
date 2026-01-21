package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ehyland/pmm2/internal/config"
	"github.com/ehyland/pmm2/internal/shim"
	"github.com/spf13/cobra"
)

func newSetupCmd(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Ensure all shims (npm, pnpm, yarn, etc.) are correctly linked",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := shim.EnsureShims(); err != nil {
				return err
			}
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			return ensurePathInBashrc(home, filepath.Join(conf.PmmDir, "bin"))
		},
	}
}

func ensurePathInBashrc(homeDir, binDir string) error {
	bashrcPath := filepath.Join(homeDir, ".bashrc")

	if _, err := exec.LookPath("pmm2"); err == nil {
		// pmm2 is already in PATH
		return nil
	}

	f, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .bashrc: %w", err)
	}
	defer f.Close()

	// Prepend a newline just in case the file doesn't end with one
	lineToAdd := fmt.Sprintf("\n# pmm2\nexport PATH=\"%s:$PATH\"\n", binDir)
	if _, err := f.WriteString(lineToAdd); err != nil {
		return err
	}

	fmt.Printf("Added %s to PATH in %s\n", binDir, bashrcPath)
	return nil
}
