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
		BackgroundJob:  cfg.BackgroundJob,
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
		OIDCFacebook:   cfg.OIDCFacebook,
		OIDCGoogle:     cfg.OIDCGoogle,
		OIDCLinkedIn:   cfg.OIDCLinkedIn,
		OIDCMicrosoft:  cfg.OIDCMicrosoft,
		OIDCTwitch:     cfg.OIDCTwitch,
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
		// Handler files exclusions
		{
			FilePaths: []string{
				"/internal/application/handler/oauth_handler.go",
				"/internal/application/service/oauth_service.go",
			},
			Check: func(p *Project) bool { return !p.HasOAuth() },
		},
		{
			FilePaths: []string{
				"/internal/application/handler/oidc_handler.go",
				"/internal/application/service/oidc_service.go",
			},
			Check: func(p *Project) bool { return !p.HasOIDC() },
		},
		// OAuth file exclusions
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
		// OIDC file exclusions
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
