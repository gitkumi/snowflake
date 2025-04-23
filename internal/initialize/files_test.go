package initialize_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gitkumi/snowflake/internal/initialize"
)

func TestRemoveEmptyDirs(t *testing.T) {
	tmpDir := t.TempDir()

	nestedDir := filepath.Join(tmpDir, "level1", "level2")
	err := os.MkdirAll(nestedDir, 0777)
	if err != nil {
		t.Fatalf("Failed to create nested directory: %v", err)
	}

	level2File := filepath.Join(nestedDir, "file.txt")
	err = os.WriteFile(level2File, []byte("test"), 0666)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	emptyDir := filepath.Join(tmpDir, "empty")
	err = os.MkdirAll(emptyDir, 0777)
	if err != nil {
		t.Fatalf("Failed to create empty directory: %v", err)
	}

	dirs := map[string]bool{
		emptyDir:  true,
		nestedDir: true,
	}

	err = initialize.RemoveEmptyDirs(dirs)
	if err != nil {
		t.Fatalf("Failed to remove empty directories: %v", err)
	}

	if _, err := os.Stat(emptyDir); err == nil {
		t.Error("Empty directory should have been removed")
	}

	if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
		t.Error("Directory with file should not have been removed")
	}
}

func TestIsDirectoryEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	emptyDir := filepath.Join(tmpDir, "empty")
	err := os.MkdirAll(emptyDir, 0777)
	if err != nil {
		t.Fatalf("Failed to create empty directory: %v", err)
	}

	nonEmptyDir := filepath.Join(tmpDir, "nonempty")
	err = os.MkdirAll(nonEmptyDir, 0777)
	if err != nil {
		t.Fatalf("Failed to create non-empty directory: %v", err)
	}

	testFile := filepath.Join(nonEmptyDir, "file.txt")
	err = os.WriteFile(testFile, []byte("test"), 0666)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	isEmpty, err := initialize.IsDirectoryEmpty(emptyDir)
	if err != nil {
		t.Fatalf("Error checking if empty directory is empty: %v", err)
	}
	if !isEmpty {
		t.Error("Empty directory should be reported as empty")
	}

	isEmpty, err = initialize.IsDirectoryEmpty(nonEmptyDir)
	if err != nil {
		t.Fatalf("Error checking if non-empty directory is empty: %v", err)
	}
	if isEmpty {
		t.Error("Non-empty directory should not be reported as empty")
	}
}
