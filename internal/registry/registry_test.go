package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ehyland/pmm2/internal/config"
	"github.com/ehyland/pmm2/internal/inspector"
)

func TestGetLatestVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pnpm" {
			t.Errorf("expected request path /pnpm, got %s", r.URL.Path)
		}
		resp := Packument{
			DistTags: map[string]string{
				"latest": "8.0.0",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	conf := &config.Config{Registry: server.URL}
	spec, err := GetLatestVersion(conf, "pnpm")
	if err != nil {
		t.Fatalf("GetLatestVersion() error = %v", err)
	}

	if spec.Name != "pnpm" || spec.Version != "8.0.0" {
		t.Errorf("expected pnpm@8.0.0, got %s@%s", spec.Name, spec.Version)
	}
}

func TestDownloadTarball(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/pnpm/-/pnpm-8.0.0.tgz"
		if r.URL.Path != expectedPath {
			t.Errorf("expected request path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "mock tarball content")
	}))
	defer server.Close()

	conf := &config.Config{Registry: server.URL}
	spec := inspector.PackageManagerSpec{Name: "pnpm", Version: "8.0.0"}
	body, err := DownloadTarball(conf, spec)
	if err != nil {
		t.Fatalf("DownloadTarball() error = %v", err)
	}
	defer body.Close()

	// Just check if we can read it
	content, _ := io.ReadAll(body)
	if string(content) != "mock tarball content" {
		t.Errorf("expected mock tarball content, got %s", string(content))
	}
}
