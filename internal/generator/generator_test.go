package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateSQLite3(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &GeneratorConfig{
		Name:      "acme",
		Database:  SQLite3,
		AppType:   API,
		InitGit:   false,
		OutputDir: tmpDir,
	}

	err = Generate(cf)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("Project directory was not created")
	}
}

func TestGeneratePostgres(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &GeneratorConfig{
		Name:      "acme",
		Database:  Postgres,
		AppType:   API,
		InitGit:   false,
		OutputDir: tmpDir,
	}

	err = Generate(cfg)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("Project directory was not created")
	}
}

func TestGenerateMySQL(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &GeneratorConfig{
		Name:      "acme",
		Database:  MySQL,
		AppType:   API,
		InitGit:   false,
		OutputDir: tmpDir,
	}

	err = Generate(cfg)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("Project directory was not created")
	}
}

