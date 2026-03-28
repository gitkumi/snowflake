package generate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateResource(t *testing.T) {
	databases := []string{"postgres", "mysql", "mariadb", "sqlite3"}

	for _, db := range databases {
		t.Run(db, func(t *testing.T) {
			projectDir := t.TempDir()

			// Setup minimal project structure
			setupProjectDir(t, projectDir, db)

			err := Run("post", []string{"title:string", "body:text", "published:bool"}, projectDir, true)
			if err != nil {
				t.Fatal(err)
			}

			// Verify generated files exist
			expectedFiles := []string{
				"cmd/app/service/post_service.go",
				"cmd/app/handlers/post_handler.go",
			}

			for _, f := range expectedFiles {
				path := filepath.Join(projectDir, f)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Errorf("expected file %s was not created", f)
				}
			}

			// Verify migration file exists (with timestamp prefix)
			migrationsDir := filepath.Join(projectDir, "cmd", "app", "sql", "migrations")
			entries, err := os.ReadDir(migrationsDir)
			if err != nil {
				t.Fatal(err)
			}

			foundMigration := false
			for _, e := range entries {
				if filepath.Ext(e.Name()) == ".sql" {
					foundMigration = true
					break
				}
			}
			if !foundMigration {
				t.Error("no migration file generated")
			}

			// Verify queries file exists
			queriesFile := filepath.Join(projectDir, "cmd", "app", "sql", "queries", "posts.sql")
			if _, err := os.Stat(queriesFile); os.IsNotExist(err) {
				t.Error("queries file not generated")
			}
		})
	}
}

func TestGenerateMigration(t *testing.T) {
	databases := []string{"postgres", "mysql", "mariadb", "sqlite3"}

	for _, db := range databases {
		t.Run(db, func(t *testing.T) {
			projectDir := t.TempDir()
			setupProjectDir(t, projectDir, db)

			err := RunMigration("create_posts", []string{"title:string", "body:text"}, projectDir, true)
			if err != nil {
				t.Fatal(err)
			}

			migrationsDir := filepath.Join(projectDir, "cmd", "app", "sql", "migrations")
			entries, err := os.ReadDir(migrationsDir)
			if err != nil {
				t.Fatal(err)
			}

			foundMigration := false
			for _, e := range entries {
				if filepath.Ext(e.Name()) == ".sql" {
					foundMigration = true
					break
				}
			}
			if !foundMigration {
				t.Error("no migration file generated")
			}

			// Verify no other files were generated
			servicePath := filepath.Join(projectDir, "cmd", "app", "service")
			entries, _ = os.ReadDir(servicePath)
			for _, e := range entries {
				t.Errorf("unexpected file in service dir: %s", e.Name())
			}
		})
	}
}

func TestGenerateResourceNoDB(t *testing.T) {
	projectDir := t.TempDir()
	setupProjectDir(t, projectDir, "none")

	err := Run("post", []string{"title:string"}, projectDir, true)
	if err == nil {
		t.Fatal("expected error for database 'none'")
	}
}

func TestGenerateResourceNoFields(t *testing.T) {
	projectDir := t.TempDir()
	setupProjectDir(t, projectDir, "postgres")

	err := Run("post", []string{}, projectDir, true)
	if err != nil {
		t.Fatal(err)
	}

	// Should still generate files (table with only id, created_at, updated_at)
	path := filepath.Join(projectDir, "cmd", "app", "handlers", "post_handler.go")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("handler not generated for empty fields")
	}
}

func TestGenerateResourceNoConfig(t *testing.T) {
	projectDir := t.TempDir()

	err := Run("post", []string{"title:string"}, projectDir, true)
	if err == nil {
		t.Fatal("expected error when .snowflake.yaml is missing")
	}
}

func setupProjectDir(t *testing.T, projectDir string, database string) {
	t.Helper()

	// Create .snowflake.yaml
	config := "name: acme\ndatabase: " + database + "\n"
	if err := os.WriteFile(filepath.Join(projectDir, ".snowflake.yaml"), []byte(config), 0666); err != nil {
		t.Fatal(err)
	}

	// Create directories
	dirs := []string{
		"cmd/app/sql/migrations",
		"cmd/app/sql/queries",
		"cmd/app/service",
		"cmd/app/handlers",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(projectDir, d), 0777); err != nil {
			t.Fatal(err)
		}
	}
}
