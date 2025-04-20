package version

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the current version",
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
