package generate

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type ProjectConfig struct {
	Module   string
	Database string
}

func LoadConfig(dir string) (*ProjectConfig, error) {
	module, err := readModule(dir)
	if err != nil {
		return nil, err
	}

	database, err := readDatabase(dir)
	if err != nil {
		return nil, err
	}

	return &ProjectConfig{
		Module:   module,
		Database: database,
	}, nil
}

func readModule(dir string) (string, error) {
	f, err := os.Open(filepath.Join(dir, "go.mod"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("go.mod not found in %s", dir)
		}
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	return "", fmt.Errorf("could not parse module name from go.mod")
}

type sqlcConfig struct {
	SQL []struct {
		Engine string `yaml:"engine"`
	} `yaml:"sql"`
}

func readDatabase(dir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(dir, "sqlc.yaml"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("sqlc.yaml not found - is this a Snowflake project with a database?")
		}
		return "", fmt.Errorf("failed to read sqlc.yaml: %w", err)
	}

	var cfg sqlcConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "", fmt.Errorf("failed to parse sqlc.yaml: %w", err)
	}

	if len(cfg.SQL) == 0 || cfg.SQL[0].Engine == "" {
		return "", fmt.Errorf("no database engine found in sqlc.yaml")
	}

	switch cfg.SQL[0].Engine {
	case "postgresql":
		return "postgres", nil
	case "mysql":
		return "mysql", nil
	case "sqlite":
		return "sqlite3", nil
	default:
		return "", fmt.Errorf("unsupported sqlc engine %q", cfg.SQL[0].Engine)
	}
}
