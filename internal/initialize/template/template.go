package template

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed base/*
var BaseFiles embed.FS

//go:embed fragments/database/*
var DatabaseFragments embed.FS

func CreateDatabaseFragments(database string) (map[string]string, error) {
	databaseFragments := make(map[string]string)
	if database == "" {
		return databaseFragments, nil
	}

	migrationsDir := filepath.Join("fragments/database", database, "migrations")
	err := fs.WalkDir(DatabaseFragments, migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		content, err := fs.ReadFile(DatabaseFragments, path)
		if err != nil {
			return err
		}

		basename := filepath.Base(path)
		templateName := "migration_" + basename
		databaseFragments[templateName] = string(content)
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load database migration fragments: %w", err)
	}

	queriesDir := filepath.Join("fragments/database", database, "queries")
	err = fs.WalkDir(DatabaseFragments, queriesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		content, err := fs.ReadFile(DatabaseFragments, path)
		if err != nil {
			return err
		}

		basename := filepath.Base(path)
		templateName := "query_" + basename
		databaseFragments[templateName] = string(content)
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load database query fragments: %w", err)
	}

	return databaseFragments, nil
}
