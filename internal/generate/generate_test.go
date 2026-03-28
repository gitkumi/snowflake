package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateResource(t *testing.T) {
	databases := []string{"postgres", "mysql", "sqlite3"}

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
				"cmd/app/generated_routes.go",
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

			routesFile := filepath.Join(projectDir, "cmd", "app", "generated_routes.go")
			routesContent, err := os.ReadFile(routesFile)
			if err != nil {
				t.Fatalf("failed to read generated routes file: %v", err)
			}
			if !strings.Contains(string(routesContent), "handlers.RegisterPostRoutes(api, postService)") {
				t.Errorf("generated routes file did not register the resource:\n%s", routesContent)
			}
		})
	}
}

func TestGenerateMigration(t *testing.T) {
	databases := []string{"postgres", "mysql", "sqlite3"}

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

	// go.mod only, no sqlc.yaml
	goMod := "module acme\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(goMod), 0666); err != nil {
		t.Fatal(err)
	}

	err := Run("post", []string{"title:string"}, projectDir, true)
	if err == nil {
		t.Fatal("expected error when sqlc.yaml is missing")
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

func TestGenerateResourceNoGoMod(t *testing.T) {
	projectDir := t.TempDir()

	err := Run("post", []string{"title:string"}, projectDir, true)
	if err == nil {
		t.Fatal("expected error when go.mod is missing")
	}
}

func setupProjectDir(t *testing.T, projectDir string, database string) {
	t.Helper()

	// Map database to sqlc engine name
	engineMap := map[string]string{
		"postgres": "postgresql",
		"mysql":    "mysql",
		"sqlite3":  "sqlite",
	}
	engine := engineMap[database]

	// Create go.mod
	goMod := "module acme\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(goMod), 0666); err != nil {
		t.Fatal(err)
	}

	// Create sqlc.yaml
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

	sqlcYaml := fmt.Sprintf(`version: "2"
sql:
- engine: "%s"
  queries: "./sql/queries/"
  schema: "./sql/migrations/"
  gen:
    go:
      package: "repo"
      out: "./repo"
`, engine)
	if err := os.WriteFile(filepath.Join(projectDir, "cmd", "app", "sqlc.yaml"), []byte(sqlcYaml), 0666); err != nil {
		t.Fatal(err)
	}
}
