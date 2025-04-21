package generate

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "gen",
		Short: "Generate CRUD from SQL",
		Run: func(cmd *cobra.Command, args []string) {
			Generate()
		},
	}
}

// func parseQueries(filePath string) ([]sqliteparser.Query, error) {
// 	content, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return []sqliteparser.Query{}, err
// 	}
//
// 	queries, err := sqliteparser.ParseQueries(string(content))
// 	if err != nil {
// 		return []sqliteparser.Query{}, err
// 	}
//
// 	return queries, nil
// }
//
// func parseTable(filePath string) ([]sqliteparser.TableSchema, error) {
// 	content, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return []sqliteparser.TableSchema{}, err
// 	}
//
// 	schemas, err := sqliteparser.ParseTable(string(content))
// 	if err != nil {
// 		return []sqliteparser.TableSchema{}, err
// 	}
//
// 	return schemas, nil
// }
//
// func readSQLCConfig() (SQLCConfig, error) {
// 	data, err := os.ReadFile("testdata/sqlc.yaml")
// 	if err != nil {
// 		return SQLCConfig{}, err
// 	}
//
// 	var config SQLCConfig
// 	if err := yaml.Unmarshal(data, &config); err != nil {
// 		return SQLCConfig{}, err
// 	}
//
// 	return config, nil
// }
