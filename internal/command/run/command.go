package run

import (
	"fmt"
	"log"

	"github.com/gitkumi/snowflake/internal/initialize"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var (
		quiet            bool
		database         string
		keyValueStore    string
		containerRuntime string
		outputDir        string
		git              bool
		smtp             bool
		storage          bool
		templ            bool
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Create a new project using command-line flags",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dbEnum, err := initialize.ParseDatabase(database)
			if err != nil {
				log.Fatal(err)
			}

			containerRuntimeEnum, err := initialize.ParseContainerRuntime(containerRuntime)
			if err != nil {
				log.Fatal(err)
			}

			kvsEnum, err := initialize.ParseKeyValueStore(keyValueStore)
			if err != nil {
				log.Fatal(err)
			}

			err = initialize.Run(&initialize.Config{
				Quiet:            quiet,
				Name:             args[0],
				Database:         dbEnum,
				KeyValueStore:    kvsEnum,
				ContainerRuntime: containerRuntimeEnum,
				Git:              git,
				OutputDir:        outputDir,
				SMTP:             smtp,
				Storage:          storage,
				Templ:            templ,
			})
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}

	cmd.Flags().StringVarP(&database, "database", "d", "none", fmt.Sprintf("Database type %v", initialize.AllDatabases))
	cmd.Flags().StringVarP(&containerRuntime, "container", "c", "podman", fmt.Sprintf("Container runtime %v", initialize.AllContainerRuntimes))
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for the generated project")
	cmd.Flags().BoolVar(&quiet, "quiet", false, "Disable project generation messages")
	cmd.Flags().BoolVar(&git, "git", true, "Initialize git")
	cmd.Flags().StringVar(&keyValueStore, "kvs", "none", fmt.Sprintf("Key-value store %v", initialize.AllKeyValueStores))
	cmd.Flags().BoolVar(&smtp, "smtp", false, "Add SMTP")
	cmd.Flags().BoolVar(&storage, "storage", false, "Add Storage (S3)")
	cmd.Flags().BoolVar(&templ, "templ", false, "Add HTML (templ)")

	return cmd
}
