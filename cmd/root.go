package cmd

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	snowflaketemplate "github.com/gitkumi/snowflake/template"

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
		Short: "Snowflake is an opinionated Go web application generator.",
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
				Name:     strings.ToLower(projectName),
				Type:     strings.ToLower(projectType),
				Database: strings.ToLower(projectDatabase),
			}

			outputPath := filepath.Join(cwd, projectName)

			templateFiles := snowflaketemplate.WebFiles
			if project.Type == "api" {
				templateFiles = snowflaketemplate.ApiFiles
			}

			err = fs.WalkDir(templateFiles, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				fileName := strings.TrimPrefix(path, project.Type)

				if d.IsDir() {
					err := os.MkdirAll(filepath.Join(outputPath, fileName), 0777)
					return err
				}

				content, err := templateFiles.ReadFile(path)
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

			command := exec.Command("go", "mod", "init", project.Name)
			command.Dir = outputPath
			err = command.Run()
			if err != nil {
				log.Fatal(err.Error())
			}

			command = exec.Command("go", "mod", "tidy")
			command.Dir = outputPath
			err = command.Run()
			if err != nil {
				log.Fatal(err.Error())
			}

			command = exec.Command("gofmt", "-w", "-s", ".")
			command.Dir = outputPath
			err = command.Run()
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}

	cmd.Flags().StringP("type", "t", "web", "Type of the project. \"web\" or \"api\".")
	cmd.Flags().StringP("database", "d", "sqlite3", "Database of the project. \"sqlite3\", \"postgres\", or \"mysql\".")

	return cmd
}

func contains(files []string, fileName string) bool {
	for _, file := range files {
		if file == fileName {
			return true
		}
	}
	return false
}
