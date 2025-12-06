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
		quiet     bool
		database  string
		queue     string
		outputDir string
		git       bool
		smtp      bool
		storage   bool
		redis     bool
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

			queueEnum := initialize.Queue(queue)
			if !queueEnum.IsValid() {
				log.Fatalf("Invalid queue type: %s. Must be one of: %v", queue, initialize.AllQueues)
			}

			err := initialize.Run(&initialize.Config{
				Quiet:     quiet,
				Name:      args[0],
				Database:  dbEnum,
				Queue:     queueEnum,
				Git:       git,
				OutputDir: outputDir,
				SMTP:      smtp,
				Storage:   storage,
				Redis:     redis,
			})
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}

	cmd.Flags().StringVarP(&database, "database", "d", "none", fmt.Sprintf("Database type %v", initialize.AllDatabases))
	cmd.Flags().StringVarP(&queue, "queue", "q", "none", fmt.Sprintf("Queue type %v", initialize.AllQueues))
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for the generated project")
	cmd.Flags().BoolVar(&quiet, "quiet", false, "Disable project generation messages")
	cmd.Flags().BoolVar(&git, "git", true, "Initialize git")
	cmd.Flags().BoolVar(&smtp, "smtp", false, "Add SMTP")
	cmd.Flags().BoolVar(&storage, "storage", false, "Add Storage (S3)")
	cmd.Flags().BoolVar(&redis, "redis", false, "Add Redis (comes with ratelimit middleware)")

	return cmd
}
