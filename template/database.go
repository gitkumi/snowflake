package template

import (
	"embed"
	"fmt"
	"path/filepath"
	"text/template"
)

//go:embed fragments/database/*
var DatabaseFragments embed.FS

type DatabaseType string

const (
	MySQL    DatabaseType = "mysql"
	Postgres DatabaseType = "postgres"
	SQLite   DatabaseType = "sqlite"
)

// DatabaseConfig holds the configuration for database templating
type DatabaseConfig struct {
	Type DatabaseType
}

// LoadDatabaseMigration loads the appropriate database migration for the given file
func LoadDatabaseMigration(dbType DatabaseType, filename string) (string, error) {
	fragmentPath := filepath.Join("fragments/database", string(dbType), "migrations", filename)
	content, err := DatabaseFragments.ReadFile(fragmentPath)
	if err != nil {
		return "", fmt.Errorf("failed to read database fragment: %w", err)
	}
	return string(content), nil
}

// ProcessMigrationTemplate processes a migration template with the appropriate database fragment
func ProcessMigrationTemplate(baseTemplate string, dbConfig DatabaseConfig) (string, error) {
	tmpl, err := template.New("migration").Parse(baseTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		DatabaseMigration string
	}{
		DatabaseMigration: "{{.DatabaseFragment}}", // This will be replaced with actual content
	}

	return ExecuteTemplate(tmpl, data)
} 