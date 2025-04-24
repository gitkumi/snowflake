package cli

import (
	"log"

	// "github.com/gitkumi/snowflake/internal/command/generate"
	"github.com/gitkumi/snowflake/internal/command/run"
	"github.com/gitkumi/snowflake/internal/command/tui"
	"github.com/gitkumi/snowflake/internal/command/version"
	"github.com/spf13/cobra"
)

func Execute() {
	cmd := &cobra.Command{
		Use:   "snowflake",
		Short: "Snowflake is an opinionated Go web application generator.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.Root().CompletionOptions.DisableDefaultCmd = true
	cmd.AddCommand(run.Command())
	cmd.AddCommand(tui.Command())
	cmd.AddCommand(version.Command())
	// cmd.AddCommand(generate.Command())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
