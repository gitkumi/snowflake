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

func (p *Project) WithAuth() bool {
	return p.Authentication != AuthenticationNone
}

func (p *Project) WithOAuth() bool {
	return p.WithAuth() && (p.OAuthGoogle || p.OAuthDiscord || p.OAuthFacebook || p.OAuthGitHub || p.OAuthInstagram || p.OAuthLinkedIn)
}

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
	project := &Project{
		Name:           cfg.Name,
		Database:       cfg.Database,
		BackgroundJob:  cfg.BackgroundJob,
		AppType:        cfg.AppType,
		SMTP:           cfg.SMTP,
		Storage:        cfg.Storage,
		Redis:          cfg.Redis || cfg.BackgroundJob == BackgroundJobAsynq,
		Authentication: cfg.Authentication,
		OAuthDiscord:   cfg.Authentication != AuthenticationNone && cfg.OAuthDiscord,
		OAuthFacebook:  cfg.Authentication != AuthenticationNone && cfg.OAuthFacebook,
		OAuthGitHub:    cfg.Authentication != AuthenticationNone && cfg.OAuthGitHub,
		OAuthGoogle:    cfg.Authentication != AuthenticationNone && cfg.OAuthGoogle,
		OAuthInstagram: cfg.Authentication != AuthenticationNone && cfg.OAuthInstagram,
		OAuthLinkedIn:  cfg.Authentication != AuthenticationNone && cfg.OAuthLinkedIn,
	}

	if cfg.Database == DatabaseNone || !cfg.SMTP {
		project.Authentication = AuthenticationNone
		project.OAuthDiscord = false
		project.OAuthFacebook = false
		project.OAuthGitHub = false
		project.OAuthGoogle = false
		project.OAuthInstagram = false
		project.OAuthLinkedIn = false
	}

	outputPath := filepath.Join(cfg.OutputDir, cfg.Name)
	templateFiles := initializetemplate.BaseFiles

	exclusions := createFileExclusions()
	renames := createFileRenames()

	databaseFragments, err := initializetemplate.CreateDatabaseFragments(string(project.Database))
	if err != nil {
		return err
	}

	if err := createFiles(project, outputPath, templateFiles, exclusions, databaseFragments, cfg.Quiet); err != nil {
		return err
	}

	if err := renameFiles(project, outputPath, renames); err != nil {
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
	exclusions *FileExclusions, databaseFragments map[string]string, quiet bool) error {

	if !quiet {
		fmt.Println("Generating files...")
	}

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

		rootTemplate := template.New(filepath.Base(templateFileName))

		for name, fragment := range databaseFragments {
			fragmentTemplate := rootTemplate.New(name)
			if _, err := fragmentTemplate.Parse(fragment); err != nil {
				return fmt.Errorf("failed to parse database fragment %s: %w", name, err)
			}
		}

		if _, err := rootTemplate.Parse(string(content)); err != nil {
			return fmt.Errorf("failed to parse template %s: %w", templateFileName, err)
		}

		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer bufPool.Put(buf)

		if err := rootTemplate.Execute(buf, project); err != nil {
			return fmt.Errorf("error executing template %s: %w", templateFileName, err)
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
