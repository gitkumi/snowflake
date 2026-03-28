package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type SnowflakeConfig struct {
	Module   string `yaml:"name"`
	Database string `yaml:"database"`
}

func LoadConfig(dir string) (*SnowflakeConfig, error) {
	data, err := os.ReadFile(filepath.Join(dir, ".snowflake.yaml"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not a Snowflake project: .snowflake.yaml not found in %s", dir)
		}
		return nil, fmt.Errorf("failed to read .snowflake.yaml: %w", err)
	}

	var cfg SnowflakeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse .snowflake.yaml: %w", err)
	}

	if cfg.Module == "" {
		return nil, fmt.Errorf("invalid .snowflake.yaml: name is required")
	}

	if cfg.Database == "" {
		return nil, fmt.Errorf("invalid .snowflake.yaml: database is required")
	}

	return &cfg, nil
}
