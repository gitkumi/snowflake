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
				"/cmd/app/devenv.yaml",
				"/cmd/app/Dockerfile",
			},
			Check: func(p *Project) bool {
				return p.Database == DatabaseSQLite3 && !p.HasKeyValueStore()
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
				"/cmd/app/sqlc.yaml",
				"/cmd/app/sql/sql.go",
				"/cmd/migrator/main.go",
				"/internal/db/db.go",
				"/cmd/app/devenv.yaml",
			},
			Check: func(p *Project) bool { return p.Database == DatabaseNone },
		},
		{
			FilePaths: []string{
				"/internal/html/html.go",
				"/internal/html/static/reset.css",
				"/internal/html/pages/index.templ",
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
