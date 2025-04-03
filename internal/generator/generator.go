package generator

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

type Project struct {
	Name     string
	Database Database
}

func Generate(projectName string, initGit bool, outputDir string, db Database) error {
	project := &Project{
		Name:     strings.ToLower(projectName),
		Database: db,
	}

	templateFuncs := template.FuncMap{
		"DatabaseMigration": func(filename string) (string, error) {
			return LoadDatabaseMigration(db, filename)
		},
		"DatabaseQuery": func(filename string) (string, error) {
			return LoadDatabaseQuery(db, filename)
		},
	}

	outputPath := filepath.Join(outputDir, projectName)
	templateFiles := snowflaketemplate.BaseFiles

	fmt.Println("Generating files...")
	err := fs.WalkDir(templateFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fileName := strings.TrimPrefix(path, "base")
		targetPath := filepath.Join(outputPath, fileName)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0777)
		}

		content, err := templateFiles.ReadFile(path)
		if err != nil {
			return err
		}

		tmpl, err := template.New(fileName).Funcs(templateFuncs).Parse(string(content))
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, project); err != nil {
			return err
		}

		newFilePath := strings.TrimSuffix(targetPath, ".templ")
		return os.WriteFile(newFilePath, buf.Bytes(), 0777)
	})
	if err != nil {
		return err
	}

	if err := runPostCommands(project, outputPath); err != nil {
		return err
	}

	if initGit {
		if err := runGitCommands(outputPath); err != nil {
			return err
		}
	}

	fmt.Println("")
	fmt.Printf(`Snowflake project generated successfully.
You can use "make" to install dependencies, run the dev server, and more:

    $ cd %s

If you don't have the required dev packages installed yet (air, sqlc, goose):

    $ make deps.get

Then start the dev server:

    $ make dev
`, project.Name)

	return nil
}

func runPostCommands(project *Project, outputPath string) error {
	commands := []struct {
		message string
		name    string
		args    []string
	}{
		{"snowflake: go mod init", "go", []string{"mod", "init", project.Name}},
		{"snowflake: go mod tidy", "go", []string{"mod", "tidy"}},
		{"snowflake: gofmt", "gofmt", []string{"-w", "-s", "."}},
		{"snowflake: make build", "make", []string{"build"}},
	}

	for _, cmdDef := range commands {
		if err := runCmd(outputPath, cmdDef.message, cmdDef.name, cmdDef.args...); err != nil {
			return err
		}
	}
	return nil
}

func runGitCommands(outputPath string) error {
	commands := []struct {
		message string
		name    string
		args    []string
	}{
		{"", "git", []string{"init"}},
		{"", "git", []string{"add", "-A"}},
		{"", "git", []string{"commit", "-m", "Initialize Snowflake project"}},
	}

	fmt.Println("snowflake: initializing git")
	for _, cmdDef := range commands {
		if err := runCmd(outputPath, cmdDef.message, cmdDef.name, cmdDef.args...); err != nil {
			return err
		}
	}
	return nil
}

func runCmd(workingDir, message, name string, args ...string) error {
	if message != "" {
		fmt.Println(message)
	}
	cmd := exec.Command(name, args...)
	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
