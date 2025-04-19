package initialize

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

	cf := &Config{
		Quiet:     true,
		Name:      "acme",
		Database:  SQLite3,
		AppType:   API,
		OutputDir: tmpDir,
		NoGit:     true,
		NoSMTP:    false,
		NoStorage: false,
		NoAuth:    false,
	}

	err = Run(cf)
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

	cfg := &Config{
		Quiet:     true,
		Name:      "acme",
		Database:  Postgres,
		AppType:   API,
		OutputDir: tmpDir,
		NoGit:     true,
		NoSMTP:    false,
		NoStorage: false,
		NoAuth:    false,
	}

	err = Run(cfg)
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

	cfg := &Config{
		Quiet:     true,
		Name:      "acme",
		Database:  MySQL,
		AppType:   API,
		OutputDir: tmpDir,
		NoGit:     true,
		NoSMTP:    false,
		NoStorage: false,
		NoAuth:    false,
	}

	err = Run(cfg)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("Project directory was not created")
	}
}

func TestGenerateWebApp(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &Config{
		Quiet:     true,
		Name:      "acme",
		Database:  SQLite3,
		AppType:   Web,
		OutputDir: tmpDir,
		NoGit:     true,
	}

	err = Run(cfg)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("Project directory was not created")
	}

	apiMainPath := filepath.Join(projectDir, "cmd", "api")
	webMainPath := filepath.Join(projectDir, "cmd", "web")
	htmlDirPath := filepath.Join(projectDir, "internal", "html")

	if _, err := os.Stat(webMainPath); os.IsNotExist(err) {
		t.Fatal("Web directory was not created")
	}

	if _, err := os.Stat(apiMainPath); err == nil {
		t.Fatal("API directory should not exist for Web app type")
	}

	if _, err := os.Stat(htmlDirPath); os.IsNotExist(err) {
		t.Fatal("HTML directory created for web app type")
	}
}

func TestGenerateNoSMTP(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &Config{
		Quiet:     true,
		Name:      "acme",
		Database:  SQLite3,
		AppType:   API,
		OutputDir: tmpDir,
		NoGit:     true,
		NoSMTP:    true,
		NoStorage: false,
		NoAuth:    false,
	}

	err = Run(cf)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("Project directory was not created")
	}
}

func TestGenerateNoStorage(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &Config{
		Quiet:     true,
		Name:      "acme",
		Database:  SQLite3,
		AppType:   API,
		OutputDir: tmpDir,
		NoGit:     true,
		NoSMTP:    false,
		NoStorage: true,
		NoAuth:    false,
	}

	err = Run(cf)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("Project directory was not created")
	}
}

func TestGenerateNoAuth(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &Config{
		Quiet:     true,
		Name:      "acme",
		Database:  SQLite3,
		AppType:   API,
		OutputDir: tmpDir,
		NoGit:     true,
		NoSMTP:    false,
		NoStorage: false,
		NoAuth:    true,
	}

	err = Run(cf)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("Project directory was not created")
	}
}
