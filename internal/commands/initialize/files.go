package initialize

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
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

func CreateTemplateFuncs(cfg *InitConfig) template.FuncMap {
	return template.FuncMap{
		"DatabaseMigration": func(filename string) (string, error) {
			return LoadDatabaseMigration(cfg.Database, filename)
		},
		"DatabaseQuery": func(filename string) (string, error) {
			return LoadDatabaseQuery(cfg.Database, filename)
		},
	}
}

func CreateFileExclusions() *FileExclusions {
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
				"/internal/html/hello.templ.templ",
				"/internal/application/handler/html_handler.go",
			},
		},
		ByDatabase: map[Database][]string{
			SQLite3: {
				"dev.yml.templ",
			},
		},
	}
}

func CreateFileRenames() *FileRenames {
	return &FileRenames{
		ByAppType: map[AppType]map[string]string{
			Web: {
				"/cmd/api/main.go": "/cmd/web/main.go",
			},
		},
	}
}

func ShouldExcludeTemplateFile(templateFileName string, project *Project, exclusions *FileExclusions) bool {
	fileName := strings.TrimSuffix(templateFileName, ".templ")

	// Check app type exclusions
	if excludedPaths, ok := exclusions.ByAppType[project.AppType]; ok {
		for _, excludedPath := range excludedPaths {
			if fileName == excludedPath {
				return true
			}
		}
	}

	// Check database type exclusions
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

func RenameFiles(project *Project, outputPath string, renames *FileRenames) error {
	renameMappings, ok := renames.ByAppType[project.AppType]
	if !ok {
		return nil
	}

	sourceDirs := make(map[string]bool)

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
			data, err := os.ReadFile(fullOldPath)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %v", fullOldPath, err)
			}

			if err := os.WriteFile(fullNewPath, data, 0666); err != nil {
				return fmt.Errorf("failed to write file %s: %v", fullNewPath, err)
			}

			if err := os.Remove(fullOldPath); err != nil {
				return fmt.Errorf("failed to remove file %s: %v", fullOldPath, err)
			}
		}

		sourceDir := filepath.Dir(fullOldPath)
		sourceDirs[sourceDir] = true
	}

	return cleanupSourceDirs(sourceDirs)
}

func cleanupSourceDirs(dirs map[string]bool) error {
	var dirList []string
	for dir := range dirs {
		dirList = append(dirList, dir)
	}

	// Sort by length in descending order to ensure child directories
	// are processed before their parents
	sort.Slice(dirList, func(i, j int) bool {
		return len(dirList[i]) > len(dirList[j])
	})

	for _, dir := range dirList {
		if err := cleanupEmptyDir(dir); err != nil {
			return err
		}
	}

	return nil
}

func cleanupEmptyDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", dir, err)
	}

	// Process subdirectories first
	for _, entry := range entries {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			if err := cleanupEmptyDir(subdir); err != nil {
				return err
			}
		}
	}

	entries, err = os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", dir, err)
	}

	if len(entries) == 0 {
		if err := os.Remove(dir); err != nil {
			return fmt.Errorf("failed to remove empty directory %s: %v", dir, err)
		}
	}

	return nil
}
