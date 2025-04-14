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
		database  string
		appType   string
		outputDir string
		noGit     bool
		noSMTP    bool
		noStorage bool
		noAuth    bool
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

			cfg := &initialize.Config{
				Name:      args[0],
				Database:  dbEnum,
				AppType:   appTypeEnum,
				NoGit:     noGit,
				OutputDir: outputDir,
				NoSMTP:    noSMTP,
				NoStorage: noStorage,
			}

			err := initialize.Initialize(cfg)
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}

	cmd.Flags().StringVarP(&appType, "appType", "t", "api", fmt.Sprintf("App type %v", initialize.AllAppTypes))
	cmd.Flags().StringVarP(&database, "database", "d", "sqlite3", fmt.Sprintf("Database type %v", initialize.AllDatabases))
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for the generated project")
	cmd.Flags().BoolVar(&noGit, "no-git", false, "Remove git")
	cmd.Flags().BoolVar(&noSMTP, "no-smtp", false, "Remove SMTP")
	cmd.Flags().BoolVar(&noStorage, "no-storage", false, "Remove Storage (S3)")
	cmd.Flags().BoolVar(&noAuth, "no-auth", false, "Remove Authentication (Authentication requires SMTP)")

	return cmd
}
