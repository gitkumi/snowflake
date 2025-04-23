package initialize

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FileExclusions struct {
	SMTP               []string
	Storage            []string
	Redis              []string
	Auth               []string
	OAuthGoogle        []string
	OAuthFacebook      []string
	OAuthGithub        []string
	OAuthLinkedIn      []string
	OAuthInstagram     []string
	OAuthDiscord       []string
	AppType            map[AppType][]string
	Database           map[Database][]string
	BackgroundJob      map[BackgroundJob][]string
	AuthenticationType map[Authentication][]string
	ExcludeFuncs       []*ExcludeFunc
}

type ExcludeFunc struct {
	FilePath string
	Check    func(*Project) bool
}

type FileRenames struct {
	ByAppType map[AppType]map[string]string
}

func createFileExclusions() *FileExclusions {
	return &FileExclusions{
		SMTP: []string{
			"/internal/smtp/mailer.go",
			"/internal/smtp/mailer_smtp.go",
			"/internal/smtp/mailer_mock.go",
		},
		Storage: []string{
			"/internal/storage/storage.go",
			"/internal/storage/storage_s3.go",
			"/internal/storage/storage_mock.go",
		},
		Redis: []string{
			"/internal/middleware/rate_limit.go",
		},
		AuthenticationType: map[Authentication][]string{
			AuthenticationNone: []string{
				"/internal/dto/auth.go",
				"/internal/password/password.go",
				"/internal/password/password_test.go",
				"/internal/middleware/auth.go",
				"/internal/middleware/auth_test.go",
				"/internal/application/handler/auth_handler_test.go",
				"/internal/application/handler/auth_handler.go",
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
		},
		AppType: map[AppType][]string{
			AppTypeAPI: {
				"/internal/html/hello.templ",
				"/internal/application/handler/html_handler.go",
			},
		},
		Database: map[Database][]string{
			DatabaseNone: {
				"/sqlc.yaml",
				"/dev.yaml",
				"/static/sql/migrations/00001_books.sql",
				"/static/sql/migrations/00002_organizations.sql",
				"/static/sql/migrations/00003_users.sql",
				"/static/sql/migrations/00004_memberships.sql",
				"/static/sql/migrations/00005_user_auth_tokens.sql",
				"/static/sql/migrations/00006_user_oauth.sql",
				"/static/sql/queries/organizations.sql",
				"/static/sql/queries/memberships.sql",
				"/static/sql/queries/users.sql",
				"/static/sql/queries/books.sql",
				"/static/sql/queries/user_auth_tokens.sql",
				"/static/sql/queries/user_oauth.sql",
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
		BackgroundJob: map[BackgroundJob][]string{
			BackgroundJobBasic: {
				"/internal/queue/queue.go",
				"/internal/queue/queue_sqs.go",
				"/internal/queue/queue_mock.go",
			},
			BackgroundJobAsynq: {
				"/internal/application/task.go",
				"/internal/queue/queue.go",
				"/internal/queue/queue_sqs.go",
				"/internal/queue/queue_mock.go",
			},
			BackgroundJobSQS: {
				"/internal/application/task.go",
			},
			BackgroundJobNone: {
				"/internal/application/task.go",
				"/internal/queue/queue.go",
				"/internal/queue/queue_sqs.go",
				"/internal/queue/queue_mock.go",
			},
		},
		OAuthGoogle: []string{
			"/internal/OAuth/google.go",
			"/internal/OAuth/google_mock.go",
			"/internal/OAuth/google_test.go",
		},
		OAuthFacebook: []string{
			"/internal/OAuth/facebook.go",
			"/internal/OAuth/facebook_mock.go",
			"/internal/OAuth/facebook_test.go",
		},
		OAuthGithub: []string{
			"/internal/OAuth/github.go",
			"/internal/OAuth/github_mock.go",
			"/internal/OAuth/github_test.go",
		},
		OAuthLinkedIn: []string{
			"/internal/OAuth/linkedin.go",
			"/internal/OAuth/linkedin_mock.go",
			"/internal/OAuth/linkedin_test.go",
		},
		OAuthInstagram: []string{
			"/internal/OAuth/instagram.go",
			"/internal/OAuth/instagram_mock.go",
			"/internal/OAuth/instagram_test.go",
		},
		OAuthDiscord: []string{
			"/internal/OAuth/discord.go",
			"/internal/OAuth/discord_mock.go",
			"/internal/OAuth/discord_test.go",
		},
		ExcludeFuncs: []*ExcludeFunc{
			{
				FilePath: "/internal/dto/oauth.go",
				Check: func(p *Project) bool {
					return !p.WithOAuth()
				},
			},
			{
				FilePath: "/internal/application/handler/oauth_handler.go",
				Check: func(p *Project) bool {
					return !p.WithOAuth()
				},
			},
			{
				FilePath: "/internal/application/service/oauth_service.go",
				Check: func(p *Project) bool {
					return !p.WithOAuth()
				},
			},
			{
				FilePath: "/static/sql/migrations/00005_user_oauth.sql",
				Check: func(p *Project) bool {
					return !p.WithOAuth()
				},
			},
			{
				FilePath: "/static/sql/queries/user_oauth.sql",
				Check: func(p *Project) bool {
					return !p.WithOAuth()
				},
			},
			{
				FilePath: "/dev.yaml",
				Check: func(p *Project) bool {
					return p.Database == DatabaseSQLite3 && !p.Redis
				},
			},
		},
	}
}

func createFileRenames() *FileRenames {
	return &FileRenames{
		ByAppType: map[AppType]map[string]string{
			AppTypeWeb: {
				"/cmd/api/main.go": "/cmd/web/main.go",
			},
		},
	}
}

func shouldExcludeTemplateFile(templateFileName string, project *Project, exclusions *FileExclusions) bool {
	fileName := strings.TrimSuffix(templateFileName, ".templ")

	for _, fn := range exclusions.ExcludeFuncs {
		if fileName == fn.FilePath && fn.Check(project) {
			return true
		}
	}

	mapExclusions := []struct {
		exclusionMap interface{}
		key          interface{}
	}{
		{exclusions.AuthenticationType, project.Authentication},
		{exclusions.AppType, project.AppType},
		{exclusions.Database, project.Database},
		{exclusions.BackgroundJob, project.BackgroundJob},
	}

	for _, mapExcl := range mapExclusions {
		switch m := mapExcl.exclusionMap.(type) {
		case map[Authentication][]string:
			if key, ok := mapExcl.key.(Authentication); ok {
				if excludedPaths, exists := m[key]; exists {
					for _, path := range excludedPaths {
						if fileName == path {
							return true
						}
					}
				}
			}
		case map[AppType][]string:
			if key, ok := mapExcl.key.(AppType); ok {
				if excludedPaths, exists := m[key]; exists {
					for _, path := range excludedPaths {
						if fileName == path {
							return true
						}
					}
				}
			}
		case map[Database][]string:
			if key, ok := mapExcl.key.(Database); ok {
				if excludedPaths, exists := m[key]; exists {
					for _, path := range excludedPaths {
						if fileName == path {
							return true
						}
					}
				}
			}
		case map[BackgroundJob][]string:
			if key, ok := mapExcl.key.(BackgroundJob); ok {
				if excludedPaths, exists := m[key]; exists {
					for _, path := range excludedPaths {
						if fileName == path {
							return true
						}
					}
				}
			}
		}
	}

	featureExclusions := []struct {
		enabled  bool
		pathList []string
	}{
		{project.Redis, exclusions.Redis},
		{project.SMTP, exclusions.SMTP},
		{project.Storage, exclusions.Storage},
		{project.OAuthGoogle, exclusions.OAuthGoogle},
		{project.OAuthFacebook, exclusions.OAuthFacebook},
		{project.OAuthGitHub, exclusions.OAuthGithub},
		{project.OAuthLinkedIn, exclusions.OAuthLinkedIn},
		{project.OAuthInstagram, exclusions.OAuthInstagram},
		{project.OAuthDiscord, exclusions.OAuthDiscord},
	}

	for _, feature := range featureExclusions {
		if !feature.enabled {
			for _, path := range feature.pathList {
				if fileName == path {
					return true
				}
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
