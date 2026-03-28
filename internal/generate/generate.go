package generate

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	generatetemplate "github.com/gitkumi/snowflake/internal/generate/template"
)

var funcMap = template.FuncMap{
	"fieldNames": func(fields []Field) string {
		names := make([]string, len(fields))
		for i, f := range fields {
			names[i] = f.Name
		}
		return strings.Join(names, ", ")
	},
	"postgresParams": func(fields []Field, start int) string {
		params := make([]string, len(fields))
		for i := range fields {
			params[i] = fmt.Sprintf("$%d", start+i)
		}
		return strings.Join(params, ", ")
	},
	"postgresSetClauses": func(fields []Field, start int) string {
		clauses := make([]string, len(fields))
		for i, f := range fields {
			clauses[i] = fmt.Sprintf("%s = $%d", f.Name, start+i)
		}
		return strings.Join(clauses, ",\n    ")
	},
	"postgresNextParam": func(fields []Field, start int) string {
		return fmt.Sprintf("$%d", start+len(fields))
	},
	"questionParams": func(fields []Field) string {
		params := make([]string, len(fields))
		for i := range fields {
			params[i] = "?"
		}
		return strings.Join(params, ", ")
	},
	"questionSetClauses": func(fields []Field) string {
		clauses := make([]string, len(fields))
		for i, f := range fields {
			clauses[i] = fmt.Sprintf("%s = ?", f.Name)
		}
		return strings.Join(clauses, ",\n    ")
	},
}

func RunMigration(name string, rawFields []string, projectDir string, quiet bool) error {
	cfg, err := LoadConfig(projectDir)
	if err != nil {
		return err
	}

	fields, err := ParseFields(rawFields, cfg.Database)
	if err != nil {
		return err
	}

	resource := NewResource(name, fields, cfg)

	migrationsDir := filepath.Join(projectDir, "cmd", "app", "sql", "migrations")
	migNum := MigrationNumber()

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(generatetemplate.Files, "*.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	outputPath := MigrationFilePath(migrationsDir, migNum, resource.NamePlural)

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, migrationTemplateName(cfg.Database), resource); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0777); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0666); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputPath, err)
	}

	if !quiet {
		rel, _ := filepath.Rel(projectDir, outputPath)
		fmt.Printf("  created %s\n", rel)
	}

	return nil
}

func Run(resourceName string, rawFields []string, projectDir string, quiet bool) error {
	cfg, err := LoadConfig(projectDir)
	if err != nil {
		return err
	}

	fields, err := ParseFields(rawFields, cfg.Database)
	if err != nil {
		return err
	}

	resource := NewResource(resourceName, fields, cfg)

	migrationsDir := filepath.Join(projectDir, "cmd", "app", "sql", "migrations")
	migNum := MigrationNumber()

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(generatetemplate.Files, "*.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	files := []struct {
		templateName string
		outputPath   string
	}{
		{
			templateName: migrationTemplateName(cfg.Database),
			outputPath:   MigrationFilePath(migrationsDir, migNum, resource.NamePlural),
		},
		{
			templateName: queriesTemplateName(cfg.Database),
			outputPath:   filepath.Join(projectDir, "cmd", "app", "sql", "queries", resource.NamePlural+".sql"),
		},
		{
			templateName: serviceTemplateName(cfg.Database),
			outputPath:   filepath.Join(projectDir, "cmd", "app", "service", resource.Name+"_service.go"),
		},
		{
			templateName: "handler.go.tmpl",
			outputPath:   filepath.Join(projectDir, "cmd", "app", "handlers", resource.Name+"_handler.go"),
		},
	}

	var buf bytes.Buffer
	for _, f := range files {
		buf.Reset()
		if err := tmpl.ExecuteTemplate(&buf, f.templateName, resource); err != nil {
			return fmt.Errorf("failed to execute template %s: %w", f.templateName, err)
		}

		if err := os.MkdirAll(filepath.Dir(f.outputPath), 0777); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		if err := os.WriteFile(f.outputPath, buf.Bytes(), 0666); err != nil {
			return fmt.Errorf("failed to write %s: %w", f.outputPath, err)
		}

		if !quiet {
			rel, _ := filepath.Rel(projectDir, f.outputPath)
			fmt.Printf("  created %s\n", rel)
		}
	}

	if err := runGenCommand("sqlc", []string{"generate"}, filepath.Join(projectDir, "cmd", "app"), quiet); err != nil {
		if !quiet {
			fmt.Println("  warning: sqlc generate failed. Run it manually: cd cmd/app && sqlc generate")
		}
	}

	goFiles := []string{}
	for _, f := range files {
		if strings.HasSuffix(f.outputPath, ".go") {
			goFiles = append(goFiles, f.outputPath)
		}
	}
	if len(goFiles) > 0 {
		args := append([]string{"-w", "-s"}, goFiles...)
		_ = runGenCommand("gofmt", args, projectDir, true)
	}

	if !quiet {
		printInstructions(resource)
	}

	return nil
}

func migrationTemplateName(database string) string {
	switch database {
	case "mysql", "mariadb":
		return "migration_mysql.sql.tmpl"
	case "sqlite3":
		return "migration_sqlite3.sql.tmpl"
	default:
		return "migration_postgres.sql.tmpl"
	}
}

func queriesTemplateName(database string) string {
	switch database {
	case "mysql", "mariadb":
		return "queries_mysql.sql.tmpl"
	case "sqlite3":
		return "queries_sqlite3.sql.tmpl"
	default:
		return "queries_postgres.sql.tmpl"
	}
}

func serviceTemplateName(database string) string {
	switch database {
	case "mysql", "mariadb":
		return "service_refetch.go.tmpl"
	default:
		return "service_returning.go.tmpl"
	}
}

func runGenCommand(name string, args []string, dir string, quiet bool) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("%s is not installed or not found in PATH", name)
	}

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if quiet {
		cmd.Stdout = io.Discard
	} else {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func printInstructions(r *Resource) {
	fmt.Printf("\nResource %q generated!\n", r.Name)
	fmt.Printf("\nAdd the following to cmd/app/server.go:\n")
	fmt.Printf("\n  1. Add import (if not already present):\n")
	fmt.Printf("       \"%s/cmd/app/repo\"\n", r.ModuleName)
	fmt.Printf("\n  2. Add to server struct:\n")
	fmt.Printf("       %sService *service.%sService\n", r.Name, r.NameTitle)
	fmt.Printf("\n  3. Add to newServer():\n")
	fmt.Printf("       queries := repo.New(db)     // if not already present\n")
	fmt.Printf("       %sService := service.New%sService(queries)\n", r.Name, r.NameTitle)
	fmt.Printf("\n  4. Add to routes():\n")
	fmt.Printf("       api.GET(\"/%s\", handlers.HandleList%s(s.%sService))\n", r.NamePlural, r.NameTitlePlural, r.Name)
	fmt.Printf("       api.GET(\"/%s/:id\", handlers.HandleGet%s(s.%sService))\n", r.NamePlural, r.NameTitle, r.Name)
	fmt.Printf("       api.POST(\"/%s\", handlers.HandleCreate%s(s.%sService))\n", r.NamePlural, r.NameTitle, r.Name)
	if len(r.Fields) > 0 {
		fmt.Printf("       api.PATCH(\"/%s/:id\", handlers.HandleUpdate%s(s.%sService))\n", r.NamePlural, r.NameTitle, r.Name)
	}
	fmt.Printf("       api.DELETE(\"/%s/:id\", handlers.HandleDelete%s(s.%sService))\n", r.NamePlural, r.NameTitle, r.Name)
}
