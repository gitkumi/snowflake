package generate

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	sqliteparser "github.com/gitkumi/snowflake/internal/parser/sqlite"
)

type SQLCConfig struct {
	Version string
	SQL     []struct {
		Engine  string
		Queries string
		Schema  string
		Gen     struct {
			Go struct {
				Package string
				out     string
			}
		}
	}
}

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "gen",
		Short: "Generate CRUD from SQL",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := readSQLCConfig()
			if err != nil {
				log.Fatal(err)
			}

			tables, err := parseTable("testdata/migration.sql")
			if err != nil {
				log.Fatal(err)
			}

			queries, err := parseQueries("testdata/queries.sql")
			if err != nil {
				log.Fatal(err)
			}

			for _, table := range tables {
				fmt.Println(table.TableName)
			}

			return

			fmt.Println(config, tables, queries)
		},
	}
}

func parseQueries(filePath string) ([]sqliteparser.Query, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return []sqliteparser.Query{}, err
	}

	queries, err := sqliteparser.ParseQueries(string(content))
	if err != nil {
		return []sqliteparser.Query{}, err
	}

	return queries, nil
}

func parseTable(filePath string) ([]sqliteparser.TableSchema, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return []sqliteparser.TableSchema{}, err
	}

	schemas, err := sqliteparser.ParseTable(string(content))
	if err != nil {
		return []sqliteparser.TableSchema{}, err
	}

	return schemas, nil
}

func readSQLCConfig() (SQLCConfig, error) {
	data, err := os.ReadFile("testdata/sqlc.yaml")
	if err != nil {
		return SQLCConfig{}, err
	}

	var config SQLCConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return SQLCConfig{}, err
	}

	return config, nil
}
