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
		quiet         bool
		database      string
		backgroundJob string
		appType       string
		outputDir     string
		withGit       bool
		withSMTP      bool
		withStorage   bool
		withRedis     bool
		withAuth      bool
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
				WithGit:       withGit,
				OutputDir:     outputDir,
				WithSMTP:      withSMTP,
				WithStorage:   withStorage,
				WithRedis:     withRedis,
				WithAuth:      withAuth,
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
	cmd.Flags().BoolVar(&withGit, "with-git", true, "Initialize git")
	cmd.Flags().BoolVar(&withSMTP, "with-smtp", false, "Add SMTP")
	cmd.Flags().BoolVar(&withStorage, "with-storage", false, "Add Storage (S3)")
	cmd.Flags().BoolVar(&withRedis, "with-redis", false, "Add Redis (comes with ratelimit middleware)")
	cmd.Flags().BoolVar(&withAuth, "with-auth", false, "Add Authentication (simple email-based)")

	return cmd
}
