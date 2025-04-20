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

type Project struct {
	Name          string
	Database      Database
	BackgroundJob BackgroundJob
	AppType       AppType
	SMTP          bool
	Storage       bool
	Auth          bool
	Redis         bool
}

type Config struct {
	Quiet         bool
	Name          string
	Database      Database
	AppType       AppType
	BackgroundJob BackgroundJob
	OutputDir     string

	NoSMTP    bool
	NoStorage bool
	NoAuth    bool
	NoGit     bool
	NoRedis   bool
}

func Run(cfg *Config) error {
	project := &Project{
		Name:          cfg.Name,
		Database:      cfg.Database,
		BackgroundJob: cfg.BackgroundJob,
		AppType:       cfg.AppType,
		SMTP:          !cfg.NoSMTP,
		Storage:       !cfg.NoStorage,
		Redis:         !cfg.NoRedis || cfg.BackgroundJob == BackgroundJobAsynq,
		Auth:          !cfg.NoAuth && !cfg.NoSMTP && cfg.Database != DatabaseNone,
	}

	outputPath := filepath.Join(cfg.OutputDir, cfg.Name)
	templateFiles := initializetemplate.BaseFiles

	templateFuncs := createTemplateFuncs(cfg)
	exclusions := createFileExclusions()
	renames := createFileRenames()

	if err := createFiles(project, outputPath, templateFiles, templateFuncs, exclusions, cfg.Quiet); err != nil {
		return err
	}

	if err := renameFiles(project, outputPath, renames); err != nil {
		return err
	}

	if err := runPostCommands(project, outputPath, cfg.Quiet); err != nil {
		return err
	}

	if !cfg.NoGit {
		if err := runGitCommands(outputPath, cfg.Quiet); err != nil {
			return err
		}
	}

	if !cfg.Quiet {
		fmt.Println("")
		successMessage := fmt.Sprintf(`✅ Snowflake project '%s' created! 🎉

Run your new project:

  $ cd %s`, project.Name, project.Name)

		if project.Database == DatabasePostgres || project.Database == DatabaseMySQL || project.Redis {
			successMessage += `
  $ make devenv # Initialize the docker dev environment
  $ make dev`
		} else {
			successMessage += `
  $ make dev`
		}

		fmt.Println(successMessage)
	}

	return nil
}

func createFiles(project *Project, outputPath string, templateFiles fs.FS,
	templateFuncs map[string]interface{}, exclusions *FileExclusions, quiet bool) error {

	if !quiet {
		fmt.Println("Generating files...")
	}

	// Use a buffer pool to reduce allocations when processing templates
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

		if shouldExcludeTemplateFile(templateFileName, project, exclusions) {
			return nil
		}

		content, err := fs.ReadFile(templateFiles, path)
		if err != nil {
			return err
		}

		tmpl, err := template.New(templateFileName).Funcs(templateFuncs).Parse(string(content))
		if err != nil {
			return err
		}

		// Get a buffer from pool and ensure it's reset
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer bufPool.Put(buf)

		if err := tmpl.Execute(buf, project); err != nil {
			return err
		}

		filePath := strings.TrimSuffix(targetPath, ".templ")
		targetDir := filepath.Dir(filePath)
		if err := os.MkdirAll(targetDir, 0777); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", targetDir, err)
		}

		return os.WriteFile(filePath, buf.Bytes(), 0666)
	})

	return err
}
