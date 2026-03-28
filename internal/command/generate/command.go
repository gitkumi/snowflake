package generate

import (
	"log"
	"os"

	"github.com/gitkumi/snowflake/internal/generate"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate code for a Snowflake project",
	}

	cmd.AddCommand(resourceCommand())
	cmd.AddCommand(migrationCommand())
	return cmd
}

func resourceCommand() *cobra.Command {
	var quiet bool

	cmd := &cobra.Command{
		Use:   "resource <name> [field:type ...]",
		Short: "Generate a CRUD resource (migration, queries, handler, service)",
		Long: `Generate a full CRUD resource with migration, SQL queries, handler, and service.

Example:
  snowflake gen resource post title:string body:text published:bool

Valid field types: string, text, int, bigint, bool, float, timestamp
If no type is specified, "string" is used as the default.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			fields := args[1:]

			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			if err := generate.Run(name, fields, cwd, quiet); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress output")
	return cmd
}

func migrationCommand() *cobra.Command {
	var quiet bool

	cmd := &cobra.Command{
		Use:   "migration <name> [field:type ...]",
		Short: "Generate a database migration",
		Long: `Generate a database migration file.

Example:
  snowflake gen migration create_posts title:string body:text published:bool

Valid field types: string, text, int, bigint, bool, float, timestamp
If no type is specified, "string" is used as the default.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			fields := args[1:]

			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			if err := generate.RunMigration(name, fields, cwd, quiet); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress output")
	return cmd
}
