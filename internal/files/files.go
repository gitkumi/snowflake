package files

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	snowflaketemplate "github.com/gitkumi/snowflake/template"
)

const (
	SQLite3  Database = "sqlite3"
	Postgres Database = "postgres"
	MySQL    Database = "mysql"
)

type Project struct {
	Name     string
	Database Database
}

type Config struct {
	Name      string
	Database  Database
	Git       bool
	OutputDir string
}

func Create(cfg *Config) error {
	project := &Project{
		Name:     strings.ToLower(cfg.Name),
		Database: cfg.Database,
	}

	outputPath := filepath.Join(cfg.OutputDir, cfg.Name)

	templateFiles := snowflaketemplate.ApiFiles

	fmt.Println("Creating files..")
	err := fs.WalkDir(templateFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fileName := strings.TrimPrefix(path, "api")

		if d.IsDir() {
			err := os.MkdirAll(filepath.Join(outputPath, fileName), 0777)
			return err
		}

		content, err := templateFiles.ReadFile(path)
		if err != nil {
			return err
		}

		tmpl, err := template.New(fileName).Parse(string(content))
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, project)
		if err != nil {
			return err
		}

		newFilePath := strings.TrimSuffix(filepath.Join(outputPath, fileName), ".templ")
		err = os.WriteFile(newFilePath, buf.Bytes(), 0777)
		return err
	})

	if err != nil {
		return err
	}

	fmt.Println("Running go mod init..")
	command := exec.Command("go", "mod", "init", project.Name)
	command.Dir = outputPath
	err = command.Run()
	if err != nil {
		return err
	}

	fmt.Println("Running go mod tidy..")
	command = exec.Command("go", "mod", "tidy")
	command.Dir = outputPath
	err = command.Run()
	if err != nil {
		return err
	}

	fmt.Println("Running gofmt..")
	command = exec.Command("gofmt", "-w", "-s", ".")
	command.Dir = outputPath
	err = command.Run()
	if err != nil {
		return err
	}

	fmt.Println("Running make build..")
	command = exec.Command("make", "build")
	command.Dir = outputPath
	err = command.Run()
	if err != nil {
		return err
	}

	if cfg.Git {
		fmt.Println("Running git init..")
		command = exec.Command("git", "init")
		command.Dir = outputPath
		err = command.Run()
		if err != nil {
			return err
		}

		fmt.Println("Running git add..")
		command = exec.Command("git", "add", "-A")
		command.Dir = outputPath
		err = command.Run()
		if err != nil {
			return err
		}

		fmt.Println("Running git commit..")
		command = exec.Command("git", "commit", "-m", "Initialize Snowflake project")
		command.Dir = outputPath
		err = command.Run()
		if err != nil {
			return err
		}
	}

	fmt.Println("Done!")
	return nil
}
