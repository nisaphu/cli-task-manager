package task

import (
	"os"
	"path/filepath"
	"testing"
)

// TempTaskFile creates a temporary file path for persistence testing and cleanup.
func TempTaskFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "tasks.json")
}

// RemoveFile removes a file if it exists, logging any error on test failure.
func RemoveFile(t *testing.T, path string) {
	t.Helper()
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		t.Logf("failed to remove temp file %s: %v", path, err)
	}
}
