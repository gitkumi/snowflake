package generate

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	sqliteparser "github.com/gitkumi/snowflake/internal/parser/sqlite"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "gen",
		Short: "Generate CRUD from SQL",
		Run: func(cmd *cobra.Command, args []string) {
			parseTable()
			parseQueries()
		},
	}
}

func parseQueries() {
	sqlFilePath := "testdata/queries.sql"
	content, err := os.ReadFile(sqlFilePath)
	if err != nil {
		fmt.Printf("Error reading SQL file: %v\n", err)
		return
	}

	queries, err := sqliteparser.ParseQueries(string(content))
	if err != nil {
		fmt.Printf("Error parsing SQL: %v\n", err)
		return
	}

	fmt.Println(queries)
}

func parseTable() {
	sqlFilePath := "testdata/migration.sql"
	content, err := os.ReadFile(sqlFilePath)
	if err != nil {
		fmt.Printf("Error reading SQL file: %v\n", err)
		return
	}

	schemas, err := sqliteparser.ParseTable(string(content))
	if err != nil {
		fmt.Printf("Error parsing SQL: %v\n", err)
		return
	}

	fmt.Println(schemas)
}
