package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"snowflake/template"

	"github.com/spf13/cobra"
)

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

			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			outputPath := cwd
			if os.Getenv("ENVIRONMENT") == "development" {
				outputPath = filepath.Join(cwd, "testdata")
			}

			err = fs.WalkDir(template.Files, "files", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				fileName := strings.Replace(path, "files", projectName, 1)

				if d.IsDir() {
					err := os.MkdirAll(filepath.Join(outputPath, fileName), 0777)
					return err
				}

				content, err := template.Files.ReadFile(path)
				if err != nil {
					return err
				}

				newFilePath := filepath.Join(outputPath, fileName)
				err = os.WriteFile(newFilePath, content, 0777)
				return err
			})

			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name of the project.")
	// cmd.Flags().StringP("type", "t", "", "Type of the project.")
	// cmd.Flags().StringP("database", "d", "", "Database of the project.")

	// type: Web or API
	// db: none, sqlite, postgres, mysql/mariadb

	return cmd
}
