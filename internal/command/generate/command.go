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
		Use:   "resource <Name> <plural> [field:type ...]",
		Short: "Generate a CRUD resource (migration, queries, handler, service)",
		Long: `Generate a full CRUD resource with migration, SQL queries, handler, and service.

The first argument is the resource name (singular, e.g. "Post").
The second argument is the plural table name (e.g. "posts").

Fields are specified as name:type pairs.

Example:
  snowflake gen resource Post posts title:string body:text published:bool

Valid field types: string, text, int, bigint, bool, float, timestamp`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			if err := generate.Run(generate.GenerateInput{
				Name:       args[0],
				Plural:     args[1],
				RawFields:  args[2:],
				ProjectDir: cwd,
				Quiet:      quiet,
			}); err != nil {
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
		Use:   "migration <Name> <plural> [field:type ...]",
		Short: "Generate a database migration",
		Long: `Generate a database migration file.

The first argument is the resource name (singular, e.g. "Post").
The second argument is the plural table name (e.g. "posts").

Fields are specified as name:type pairs.

Example:
  snowflake gen migration Post posts title:string body:text published:bool

Valid field types: string, text, int, bigint, bool, float, timestamp`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			if err := generate.RunMigration(generate.GenerateInput{
				Name:       args[0],
				Plural:     args[1],
				RawFields:  args[2:],
				ProjectDir: cwd,
				Quiet:      quiet,
			}); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress output")
	return cmd
}
