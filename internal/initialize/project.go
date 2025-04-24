package initialize

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Project struct {
	Name          string
	AppType       AppType
	Database      Database
	BackgroundJob BackgroundJob

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
		Name:          cfg.Name,
		Database:      cfg.Database,
		BackgroundJob: cfg.BackgroundJob,
		AppType:       cfg.AppType,
		SMTP:          cfg.SMTP,
		Storage:       cfg.Storage,
		Redis:         cfg.Redis,
	}

	project.fileExclusions = []*FileExclusion{
		{
			FilePaths: []string{"/dev.yaml"},
			Check: func(p *Project) bool {
				return p.Database == DatabaseSQLite3 && !p.Redis
			},
		},
		{
			FilePaths: []string{
				"/internal/smtp/mailer.go",
				"/internal/smtp/mailer_smtp.go",
				"/internal/smtp/mailer_mock.go",
			},
			Check: func(p *Project) bool { return !p.SMTP },
		},
		{
			FilePaths: []string{
				"/internal/storage/storage.go",
				"/internal/storage/storage_s3.go",
				"/internal/storage/storage_mock.go",
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
				"/internal/html/hello.templ",
				"/internal/application/handler/html_handler.go",
			},
			Check: func(p *Project) bool { return p.AppType == AppTypeAPI },
		},
		{
			FilePaths: []string{
				"/sqlc.yaml",
				"/dev.yaml",
				"/static/sql/migrations/00001_books.sql",
				"/static/sql/queries/books.sql",
				"/static/static.go",
				"/internal/application/db.go",
				"/test/fixtures.go",
				"/internal/application/handler/book_handler.go",
				"/internal/application/handler/book_handler_test.go",
				"/internal/application/service/book_service.go",
			},
			Check: func(p *Project) bool { return p.Database == DatabaseNone },
		},
		{
			FilePaths: []string{
				"/internal/queue/queue.go",
				"/internal/queue/queue_sqs.go",
				"/internal/queue/queue_mock.go",
			},
			Check: func(p *Project) bool { return p.BackgroundJob == BackgroundJobBasic },
		},
		{
			FilePaths: []string{
				"/internal/application/task.go",
			},
			Check: func(p *Project) bool { return p.BackgroundJob == BackgroundJobSQS },
		},
		{
			FilePaths: []string{
				"/internal/application/task.go",
				"/internal/queue/queue.go",
				"/internal/queue/queue_sqs.go",
				"/internal/queue/queue_mock.go",
			},
			Check: func(p *Project) bool { return p.BackgroundJob == BackgroundJobNone },
		},
	}

	project.fileRenames = []*FileRename{
		{
			OldPath: "/cmd/api/main.go",
			NewPath: "/cmd/web/main.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
	}

	return project
}

func (p *Project) UsesDockerOnDev() bool {
	return p.Redis || p.Database != DatabaseNone
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

		targetDir := filepath.Dir(fullNewPath)
		if err := os.MkdirAll(targetDir, 0777); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", targetDir, err)
		}

		if err := os.Rename(fullOldPath, fullNewPath); err != nil {
			return fmt.Errorf("failed to rename file %s to %s: %v", fullOldPath, fullNewPath, err)
		}

		// Track the old directory for potential removal if empty later
		oldDir := path.Dir(fullOldPath)
		oldDirs[oldDir] = true
	}

	return RemoveEmptyDirs(oldDirs)
}
