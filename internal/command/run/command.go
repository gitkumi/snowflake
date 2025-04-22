package run

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gitkumi/snowflake/internal/initialize"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var (
		quiet          bool
		database       string
		backgroundJob  string
		appType        string
		outputDir      string
		git            bool
		smtp           bool
		storage        bool
		redis          bool
		auth           bool
		oauthDiscord   bool
		oauthFacebook  bool
		oauthGitHub    bool
		oauthGoogle    bool
		oauthInstagram bool
		oauthLinkedIn  bool
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Create a new project using command-line flags",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if outputDir == "" {
				cwd, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				outputDir = cwd
			}

			if !filepath.IsAbs(outputDir) {
				cwd, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				outputDir = filepath.Join(cwd, outputDir)
			}

			// Ensure output directory exists
			if _, err := os.Stat(outputDir); os.IsNotExist(err) {
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					log.Fatalf("Failed to create output directory: %v", err)
				}
			}

			dbEnum := initialize.Database(database)
			if !dbEnum.IsValid() {
				log.Fatalf("Invalid database type: %s. Must be one of: %v", database, initialize.AllDatabases)
			}

			appTypeEnum := initialize.AppType(appType)
			if !appTypeEnum.IsValid() {
				log.Fatalf("Invalid app type: %s. Must be one of: %v", appType, initialize.AllAppTypes)
			}

			backgroundJobEnum := initialize.BackgroundJob(backgroundJob)
			if !backgroundJobEnum.IsValid() {
				log.Fatalf("Invalid app type: %s. Must be one of: %v", backgroundJob, initialize.AllBackgroundJobs)
			}

			err := initialize.Run(&initialize.Config{
				Quiet:         quiet,
				Name:          args[0],
				Database:      dbEnum,
				BackgroundJob: backgroundJobEnum,
				AppType:       appTypeEnum,
				Git:           git,
				OutputDir:     outputDir,
				SMTP:          smtp,
				Storage:       storage,
				Redis:         redis,
				Auth:          auth,
			})
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}

	cmd.Flags().StringVarP(&appType, "app-type", "t", "api", fmt.Sprintf("App type %v", initialize.AllAppTypes))
	cmd.Flags().StringVarP(&database, "database", "d", "sqlite3", fmt.Sprintf("Database type %v", initialize.AllDatabases))

	cmd.Flags().StringVarP(&backgroundJob, "background-job", "b", "basic", fmt.Sprintf("Background Job type %v", initialize.AllBackgroundJobs))
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for the generated project")
	cmd.Flags().BoolVar(&quiet, "quiet", false, "Disable project generation messages")
	cmd.Flags().BoolVar(&git, "git", true, "Initialize git")
	cmd.Flags().BoolVar(&smtp, "smtp", false, "Add SMTP")
	cmd.Flags().BoolVar(&storage, "storage", false, "Add Storage (S3)")
	cmd.Flags().BoolVar(&redis, "redis", false, "Add Redis (comes with ratelimit middleware)")
	cmd.Flags().BoolVar(&auth, "auth", false, "Add Authentication (simple email-based)")
	cmd.Flags().BoolVar(&oauthDiscord, "oauth-discord", false, "Add Discord OAuth")
	cmd.Flags().BoolVar(&oauthFacebook, "oauth-facebook", false, "Add Facebook OAuth")
	cmd.Flags().BoolVar(&oauthGitHub, "oauth-github", false, "Add GitHub OAuth")
	cmd.Flags().BoolVar(&oauthGoogle, "oauth-google", false, "Add Google OAuth")
	cmd.Flags().BoolVar(&oauthInstagram, "oauth-instagram", false, "Add Instagram OAuth")
	cmd.Flags().BoolVar(&oauthLinkedIn, "oauth-linkedin", false, "Add LinkedIn OAuth")

	return cmd
}
