package initialize_test

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gitkumi/snowflake/internal/initialize"
)

func TestGenerateVariants(t *testing.T) {
	tests := []struct {
		name string
		cfg  initialize.Config
	}{
		{
			name: "no_db",
			cfg: initialize.Config{
				Quiet:         true,
				Name:          "acme",
				Database:      initialize.DatabaseNone,
				Git:           false,
				SMTP:          true,
				Storage:       true,
				KeyValueStore: initialize.KeyValueStoreRedis,
			},
		},
		{
			name: "sqlite3",
			cfg: initialize.Config{
				Quiet:         true,
				Name:          "acme",
				Database:      initialize.DatabaseSQLite3,
				Git:           false,
				SMTP:          true,
				Storage:       true,
				KeyValueStore: initialize.KeyValueStoreRedis,
			},
		},
		{
			name: "postgres",
			cfg: initialize.Config{
				Quiet:         true,
				Name:          "acme",
				Database:      initialize.DatabasePostgres,
				Git:           false,
				SMTP:          true,
				Storage:       true,
				KeyValueStore: initialize.KeyValueStoreRedis,
			},
		},
		{
			name: "mysql",
			cfg: initialize.Config{
				Quiet:         true,
				Name:          "acme",
				Database:      initialize.DatabaseMySQL,
				Git:           false,
				SMTP:          true,
				Storage:       true,
				KeyValueStore: initialize.KeyValueStoreRedis,
			},
		},
		{
			name: "mariadb",
			cfg: initialize.Config{
				Quiet:         true,
				Name:          "acme",
				Database:      initialize.DatabaseMariaDB,
				Git:           false,
				SMTP:          true,
				Storage:       true,
				KeyValueStore: initialize.KeyValueStoreRedis,
			},
		},
		{
			name: "no_smtp",
			cfg: initialize.Config{
				Quiet:         true,
				Name:          "acme",
				Database:      initialize.DatabaseSQLite3,
				Git:           false,
				SMTP:          false,
				Storage:       true,
				KeyValueStore: initialize.KeyValueStoreRedis,
			},
		},
		{
			name: "no_storage",
			cfg: initialize.Config{
				Quiet:         true,
				Name:          "acme",
				Database:      initialize.DatabaseSQLite3,
				Git:           false,
				SMTP:          true,
				Storage:       false,
				KeyValueStore: initialize.KeyValueStoreRedis,
			},
		},
		{
			name: "valkey",
			cfg: initialize.Config{
				Quiet:         true,
				Name:          "acme",
				Database:      initialize.DatabaseSQLite3,
				Git:           false,
				SMTP:          true,
				Storage:       true,
				KeyValueStore: initialize.KeyValueStoreValkey,
			},
		},
		{
			name: "no_redis",
			cfg: initialize.Config{
				Quiet:         true,
				Name:          "acme",
				Database:      initialize.DatabaseSQLite3,
				Git:           false,
				SMTP:          true,
				Storage:       true,
				KeyValueStore: initialize.KeyValueStoreNone,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.cfg
			cfg.OutputDir = t.TempDir()

			if err := initialize.Generate(&cfg); err != nil {
				t.Fatal(err)
			}

			projectDir := filepath.Join(cfg.OutputDir, cfg.Name)
			if _, err := os.Stat(projectDir); os.IsNotExist(err) {
				t.Fatal("project directory not created")
			}
		})
	}
}

func TestGenerateTempl(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Generate(&initialize.Config{
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

	err := initialize.Generate(&initialize.Config{
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

	err := initialize.Generate(&initialize.Config{
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
	envPath := filepath.Join(projectDir, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		t.Fatalf(".env file not created at %s", envPath)
	}

	// Check for .env.test file
	envTestPath := filepath.Join(projectDir, ".env.test")
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

func TestGenerateDatabaseScaffoldInternalImportsResolve(t *testing.T) {
	cfg := initialize.Config{
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseSQLite3,
		OutputDir: t.TempDir(),
		Git:       false,
	}

	if err := initialize.Generate(&cfg); err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(cfg.OutputDir, cfg.Name)
	if err := assertInternalImportsResolve(projectDir, cfg.Name); err != nil {
		t.Fatal(err)
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

		err := initialize.Generate(&initialize.Config{
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
			t.Logf("initialize.Generate returned error: %v", err)
			return
		}

		projectDir := filepath.Join(tmpDir, "acme")
		if _, err := os.Stat(projectDir); os.IsNotExist(err) {
			t.Fatal("project directory not created")
		}
	})
}

func assertInternalImportsResolve(projectDir string, module string) error {
	fset := token.NewFileSet()

	return filepath.WalkDir(projectDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}

		for _, spec := range file.Imports {
			importPath, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return fmt.Errorf("unquote import in %s: %w", path, err)
			}
			if !strings.HasPrefix(importPath, module+"/") {
				continue
			}

			pkgDir := filepath.Join(projectDir, filepath.FromSlash(strings.TrimPrefix(importPath, module+"/")))
			if err := assertPackageExists(pkgDir); err != nil {
				relPath, relErr := filepath.Rel(projectDir, path)
				if relErr != nil {
					relPath = path
				}
				return fmt.Errorf("%s imports missing package %q: %w", relPath, importPath, err)
			}
		}

		return nil
	})
}

func assertPackageExists(pkgDir string) error {
	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".go") {
			return nil
		}
	}

	return fmt.Errorf("no Go files found in %s", pkgDir)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
