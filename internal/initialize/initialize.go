package initialize

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	initializetemplate "github.com/gitkumi/snowflake/internal/initialize/template"
)

type Config struct {
	Quiet     bool
	OutputDir string
	Git       bool

	Name           string
	AppType        AppType
	Database       Database
	BackgroundJob  BackgroundJob
	Authentication Authentication

	SMTP           bool
	Storage        bool
	Redis          bool
	OAuthDiscord   bool
	OAuthFacebook  bool
	OAuthGitHub    bool
	OAuthGoogle    bool
	OAuthInstagram bool
	OAuthLinkedIn  bool
}

func Run(cfg *Config) error {
	project := NewProject(cfg)
	outputPath := filepath.Join(cfg.OutputDir, cfg.Name)

	templateFiles := initializetemplate.BaseFiles
	exclusions := NewFileExclusions()
	renames := NewFileRenames()

	databaseFragments, err := initializetemplate.CreateDatabaseFragments(string(project.Database))
	if err != nil {
		return err
	}

	if err := createFiles(project, outputPath, templateFiles, exclusions, databaseFragments, cfg.Quiet); err != nil {
		return err
	}

	if err := RenameFiles(project, outputPath, renames); err != nil {
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

	printSuccessMessage(project.Name, project.Database, project.Redis, cfg.Quiet)

	return nil
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
	exclusions *FileExclusions, databaseFragments map[string]string, quiet bool) error {

	if !quiet {
		fmt.Println("Generating files...")
	}

	// Create a buffer pool for template rendering
	bufPool := sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
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

		if ExcludeTemplateFile(templateFileName, project, exclusions) {
			return nil
		}

		content, err := fs.ReadFile(templateFiles, path)
		if err != nil {
			return err
		}

		buf := bufPool.Get().(*bytes.Buffer)
		defer bufPool.Put(buf)

		processedContent, err := processTemplate(content, templateFileName, databaseFragments, project, buf)
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

func printSuccessMessage(projectName string, database Database, redis bool, quiet bool) {
	if quiet {
		return
	}

	fmt.Println("")
	successMessage := fmt.Sprintf(`✅ Snowflake project '%s' created! 🎉

Run your new project:

  $ cd %s`, projectName, projectName)

	if database == DatabasePostgres || database == DatabaseMySQL || redis {
		successMessage += `
  $ make devenv # Initialize the docker dev environment
  $ make dev`
	} else {
		successMessage += `
  $ make dev`
	}

	fmt.Println(successMessage)
}
