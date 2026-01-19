package executor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/ehyland/pmm2/internal/config"
	"github.com/ehyland/pmm2/internal/defaults"
	"github.com/ehyland/pmm2/internal/inspector"
	"github.com/ehyland/pmm2/internal/installer"
)

func RunPackageManager(packageManagerName string, executableName string, args []string) error {
	conf := config.LoadConfig()

	if !config.IsSupported(packageManagerName) {
		return fmt.Errorf("unsupported package manager: %s", packageManagerName)
	}

	found, err := inspector.FindPackageManagerSpec()
	if err != nil {
		return err
	}

	var spec *inspector.PackageManagerSpec
	if found != nil {
		if found.Spec.Name != packageManagerName {
			if conf.IgnoreSpecMismatch {
				spec = nil
			} else {
				relPath, _ := filepath.Rel(".", found.PackageJSONPath)
				fmt.Fprintf(os.Stderr, "⚠️  This project is configured to use %s.\n", found.Spec.Name)
				fmt.Fprintf(os.Stderr, "See \"packageManager\" field in ./%s\n", relPath)
				os.Exit(1)
			}
		} else {
			spec = &found.Spec
		}
	}

	if spec == nil {
		version, err := defaults.GetDefaultVersion(conf, packageManagerName)
		if err != nil {
			return err
		}
		spec = &inspector.PackageManagerSpec{
			Name:    packageManagerName,
			Version: version,
		}
	}

	if err := installer.Install(conf, *spec); err != nil {
		return err
	}

	exePath, err := installer.GetExecutablePath(conf, *spec, executableName)
	if err != nil {
		return err
	}

	nodePath, err := exec.LookPath("node")
	if err != nil {
		return fmt.Errorf("node not found in PATH: %w", err)
	}

	cmdArgs := append([]string{exePath}, args...)
	env := os.Environ()
	env = append(env, "PMM_IGNORE_SPEC_MISS_MATCH=1")

	return syscall.Exec(nodePath, append([]string{"node"}, cmdArgs...), env)
}
