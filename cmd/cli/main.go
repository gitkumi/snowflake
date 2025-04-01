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
	var initGit bool

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err.Error())
			}

			err = generator.Generate(args[0], initGit, cwd)
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}

	cmd.Flags().BoolVarP(&initGit, "git", "g", true, "Initialize git")

	return cmd
}
