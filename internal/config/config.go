package config

import (
	"os"
	"path/filepath"
	"strings"
)

var SupportedPackageManagers = []string{"pnpm", "npm", "yarn"}

type Config struct {
	Registry             string
	PmmDir               string
	IgnoreSpecMismatch   bool
}

func LoadConfig() *Config {
	registry := os.Getenv("PMM_NPM_REGISTRY")
	if registry == "" {
		registry = "https://registry.npmjs.org"
	}

	pmmDir := os.Getenv("PMM_DIR")
	if pmmDir == "" {
		home, _ := os.UserHomeDir()
		pmmDir = filepath.Join(home, ".pmm")
	}

	ignoreStr := os.Getenv("PMM_IGNORE_SPEC_MISS_MATCH")
	ignore := false
	if strings.ToLower(ignoreStr) == "yes" || strings.ToLower(ignoreStr) == "true" || ignoreStr == "1" {
		ignore = true
	}

	return &Config{
		Registry:           registry,
		PmmDir:             pmmDir,
		IgnoreSpecMismatch: ignore,
	}
}

func IsSupported(name string) bool {
	for _, supported := range SupportedPackageManagers {
		if name == supported {
			return true
		}
	}
	return false
}
