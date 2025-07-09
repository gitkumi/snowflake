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
	AppType  AppType
	Database Database
	Queue    Queue

	SMTP    bool
	Storage bool
	Redis   bool

	OAuthGoogle    bool
	OAuthDiscord   bool
	OAuthGitHub    bool
	OAuthInstagram bool
	OAuthMicrosoft bool
	OAuthReddit    bool
	OAuthSpotify   bool
	OAuthTwitch    bool
	OAuthFacebook  bool
	OAuthLinkedIn  bool
	OAuthSlack     bool
	OAuthStripe    bool
	OAuthX         bool

	OIDCFacebook  bool
	OIDCGoogle    bool
	OIDCLinkedIn  bool
	OIDCMicrosoft bool
	OIDCTwitch    bool
	OIDCDiscord   bool

	fileExclusions []*FileExclusion
	fileRenames    []*FileRename
}

func (p *Project) HasOAuth() bool {
	return p.OAuthGoogle ||
		p.OAuthGitHub ||
		p.OAuthFacebook ||
		p.OAuthInstagram ||
		p.OAuthDiscord ||
		p.OAuthLinkedIn ||
		p.OAuthReddit ||
		p.OAuthTwitch ||
		p.OAuthStripe ||
		p.OAuthX ||
		p.OAuthMicrosoft ||
		p.OAuthSlack ||
		p.OAuthSpotify
}

