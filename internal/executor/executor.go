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

type SpecMismatchError struct {
	Expected string
	Path     string
}

func (e *SpecMismatchError) Error() string {
	relPath, _ := filepath.Rel(".", e.Path)
	return fmt.Sprintf("⚠️  This project is configured to use %s.\nSee \"packageManager\" field in ./%s\n\nYou can ignore this error by setting the environment variable PMM_IGNORE_SPEC_MISS_MATCH=1", e.Expected, relPath)
}

func RunPackageManager(conf *config.Config, packageManagerName string, executableName string, args []string) error {
	if !config.IsSupported(packageManagerName) {
		return fmt.Errorf("unsupported package manager: %s", packageManagerName)
	}

	found, err := inspector.FindPackageManagerSpec()
	if err != nil {
		return fmt.Errorf("failed to find package manager spec: %w", err)
	}

	var spec *inspector.PackageManagerSpec
	if found != nil {
		if found.Spec.Name != packageManagerName {
			// TODO: move bun exception to config
			if packageManagerName == "bun" || conf.IgnoreSpecMismatch {
				spec = nil
			} else {
				return &SpecMismatchError{
					Expected: found.Spec.Name,
					Path:     found.PackageJSONPath,
				}
			}
		} else {
			spec = &found.Spec
		}
	}

	if spec == nil {
		version, err := defaults.GetDefaultVersion(conf, packageManagerName)
		if err != nil {
			return fmt.Errorf("failed to get default version: %w", err)
		}
		spec = &inspector.PackageManagerSpec{
			Name:    packageManagerName,
			Version: version,
		}
	}

	if err := installer.Install(conf, *spec); err != nil {
		return fmt.Errorf("failed to install: %w", err)
	}

	exePath, err := installer.GetExecutablePath(conf, *spec, executableName)
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	cmdArgs := append([]string{exePath}, args...)
	env := os.Environ()
	env = append(env, "PMM_IGNORE_SPEC_MISS_MATCH=1")

	if packageManagerName == "bun" {
		return syscall.Exec(exePath, append([]string{executableName}, args...), env)
	}

	nodePath, err := exec.LookPath("node")
	if err != nil {
		return fmt.Errorf("node not found in PATH: %w", err)
	}

	return syscall.Exec(nodePath, append([]string{"node"}, cmdArgs...), env)
}
