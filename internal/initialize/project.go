package initialize

import (
	"strings"
)

type Project struct {
	*Config
	fileExclusions []FileExclusion
}

type FileExclusion struct {
	FilePaths []string
	Check     func(*Project) bool
}

func NewProject(cfg *Config) *Project {
	project := &Project{
		Config: cfg,
	}

	project.fileExclusions = []FileExclusion{
		{
			FilePaths: []string{
				"/devenv.yaml",
				"/Dockerfile",
			},
			Check: func(p *Project) bool {
				return p.Database == DatabaseSQLite3 && !p.HasKeyValueStore()
			},
		},
		{
			FilePaths: []string{
				"/internal/smtp/mailer.go",
				"/internal/smtp/smtp.go",
				"/internal/smtp/dev_mailbox.go",
				"/internal/smtp/handler.go",
				"/internal/smtp/handler_test.go",
				"/internal/smtp/layout.templ",
				"/internal/smtp/list.templ",
				"/internal/smtp/show.templ",
				"/cmd/app/handlers/send_handler.go",
			},
			Check: func(p *Project) bool { return !p.SMTP },
		},
		{
			FilePaths: []string{
				"/internal/storage/storage.go",
				"/internal/storage/s3.go",
				"/internal/storage/dev_storage.go",
				"/internal/storage/handler.go",
				"/internal/storage/handler_test.go",
				"/internal/storage/layout.templ",
				"/internal/storage/list.templ",
				"/internal/storage/show.templ",
			},
			Check: func(p *Project) bool { return !p.Storage },
		},
		{
			FilePaths: []string{
				"/sqlc.yaml",
				"/cmd/app/sql/sql.go",
				"/cmd/migrator/main.go",
				"/internal/db/db.go",
				"/devenv.yaml",
			},
			Check: func(p *Project) bool { return p.Database == DatabaseNone },
		},
		{
			FilePaths: []string{
				"/internal/html/html.go",
				"/internal/html/static/reset.css",
				"/internal/html/pages/index.templ",
				"/internal/html/ui/button.go",
				"/internal/html/ui/button.templ",
				"/internal/html/ui/dev_page.templ",
				"/internal/html/ui/field.templ",
				"/cmd/app/handlers/page_handler.go",
			},
			Check: func(p *Project) bool { return !p.Templ },
		},
	}

	return project
}

func (p *Project) HasKeyValueStore() bool {
	return p.KeyValueStore != KeyValueStoreNone
}

func (p *Project) HasDevEnv() bool {
	return !(p.Database == DatabaseSQLite3 && !p.HasKeyValueStore())
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
