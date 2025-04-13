package cli

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/gitkumi/snowflake/internal/commands/initialize"
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
	cmd.AddCommand(initialize.InitProject())
	cmd.AddCommand(showVersion())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func showVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show current version",
		Run: func(cmd *cobra.Command, args []string) {
			info, ok := debug.ReadBuildInfo()
			if ok {
				fmt.Println(info.Main.Version)
				return
			}

			fmt.Println("dev")
		},
	}
}
