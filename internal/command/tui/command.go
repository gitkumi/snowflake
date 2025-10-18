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
			database := initialize.AllDatabases[0]
			queue := initialize.AllQueues[0]
			billing := initialize.BillingNone
			selectedFeatures := []string{"Git"}
			selectedOAuth := []string{}
			selectedOIDC := []string{}

			projectNameGroup := huh.NewGroup(
				huh.NewInput().
					Title("Enter project name").
					Placeholder("acme").
					Value(&projectName),
			)

			databaseGroup := huh.NewGroup(
				huh.NewSelect[initialize.Database]().
					Title("Select database").
					Options(
						huh.NewOption("None", initialize.DatabaseNone),
						huh.NewOption("SQLite3", initialize.DatabaseSQLite3),
						huh.NewOption("Postgres", initialize.DatabasePostgres),
						huh.NewOption("MySQL", initialize.DatabaseMySQL),
						huh.NewOption("MariaDB", initialize.DatabaseMariaDB),
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
						huh.NewOption("HTML", "HTML"),
					).
					Value(&selectedFeatures),
			)

			queueGroup := huh.NewGroup(
				huh.NewSelect[initialize.Queue]().
					Title("Select queue").
					Options(
						huh.NewOption("None", initialize.QueueNone),
						huh.NewOption("SQS", initialize.QueueSQS),
					).
					Value(&queue),
			)

			billingGroup := huh.NewGroup(
				huh.NewSelect[initialize.Billing]().
					Title("Add billing").
					Options(
						huh.NewOption("None", initialize.BillingNone),
						huh.NewOption("Stripe", initialize.BillingStripe),
					).
					Value(&billing),
			)

			oauthGroup := huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Add OAuth providers").
					Options(
						huh.NewOption("Google", "Google"),
						huh.NewOption("Discord", "Discord"),
						huh.NewOption("GitHub", "GitHub"),
						huh.NewOption("Instagram", "Instagram"),
						huh.NewOption("Microsoft", "Microsoft"),
						huh.NewOption("Reddit", "Reddit"),
						huh.NewOption("Spotify", "Spotify"),
						huh.NewOption("Twitch", "Twitch"),
						huh.NewOption("Facebook", "Facebook"),
						huh.NewOption("LinkedIn", "LinkedIn"),
						huh.NewOption("Slack", "Slack"),
						huh.NewOption("Stripe", "Stripe"),
						huh.NewOption("X", "X"),
					).
					Value(&selectedOAuth),
			)

			initialForm := huh.NewForm(
				projectNameGroup,
				databaseGroup,
				featuresGroup,
				queueGroup,
				billingGroup,
				oauthGroup,
			)

			if err := initialForm.Run(); err != nil {
				fmt.Printf("error running initial form: %v\n", err)
				return
			}

			oidcProviders := []string{"Facebook", "Google", "LinkedIn", "Microsoft", "Twitch", "Discord"}
			var oidcOptions []huh.Option[string]

			for _, provider := range oidcProviders {
				if contains(selectedOAuth, provider) {
					oidcOptions = append(oidcOptions, huh.NewOption(provider, provider))
				}
			}

			if len(oidcOptions) > 0 {
				oidcForm := huh.NewForm(
					huh.NewGroup(
						huh.NewMultiSelect[string]().
							Title("Add OIDC providers").
							Options(oidcOptions...).
							Value(&selectedOIDC),
					),
				)

				if err := oidcForm.Run(); err != nil {
					fmt.Printf("error running OIDC form: %v\n", err)
					return
				}
			}

			cfg.Name = projectName
			cfg.Database = database
			cfg.Queue = queue
			cfg.Billing = billing

			cfg.Git = contains(selectedFeatures, "Git")
			cfg.SMTP = contains(selectedFeatures, "SMTP")
			cfg.Storage = contains(selectedFeatures, "Storage")
			cfg.Redis = contains(selectedFeatures, "Redis")
			cfg.ServeHTML = contains(selectedFeatures, "HTML")

			cfg.OAuthGoogle = contains(selectedOAuth, "Google")
			cfg.OAuthDiscord = contains(selectedOAuth, "Discord")
			cfg.OAuthGitHub = contains(selectedOAuth, "GitHub")
			cfg.OAuthInstagram = contains(selectedOAuth, "Instagram")
			cfg.OAuthMicrosoft = contains(selectedOAuth, "Microsoft")
			cfg.OAuthReddit = contains(selectedOAuth, "Reddit")
			cfg.OAuthSpotify = contains(selectedOAuth, "Spotify")
			cfg.OAuthTwitch = contains(selectedOAuth, "Twitch")
			cfg.OAuthFacebook = contains(selectedOAuth, "Facebook")
			cfg.OAuthLinkedIn = contains(selectedOAuth, "LinkedIn")
			cfg.OAuthSlack = contains(selectedOAuth, "Slack")
			cfg.OAuthStripe = contains(selectedOAuth, "Stripe")
			cfg.OAuthX = contains(selectedOAuth, "X")

			cfg.OIDCFacebook = contains(selectedOIDC, "Facebook")
			cfg.OIDCGoogle = contains(selectedOIDC, "Google")
			cfg.OIDCLinkedIn = contains(selectedOIDC, "LinkedIn")
			cfg.OIDCMicrosoft = contains(selectedOIDC, "Microsoft")
			cfg.OIDCTwitch = contains(selectedOIDC, "Twitch")
			cfg.OIDCDiscord = contains(selectedOIDC, "Discord")

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
