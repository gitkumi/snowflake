package initialize

type Project struct {
	Name           string
	AppType        AppType
	Database       Database
	BackgroundJob  BackgroundJob
	Authentication Authentication

	SMTP           bool
	Storage        bool
	Redis          bool
	OAuthDiscord   bool
	OAuthFacebook  bool
	OAuthGitHub    bool
	OAuthGoogle    bool
	OAuthInstagram bool
	OAuthLinkedIn  bool
}

func (p *Project) WithAuth() bool {
	return p.Authentication != AuthenticationNone
}

func (p *Project) WithOAuth() bool {
	return p.WithAuth() && (p.OAuthGoogle || p.OAuthDiscord || p.OAuthFacebook || p.OAuthGitHub || p.OAuthInstagram || p.OAuthLinkedIn)
}

func NewProject(cfg *Config) *Project {
	project := &Project{
		Name:           cfg.Name,
		Database:       cfg.Database,
		BackgroundJob:  cfg.BackgroundJob,
		AppType:        cfg.AppType,
		SMTP:           cfg.SMTP,
		Storage:        cfg.Storage,
		Redis:          cfg.Redis || cfg.BackgroundJob == BackgroundJobAsynq,
		Authentication: cfg.Authentication,
		OAuthDiscord:   cfg.Authentication != AuthenticationNone && cfg.OAuthDiscord,
		OAuthFacebook:  cfg.Authentication != AuthenticationNone && cfg.OAuthFacebook,
		OAuthGitHub:    cfg.Authentication != AuthenticationNone && cfg.OAuthGitHub,
		OAuthGoogle:    cfg.Authentication != AuthenticationNone && cfg.OAuthGoogle,
		OAuthInstagram: cfg.Authentication != AuthenticationNone && cfg.OAuthInstagram,
		OAuthLinkedIn:  cfg.Authentication != AuthenticationNone && cfg.OAuthLinkedIn,
	}

	if cfg.Database == DatabaseNone || !cfg.SMTP {
		project.Authentication = AuthenticationNone
		project.OAuthDiscord = false
		project.OAuthFacebook = false
		project.OAuthGitHub = false
		project.OAuthGoogle = false
		project.OAuthInstagram = false
		project.OAuthLinkedIn = false
	}

	return project
}
