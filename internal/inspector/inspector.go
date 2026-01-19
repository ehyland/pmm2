package inspector

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ehyland/pmm2/internal/config"
)

type PackageManagerSpec struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type PackageJSON struct {
	PackageManager string `json:"packageManager"`
}

type FoundSpec struct {
	PackageJSONPath string
	Spec            PackageManagerSpec
}

func ParseSpecString(specString string) (PackageManagerSpec, error) {
	parts := strings.Split(specString, "@")
	if len(parts) != 2 {
		return PackageManagerSpec{}, fmt.Errorf("invalid spec format: %s", specString)
	}

	name := parts[0]
	version := parts[1]

	if !config.IsSupported(name) {
		return PackageManagerSpec{}, fmt.Errorf("unsupported package manager: %s", name)
	}

	return PackageManagerSpec{
		Name:    name,
		Version: version,
	}, nil
}

func FindPackageManagerSpec() (*FoundSpec, error) {
	current, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	for {
		pkgJSONPath := filepath.Join(current, "package.json")
		if _, err := os.Stat(pkgJSONPath); err == nil {
			spec, err := loadSpecFromPkgJSON(pkgJSONPath)
			if err != nil {
				return nil, fmt.Errorf("failed to load spec from %s: %w", pkgJSONPath, err)
			}
			if spec != nil {
				return &FoundSpec{
					PackageJSONPath: pkgJSONPath,
					Spec:            *spec,
				}, nil
			}
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return nil, nil
}

func loadSpecFromPkgJSON(path string) (*PackageManagerSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	if pkg.PackageManager == "" {
		return nil, nil
	}

	spec, err := ParseSpecString(pkg.PackageManager)
	if err != nil {
		return nil, err
	}

	return &spec, nil
}

func UpdateSpecInPackageJSON(path string, spec PackageManagerSpec) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return err
	}

	pkg["packageManager"] = fmt.Sprintf("%s@%s", spec.Name, spec.Version)

	newData, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, newData, 0644)
}
