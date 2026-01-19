package installer

import (
	"path/filepath"
	"testing"

	"github.com/ehyland/pmm2/internal/config"
	"github.com/ehyland/pmm2/internal/inspector"
)

func TestGetInstallPath(t *testing.T) {
	conf := &config.Config{PmmDir: "/tmp/.pmm"}
	spec := inspector.PackageManagerSpec{Name: "pnpm", Version: "8.0.0"}
	path := GetInstallPath(conf, spec)
	expected := filepath.Join("/tmp/.pmm", "installed-versions", "pnpm-8.0.0")
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)

	}
}
