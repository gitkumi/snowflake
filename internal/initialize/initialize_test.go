package initialize_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gitkumi/snowflake/internal/initialize"
)

func TestGenerateNoDB(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseNone,
		BackgroundJob: initialize.BackgroundJobNone,
		AppType:       initialize.AppTypeAPI,
		OutputDir:     tmpDir,
		NoGit:         true,
		NoSMTP:        false,
		NoStorage:     false,
		NoAuth:        false,
		NoRedis:       false,
	}

	err = initialize.Run(cf)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateSQLite3(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		BackgroundJob: initialize.BackgroundJobNone,
		AppType:       initialize.AppTypeAPI,
		OutputDir:     tmpDir,
		NoGit:         true,
		NoSMTP:        false,
		NoStorage:     false,
		NoAuth:        false,
		NoRedis:       false,
	}

	err = initialize.Run(cf)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGeneratePostgres(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabasePostgres,
		BackgroundJob: initialize.BackgroundJobNone,
		AppType:       initialize.AppTypeAPI,
		OutputDir:     tmpDir,
		NoGit:         true,
		NoSMTP:        false,
		NoStorage:     false,
		NoAuth:        false,
		NoRedis:       false,
	}

	err = initialize.Run(cfg)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateMySQL(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseMySQL,
		BackgroundJob: initialize.BackgroundJobNone,
		AppType:       initialize.AppTypeAPI,
		OutputDir:     tmpDir,
		NoGit:         true,
		NoSMTP:        false,
		NoStorage:     false,
		NoAuth:        false,
		NoRedis:       false,
	}

	err = initialize.Run(cfg)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateWebApp(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		BackgroundJob: initialize.BackgroundJobNone,
		AppType:       initialize.AppTypeWeb,
		OutputDir:     tmpDir,
		NoGit:         true,
		NoSMTP:        false,
		NoStorage:     false,
		NoAuth:        false,
		NoRedis:       false,
	}

	err = initialize.Run(cfg)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
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

	cf := &initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		BackgroundJob: initialize.BackgroundJobNone,
		AppType:       initialize.AppTypeAPI,
		OutputDir:     tmpDir,
		NoGit:         true,
		NoSMTP:        true,
		NoStorage:     false,
		NoAuth:        false,
		NoRedis:       false,
	}

	err = initialize.Run(cf)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateNoStorage(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		BackgroundJob: initialize.BackgroundJobNone,
		AppType:       initialize.AppTypeAPI,
		OutputDir:     tmpDir,
		NoGit:         true,
		NoSMTP:        false,
		NoStorage:     true,
		NoAuth:        false,
		NoRedis:       false,
	}

	err = initialize.Run(cf)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateNoAuth(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		BackgroundJob: initialize.BackgroundJobNone,
		AppType:       initialize.AppTypeAPI,
		OutputDir:     tmpDir,
		NoGit:         true,
		NoSMTP:        false,
		NoStorage:     false,
		NoAuth:        true,
		NoRedis:       false,
	}

	err = initialize.Run(cf)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateNoRedis(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		BackgroundJob: initialize.BackgroundJobNone,
		AppType:       initialize.AppTypeAPI,
		OutputDir:     tmpDir,
		NoGit:         true,
		NoSMTP:        false,
		NoStorage:     false,
		NoAuth:        false,
		NoRedis:       true,
	}

	err = initialize.Run(cf)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateBackgroundJobBasic(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		BackgroundJob: initialize.BackgroundJobBasic,
		AppType:       initialize.AppTypeAPI,
		OutputDir:     tmpDir,
		NoGit:         true,
		NoSMTP:        false,
		NoStorage:     false,
		NoAuth:        false,
		NoRedis:       false,
	}

	err = initialize.Run(cf)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateBackgroundJobSQS(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		BackgroundJob: initialize.BackgroundJobSQS,
		AppType:       initialize.AppTypeAPI,
		OutputDir:     tmpDir,
		NoGit:         true,
		NoSMTP:        false,
		NoStorage:     false,
		NoAuth:        false,
		NoRedis:       false,
	}

	err = initialize.Run(cf)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateBackgroundJobAsynq(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "snowflake_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cf := &initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		BackgroundJob: initialize.BackgroundJobAsynq,
		AppType:       initialize.AppTypeAPI,
		OutputDir:     tmpDir,
		NoGit:         true,
		NoSMTP:        false,
		NoStorage:     false,
		NoAuth:        false,
		NoRedis:       false,
	}

	err = initialize.Run(cf)
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}
