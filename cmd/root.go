package cmd

import (
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
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]

			err := fs.WalkDir(template.Files, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				cwd, err := os.Getwd()
				if err != nil {
					return err
				}

				relativePath := strings.Replace(path, "files", projectName, 1)
				outputPath := filepath.Join(cwd, "testdata")

				if d.IsDir() {
					err := os.MkdirAll(filepath.Join(outputPath, relativePath), 0777)

					if err != nil {
						return err
					}

					return nil
				}

				content, err := template.Files.ReadFile(path)
				if err != nil {
					return err
				}

				if err != nil {
					return err
				}

				newFileName := filepath.Join(cwd, "testdata", relativePath)
				err = os.WriteFile(newFileName, content, 0777)
				if err != nil {
					return err
				}

				return nil
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
