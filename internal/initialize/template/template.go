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

func loadFragments(fsys fs.FS, dir string, prefix string) (map[string]string, error) {
	fragments := make(map[string]string)

	err := fs.WalkDir(fsys, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}

		basename := filepath.Base(path)
		templateName := prefix + "_" + basename
		fragments[templateName] = string(content)
		return nil
	})

	// Ignore directory not found errors
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load fragments from directory %s: %w", dir, err)
	}

	return fragments, nil
}

func CreateDatabaseFragments(database string) (map[string]string, error) {
	databaseFragments := make(map[string]string)
	if database == "" {
		return databaseFragments, nil
	}

	// Load migration fragments
	migrationsDir := filepath.Join("fragments/database", database, "migrations")
	migrationFragments, err := loadFragments(DatabaseFragments, migrationsDir, "migration")
	if err != nil {
		return nil, err
	}

	// Add migration fragments to result
	for k, v := range migrationFragments {
		databaseFragments[k] = v
	}

	// Load query fragments
	queriesDir := filepath.Join("fragments/database", database, "queries")
	queryFragments, err := loadFragments(DatabaseFragments, queriesDir, "query")
	if err != nil {
		return nil, err
	}

	// Add query fragments to result
	for k, v := range queryFragments {
		databaseFragments[k] = v
	}

	return databaseFragments, nil
}
