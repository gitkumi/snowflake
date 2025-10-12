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
		queue          string
		outputDir      string
		git            bool
		smtp           bool
		storage        bool
		redis          bool
		serveHTML      bool
		oauthGoogle    bool
		oauthDiscord   bool
		oauthGitHub    bool
		oauthInstagram bool
		oauthMicrosoft bool
		oauthReddit    bool
		oauthSpotify   bool
		oauthTwitch    bool
		oauthFacebook  bool
		oauthLinkedIn  bool
		oauthSlack     bool
		oauthStripe    bool
		oauthX         bool
		oidcFacebook   bool
		oidcGoogle     bool
		oidcLinkedIn   bool
		oidcMicrosoft  bool
		oidcTwitch     bool
		oidcDiscord    bool
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
				Quiet:          quiet,
				Name:           args[0],
				Database:       dbEnum,
				Queue:          queueEnum,
				Git:            git,
				OutputDir:      outputDir,
				SMTP:           smtp,
				Storage:        storage,
				Redis:          redis,
				ServeHTML:      serveHTML,
				OAuthGoogle:    oauthGoogle,
				OAuthDiscord:   oauthDiscord,
				OAuthGitHub:    oauthGitHub,
				OAuthInstagram: oauthInstagram,
				OAuthMicrosoft: oauthMicrosoft,
				OAuthReddit:    oauthReddit,
				OAuthSpotify:   oauthSpotify,
				OAuthTwitch:    oauthTwitch,
				OAuthFacebook:  oauthFacebook,
				OAuthLinkedIn:  oauthLinkedIn,
				OAuthSlack:     oauthSlack,
				OAuthStripe:    oauthStripe,
				OAuthX:         oauthX,
				OIDCFacebook:   oidcFacebook,
				OIDCGoogle:     oidcGoogle,
				OIDCLinkedIn:   oidcLinkedIn,
				OIDCMicrosoft:  oidcMicrosoft,
				OIDCTwitch:     oidcTwitch,
				OIDCDiscord:    oidcDiscord,
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
	cmd.Flags().BoolVar(&serveHTML, "html", false, "Serve HTML with templ")

	cmd.Flags().BoolVar(&oauthGoogle, "oauth-google", false, "Add Google OAuth")
	cmd.Flags().BoolVar(&oauthDiscord, "oauth-discord", false, "Add Discord OAuth")
	cmd.Flags().BoolVar(&oauthGitHub, "oauth-github", false, "Add GitHub OAuth")
	cmd.Flags().BoolVar(&oauthInstagram, "oauth-instagram", false, "Add Instagram OAuth")
	cmd.Flags().BoolVar(&oauthMicrosoft, "oauth-microsoft", false, "Add Microsoft OAuth")
	cmd.Flags().BoolVar(&oauthReddit, "oauth-reddit", false, "Add Reddit OAuth")
	cmd.Flags().BoolVar(&oauthSpotify, "oauth-spotify", false, "Add Spotify OAuth")
	cmd.Flags().BoolVar(&oauthTwitch, "oauth-twitch", false, "Add Twitch OAuth")
	cmd.Flags().BoolVar(&oauthFacebook, "oauth-facebook", false, "Add Facebook OAuth")
	cmd.Flags().BoolVar(&oauthLinkedIn, "oauth-linkedin", false, "Add LinkedIn OAuth")
	cmd.Flags().BoolVar(&oauthSlack, "oauth-slack", false, "Add Slack OAuth")
	cmd.Flags().BoolVar(&oauthStripe, "oauth-stripe", false, "Add Stripe OAuth")
	cmd.Flags().BoolVar(&oauthX, "oauth-x", false, "Add X OAuth")
	cmd.Flags().BoolVar(&oidcFacebook, "oidc-facebook", false, "Add Facebook OIDC")
	cmd.Flags().BoolVar(&oidcGoogle, "oidc-google", false, "Add Google OIDC")
	cmd.Flags().BoolVar(&oidcLinkedIn, "oidc-linkedin", false, "Add LinkedIn OIDC")
	cmd.Flags().BoolVar(&oidcMicrosoft, "oidc-microsoft", false, "Add Microsoft OIDC")
	cmd.Flags().BoolVar(&oidcTwitch, "oidc-twitch", false, "Add Twitch OIDC")
	cmd.Flags().BoolVar(&oidcDiscord, "oidc-discord", false, "Add Discord OIDC")

	return cmd
}
