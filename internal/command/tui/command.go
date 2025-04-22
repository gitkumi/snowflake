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
			authEnabled := false
			selectedAuthProviders := []string{}

			// Create the initial form groups
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
					Title("Select additional features").
					Options(
						huh.NewOption("Git", "Git"),
						huh.NewOption("SMTP", "SMTP"),
						huh.NewOption("Storage (S3)", "Storage"),
						huh.NewOption("Redis", "Redis"),
					).
					Value(&selectedFeatures),
			)

			backgroundJobGroup := huh.NewGroup(
				huh.NewSelect[initialize.BackgroundJob]().
					Title("Select background job").
					Options(huh.NewOptions(initialize.AllBackgroundJobs...)...).
					Value(&backgroundJob),
			)

			authGroup := huh.NewGroup(
				huh.NewConfirm().
					Title("Enable Email Authentication").
					Value(&authEnabled),
			)

			oauthProvidersGroup := huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Add OAuth Provider").
					Options(
						huh.NewOption("Discord OAuth", "OAuthDiscord"),
						huh.NewOption("Facebook OAuth", "OAuthFacebook"),
						huh.NewOption("GitHub OAuth", "OAuthGitHub"),
						huh.NewOption("Google OAuth", "OAuthGoogle"),
						huh.NewOption("Instagram OAuth", "OAuthInstagram"),
						huh.NewOption("LinkedIn OAuth", "OAuthLinkedIn"),
					).
					Value(&selectedAuthProviders),
			)

			initialForm := huh.NewForm(
				projectNameGroup,
				appTypeGroup,
				databaseGroup,
				featuresGroup,
				backgroundJobGroup,
				authGroup,
			)

			if err := initialForm.Run(); err != nil {
				fmt.Printf("error running form: %v\n", err)
				return
			}

			// If auth is enabled, show the OAuth providers form
			if authEnabled {
				oauthForm := huh.NewForm(oauthProvidersGroup)
				if err := oauthForm.Run(); err != nil {
					fmt.Printf("error running OAuth form: %v\n", err)
					return
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
			cfg.Git = featureEnabled("Git")
			cfg.SMTP = featureEnabled("SMTP")
			cfg.Storage = featureEnabled("Storage")
			cfg.Redis = featureEnabled("Redis")
			cfg.Auth = authEnabled
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
