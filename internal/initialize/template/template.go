package template

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:base
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

	fragmentTypes := []struct {
		subDir string
		prefix string
	}{
		{"migrations", "migration"},
		{"queries", "query"},
	}

	for _, ft := range fragmentTypes {
		dir := filepath.Join("fragments/database", database, ft.subDir)
		fragments, err := loadFragments(DatabaseFragments, dir, ft.prefix)
		if err != nil {
			return nil, err
		}

		for k, v := range fragments {
			databaseFragments[k] = v
		}
	}

	return databaseFragments, nil
}
