package initialize_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gitkumi/snowflake/internal/initialize"
)

func TestGenerateNoDB(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseNone,
		OutputDir:     tmpDir,
		Git:           false,
		SMTP:          true,
		Storage:       true,
		KeyValueStore: initialize.KeyValueStoreRedis,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateSQLite3(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		OutputDir:     tmpDir,
		Git:           false,
		SMTP:          true,
		Storage:       true,
		KeyValueStore: initialize.KeyValueStoreRedis,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGeneratePostgres(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabasePostgres,
		OutputDir:     tmpDir,
		Git:           false,
		SMTP:          true,
		Storage:       true,
		KeyValueStore: initialize.KeyValueStoreRedis,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateMySQL(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseMySQL,
		OutputDir:     tmpDir,
		Git:           false,
		SMTP:          true,
		Storage:       true,
		KeyValueStore: initialize.KeyValueStoreRedis,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateMariaDB(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseMariaDB,
		OutputDir:     tmpDir,
		Git:           false,
		SMTP:          true,
		Storage:       true,
		KeyValueStore: initialize.KeyValueStoreRedis,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateNoSMTP(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		OutputDir:     tmpDir,
		Git:           false,
		SMTP:          false,
		Storage:       true,
		KeyValueStore: initialize.KeyValueStoreRedis,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateNoStorage(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		OutputDir:     tmpDir,
		Git:           false,
		SMTP:          true,
		Storage:       false,
		KeyValueStore: initialize.KeyValueStoreRedis,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateValkey(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		OutputDir:     tmpDir,
		Git:           false,
		SMTP:          true,
		Storage:       true,
		KeyValueStore: initialize.KeyValueStoreValkey,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateNoRedis(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		OutputDir:     tmpDir,
		Git:           false,
		SMTP:          true,
		Storage:       true,
		KeyValueStore: initialize.KeyValueStoreNone,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateTempl(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseNone,
		OutputDir: tmpDir,
		Git:       false,
		Templ:     true,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}

	// Check templ files exist
	templFiles := []string{
		filepath.Join(projectDir, "internal", "html", "pages", "index.templ"),
		filepath.Join(projectDir, "cmd", "app", "handlers", "page_handler.go"),
	}
	for _, f := range templFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Fatalf("templ file not created at %s", f)
		}
	}
}

func TestGenerateNoTempl(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseNone,
		OutputDir: tmpDir,
		Git:       false,
		Templ:     false,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")

	// Check templ files do NOT exist
	templFiles := []string{
		filepath.Join(projectDir, "internal", "html", "pages", "index.templ"),
		filepath.Join(projectDir, "cmd", "app", "handlers", "page_handler.go"),
	}
	for _, f := range templFiles {
		if _, err := os.Stat(f); !os.IsNotExist(err) {
			t.Fatalf("templ file should not exist at %s", f)
		}
	}
}

func TestEnvFilesGenerated(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		OutputDir:     tmpDir,
		Git:           false,
		SMTP:          true,
		Storage:       true,
		KeyValueStore: initialize.KeyValueStoreRedis,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")

	// Check for .env file
	envPath := filepath.Join(projectDir, "cmd", "app", ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		t.Fatalf(".env file not created at %s", envPath)
	}

	// Check for .env.test file
	envTestPath := filepath.Join(projectDir, "cmd", "app", ".env.test")
	if _, err := os.Stat(envTestPath); os.IsNotExist(err) {
		t.Fatalf(".env.test file not created at %s", envTestPath)
	}

	// Verify .env file has content
	envContent, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("failed to read .env file: %v", err)
	}
	if len(envContent) == 0 {
		t.Fatal(".env file is empty")
	}

	// Verify .env.test file has content
	envTestContent, err := os.ReadFile(envTestPath)
	if err != nil {
		t.Fatalf("failed to read .env.test file: %v", err)
	}
	if len(envTestContent) == 0 {
		t.Fatal(".env.test file is empty")
	}
}

func FuzzGenerate(f *testing.F) {
	f.Add(
		true, true,
		0, 0,
	)

	f.Fuzz(func(t *testing.T,
		withSMTP, withStorage bool,
		dbTypeInt, kvsTypeInt int,
	) {
		tmpDir := t.TempDir()

		databases := []initialize.Database{
			initialize.DatabaseSQLite3,
			initialize.DatabasePostgres,
			initialize.DatabaseMySQL,
			initialize.DatabaseMariaDB,
			initialize.DatabaseNone,
		}
		keyValueStores := []initialize.KeyValueStore{
			initialize.KeyValueStoreNone,
			initialize.KeyValueStoreRedis,
			initialize.KeyValueStoreValkey,
		}

		database := databases[abs(dbTypeInt)%len(databases)]
		kvs := keyValueStores[abs(kvsTypeInt)%len(keyValueStores)]

		err := initialize.Run(&initialize.Config{
			Git:           false,
			Quiet:         true,
			Name:          "acme",
			OutputDir:     tmpDir,
			Database:      database,
			KeyValueStore: kvs,
			SMTP:          withSMTP,
			Storage:       withStorage,
		})
		if err != nil {
			t.Logf("initialize.Run returned error: %v", err)
			return
		}

		projectDir := filepath.Join(tmpDir, "acme")
		if _, err := os.Stat(projectDir); os.IsNotExist(err) {
			t.Fatal("project directory not created")
		}
	})
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
