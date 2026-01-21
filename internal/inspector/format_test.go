package inspector

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdateSpecInPackageJSON_PreservesFormat(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "package.json")

	// Notice:
	// 1. "version" is before "name" (unusual, but order should be preserved)
	// 2. "scripts" has "&&" inside
	// 3. Indentation is 4 spaces (standard MarshalIndent uses 2, or we pass it)
	// 4. Trailing newline
	initialContent := `{
    "version": "1.0.0",
    "name": "test-pkg",
    "scripts": {
        "test": "echo 'hello' && exit 0"
    },
    "packageManager": "npm@6.0.0"
}
`
	if err := os.WriteFile(pkgPath, []byte(initialContent), 0644); err != nil {
		t.Fatal(err)
	}

	newSpec := PackageManagerSpec{Name: "pnpm", Version: "8.0.0"}
	if err := UpdateSpecInPackageJSON(pkgPath, newSpec); err != nil {
		t.Fatalf("UpdateSpecInPackageJSON failed: %v", err)
	}

	contentBytes, err := os.ReadFile(pkgPath)
	if err != nil {
		t.Fatal(err)
	}
	content := string(contentBytes)

	// Check 1: HTML escaping
	if strings.Contains(content, "\\u0026") {
		t.Errorf("found unicode escape sequence for &: %s", content)
	}
	if !strings.Contains(content, "&&") {
		t.Errorf("expected && to be preserved, got: %s", content)
	}

	// Check 2: Order preservation
	// We expect "version" to appear before "name"
	vIndex := strings.Index(content, "\"version\"")
	nIndex := strings.Index(content, "\"name\"")
	if vIndex == -1 || nIndex == -1 {
		t.Fatal("missing keys")
	}
	if vIndex > nIndex {
		t.Errorf("keys were reordered! version=%d, name=%d\nContent:\n%s", vIndex, nIndex, content)
	}
}
