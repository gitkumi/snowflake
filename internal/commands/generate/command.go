package generate

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "gen",
		Short: "Generated CRUD",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
}
