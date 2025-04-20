package generate

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type SQLCConfig struct {
	Version string `yaml:"version"`
	SQL     []struct {
		Engine  string `yaml:"engine"`
		Schema  string `yaml:"schema"`
		Queries string `yaml:"queries"`
		Gen     struct {
			Go struct {
				Package string `yaml:"package"`
				Out     string `yaml:"out"`
			} `yaml:"go"`
		} `yaml:"gen"`
	} `yaml:"sql"`
}

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "gen",
		Short: "Generated CRUD",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
}

func run() {
	data, err := os.ReadFile("testdata/sqlc.yaml")
	if err != nil {
		fmt.Println("Failed to read file:", err)
		return
	}

	var config SQLCConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Println("Failed to parse YAML:", err)
		return
	}

	fmt.Println(config.SQL[0].Gen.Go.Out)
}
