package initialize

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateNoDB(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &Config{
		Quiet:     true,
		Name:      "acme",
		Database:  DatabaseNone,
		AppType:   AppTypeAPI,
		OutputDir: tmpDir,
		NoGit:     true,
		NoSMTP:    false,
		NoStorage: false,
		NoAuth:    false,
		NoRedis:   false,
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

func TestGenerateSQLite3(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &Config{
		Quiet:     true,
		Name:      "acme",
		Database:  DatabaseSQLite3,
		AppType:   AppTypeAPI,
		OutputDir: tmpDir,
		NoGit:     true,
		NoSMTP:    false,
		NoStorage: false,
		NoAuth:    false,
		NoRedis:   false,
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
		Database:  DatabasePostgres,
		AppType:   AppTypeAPI,
		OutputDir: tmpDir,
		NoGit:     true,
		NoSMTP:    false,
		NoStorage: false,
		NoAuth:    false,
		NoRedis:   false,
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
		Database:  DatabaseMySQL,
		AppType:   AppTypeAPI,
		OutputDir: tmpDir,
		NoGit:     true,
		NoSMTP:    false,
		NoStorage: false,
		NoAuth:    false,
		NoRedis:   false,
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
		Database:  DatabaseSQLite3,
		AppType:   AppTypeWeb,
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
		Database:  DatabaseSQLite3,
		AppType:   AppTypeAPI,
		OutputDir: tmpDir,
		NoGit:     true,
		NoSMTP:    true,
		NoStorage: false,
		NoAuth:    false,
		NoRedis:   false,
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
		Database:  DatabaseSQLite3,
		AppType:   AppTypeAPI,
		OutputDir: tmpDir,
		NoGit:     true,
		NoSMTP:    false,
		NoStorage: true,
		NoAuth:    false,
		NoRedis:   false,
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
		Database:  DatabaseSQLite3,
		AppType:   AppTypeAPI,
		OutputDir: tmpDir,
		NoGit:     true,
		NoSMTP:    false,
		NoStorage: false,
		NoAuth:    true,
		NoRedis:   false,
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

func TestGenerateNoRedis(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &Config{
		Quiet:     true,
		Name:      "acme",
		Database:  DatabaseSQLite3,
		AppType:   AppTypeAPI,
		OutputDir: tmpDir,
		NoGit:     true,
		NoSMTP:    false,
		NoStorage: false,
		NoAuth:    false,
		NoRedis:   true,
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
