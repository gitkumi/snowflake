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
			selectedFeatures := []string{"Git", "SMTP", "Storage"}

			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Enter project name").
						Placeholder("acme").
						Value(&projectName),

					huh.NewSelect[initialize.AppType]().
						Title("Select application type").
						Options(huh.NewOptions(initialize.AllAppTypes...)...).
						Value(&appType),

					huh.NewSelect[initialize.Database]().
						Title("Select database").
						Options(huh.NewOptions(initialize.AllDatabases...)...).
						Value(&database),

					huh.NewSelect[initialize.BackgroundJob]().
						Title("Select background job").
						Options(huh.NewOptions(initialize.AllBackgroundJobs[]...)...).
						Value(&backgroundJob),

					huh.NewMultiSelect[string]().
						Title("Select features").
						Options(
							huh.NewOption("Git", "Git"),
							huh.NewOption("SMTP", "SMTP"),
							huh.NewOption("Storage", "Storage"),
							huh.NewOption("Redis", "Redis"),
							huh.NewOption("Auth", "Auth"),
						).
						Value(&selectedFeatures),
				),
			)

			if err := form.Run(); err != nil {
				fmt.Printf("error running form: %v\n", err)
				return
			}

			// Map selected features to config flags
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
			cfg.NoGit = !featureEnabled("Git")
			cfg.NoSMTP = !featureEnabled("SMTP")
			cfg.NoStorage = !featureEnabled("Storage")
			cfg.NoRedis = !featureEnabled("Redis")
			cfg.NoAuth = !featureEnabled("Auth")

			if err := initialize.Run(cfg); err != nil {
				fmt.Printf("error creating project: %v\n", err)
			}
		},
	}

	return cmd
}
