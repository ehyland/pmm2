package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ehyland/pmm2/internal/config"
	"github.com/ehyland/pmm2/internal/inspector"
)

type Packument struct {
	DistTags map[string]string `json:"dist-tags"`
}

func GetLatestVersion(conf *config.Config, name string) (*inspector.PackageManagerSpec, error) {
	url := fmt.Sprintf("%s/%s", conf.Registry, name)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch packument: %s", resp.Status)
	}

	var packument Packument
	if err := json.NewDecoder(resp.Body).Decode(&packument); err != nil {
		return nil, err
	}

	version, ok := packument.DistTags["latest"]
	if !ok {
		return nil, fmt.Errorf("latest dist-tag not found for %s", name)
	}

	return &inspector.PackageManagerSpec{
		Name:    name,
		Version: version,
	}, nil
}

func DownloadTarball(conf *config.Config, spec inspector.PackageManagerSpec) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/%s/-/%s-%s.tgz", conf.Registry, spec.Name, spec.Name, spec.Version)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to download tarball: %s", resp.Status)
	}

	return resp.Body, nil
}
