package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/gitkumi/snowflake/internal/generator"
	"github.com/spf13/cobra"
)

func Execute() {
	var showVersion bool

	cmd := &cobra.Command{
		Use:   "snowflake",
		Short: "Snowflake is an opinionated Go REST API application generator.",
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Println("v0.16.0")
				return
			}

			cmd.Help()
		},
	}

	cmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show snowflake version")

	cmd.Root().CompletionOptions.DisableDefaultCmd = true

	cmd.AddCommand(new())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func new() *cobra.Command {
	var (
		initGit  bool
		database string
	)

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err.Error())
			}

			dbEnum := generator.Database(database)
			if !dbEnum.IsValid() {
				log.Fatalf("Invalid database type: %s. Must be one of: %v", database, generator.AllDatabases)
			}

			err = generator.Generate(args[0], initGit, cwd, dbEnum)
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}

	cmd.Flags().BoolVarP(&initGit, "git", "g", true, "Initialize git")
	cmd.Flags().StringVarP(&database, "database", "d", "sqlite3", fmt.Sprintf("Database type %v", generator.AllDatabases))

	return cmd
}
