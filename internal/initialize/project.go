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
				"/Dockerfile",
			},
			Check: func(p *Project) bool {
				return p.Database == DatabaseSQLite3 && !p.HasKeyValueStore()
			},
		},
		{
			FilePaths: []string{
				"/.devcontainer/devcontainer.json",
				"/.devcontainer/compose.yaml",
			},
			Check: func(p *Project) bool { return !p.HasDevContainer() },
		},
		{
			FilePaths: []string{
				"/internal/smtp/mailer.go",
				"/internal/smtp/smtp.go",
				"/internal/smtp/dev_mailbox.go",
				"/cmd/app/handlers/send_handler.go",
			},
			Check: func(p *Project) bool { return !p.SMTP },
		},
		{
			FilePaths: []string{
				"/internal/smtp/handler.go",
				"/internal/smtp/handler_test.go",
				"/internal/smtp/layout.templ",
				"/internal/smtp/list.templ",
				"/internal/smtp/show.templ",
			},
			Check: func(p *Project) bool { return !p.SMTP || !p.DevMailboxDashboard },
		},
		{
			FilePaths: []string{
				"/internal/storage/storage.go",
				"/internal/storage/s3.go",
				"/internal/storage/dev_storage.go",
				"/cmd/app/handlers/storage_handler.go",
			},
			Check: func(p *Project) bool { return !p.Storage },
		},
		{
			FilePaths: []string{
				"/internal/storage/handler.go",
				"/internal/storage/handler_test.go",
				"/internal/storage/layout.templ",
				"/internal/storage/list.templ",
				"/internal/storage/show.templ",
			},
			Check: func(p *Project) bool { return !p.Storage || !p.DevStorageDashboard },
		},
		{
			FilePaths: []string{
				"/internal/db/dev_db.go",
				"/internal/db/dev_db_handler.go",
				"/internal/db/dev_db_queries.go",
				"/internal/db/dev_db_layout.templ",
				"/internal/db/dev_db_row_form.templ",
				"/internal/db/dev_db_rows.templ",
				"/internal/db/dev_db_tables.templ",
			},
			Check: func(p *Project) bool { return !p.DevDBDashboard },
		},
		{
			FilePaths: []string{
				"/sqlc.yaml",
				"/cmd/app/sql/sql.go",
				"/cmd/migrator/main.go",
				"/internal/db/db.go",
			},
			Check: func(p *Project) bool { return p.Database == DatabaseNone },
		},
		{
			FilePaths: []string{
				"/internal/jobs/jobs.go",
				"/internal/jobs/tasks.go",
				"/internal/jobs/absurd.sql",
				"/cmd/app/handlers/jobs_handler.go",
			},
			Check: func(p *Project) bool { return !p.HasJobs() },
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

func (p *Project) HasJobs() bool {
	return p.JobProcessor != JobProcessorNone
}

// HasDatabaseService reports whether the database runs as its own container.
// SQLite is a file, so it needs no service.
func (p *Project) HasDatabaseService() bool {
	return p.Database == DatabasePostgres || p.Database == DatabaseMySQL || p.Database == DatabaseMariaDB
}

// HasDevContainer reports whether the project gets a .devcontainer/. There is
// nothing to run when the project has neither a database service nor a
// key-value store, so the devcontainer is skipped in that case.
func (p *Project) HasDevContainer() bool {
	return p.HasDatabaseService() || p.HasKeyValueStore()
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
