package tui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/gitkumi/snowflake/internal/initialize"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new project using the TUI",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := &initialize.Config{}
			projectName := ""
			appType := initialize.AllAppTypes[0]
			database := initialize.AllDatabases[0]
			backgroundJob := initialize.AllBackgroundJobs[0]
			selectedFeatures := []string{"Git"}

			projectNameGroup := huh.NewGroup(
				huh.NewInput().
					Title("Enter project name").
					Placeholder("acme").
					Value(&projectName),
			)

			appTypeGroup := huh.NewGroup(
				huh.NewSelect[initialize.AppType]().
					Title("Select application type").
					Options(huh.NewOptions(initialize.AllAppTypes...)...).
					Value(&appType),
			)

			databaseGroup := huh.NewGroup(
				huh.NewSelect[initialize.Database]().
					Title("Select database").
					Options(
						huh.NewOption("None", initialize.DatabaseNone),
						huh.NewOption("SQLite3", initialize.DatabaseSQLite3),
						huh.NewOption("Postgres", initialize.DatabasePostgres),
						huh.NewOption("MySQL", initialize.DatabaseMySQL),
					).
					Value(&database),
			)

			featuresGroup := huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Add features").
					Options(
						huh.NewOption("Git", "Git"),
						huh.NewOption("SMTP", "SMTP"),
						huh.NewOption("S3", "Storage"),
						huh.NewOption("Redis", "Redis"),
					).
					Value(&selectedFeatures),
			)

			backgroundJobGroup := huh.NewGroup(
				huh.NewSelect[initialize.BackgroundJob]().
					Title("Select queue").
					Options(
						huh.NewOption("None", initialize.BackgroundJobNone),
						huh.NewOption("Basic (sync.WaitGroup)", initialize.BackgroundJobBasic),
						huh.NewOption("SQS", initialize.BackgroundJobSQS),
					).
					Value(&backgroundJob),
			)

			initialForm := huh.NewForm(
				projectNameGroup,
				appTypeGroup,
				databaseGroup,
				featuresGroup,
				backgroundJobGroup,
			)

			if err := initialForm.Run(); err != nil {
				fmt.Printf("error running form: %v\n", err)
				return
			}

			featureEnabled := func(name string) bool {
				for _, f := range selectedFeatures {
					if f == name {
						return true
					}
				}
				return false
			}

			cfg.Name = projectName
			cfg.AppType = appType
			cfg.Database = database
			cfg.BackgroundJob = backgroundJob
			cfg.Git = featureEnabled("Git")
			cfg.SMTP = featureEnabled("SMTP")
			cfg.Storage = featureEnabled("Storage")
			cfg.Redis = featureEnabled("Redis")

			if err := initialize.Run(cfg); err != nil {
				fmt.Printf("error creating project: %v\n", err)
			}
		},
	}

	return cmd
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
