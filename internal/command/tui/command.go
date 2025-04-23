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
			selectedAuthProviders := []string{}
			authType := initialize.AllAuthentications[0]

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
					Options(huh.NewOptions(initialize.AllDatabases...)...).
					Value(&database),
			)

			featuresGroup := huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Add additional features").
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
					Title("Select Queue").
					Options(
						huh.NewOption("None", initialize.BackgroundJobNone),
						huh.NewOption("Basic (sync.WaitGroup)", initialize.BackgroundJobBasic),
						huh.NewOption("SQS", initialize.BackgroundJobSQS),
					).
					Value(&backgroundJob),
			)

			authGroup := huh.NewGroup(
				huh.NewSelect[initialize.Authentication]().
					Title("Select Authentication").
					Options(
						huh.NewOption("None", initialize.AuthenticationNone),
						huh.NewOption("Email", initialize.AuthenticationEmail),
						huh.NewOption("Email with username", initialize.AuthenticationEmailWithUsername),
					).
					Value(&authType),
			)

			oauthProvidersGroup := huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Add OAuth Provider").
					Options(
						huh.NewOption("Google", "OAuthGoogle"),
						huh.NewOption("Facebook", "OAuthFacebook"),
						huh.NewOption("GitHub", "OAuthGitHub"),
						huh.NewOption("Discord", "OAuthDiscord"),
						huh.NewOption("Instagram", "OAuthInstagram"),
						huh.NewOption("LinkedIn", "OAuthLinkedIn"),
					).
					Value(&selectedAuthProviders),
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

			// Only ask for authentication if database is not 'none'
			if database != initialize.DatabaseNone {
				authForm := huh.NewForm(authGroup)
				if err := authForm.Run(); err != nil {
					fmt.Printf("error running auth form: %v\n", err)
					return
				}

				// Only ask for OAuth providers if authentication is not 'none'
				if authType != initialize.AuthenticationNone {
					oauthForm := huh.NewForm(oauthProvidersGroup)
					if err := oauthForm.Run(); err != nil {
						fmt.Printf("error running OAuth form: %v\n", err)
						return
					}
				}
			}

			featureEnabled := func(name string) bool {
				for _, f := range selectedFeatures {
					if f == name {
						return true
					}
				}
				return false
			}

			authProviderEnabled := func(name string) bool {
				for _, p := range selectedAuthProviders {
					if p == name {
						return true
					}
				}
				return false
			}

			cfg.Name = projectName
			cfg.AppType = appType
			cfg.Database = database
			cfg.BackgroundJob = backgroundJob
			cfg.Authentication = authType
			cfg.Git = featureEnabled("Git")
			cfg.SMTP = featureEnabled("SMTP")
			cfg.Storage = featureEnabled("Storage")
			cfg.Redis = featureEnabled("Redis")
			cfg.OAuthDiscord = authProviderEnabled("OAuthDiscord")
			cfg.OAuthFacebook = authProviderEnabled("OAuthFacebook")
			cfg.OAuthGitHub = authProviderEnabled("OAuthGitHub")
			cfg.OAuthGoogle = authProviderEnabled("OAuthGoogle")
			cfg.OAuthInstagram = authProviderEnabled("OAuthInstagram")
			cfg.OAuthLinkedIn = authProviderEnabled("OAuthLinkedIn")

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
