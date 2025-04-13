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

	snowflaketemplate "github.com/gitkumi/snowflake/template"
)

type Project struct {
	Name     string
	Database Database
	AppType  AppType
	SMTP     bool
	Storage  bool
	Auth     bool
}

type GeneratorConfig struct {
	Name     string
	Database Database
	AppType  AppType
	SMTP     bool
	Storage  bool
	Auth     bool

	InitGit   bool
	OutputDir string
}

func InitProject() *cobra.Command {
	var (
		initGit   bool
		database  string
		appType   string
		outputDir string
		smtp      bool
		storage   bool
		auth      bool
	)

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// If no output directory is specified, use the current working directory
			if outputDir == "" {
				var err error
				outputDir, err = os.Getwd()
				if err != nil {
					log.Fatal(err.Error())
				}
			} else {
				// If the provided path is not absolute, make it absolute
				if !filepath.IsAbs(outputDir) {
					cwd, err := os.Getwd()
					if err != nil {
						log.Fatal(err.Error())
					}
					outputDir = filepath.Join(cwd, outputDir)
				}

				// Ensure output directory exists
				if _, err := os.Stat(outputDir); os.IsNotExist(err) {
					if err := os.MkdirAll(outputDir, 0755); err != nil {
						log.Fatalf("Failed to create output directory: %v", err)
					}
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

			cfg := &GeneratorConfig{
				Name:      args[0],
				Database:  dbEnum,
				AppType:   appTypeEnum,
				InitGit:   initGit,
				OutputDir: outputDir,
				SMTP:      smtp,
				Storage:   storage,
			}

			err := Generate(cfg)
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}

	cmd.Flags().StringVarP(&appType, "appType", "t", "api", fmt.Sprintf("App type %v", AllAppTypes))
	cmd.Flags().StringVarP(&database, "database", "d", "sqlite3", fmt.Sprintf("Database type %v", AllDatabases))
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for the generated project")
	cmd.Flags().BoolVar(&initGit, "git", true, "Initialize git")
	cmd.Flags().BoolVar(&smtp, "smtp", true, "Add smtp feature")
	cmd.Flags().BoolVar(&storage, "storage", true, "Add storage feature (S3)")
	cmd.Flags().BoolVar(&auth, "auth", true, "Add authentication feature")

	return cmd
}

func Generate(cfg *GeneratorConfig) error {
	project := &Project{
		Name:     cfg.Name,
		Database: cfg.Database,
		AppType:  cfg.AppType,
		SMTP:     cfg.SMTP,
		Storage:  cfg.Storage,
		Auth:     cfg.Auth && cfg.SMTP,
	}

	outputPath := filepath.Join(cfg.OutputDir, cfg.Name)
	templateFiles := snowflaketemplate.BaseFiles

	templateFuncs := CreateTemplateFuncs(cfg)
	exclusions := CreateFileExclusions()
	renames := CreateFileRenames()

	if err := generateFromTemplates(project, outputPath, templateFiles, templateFuncs, exclusions); err != nil {
		return err
	}

	if err := ProcessFileRenames(project, outputPath, renames); err != nil {
		return err
	}

	if err := RunPostCommands(project, outputPath); err != nil {
		return err
	}

	if cfg.InitGit {
		if err := RunGitCommands(outputPath); err != nil {
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

func generateFromTemplates(project *Project, outputPath string, templateFiles fs.FS,
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

		if ShouldExcludeTemplateFile(templateFileName, project, exclusions) {
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
