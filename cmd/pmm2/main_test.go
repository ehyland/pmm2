package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsurePathInBashrc(t *testing.T) {
	// Create a temporary directory for home
	tmpHome, err := os.MkdirTemp("", "pmm2-test-home")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	// Create a dummy bin dir path
	binDir := "/tmp/pmm2/bin"

	// Mocking PATH to ensure pmm2 is NOT found
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)
	os.Setenv("PATH", "") // clear path so LookPath("pmm2") fails

	// Case 1: .bashrc doesn't exist
	err = ensurePathInBashrc(tmpHome, binDir)
	if err != nil {
		t.Errorf("ensurePathInBashrc failed when file missing: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpHome, ".bashrc"))
	if err != nil {
		t.Fatalf("Failed to read .bashrc: %v", err)
	}
	if !strings.Contains(string(content), binDir) {
		t.Errorf("Expected .bashrc to contain %s, got:\n%s", binDir, string(content))
	}

	// Case 2: .bashrc exists and has content
	// We run it again. Since our logic checks LookPath (which fails), it should append again.
	err = ensurePathInBashrc(tmpHome, binDir)
	if err != nil {
		t.Errorf("ensurePathInBashrc failed on second run: %v", err)
	}
	content, _ = os.ReadFile(filepath.Join(tmpHome, ".bashrc"))
	if strings.Count(string(content), binDir) != 2 {
		t.Errorf("Expected .bashrc to contain %s twice, got count %d:\n%s", binDir, strings.Count(string(content), binDir), string(content))
	}

	// Case 3: pmm2 IS in PATH
	// We put pmm2 in a directory and add it to PATH
	tmpBin := filepath.Join(tmpHome, "bin")
	os.MkdirAll(tmpBin, 0755)

	// Create a dummy pmm2 executable
	dummyExe := filepath.Join(tmpBin, "pmm2")
	f, err := os.Create(dummyExe)
	if err != nil {
		t.Fatalf("Failed to create dummy exe: %v", err)
	}
	f.Chmod(0755)
	f.Close()

	os.Setenv("PATH", tmpBin)

	err = ensurePathInBashrc(tmpHome, binDir)
	if err != nil {
		t.Errorf("ensurePathInBashrc failed when in PATH: %v", err)
	}

	// Should NOT have appended a 3rd time
	content, _ = os.ReadFile(filepath.Join(tmpHome, ".bashrc"))
	if strings.Count(string(content), binDir) != 2 {
		t.Errorf("Expected .bashrc to still contain %s twice, got count %d", binDir, strings.Count(string(content), binDir))
	}
}
