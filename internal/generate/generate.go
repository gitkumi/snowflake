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
	"hasParamsStruct": func(fields []Field) bool {
		return len(fields) > 1
	},
}

type generatedTarget struct {
	templateName string
	outputPath   string
}

type generationContext struct {
	config    *ProjectConfig
	resource  *Resource
	templates *template.Template
}

func RunMigration(name string, rawFields []string, projectDir string, quiet bool) error {
	ctx, err := prepareGeneration(name, rawFields, projectDir)
	if err != nil {
		return err
	}

	migrationsDir := filepath.Join(projectDir, "cmd", "app", "sql", "migrations")
	migNum := MigrationNumber()
	_, err = renderTargets(ctx.templates, ctx.resource, []generatedTarget{
		{
			templateName: migrationTemplateName(ctx.config.Database),
			outputPath:   MigrationFilePath(migrationsDir, migNum, ctx.resource.PluralName),
		},
	}, projectDir, quiet)
	return err
}

func Run(resourceName string, rawFields []string, projectDir string, quiet bool) error {
	ctx, err := prepareGeneration(resourceName, rawFields, projectDir)
	if err != nil {
		return err
	}

	migrationsDir := filepath.Join(projectDir, "cmd", "app", "sql", "migrations")
	migNum := MigrationNumber()
	files := []generatedTarget{
		{
			templateName: migrationTemplateName(ctx.config.Database),
			outputPath:   MigrationFilePath(migrationsDir, migNum, ctx.resource.PluralName),
		},
		{
			templateName: queriesTemplateName(ctx.config.Database),
			outputPath:   filepath.Join(projectDir, "cmd", "app", "sql", "queries", ctx.resource.PluralName+".sql"),
		},
		{
			templateName: serviceTemplateName(ctx.config.Database),
			outputPath:   filepath.Join(projectDir, "cmd", "app", "service", ctx.resource.Name+"_service.go"),
		},
		{
			templateName: "handler.go.tmpl",
			outputPath:   filepath.Join(projectDir, "cmd", "app", "handlers", ctx.resource.Name+"_handler.go"),
		},
	}

	goFiles, err := renderTargets(ctx.templates, ctx.resource, files, projectDir, quiet)
	if err != nil {
		return err
	}

	if routesFile, err := wireGeneratedRoutes(projectDir, ctx.config, ctx.resource); err != nil {
		return err
	} else if routesFile != "" {
		goFiles = append(goFiles, routesFile)
	}

	if serverFile, err := ensureServerUsesGeneratedRoutes(projectDir); err != nil {
		return err
	} else if serverFile != "" {
		goFiles = append(goFiles, serverFile)
	}

	if err := runGenCommand("sqlc", []string{"generate"}, filepath.Join(projectDir, "cmd", "app"), quiet); err != nil {
		if !quiet {
			fmt.Println("  warning: sqlc generate failed. Run it manually: cd cmd/app && sqlc generate")
		}
	}

	if len(goFiles) > 0 {
		args := append([]string{"-w", "-s"}, uniquePaths(goFiles)...)
		_ = runGenCommand("gofmt", args, projectDir, true)
	}

	if !quiet {
		printSuccess(ctx.resource)
	}

	return nil
}

func prepareGeneration(name string, rawFields []string, projectDir string) (*generationContext, error) {
	cfg, err := LoadConfig(projectDir)
	if err != nil {
		return nil, err
	}

	fields, err := ParseFields(rawFields, cfg.Database)
	if err != nil {
		return nil, err
	}

	templates, err := parseTemplates()
	if err != nil {
		return nil, err
	}

	return &generationContext{
		config:    cfg,
		resource:  NewResource(name, fields, cfg),
		templates: templates,
	}, nil
}

