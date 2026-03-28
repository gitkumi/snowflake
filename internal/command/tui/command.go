package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
			projectPath := ""
			database := initialize.AllDatabases[0]
			keyValueStore := initialize.AllKeyValueStores[0]
			queue := initialize.AllQueues[0]
			containerRuntime := initialize.AllContainerRuntimes[0] // Defaults to Podman
			selectedFeatures := []string{"Git"}

			projectPathGroup := huh.NewGroup(
				huh.NewInput().
					Title("Project path").
					Placeholder("./acme").
					Value(&projectPath),
			)

			featuresGroup := huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Add features").
					Options(
						huh.NewOption("Git", "Git"),
						huh.NewOption("SMTP", "SMTP"),
						huh.NewOption("S3-compatible", "Storage"),
						huh.NewOption("HTML (templ)", "Templ"),
					).
					Value(&selectedFeatures),
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

			keyValueStoreGroup := huh.NewGroup(
				huh.NewSelect[initialize.KeyValueStore]().
					Title("Select key-value store").
					Options(
						huh.NewOption("None", initialize.KeyValueStoreNone),
						huh.NewOption("Redis", initialize.KeyValueStoreRedis),
						huh.NewOption("Valkey", initialize.KeyValueStoreValkey),
					).
					Value(&keyValueStore),
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

			containerRuntimeGroup := huh.NewGroup(
				huh.NewSelect[initialize.ContainerRuntime]().
					Title("Select container runtime").
					Options(
						huh.NewOption("Podman", initialize.ContainerRuntimePodman),
						huh.NewOption("Docker", initialize.ContainerRuntimeDocker),
					).
					Value(&containerRuntime),
			)

			initialForm := huh.NewForm(
				projectPathGroup,
				databaseGroup,
				keyValueStoreGroup,
				featuresGroup,
				queueGroup,
				containerRuntimeGroup,
			)

			if err := initialForm.Run(); err != nil {
				fmt.Printf("error running initial form: %v\n", err)
				return
			}

			name, outputDir, err := ParseProjectPath(projectPath)
			if err != nil {
				fmt.Printf("error parsing project path: %v\n", err)
				return
			}

			cfg.Name = name
			cfg.OutputDir = outputDir
			cfg.Database = database
			cfg.KeyValueStore = keyValueStore
			cfg.Queue = queue
			cfg.ContainerRuntime = containerRuntime

			cfg.Git = contains(selectedFeatures, "Git")
			cfg.SMTP = contains(selectedFeatures, "SMTP")
			cfg.Storage = contains(selectedFeatures, "Storage")
			cfg.Templ = contains(selectedFeatures, "Templ")

			if err := initialize.Run(cfg); err != nil {
				fmt.Printf("error creating project: %v\n", err)
			}
		},
	}

	return cmd
}

func ParseProjectPath(input string) (name string, outputDir string, err error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", "", fmt.Errorf("project path cannot be empty")
	}

	cleaned := filepath.Clean(input)
	if cleaned == "." || cleaned == "/" {
		return "", "", fmt.Errorf("project path must include a project name")
	}

	resolved := cleaned
	if !filepath.IsAbs(resolved) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", "", fmt.Errorf("failed to get current directory: %w", err)
		}
		resolved = filepath.Join(cwd, resolved)
	}

	return filepath.Base(resolved), filepath.Dir(resolved), nil
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
