package initialize

type Project struct {
	Name           string
	AppType        AppType
	Database       Database
	BackgroundJob  BackgroundJob

	SMTP           bool
	Storage        bool
	Redis          bool
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
	}

	return project
}
