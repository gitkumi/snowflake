package cmd

import (
	"bytes"
	"fmt"
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

type Project struct {
	Name string
}

func Execute() {
	cmd := &cobra.Command{
		Use:   "snowflake",
		Short: "Snowflake is an opinionated Go REST API application generator.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.Root().CompletionOptions.DisableDefaultCmd = true

	cmd.AddCommand(new())
	cmd.AddCommand(version())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func version() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version",
		Run: func(_cmd *cobra.Command, _args []string) {
			fmt.Println("v0.9.1")
		},
	}

	return cmd
}

func new() *cobra.Command {
	var initGit bool

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]

			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err.Error())
			}

			project := &Project{
				Name: strings.ToLower(projectName),
			}

			outputPath := filepath.Join(cwd, projectName)

			templateFiles := snowflaketemplate.ApiFiles

			err = fs.WalkDir(templateFiles, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				fileName := strings.TrimPrefix(path, "api")

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

			command = exec.Command("make", "build")
			command.Dir = outputPath
			err = command.Run()
			if err != nil {
				log.Fatal(err.Error())
			}

			if initGit {
				command = exec.Command("git", "init")
				command.Dir = outputPath
				err = command.Run()
				if err != nil {
					log.Fatal(err.Error())
				}

				command = exec.Command("git", "add", "-A")
				command.Dir = outputPath
				err = command.Run()
				if err != nil {
					log.Fatal(err.Error())
				}

				command = exec.Command("git", "commit", "-m", "Initialize Snowflake project")
				command.Dir = outputPath
				err = command.Run()
				if err != nil {
					log.Fatal(err.Error())
				}
			}

			fmt.Println("Initialized Snowflake project.")
		},
	}

	cmd.Flags().BoolVarP(&initGit, "git", "g", true, "Initialize git")

	return cmd
}
