package generator

import (
	"fmt"
	"path/filepath"

	snowflaketemplate "github.com/gitkumi/snowflake/template"
)

func LoadDatabaseMigration(db Database, filename string) (string, error) {
	fragmentPath := filepath.Join("fragments/database", string(db), "migrations", filename)
	content, err := snowflaketemplate.DatabaseFragments.ReadFile(fragmentPath)
	if err != nil {
		return "", fmt.Errorf("failed to read database fragment: %w", err)
	}
	return string(content), nil
}
