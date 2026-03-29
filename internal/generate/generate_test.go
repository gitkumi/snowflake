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

			err := Run(GenerateInput{
				Name:       "post",
				Plural:     "posts",
				RawFields:  []string{"title:string", "body:text", "published:bool"},
				ProjectDir: projectDir,
				Quiet:      true,
			})
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

			routesFile := filepath.Join(projectDir, "cmd", "app", "generated_routes.go")
			if _, err := os.Stat(routesFile); !os.IsNotExist(err) {
				t.Errorf("unexpected generated routes file created: %s", routesFile)
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

			err := RunMigration(GenerateInput{
				Name:       "create_posts",
				Plural:     "posts",
				RawFields:  []string{"title:string", "body:text"},
				ProjectDir: projectDir,
				Quiet:      true,
			})
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

	err := Run(GenerateInput{
		Name:       "post",
		Plural:     "posts",
		RawFields:  []string{"title:string"},
		ProjectDir: projectDir,
		Quiet:      true,
	})
	if err == nil {
		t.Fatal("expected error when sqlc.yaml is missing")
	}
}

func TestGenerateResourceNoFields(t *testing.T) {
	projectDir := t.TempDir()
	setupProjectDir(t, projectDir, "postgres")

	err := Run(GenerateInput{
		Name:       "post",
		Plural:     "posts",
		ProjectDir: projectDir,
		Quiet:      true,
	})
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

	err := Run(GenerateInput{
		Name:       "post",
		Plural:     "posts",
		RawFields:  []string{"title:string"},
		ProjectDir: projectDir,
		Quiet:      true,
	})
	if err == nil {
		t.Fatal("expected error when go.mod is missing")
	}
}

func TestRouteInstructions(t *testing.T) {
	projectDir := t.TempDir()
	setupProjectDir(t, projectDir, "postgres")

	cfg, err := LoadConfig(projectDir)
	if err != nil {
		t.Fatal(err)
	}

	resource := NewResource("post", "posts", nil, cfg)
	instructions := routeInstructions(projectDir, cfg, resource)

	for _, want := range []string{
		`Add these imports to cmd/app/routes.go:`,
		`"acme/cmd/app/handlers"`,
		`"acme/cmd/app/repo"`,
		`"acme/cmd/app/service"`,
		`queries := repo.New(db)`,
		`postService := service.NewPostService(queries)`,
		`handlers.RegisterPostRoutes(api, postService)`,
	} {
		if !strings.Contains(instructions, want) {
			t.Errorf("expected route instructions to contain %q, got:\n%s", want, instructions)
		}
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

	routesGo := `package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func registerRoutes(api *gin.RouterGroup, db *sql.DB) {
	_ = api
	_ = db
}
`
	if err := os.WriteFile(filepath.Join(projectDir, "cmd", "app", "routes.go"), []byte(routesGo), 0666); err != nil {
		t.Fatal(err)
	}
}
