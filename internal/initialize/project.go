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

	SMTP      bool
	Storage   bool
	Redis     bool
	ServeHTML bool

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

	Billing Billing

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
		SMTP:           cfg.SMTP,
		Storage:        cfg.Storage,
		Redis:          cfg.Redis,
		ServeHTML:      cfg.ServeHTML,
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
		Billing:        cfg.Billing,
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
			FilePaths: []string{
				"/cmd/app/dev.yaml",
				"/cmd/app/Dockerfile",
			},
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
				"/cmd/app/handler/html_handler.go",
				"/cmd/app/html/hello.templ",
			},
			Check: func(p *Project) bool { return !p.ServeHTML },
		},
		{
			FilePaths: []string{
				"/cmd/app/sqlc.yaml",
				"/cmd/app/dev.yaml",
				"/cmd/app/static/sql/migrations/00001_books.sql",
				"/cmd/app/static/sql/queries/books.sql",
				"/cmd/app/static/static.go",
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
				"/internal/queue/queue_sqs.go",
				"/internal/queue/queue_mock.go",
			},
			Check: func(p *Project) bool { return p.Queue == QueueBasic },
		},
		{
			FilePaths: []string{
				"/cmd/app/application/task.go",
			},
			Check: func(p *Project) bool { return p.Queue == QueueSQS },
		},
		{
			FilePaths: []string{
				"/cmd/app/application/task.go",
				"/internal/queue/queue.go",
				"/internal/queue/queue_sqs.go",
				"/internal/queue/queue_mock.go",
			},
			Check: func(p *Project) bool { return p.Queue == QueueNone },
		},
		{
			FilePaths: []string{
				"/internal/util/http.go",
				"/cmd/app/handler/oauth_handler.go",
				"/cmd/app/service/oauth_service.go",
			},
			Check: func(p *Project) bool { return !p.HasOAuth() },
		},
		{
			FilePaths: []string{
				"/cmd/app/handler/oidc_handler.go",
				"/cmd/app/service/oidc_service.go",
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
		{
			FilePaths: []string{
				"/internal/billing/billing.go",
				"/internal/billing/billing_stripe.go",
				"/cmd/app/dto/billing.go",
				"/cmd/app/service/billing_service.go",
				"/cmd/app/handler/billing_handler.go",
				"/cmd/app/handler/billing_webhook_handler.go",
			},
			Check: func(p *Project) bool { return p.Billing == BillingNone },
		},
	}

	project.fileRenames = []*FileRename{}

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
