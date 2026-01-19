package shim

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ehyland/pmm2/internal/config"
)

func EnsureShims() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	binDir := filepath.Dir(exePath)
	exeName := filepath.Base(exePath)

	for _, shimName := range config.GetShims() {
		shimPath := filepath.Join(binDir, shimName)
		// Remove if exists
		if _, err := os.Lstat(shimPath); err == nil {
			os.Remove(shimPath)
		}
		fmt.Printf("Creating shim: %s -> %s\n", shimName, exeName)
		if err := os.Symlink(exeName, shimPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to create shim %s: %v\n", shimName, err)
		}
	}
	return nil
}
