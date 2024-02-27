package cmd

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	snowflaketemplate "snowflake/template"

	"github.com/spf13/cobra"
)

type project struct {
	Name     string
	Type     string
	Database string
}

func Execute() {
	cmd := &cobra.Command{
		Use:   "snowflake",
		Short: "Snowflake is an opinionated Go application generator.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.Root().CompletionOptions.DisableDefaultCmd = true

	cmd.AddCommand(new())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func new() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]
			projectType, err := cmd.Flags().GetString("type")
			if err != nil {
				log.Fatal(err.Error())
			}

			projectDatabase, err := cmd.Flags().GetString("database")
			if err != nil {
				log.Fatal(err.Error())
			}

			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err.Error())
			}

			project := &project{
				Name:     projectName,
				Type:     projectType,
				Database: projectDatabase,
			}

			outputPath := filepath.Join(cwd, projectName)
			if os.Getenv("ENVIRONMENT") != "production" {
				outputPath = filepath.Join(cwd, "testdata", projectName)
			}

			err = fs.WalkDir(snowflaketemplate.Files, "files", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				fileName := strings.TrimPrefix(path, "files")

				switch project.Type {
				case "web":
					if isWebFile(fileName) {
						return nil
					}
				case "api":
					if isAPIFile(fileName) {
						return nil
					}
				}

				if d.IsDir() {
					err := os.MkdirAll(filepath.Join(outputPath, fileName), 0777)
					return err
				}

				content, err := snowflaketemplate.Files.ReadFile(path)
				if err != nil {
					return err
				}

				tmpl, err := template.New(fileName).Parse(string(content))
				if err != nil {
					return err
				}

				var buf bytes.Buffer
				err = tmpl.Execute(&buf, project)
				if err != nil {
					return err
				}

				newFilePath := strings.TrimSuffix(filepath.Join(outputPath, fileName), ".templ")
				err = os.WriteFile(newFilePath, buf.Bytes(), 0777)
				return err
			})

			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}

	cmd.Flags().StringP("type", "t", "web", "Type of the project.")
	cmd.Flags().StringP("database", "d", "sqlite3", "Database of the project.")

	return cmd
}

func isWebFile(fileName string) bool {
	webFiles := []string{
		"/cmd/web",
		"/cmd/web/main.go.templ",
		"/tygo.yaml.templ",
		"/package.json.templ",
		"/tailwind.config.js.templ",
		"/internal/pages",
		"/internal/pages/error.templ.templ",
		"/internal/pages/home.templ.templ",
		"/internal/pages/home.ts.templ",
		"/internal/gintemplrenderer",
		"/internal/gintemplrenderer/renderer.go.templ",
		"/static",
		"/static/public",
		"/static/public/assets",
		"/static/public/assets/home.js.templ",
		"/static/public/assets/style.css.templ",
	}
	return contains(webFiles, fileName)
}

func isAPIFile(fileName string) bool {
	apiFiles := []string{
		"/cmd/api",
		"/cmd/api/main.go.templ",
	}
	return contains(apiFiles, fileName)
}

func contains(files []string, fileName string) bool {
	for _, file := range files {
		if file == fileName {
			return true
		}
	}
	return false
}
