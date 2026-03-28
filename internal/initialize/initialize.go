package initialize

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	initializetemplate "github.com/gitkumi/snowflake/internal/initialize/template"
)

type Config struct {
	Quiet     bool
	OutputDir string
	Git       bool

	Name             string
	Database         Database
	ContainerRuntime ContainerRuntime

	SMTP          bool
	Storage       bool
	KeyValueStore KeyValueStore
	Templ         bool
}

// Generate creates the project files without running any external commands.
func Generate(cfg *Config) error {
	project, outputPath, err := prepare(cfg)
	if err != nil {
		return err
	}

	databaseFragments, err := initializetemplate.CreateDatabaseFragments(string(project.Database))
	if err != nil {
		return err
	}

	return createFiles(project, outputPath, initializetemplate.BaseFiles, databaseFragments, cfg.Quiet)
}

// Finalize runs post-generation commands for a previously generated project.
func Finalize(cfg *Config) error {
	project, outputPath, err := prepare(cfg)
	if err != nil {
		return err
	}

	if err := runPostCommands(project, outputPath, cfg.Quiet); err != nil {
		return err
	}

	if cfg.Git {
		if err := runGitCommands(outputPath, cfg.Quiet); err != nil {
			return err
		}
	}

	printSuccessMessage(project.Name, project.Database, project.HasKeyValueStore(), cfg.Quiet)

	return nil
}

func Run(cfg *Config) error {
	if err := Generate(cfg); err != nil {
		return err
	}

	return Finalize(cfg)
}

func prepare(cfg *Config) (*Project, string, error) {
	if err := normalizeConfig(cfg); err != nil {
		return nil, "", err
	}

	project := NewProject(cfg)
	outputPath := filepath.Join(cfg.OutputDir, cfg.Name)

	return project, outputPath, nil
}

func processTemplate(templateContent []byte, templateFileName string,
	databaseFragments map[string]string, project *Project, buf *bytes.Buffer) ([]byte, error) {

	rootTemplate := template.New(filepath.Base(templateFileName))

	// Add database fragments as sub-templates
	for name, fragment := range databaseFragments {
		fragmentTemplate := rootTemplate.New(name)
		if _, err := fragmentTemplate.Parse(fragment); err != nil {
			return nil, fmt.Errorf("failed to parse database fragment %s: %w", name, err)
		}
	}

	// Parse the main template
	if _, err := rootTemplate.Parse(string(templateContent)); err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", templateFileName, err)
	}

	// Execute the template with the project data
	buf.Reset()
	if err := rootTemplate.Execute(buf, project); err != nil {
		return nil, fmt.Errorf("error executing template %s: %w", templateFileName, err)
	}

	return buf.Bytes(), nil
}

func createFiles(project *Project, outputPath string, templateFiles fs.FS,
	databaseFragments map[string]string, quiet bool) error {

	if !quiet {
		fmt.Println("Generating files...")
	}

	err := fs.WalkDir(templateFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		templateFileName := strings.TrimPrefix(path, "base")
		targetPath := filepath.Join(outputPath, templateFileName)

		if project.ExcludeFile(templateFileName) {
			return nil
		}

		content, err := fs.ReadFile(templateFiles, path)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		processedContent, err := processTemplate(content, templateFileName, databaseFragments, project, &buf)
		if err != nil {
			return err
		}

		filePath := strings.TrimSuffix(targetPath, ".templ")
		targetDir := filepath.Dir(filePath)
		if err := os.MkdirAll(targetDir, 0777); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", targetDir, err)
		}

		return os.WriteFile(filePath, processedContent, 0666)
	})

	return err
}

func normalizeConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}

	cfg.Name = strings.TrimSpace(cfg.Name)
	if cfg.Name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if cfg.Database == "" {
		cfg.Database = DatabaseNone
	}
	if cfg.KeyValueStore == "" {
		cfg.KeyValueStore = KeyValueStoreNone
	}
	if cfg.ContainerRuntime == "" {
		cfg.ContainerRuntime = ContainerRuntimePodman
	}

	if !cfg.Database.IsValid() {
		return fmt.Errorf("invalid database type: %s. Must be one of: %v", cfg.Database, AllDatabases)
	}
	if !cfg.KeyValueStore.IsValid() {
		return fmt.Errorf("invalid key-value store: %s. Must be one of: %v", cfg.KeyValueStore, AllKeyValueStores)
	}
	if !cfg.ContainerRuntime.IsValid() {
		return fmt.Errorf("invalid container runtime: %s. Must be one of: %v", cfg.ContainerRuntime, AllContainerRuntimes)
	}

	if strings.TrimSpace(cfg.OutputDir) == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		cfg.OutputDir = cwd
	} else if cfg.OutputDir = filepath.Clean(cfg.OutputDir); !filepath.IsAbs(cfg.OutputDir) {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		cfg.OutputDir = filepath.Join(cwd, cfg.OutputDir)
	}

	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", cfg.OutputDir, err)
	}

	return nil
}

func printSuccessMessage(projectName string, database Database, hasKeyValueStore bool, quiet bool) {
	if quiet {
		return
	}

	fmt.Println("")
	successMessage := fmt.Sprintf(`✅ Snowflake project '%s' created! 🎉

Run your new project:

  $ cd %s`, projectName, projectName)

	if database == DatabasePostgres || database == DatabaseMySQL || database == DatabaseMariaDB || hasKeyValueStore {
		successMessage += `
  $ make app.devenv.up # Initialize the dev environment
  $ make app.dev`
	} else {
		successMessage += `
  $ make app.dev`
	}

	fmt.Println(successMessage)
}
