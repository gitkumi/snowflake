package initialize

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	initializetemplate "github.com/gitkumi/snowflake/internal/initialize/template"
)

type FileExclusions struct {
	NoSMTP     []string
	NoStorage  []string
	NoAuth     []string
	ByAppType  map[AppType][]string
	ByDatabase map[Database][]string
}

type FileRenames struct {
	ByAppType map[AppType]map[string]string
}

func createTemplateFuncs(cfg *Config) template.FuncMap {
	return template.FuncMap{
		"DatabaseMigration": func(filename string) (string, error) {
			return loadDatabaseMigration(cfg.Database, filename)
		},
		"DatabaseQuery": func(filename string) (string, error) {
			return loadDatabaseQuery(cfg.Database, filename)
		},
	}
}

func loadDatabaseMigration(db Database, filename string) (string, error) {
	if db == None {
		return "", nil
	}

	fragmentPath := filepath.Join("fragments/database", string(db), "migrations", filename)
	content, err := initializetemplate.DatabaseFragments.ReadFile(fragmentPath)
	if err != nil {
		return "", fmt.Errorf("failed to read database fragment: %w", err)
	}
	return string(content), nil
}

func loadDatabaseQuery(db Database, filename string) (string, error) {
	if db == None {
		return "", nil
	}

	fragmentPath := filepath.Join("fragments/database", string(db), "queries", filename)
	content, err := initializetemplate.DatabaseFragments.ReadFile(fragmentPath)
	if err != nil {
		return "", fmt.Errorf("failed to read database query fragment: %w", err)
	}
	return string(content), nil
}

func createFileExclusions() *FileExclusions {
	return &FileExclusions{
		NoSMTP: []string{
			"/internal/smtp/mailer.go",
			"/internal/smtp/mailer_smtp.go",
			"/internal/smtp/mailer_mock.go",
		},
		NoStorage: []string{
			"/internal/storage/storage.go",
			"/internal/storage/storage_s3.go",
			"/internal/storage/storage_mock.go",
		},
		NoAuth: []string{
			"/internal/password/password.go",
			"/internal/password/password_test.go",
			"/internal/middleware/auth.go",
			"/internal/middleware/auth_test.go",
			"/internal/application/handler/auth_handler_test.go",
			"/internal/application/handler/auth_handler_types.go",
			"/internal/application/handler/auth_handler.go",
			"/internal/application/service/auth_service_types.go",
			"/internal/application/service/auth_service.go",
			"/static/sql/migrations/00002_organizations.sql",
			"/static/sql/migrations/00003_users.sql",
			"/static/sql/migrations/00004_memberships.sql",
			"/static/sql/migrations/00005_user_auth_tokens.sql",
			"/static/sql/queries/memberships.sql",
			"/static/sql/queries/organizations.sql",
			"/static/sql/queries/user_auth_tokens.sql",
			"/static/sql/queries/users.sql",
		},
		ByAppType: map[AppType][]string{
			API: {
				"/internal/html/hello.templ",
				"/internal/application/handler/html_handler.go",
			},
		},
		ByDatabase: map[Database][]string{
			SQLite3: {
				"/dev.yaml",
			},
			None: {
				"/sqlc.yaml",
				"/dev.yaml",
				"/static/sql/migrations/00001_books.sql",
				"/static/sql/migrations/00002_organizations.sql",
				"/static/sql/migrations/00003_users.sql",
				"/static/sql/migrations/00004_memberships.sql",
				"/static/sql/migrations/00005_user_auth_tokens.sql",
				"/static/sql/queries/organizations.sql",
				"/static/sql/queries/memberships.sql",
				"/static/sql/queries/users.sql",
				"/static/sql/queries/books.sql",
				"/static/sql/queries/user_auth_tokens.sql",
				"/static/static.go",
				"/internal/application/db.go",
				"/test/fixtures.go",
				"/internal/application/handler/book_handler.go",
				"/internal/application/handler/book_handler_test.go",
				"/internal/application/handler/auth_handler.go",
				"/internal/application/handler/auth_handler_test.go",
				"/internal/application/handler/auth_handler_types.go",
				"/internal/application/service/auth_service.go",
				"/internal/application/service/auth_service_types.go",
				"/internal/application/service/book_service.go",
			},
		},
	}
}

func createFileRenames() *FileRenames {
	return &FileRenames{
		ByAppType: map[AppType]map[string]string{
			Web: {
				"/cmd/api/main.go": "/cmd/web/main.go",
			},
		},
	}
}

func shouldExcludeTemplateFile(templateFileName string, project *Project, exclusions *FileExclusions) bool {
	fileName := strings.TrimSuffix(templateFileName, ".templ")

	if excludedPaths, ok := exclusions.ByAppType[project.AppType]; ok {
		for _, excludedPath := range excludedPaths {
			if fileName == excludedPath {
				return true
			}
		}
	}

	if excludedPaths, ok := exclusions.ByDatabase[project.Database]; ok {
		for _, excludedPath := range excludedPaths {
			if fileName == excludedPath {
				return true
			}
		}
	}

	if !project.SMTP {
		for _, excludedPath := range exclusions.NoSMTP {
			if fileName == excludedPath {
				return true
			}
		}
	}

	if !project.Storage {
		for _, excludedPath := range exclusions.NoStorage {
			if fileName == excludedPath {
				return true
			}
		}
	}

	if !project.Auth {
		for _, excludedPath := range exclusions.NoAuth {
			if fileName == excludedPath {
				return true
			}
		}
	}

	return false
}

func renameFiles(project *Project, outputPath string, renames *FileRenames) error {
	oldDirs := make(map[string]bool)

	renameMappings, ok := renames.ByAppType[project.AppType]
	if !ok {
		return nil
	}

	for oldPath, newPath := range renameMappings {
		fullOldPath := filepath.Join(outputPath, oldPath)
		fullNewPath := filepath.Join(outputPath, newPath)

		if _, err := os.Stat(fullOldPath); os.IsNotExist(err) {
			continue
		}

		targetDir := filepath.Dir(fullNewPath)
		if err := os.MkdirAll(targetDir, 0777); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", targetDir, err)
		}

		if err := os.Rename(fullOldPath, fullNewPath); err != nil {
			return fmt.Errorf("failed to rename file %s: %v", fullOldPath, fullNewPath)
		}

		oldDir := path.Dir(fullOldPath)
		oldDirs[oldDir] = true
	}

	return removeEmptyDirs(oldDirs)
}

func removeEmptyDirs(paths map[string]bool) error {
	for dir := range paths {
		isEmpty, err := isDirectoryEmpty(dir)
		if err != nil {
			return fmt.Errorf("failed to check if directory %s is empty: %v", dir, err)
		}
		if isEmpty {
			if err := os.Remove(dir); err != nil {
				return fmt.Errorf("failed to remove empty directory %s: %v", dir, err)
			}
		}
	}
	return nil
}

func isDirectoryEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err != nil {
		return true, nil
	}

	return false, nil
}