func (p *Project) HasOIDC() bool {
	return p.OIDCGoogle ||
		p.OIDCMicrosoft ||
		p.OIDCFacebook ||
		p.OIDCLinkedIn ||
		p.OIDCTwitch ||
		p.OIDCDiscord
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
		Name:           cfg.Name,
		Database:       cfg.Database,
		Queue:          cfg.Queue,
		AppType:        cfg.AppType,
		SMTP:           cfg.SMTP,
		Storage:        cfg.Storage,
		Redis:          cfg.Redis,
		OAuthGoogle:    cfg.OAuthGoogle,
		OAuthDiscord:   cfg.OAuthDiscord,
		OAuthGitHub:    cfg.OAuthGitHub,
		OAuthInstagram: cfg.OAuthInstagram,
		OAuthMicrosoft: cfg.OAuthMicrosoft,
		OAuthReddit:    cfg.OAuthReddit,
		OAuthSpotify:   cfg.OAuthSpotify,
		OAuthTwitch:    cfg.OAuthTwitch,
		OAuthFacebook:  cfg.OAuthFacebook,
		OAuthLinkedIn:  cfg.OAuthLinkedIn,
		OAuthSlack:     cfg.OAuthSlack,
		OAuthStripe:    cfg.OAuthStripe,
		OAuthX:         cfg.OAuthX,
		OIDCGoogle:     cfg.OIDCGoogle,
		OIDCFacebook:   cfg.OIDCFacebook,
		OIDCLinkedIn:   cfg.OIDCLinkedIn,
		OIDCMicrosoft:  cfg.OIDCMicrosoft,
		OIDCTwitch:     cfg.OIDCTwitch,
		OIDCDiscord:    cfg.OIDCDiscord,
	}

	if project.OIDCGoogle {
		project.OAuthGoogle = true
	}
	if project.OIDCFacebook {
		project.OAuthFacebook = true
	}
	if project.OIDCLinkedIn {
		project.OAuthLinkedIn = true
	}
	if project.OIDCMicrosoft {
		project.OAuthMicrosoft = true
	}
	if project.OIDCTwitch {
		project.OAuthTwitch = true
	}
	if project.OIDCDiscord {
		project.OAuthDiscord = true
	}

	if project.HasOAuth() || project.HasOIDC() {
		project.Redis = true
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
				"/cmd/api/handler/html_handler.go",
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
				"/cmd/api/db.go",
				"/test/fixtures.go",
				"/cmd/api/handler/book_handler.go",
				"/cmd/api/handler/book_handler_test.go",
				"/cmd/api/service/book_service.go",
				"/internal/dto/book.go",
			},
			Check: func(p *Project) bool { return p.Database == DatabaseNone },
		},
		{
			FilePaths: []string{
				"/internal/queue/queue.go",
				"/internal/queue/queue_sqs.go",
				"/internal/queue/queue_mock.go",
			},
			Check: func(p *Project) bool { return p.Queue == QueueBasic },
		},
		{
			FilePaths: []string{
				"/cmd/api/task.go",
			},
			Check: func(p *Project) bool { return p.Queue == QueueSQS },
		},
		{
			FilePaths: []string{
				"/cmd/api/task.go",
				"/internal/queue/queue.go",
				"/internal/queue/queue_sqs.go",
				"/internal/queue/queue_mock.go",
			},
			Check: func(p *Project) bool { return p.Queue == QueueNone },
		},
		{
			FilePaths: []string{
				"/internal/util/http.go",
				"/cmd/api/handler/oauth_handler.go",
				"/cmd/api/service/oauth_service.go",
			},
			Check: func(p *Project) bool { return !p.HasOAuth() },
		},
		{
			FilePaths: []string{
				"/cmd/api/handler/oidc_handler.go",
				"/cmd/api/service/oidc_service.go",
			},
			Check: func(p *Project) bool { return !p.HasOIDC() },
		},
		{
			FilePaths: []string{
				"/internal/oauth/google.go",
			},
			Check: func(p *Project) bool { return !p.OAuthGoogle },
		},
		{
			FilePaths: []string{
				"/internal/oauth/facebook.go",
			},
			Check: func(p *Project) bool { return !p.OAuthFacebook },
		},
		{
			FilePaths: []string{
				"/internal/oauth/github.go",
			},
			Check: func(p *Project) bool { return !p.OAuthGitHub },
		},
		{
			FilePaths: []string{
				"/internal/oauth/discord.go",
			},
			Check: func(p *Project) bool { return !p.OAuthDiscord },
		},
		{
			FilePaths: []string{
				"/internal/oauth/instagram.go",
			},
			Check: func(p *Project) bool { return !p.OAuthInstagram },
		},
		{
			FilePaths: []string{
				"/internal/oauth/linkedin.go",
			},
			Check: func(p *Project) bool { return !p.OAuthLinkedIn },
		},
		{
			FilePaths: []string{
				"/internal/oauth/microsoft.go",
			},
			Check: func(p *Project) bool { return !p.OAuthMicrosoft },
		},
		{
			FilePaths: []string{
				"/internal/oauth/reddit.go",
			},
			Check: func(p *Project) bool { return !p.OAuthReddit },
		},
		{
			FilePaths: []string{
				"/internal/oauth/slack.go",
			},
			Check: func(p *Project) bool { return !p.OAuthSlack },
		},
		{
			FilePaths: []string{
				"/internal/oauth/spotify.go",
			},
			Check: func(p *Project) bool { return !p.OAuthSpotify },
		},
		{
			FilePaths: []string{
				"/internal/oauth/stripe.go",
			},
			Check: func(p *Project) bool { return !p.OAuthStripe },
		},
		{
			FilePaths: []string{
				"/internal/oauth/twitch.go",
			},
			Check: func(p *Project) bool { return !p.OAuthTwitch },
		},
		{
			FilePaths: []string{
				"/internal/oauth/x.go",
			},
			Check: func(p *Project) bool { return !p.OAuthX },
		},
		{
			FilePaths: []string{
				"/internal/oauth/oauth.go",
			},
			Check: func(p *Project) bool { return !p.HasOAuth() },
		},
		{
			FilePaths: []string{
				"/internal/oidc/google.go",
			},
			Check: func(p *Project) bool { return !p.OIDCGoogle },
		},
		{
			FilePaths: []string{
				"/internal/oidc/facebook.go",
			},
			Check: func(p *Project) bool { return !p.OIDCFacebook },
		},
		{
			FilePaths: []string{
				"/internal/oidc/linkedin.go",
			},
			Check: func(p *Project) bool { return !p.OIDCLinkedIn },
		},
		{
			FilePaths: []string{
				"/internal/oidc/microsoft.go",
			},
			Check: func(p *Project) bool { return !p.OIDCMicrosoft },
		},
		{
			FilePaths: []string{
				"/internal/oidc/twitch.go",
			},
			Check: func(p *Project) bool { return !p.OIDCTwitch },
		},
		{
			FilePaths: []string{
				"/internal/oidc/discord.go",
			},
			Check: func(p *Project) bool { return !p.OIDCDiscord },
		},
		{
			FilePaths: []string{
				"/internal/oidc/oidc.go",
				"/internal/oidc/token.go",
			},
			Check: func(p *Project) bool { return !p.HasOIDC() },
		},
	}

	project.fileRenames = []*FileRename{
		{
			OldPath: "/cmd/api/main.go",
			NewPath: "/cmd/web/main.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/application.go",
			NewPath: "/cmd/web/application.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/task.go",
			NewPath: "/cmd/web/task.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/router.go",
			NewPath: "/cmd/web/router.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/db.go",
			NewPath: "/cmd/web/db.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/service/health_service.go",
			NewPath: "/cmd/web/service/health_service.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/service/oauth_service.go",
			NewPath: "/cmd/web/service/oauth_service.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/service/oidc_service.go",
			NewPath: "/cmd/web/service/oidc_service.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/service/book_service.go",
			NewPath: "/cmd/web/service/book_service.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/handler/book_handler.go",
			NewPath: "/cmd/web/handler/book_handler.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/handler/health_handler.go",
			NewPath: "/cmd/web/handler/health_handler.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/handler/oauth_handler.go",
			NewPath: "/cmd/web/handler/oauth_handler.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/handler/oidc_handler.go",
			NewPath: "/cmd/web/handler/oidc_handler.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/handler/pagination.go",
			NewPath: "/cmd/web/handler/pagination.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/handler/validator.go",
			NewPath: "/cmd/web/handler/validator.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/handler/response.go",
			NewPath: "/cmd/web/handler/response.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/handler/html_handler.go",
			NewPath: "/cmd/web/handler/html_handler.go",
			Check:   func(p *Project) bool { return p.AppType == AppTypeWeb },
		},
		{
			OldPath: "/cmd/api/handler/book_handler_test.go",
			NewPath: "/cmd/web/handler/book_handler_test.go",
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
