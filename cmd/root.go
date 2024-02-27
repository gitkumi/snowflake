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
	Name string
	Type string
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
				log.Fatal(err)
			}

			projectDatabase, err := cmd.Flags().GetString("database")
			if err != nil {
				log.Fatal(err)
			}

			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			project := &project{
				Name: projectName,
				Type: projectType,
				Database: projectDatabase,
			}

			outputPath := cwd
			if os.Getenv("ENVIRONMENT") != "production" {
				outputPath = filepath.Join(cwd, "testdata")
			}
			outputPath = filepath.Join(outputPath, projectName)

			err = fs.WalkDir(snowflaketemplate.Files, "files", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				fileName := strings.Replace(path, "files", "", 1)

				if project.Type == "api" {
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

					isWeb := false
					for _, file := range webFiles {
						if file == fileName {
							isWeb = true
							break
						}
					}

					if isWeb {
						return nil
					}
				}

				if project.Type == "web" {
					apiFiles := []string{
						"cmd/api",
						"cmd/api/main.go.templ",
					}

					isApi := false
					for _, file := range apiFiles {
						if file == fileName {
							isApi = true
							break
						}
					}

					if isApi {
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

				temp, err := template.New("test").Parse(string(content))
				if err != nil {
					return err
				}

				var buf bytes.Buffer
				err = temp.Execute(&buf, project)

				newFilePath := strings.TrimSuffix(filepath.Join(outputPath, fileName), ".templ") 
				err = os.WriteFile(newFilePath, buf.Bytes(), 0777)
				return err
			})

			if err != nil {
				log.Fatal(err)
			}
		},
	}

	// type: Web or API
	cmd.Flags().StringP("type", "t", "", "Type of the project.")

	// db: none, sqlite, postgres, mysql/mariadb
	cmd.Flags().StringP("database", "d", "", "Database of the project.")

	return cmd
}
