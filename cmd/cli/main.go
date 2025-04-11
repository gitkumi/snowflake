package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gitkumi/snowflake/internal/generator"
	"github.com/spf13/cobra"
)

// Will be set at build time
var (
	version = "dev"
)

func Execute() {
	var showVersion bool

	cmd := &cobra.Command{
		Use:   "snowflake",
		Short: "Snowflake is an opinionated Go web application generator.",
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Println(version)
				return
			}

			cmd.Help()
		},
	}

	cmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show snowflake version")

	cmd.Root().CompletionOptions.DisableDefaultCmd = true

	cmd.AddCommand(new())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func new() *cobra.Command {
	var (
		initGit   bool
		database  string
		appType   string
		outputDir string
		smtp      bool
		storage   bool
		auth      bool
	)

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// If no output directory is specified, use the current working directory
			if outputDir == "" {
				var err error
				outputDir, err = os.Getwd()
				if err != nil {
					log.Fatal(err.Error())
				}
			} else {
				// If the provided path is not absolute, make it absolute
				if !filepath.IsAbs(outputDir) {
					cwd, err := os.Getwd()
					if err != nil {
						log.Fatal(err.Error())
					}
					outputDir = filepath.Join(cwd, outputDir)
				}

				// Ensure output directory exists
				if _, err := os.Stat(outputDir); os.IsNotExist(err) {
					if err := os.MkdirAll(outputDir, 0755); err != nil {
						log.Fatalf("Failed to create output directory: %v", err)
					}
				}
			}

			dbEnum := generator.Database(database)
			if !dbEnum.IsValid() {
				log.Fatalf("Invalid database type: %s. Must be one of: %v", database, generator.AllDatabases)
			}

			appTypeEnum := generator.AppType(appType)
			if !appTypeEnum.IsValid() {
				log.Fatalf("Invalid app type: %s. Must be one of: %v", appType, generator.AllAppTypes)
			}

			cfg := &generator.GeneratorConfig{
				Name:      args[0],
				Database:  dbEnum,
				AppType:   appTypeEnum,
				InitGit:   initGit,
				OutputDir: outputDir,
				SMTP:      smtp,
				Storage:   storage,
			}

			err := generator.Generate(cfg)
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}

	cmd.Flags().StringVarP(&appType, "appType", "t", "api", fmt.Sprintf("App type %v", generator.AllAppTypes))
	cmd.Flags().StringVarP(&database, "database", "d", "sqlite3", fmt.Sprintf("Database type %v", generator.AllDatabases))
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for the generated project")
	cmd.Flags().BoolVar(&initGit, "git", true, "Initialize git")
	cmd.Flags().BoolVar(&smtp, "smtp", true, "Add smtp feature")
	cmd.Flags().BoolVar(&storage, "storage", true, "Add storage feature")
	cmd.Flags().BoolVar(&auth, "auth", true, "Add authentication feature")

	return cmd
}
