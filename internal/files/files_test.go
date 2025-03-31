package files_test

import (
	"os"
	"testing"

	"github.com/gitkumi/snowflake/internal/files"
)

func TestCreateFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "files_test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(tmpDir)

	if err := files.Create("acme", false, tmpDir); err != nil {
		t.Fatalf("files.Create returned an error: %v", err)
	}
}
