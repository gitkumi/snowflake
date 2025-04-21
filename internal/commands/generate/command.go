package generate

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "gen",
		Short: "Generate CRUD from SQL",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("wip")
		},
	}
}
