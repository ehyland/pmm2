package installer

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ehyland/pmm2/internal/config"
	"github.com/ehyland/pmm2/internal/inspector"
	"github.com/ehyland/pmm2/internal/registry"
)

type PackageJSON struct {
	Name string            `json:"name"`
	Bin  map[string]string `json:"bin"`
}

func GetInstallPath(conf *config.Config, spec inspector.PackageManagerSpec) string {
	return filepath.Join(conf.PmmDir, "installed-versions", fmt.Sprintf("%s-%s", spec.Name, spec.Version))
}

func IsInstalled(conf *config.Config, spec inspector.PackageManagerSpec) bool {
	path := filepath.Join(GetInstallPath(conf, spec), "package.json")
	_, err := os.Stat(path)
	return err == nil
}

func Install(conf *config.Config, spec inspector.PackageManagerSpec) error {
	if IsInstalled(conf, spec) {
		return nil
	}

	fmt.Printf("Installing %s@%s...\n", spec.Name, spec.Version)

	body, err := registry.DownloadTarball(conf, spec)
	if err != nil {
		return err
	}
	defer body.Close()

	installPath := GetInstallPath(conf, spec)
	if err := os.RemoveAll(installPath); err != nil {
		return err
	}
	if err := os.MkdirAll(installPath, 0755); err != nil {
		return err
	}

	if err := extractTarGz(body, installPath); err != nil {
		return err
	}

	return nil
}

func extractTarGz(gzipStream io.Reader, dest string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Most npm packages contain a "package/" prefix in the tarball
		// we want to strip the first component
		parts := strings.Split(header.Name, "/")
		if len(parts) <= 1 {
			continue
		}
		relPath := filepath.Join(parts[1:]...)
		target := filepath.Join(dest, relPath)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tarReader); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
	return nil
}

func GetExecutablePath(conf *config.Config, spec inspector.PackageManagerSpec, executableName string) (string, error) {
	installPath := GetInstallPath(conf, spec)
	pkgJSONPath := filepath.Join(installPath, "package.json")

	data, err := os.ReadFile(pkgJSONPath)
	if err != nil {
		return "", err
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "", err
	}

	relPath, ok := pkg.Bin[executableName]
	if !ok {
		// Sometimes 'bin' is just a string
		// But for pnpm/npm/yarn it's usually an object
		return "", fmt.Errorf("executable %s not found in package.json", executableName)
	}

	return filepath.Join(installPath, relPath), nil
}
