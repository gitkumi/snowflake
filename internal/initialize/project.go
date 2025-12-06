package initialize

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Project struct {
	Name     string
	Database Database
	Queue    Queue

	SMTP    bool
	Storage bool
	Redis   bool

	fileExclusions []*FileExclusion
	fileRenames    []*FileRename
}

type FileRename struct {
	OldPath string
	NewPath string
	Check   func(*Project) bool
}

type FileExclusion struct {
	FilePaths []string
	Check     func(*Project) bool
}

func NewProject(cfg *Config) *Project {
	project := &Project{
		Name:     cfg.Name,
		Database: cfg.Database,
		Queue:    cfg.Queue,
		SMTP:     cfg.SMTP,
		Storage:  cfg.Storage,
		Redis:    cfg.Redis,
	}

	project.fileExclusions = []*FileExclusion{
		{
			FilePaths: []string{
				"/cmd/app/devenv.yaml",
				"/cmd/app/Dockerfile",
			},
			Check: func(p *Project) bool {
				return p.Database == DatabaseSQLite3 && !p.Redis
			},
		},
		{
			FilePaths: []string{
				"/internal/smtp/mailer.go",
				"/internal/smtp/smtp.go",
				"/internal/smtp/mock.go",
			},
			Check: func(p *Project) bool { return !p.SMTP },
		},
		{
			FilePaths: []string{
				"/internal/storage/storage.go",
				"/internal/storage/s3.go",
				"/internal/storage/mock.go",
			},
			Check: func(p *Project) bool { return !p.Storage },
		},
		{
			FilePaths: []string{
				"/internal/middleware/rate_limit.go",
			},
			Check: func(p *Project) bool { return !p.Redis },
		},
		{
			FilePaths: []string{
				"/sqlc.yaml",
				"/sql/migrations/00001_books.sql",
				"/sql/queries/books.sql",
				"/sql/sql.go",
				"/cmd/migrator/main.go",
				"/internal/db/db.go",
				"/cmd/app/devenv.yaml",
				"/cmd/app/application/db.go",
				"/cmd/app/handler/book_handler.go",
				"/cmd/app/handler/book_handler_test.go",
				"/cmd/app/service/book_service.go",
				"/cmd/app/dto/book.go",
				"/cmd/app/dto/dto.go",
			},
			Check: func(p *Project) bool { return p.Database == DatabaseNone },
		},
		{
			FilePaths: []string{
				"/internal/queue/queue.go",
				"/internal/queue/sqs.go",
				"/internal/queue/mock.go",
			},
			Check: func(p *Project) bool { return p.Queue == QueueNone },
		},
	}

	project.fileRenames = []*FileRename{}

	return project
}

func (p *Project) UsesDockerOnDev() bool {
	return p.Redis || p.Database != DatabaseNone
}

func (p *Project) HasDevEnv() bool {
	return !(p.Database == DatabaseSQLite3 && !p.Redis)
}

func (p *Project) ExcludeFile(templateFileName string) bool {
	fileName := strings.TrimSuffix(templateFileName, ".templ")

	for _, exclusion := range p.fileExclusions {
		for _, filePath := range exclusion.FilePaths {
			if fileName == filePath && exclusion.Check(p) {
				return true
			}
		}
	}

	return false
}

func (p *Project) RenameFiles(outputPath string) error {
	oldDirs := make(map[string]bool)

	for _, rename := range p.fileRenames {
		if !rename.Check(p) {
			continue
		}

		fullOldPath := filepath.Join(outputPath, rename.OldPath)
		fullNewPath := filepath.Join(outputPath, rename.NewPath)

		// Track the old directory for potential removal if empty later
		oldDir := path.Dir(fullOldPath)
		oldDirs[oldDir] = true

		// Check if source file exists, skip if it doesn't (could be excluded)
		if _, err := os.Stat(fullOldPath); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return fmt.Errorf("failed to check if file exists %s: %v", fullOldPath, err)
		}

		if err := os.MkdirAll(filepath.Dir(fullNewPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", filepath.Dir(fullNewPath), err)
		}

		if err := os.Rename(fullOldPath, fullNewPath); err != nil {
			return fmt.Errorf("failed to rename file %s to %s: %v", fullOldPath, fullNewPath, err)
		}
	}

	return RemoveEmptyDirs(oldDirs)
}
