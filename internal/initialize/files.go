package initialize

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FileExclusions struct {
	ExcludeFuncs []*ExcludeFunc
}

type ExcludeFunc struct {
	FilePaths []string
	Check     func(*Project) bool
}

type FileRenames struct {
	ByAppType map[AppType]map[string]string
}

func NewFileExclusions() *FileExclusions {
	return &FileExclusions{
		ExcludeFuncs: []*ExcludeFunc{
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
		},
	}
}

func NewFileRenames() *FileRenames {
	return &FileRenames{
		ByAppType: map[AppType]map[string]string{
			AppTypeWeb: {
				"/cmd/api/main.go": "/cmd/web/main.go",
			},
		},
	}
}

func ExcludeTemplateFile(templateFileName string, project *Project, exclusions *FileExclusions) bool {
	fileName := strings.TrimSuffix(templateFileName, ".templ")

	for _, fn := range exclusions.ExcludeFuncs {
		for _, filePath := range fn.FilePaths {
			if fileName == filePath && fn.Check(project) {
				return true
			}
		}
	}

	return false
}

func RenameFiles(project *Project, outputPath string, renames *FileRenames) error {
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

	return RemoveEmptyDirs(oldDirs)
}

func RemoveEmptyDirs(paths map[string]bool) error {
	for dir := range paths {
		isEmpty, err := IsDirectoryEmpty(dir)
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

func IsDirectoryEmpty(name string) (bool, error) {
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
