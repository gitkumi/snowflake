package generator_test

import (
	"os"
	"testing"

	"github.com/gitkumi/snowflake/internal/generator"
)

func TestGenerateFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "generator_test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(tmpDir)

	if err := generator.Generate("acme", false, tmpDir); err != nil {
		t.Fatalf("generator.Create returned an error: %v", err)
	}
}
