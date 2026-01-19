package defaults

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ehyland/pmm2/internal/config"
	"github.com/ehyland/pmm2/internal/inspector"
	"github.com/ehyland/pmm2/internal/registry"
)

func GetDefaultFilePath(conf *config.Config, name string) string {
	return filepath.Join(conf.PmmDir, "installed-versions", ".defaults", name+"-version")
}

func GetDefaultVersion(conf *config.Config, name string) (string, error) {
	path := GetDefaultFilePath(conf, name)
	data, err := os.ReadFile(path)
	if err == nil {
		version := strings.TrimSpace(string(data))
		if version != "" {
			return version, nil
		}
	}

	latest, err := registry.GetLatestVersion(conf, name)
	if err != nil {
		return "", err
	}

	if err := UpdateDefault(conf, *latest); err != nil {
		return "", err
	}

	return latest.Version, nil
}

func UpdateDefault(conf *config.Config, spec inspector.PackageManagerSpec) error {
	path := GetDefaultFilePath(conf, spec.Name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(spec.Version), 0644)
}