func parseTemplates() (*template.Template, error) {
	tmpl, err := template.New("").Funcs(funcMap).ParseFS(generatetemplate.Files, "*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}
	return tmpl, nil
}

func renderTargets(tmpl *template.Template, data any, targets []generatedTarget, projectDir string, quiet bool) ([]string, error) {
	var (
		buf     bytes.Buffer
		goFiles []string
	)

	for _, target := range targets {
		if err := renderTarget(tmpl, data, target, projectDir, quiet, &buf); err != nil {
			return nil, err
		}
		if strings.HasSuffix(target.outputPath, ".go") {
			goFiles = append(goFiles, target.outputPath)
		}
	}

	return goFiles, nil
}

func renderTarget(tmpl *template.Template, data any, target generatedTarget, projectDir string, quiet bool, buf *bytes.Buffer) error {
	buf.Reset()
	if err := tmpl.ExecuteTemplate(buf, target.templateName, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", target.templateName, err)
	}

	if err := os.MkdirAll(filepath.Dir(target.outputPath), 0777); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(target.outputPath, buf.Bytes(), 0666); err != nil {
		return fmt.Errorf("failed to write %s: %w", target.outputPath, err)
	}

	if !quiet {
		rel, _ := filepath.Rel(projectDir, target.outputPath)
		fmt.Printf("  created %s\n", rel)
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

func printSuccess(r *Resource) {
	fmt.Printf("\nResource %q generated and registered.\n", r.Name)
}

func wireGeneratedRoutes(projectDir string, cfg *ProjectConfig, resource *Resource) (string, error) {
	routesFile, err := ensureGeneratedRoutesFile(projectDir, cfg)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(routesFile)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", routesFile, err)
	}

	updated, err := appendBeforeMarker(string(content), "\t// snowflake:imports",
		fmt.Sprintf("\t%q\n", cfg.Module+"/cmd/app/handlers"))
	if err != nil {
		return "", err
	}
	updated, err = appendBeforeMarker(updated, "\t// snowflake:imports",
		fmt.Sprintf("\t%q\n", cfg.Module+"/cmd/app/service"))
	if err != nil {
		return "", err
	}

	routeBlock := fmt.Sprintf(
		"\t%sService := service.New%sService(queries)\n\thandlers.Register%sRoutes(api, %sService)\n",
		resource.Name,
		resource.NameTitle,
		resource.NameTitle,
		resource.Name,
	)
	updated, err = appendBeforeMarker(updated, "\t// snowflake:routes", routeBlock)
	if err != nil {
		return "", err
	}

	if updated == string(content) {
		return "", nil
	}

	if err := os.WriteFile(routesFile, []byte(updated), 0666); err != nil {
		return "", fmt.Errorf("failed to update %s: %w", routesFile, err)
	}

	return routesFile, nil
}

func ensureGeneratedRoutesFile(projectDir string, cfg *ProjectConfig) (string, error) {
	routesFile := filepath.Join(projectDir, "cmd", "app", "generated_routes.go")
	if _, err := os.Stat(routesFile); err == nil {
		return routesFile, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to stat %s: %w", routesFile, err)
	}

	content := fmt.Sprintf(`package main

import (
	"database/sql"

	"%s/cmd/app/repo"
	// snowflake:imports

	"github.com/gin-gonic/gin"
)

func registerGeneratedRoutes(api *gin.RouterGroup, db *sql.DB) {
	queries := repo.New(db)
	_ = queries
	_ = api
	// snowflake:routes
}
`, cfg.Module)

	if err := os.MkdirAll(filepath.Dir(routesFile), 0777); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(routesFile, []byte(content), 0666); err != nil {
		return "", fmt.Errorf("failed to create %s: %w", routesFile, err)
	}

	return routesFile, nil
}

func ensureServerUsesGeneratedRoutes(projectDir string) (string, error) {
	serverFile := filepath.Join(projectDir, "cmd", "app", "server.go")
	content, err := os.ReadFile(serverFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("failed to read %s: %w", serverFile, err)
	}

	if strings.Contains(string(content), "registerGeneratedRoutes(") {
		return "", nil
	}

	legacyHealthLine := "\t\tapi.GET(\"/health\", handlers.HandleHealth(s.healthService))"
	if !strings.Contains(string(content), legacyHealthLine) {
		return "", nil
	}

	updated := strings.Replace(string(content), legacyHealthLine,
		legacyHealthLine+"\n\t\tregisterGeneratedRoutes(api, s.db)", 1)

	if err := os.WriteFile(serverFile, []byte(updated), 0666); err != nil {
		return "", fmt.Errorf("failed to update %s: %w", serverFile, err)
	}

	return serverFile, nil
}

func appendBeforeMarker(content string, marker string, block string) (string, error) {
	if strings.Contains(content, block) {
		return content, nil
	}
	if !strings.Contains(content, marker) {
		return "", fmt.Errorf("marker %q not found", marker)
	}
	return strings.Replace(content, marker, block+marker, 1), nil
}

func uniquePaths(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		if _, ok := seen[path]; ok {
			continue
		}
		seen[path] = struct{}{}
		result = append(result, path)
	}
	return result
}
