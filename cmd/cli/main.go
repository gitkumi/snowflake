package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/gitkumi/snowflake/internal/files"
	"github.com/spf13/cobra"
)

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
			fmt.Println("v0.13.0")
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
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err.Error())
			}

			err = files.Create(args[0], initGit, cwd)
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}

	cmd.Flags().BoolVarP(&initGit, "git", "g", true, "Initialize git")

	return cmd
}
