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

type GenerateInput struct {
	Name       string
	Plural     string
	RawFields  []string
	ProjectDir string
	Quiet      bool
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

func RunMigration(input GenerateInput) error {
	ctx, err := prepareGeneration(input)
	if err != nil {
		return err
	}

	migrationsDir := filepath.Join(input.ProjectDir, "cmd", "app", "sql", "migrations")
	migNum := MigrationNumber()
	_, err = renderTargets(ctx.templates, ctx.resource, []generatedTarget{
		{
			templateName: migrationTemplateName(ctx.config.Database),
			outputPath:   MigrationFilePath(migrationsDir, migNum, ctx.resource.PluralName),
		},
	}, input.ProjectDir, input.Quiet)
	return err
}

func Run(input GenerateInput) error {
	ctx, err := prepareGeneration(input)
	if err != nil {
		return err
	}

	migrationsDir := filepath.Join(input.ProjectDir, "cmd", "app", "sql", "migrations")
	migNum := MigrationNumber()
	files := []generatedTarget{
		{
			templateName: migrationTemplateName(ctx.config.Database),
			outputPath:   MigrationFilePath(migrationsDir, migNum, ctx.resource.PluralName),
		},
		{
			templateName: queriesTemplateName(ctx.config.Database),
			outputPath:   filepath.Join(input.ProjectDir, "cmd", "app", "sql", "queries", ctx.resource.PluralName+".sql"),
		},
		{
			templateName: serviceTemplateName(ctx.config.Database),
			outputPath:   filepath.Join(input.ProjectDir, "cmd", "app", "service", ctx.resource.Name+"_service.go"),
		},
		{
			templateName: "handler.go.tmpl",
			outputPath:   filepath.Join(input.ProjectDir, "cmd", "app", "handlers", ctx.resource.Name+"_handler.go"),
		},
	}

	goFiles, err := renderTargets(ctx.templates, ctx.resource, files, input.ProjectDir, input.Quiet)
	if err != nil {
		return err
	}

	if err := runGenCommand("sqlc", []string{"generate"}, filepath.Join(input.ProjectDir, "cmd", "app"), input.Quiet); err != nil {
		if !input.Quiet {
			fmt.Println("  warning: sqlc generate failed. Run it manually: cd cmd/app && sqlc generate")
		}
	}

	if len(goFiles) > 0 {
		args := append([]string{"-w", "-s"}, uniquePaths(goFiles)...)
		_ = runGenCommand("gofmt", args, input.ProjectDir, true)
	}

	if !input.Quiet {
		printSuccess(input.ProjectDir, ctx.config, ctx.resource)
	}

	return nil
}

func prepareGeneration(input GenerateInput) (*generationContext, error) {
	cfg, err := LoadConfig(input.ProjectDir)
	if err != nil {
		return nil, err
	}

	fields, err := ParseFields(input.RawFields, cfg.Database)
	if err != nil {
		return nil, err
	}

	templates, err := parseTemplates()
	if err != nil {
		return nil, err
	}

	return &generationContext{
		config:    cfg,
		resource:  NewResource(input.Name, input.Plural, fields, cfg),
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

func printSuccess(projectDir string, cfg *ProjectConfig, r *Resource) {
	fmt.Println()
	fmt.Printf("Generated %s resource with %d field(s).\n", r.Name, len(r.Fields))
	if len(r.Fields) > 0 {
		fmt.Println()
		for _, f := range r.Fields {
			fmt.Printf("    %s:%s\n", f.Name, f.Type)
		}
	}
	fmt.Printf("\n%s\n", routeInstructions(projectDir, cfg, r))
}

func routeInstructions(projectDir string, cfg *ProjectConfig, resource *Resource) string {
	content, err := os.ReadFile(filepath.Join(projectDir, "cmd", "app", "router.go"))
	if err != nil {
		content = nil
	}
	return buildRouteInstructions(string(content), cfg, resource)
}

func buildRouteInstructions(content string, cfg *ProjectConfig, resource *Resource) string {
	queriesLine := "queries := repo.New(db)"
	serviceLine := fmt.Sprintf("%sService := service.New%sService(queries)", resource.Name, resource.NameTitle)
	registerLine := fmt.Sprintf("handlers.Register%sRoutes(api, %sService)", resource.NameTitle, resource.Name)

	needsQueries := !strings.Contains(content, queriesLine)
	needsService := !strings.Contains(content, serviceLine)
	needsRegister := !strings.Contains(content, registerLine)

	var imports []string
	if needsRegister && !hasImport(content, cfg.Module+"/cmd/app/handlers") {
		imports = append(imports, fmt.Sprintf("%q", cfg.Module+"/cmd/app/handlers"))
	}
	if needsQueries && !hasImport(content, cfg.Module+"/cmd/app/repo") {
		imports = append(imports, fmt.Sprintf("%q", cfg.Module+"/cmd/app/repo"))
	}
	if needsService && !hasImport(content, cfg.Module+"/cmd/app/service") {
		imports = append(imports, fmt.Sprintf("%q", cfg.Module+"/cmd/app/service"))
	}

	var sections []string
	if len(imports) > 0 {
		sections = append(sections, "Add these imports to cmd/app/router.go:\n"+indentLines(imports))
	}

	var lines []string
	if needsQueries {
		lines = append(lines, queriesLine)
	}
	if needsService {
		lines = append(lines, serviceLine)
	}
	if needsRegister {
		lines = append(lines, registerLine)
	}
	if len(lines) > 0 {
		sections = append(sections, "Add this inside registerRoutes in cmd/app/router.go:\n"+indentLines(lines))
	}

	if len(sections) == 0 {
		return "Routes for this resource already appear to be declared in cmd/app/router.go."
	}

	return strings.Join(sections, "\n\n")
}

func hasImport(content string, path string) bool {
	return strings.Contains(content, fmt.Sprintf("%q", path))
}

func indentLines(lines []string) string {
	indented := make([]string, len(lines))
	for i, line := range lines {
		indented[i] = "    " + line
	}
	return strings.Join(indented, "\n")
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
