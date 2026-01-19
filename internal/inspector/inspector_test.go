package inspector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSpecString(t *testing.T) {
	tests := []struct {
		input    string
		expected PackageManagerSpec
		wantErr  bool
	}{
		{"pnpm@8.0.0", PackageManagerSpec{"pnpm", "8.0.0"}, false},
		{"npm@6.14.15", PackageManagerSpec{"npm", "6.14.15"}, false},
		{"yarn@1.22.19", PackageManagerSpec{"yarn", "1.22.19"}, false},
		{"invalid", PackageManagerSpec{}, true},
		{"pnpm@latest", PackageManagerSpec{"pnpm", "latest"}, false},
		{"bun@1.0.0", PackageManagerSpec{}, true},
	}

	for _, tt := range tests {
		got, err := ParseSpecString(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseSpecString(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.expected {
			t.Errorf("ParseSpecString(%s) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestFindPackageManagerSpec(t *testing.T) {
	tmpDir := t.TempDir()

	pkgJSONPath := filepath.Join(tmpDir, "package.json")
	content := `{"packageManager": "pnpm@8.0.0"}`
	if err := os.WriteFile(pkgJSONPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	found, err := FindPackageManagerSpec()
	if err != nil {
		t.Fatalf("FindPackageManagerSpec() error = %v", err)
	}
	if found == nil {
		t.Fatal("expected to find spec, got nil")
	}
	if found.Spec.Name != "pnpm" || found.Spec.Version != "8.0.0" {
		t.Errorf("expected pnpm@8.0.0, got %s@%s", found.Spec.Name, found.Spec.Version)
	}
}
