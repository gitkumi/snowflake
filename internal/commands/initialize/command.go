package initialize

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/spf13/cobra"

	initializetemplate "github.com/gitkumi/snowflake/internal/commands/initialize/template"
)

type Project struct {
	Name     string
	Database Database
	AppType  AppType
	SMTP     bool
	Storage  bool
	Auth     bool
}

type InitConfig struct {
	Name      string
	Database  Database
	AppType   AppType
	OutputDir string

	NoSMTP    bool
	NoStorage bool
	NoAuth    bool
	NoGit     bool
}

func InitProject() *cobra.Command {
	var (
		database  string
		appType   string
		outputDir string
		noGit     bool
		noSMTP    bool
		noStorage bool
		noAuth    bool
	)

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Use current working directory if output dir is not provided
			if outputDir == "" {
				cwd, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				outputDir = cwd
			}

			// If the provided path is not absolute, make it absolute
			if !filepath.IsAbs(outputDir) {
				cwd, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				outputDir = filepath.Join(cwd, outputDir)
			}

			// Ensure output directory exists
			if _, err := os.Stat(outputDir); os.IsNotExist(err) {
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					log.Fatalf("Failed to create output directory: %v", err)
				}
			}

			dbEnum := Database(database)
			if !dbEnum.IsValid() {
				log.Fatalf("Invalid database type: %s. Must be one of: %v", database, AllDatabases)
			}

			appTypeEnum := AppType(appType)
			if !appTypeEnum.IsValid() {
				log.Fatalf("Invalid app type: %s. Must be one of: %v", appType, AllAppTypes)
			}

			cfg := &InitConfig{
				Name:      args[0],
				Database:  dbEnum,
				AppType:   appTypeEnum,
				NoGit:     noGit,
				OutputDir: outputDir,
				NoSMTP:    noSMTP,
				NoStorage: noStorage,
			}

			err := Initialize(cfg)
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}

	cmd.Flags().StringVarP(&appType, "appType", "t", "api", fmt.Sprintf("App type %v", AllAppTypes))
	cmd.Flags().StringVarP(&database, "database", "d", "sqlite3", fmt.Sprintf("Database type %v", AllDatabases))
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for the generated project")
	cmd.Flags().BoolVar(&noGit, "no-git", false, "Remove git")
	cmd.Flags().BoolVar(&noSMTP, "no-smtp", false, "Remove SMTP")
	cmd.Flags().BoolVar(&noStorage, "no-storage", false, "Remove Storage (S3)")
	cmd.Flags().BoolVar(&noAuth, "no-auth", false, "Remove Authentication (Authentication requires SMTP)")

	return cmd
}

func Initialize(cfg *InitConfig) error {
	project := &Project{
		Name:     cfg.Name,
		Database: cfg.Database,
		AppType:  cfg.AppType,
		SMTP:     !cfg.NoSMTP,
		Storage:  !cfg.NoStorage,
		Auth:     !cfg.NoAuth && !cfg.NoSMTP,
	}

	outputPath := filepath.Join(cfg.OutputDir, cfg.Name)
	templateFiles := initializetemplate.BaseFiles

	templateFuncs := createTemplateFuncs(cfg)
	exclusions := createFileExclusions()
	renames := createFileRenames()

	if err := createFiles(project, outputPath, templateFiles, templateFuncs, exclusions); err != nil {
		return err
	}

	if err := renameFiles(project, outputPath, renames); err != nil {
		return err
	}

	if err := runPostCommands(project, outputPath); err != nil {
		return err
	}

	if !cfg.NoGit {
		if err := runGitCommands(outputPath); err != nil {
			return err
		}
	}

	fmt.Println("")
	successMessage := fmt.Sprintf(`✅ Snowflake project '%s' generated successfully! 🎉

Run your new project:

  $ cd %s`, project.Name, project.Name)

	if project.Database == Postgres || project.Database == MySQL {
		successMessage += `
  $ make db.init  # Initialize the docker database dev environment
  $ make dev`
	} else {
		successMessage += `
  $ make dev`
	}

	fmt.Println(successMessage)

	return nil
}

func createFiles(project *Project, outputPath string, templateFiles fs.FS,
	templateFuncs map[string]interface{}, exclusions *FileExclusions) error {

	fmt.Println("Generating files...")

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
