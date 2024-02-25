package cmd

import (
	"log"

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
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	cmd.Flags().StringP("name", "n", "", "Name of the project.")

	// Web or API
	// sqlite, postgres, mysql/mariadb

	return cmd
}
