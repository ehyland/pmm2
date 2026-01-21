package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("PMM_NPM_REGISTRY", "https://test.registry.org")
	os.Setenv("PMM2_DIR", "/tmp/.pmm-test")
	os.Setenv("PMM_IGNORE_SPEC_MISS_MATCH", "true")

	conf := LoadConfig()

	if conf.Registry != "https://test.registry.org" {
		t.Errorf("expected registry https://test.registry.org, got %s", conf.Registry)
	}

	if conf.PmmDir != "/tmp/.pmm-test" {
		t.Errorf("expected pmmDir /tmp/.pmm-test, got %s", conf.PmmDir)
	}

	if !conf.IgnoreSpecMismatch {
		t.Errorf("expected IgnoreSpecMismatch true, got false")
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	os.Unsetenv("PMM_NPM_REGISTRY")
	os.Unsetenv("PMM2_DIR")
	os.Unsetenv("PMM_IGNORE_SPEC_MISS_MATCH")

	conf := LoadConfig()

	if conf.Registry != "https://registry.npmjs.org" {
		t.Errorf("expected default registry, got %s", conf.Registry)
	}

	home, _ := os.UserHomeDir()
	expectedDir := filepath.Join(home, ".pmm2")
	if conf.PmmDir != expectedDir {
		t.Errorf("expected default pmmDir %s, got %s", expectedDir, conf.PmmDir)
	}

	if conf.IgnoreSpecMismatch {
		t.Errorf("expected default IgnoreSpecMismatch false, got true")
	}
}

func TestIsSupported(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"npm", true},
		{"pnpm", true},
		{"yarn", true},
		{"bun", true},
	}

	for _, tt := range tests {
		if got := IsSupported(tt.name); got != tt.expected {
			t.Errorf("IsSupported(%s) = %v; want %v", tt.name, got, tt.expected)
		}
	}
}
