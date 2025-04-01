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
	Name string
}

func Generate(projectName string, initGit bool, outputDir string) error {
	project := &Project{
		Name: strings.ToLower(projectName),
	}
	outputPath := filepath.Join(outputDir, projectName)
	templateFiles := snowflaketemplate.ProjectFiles

	fmt.Println("Generating files...")
	err := fs.WalkDir(templateFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Remove the "project" prefix from the template path.
		fileName := strings.TrimPrefix(path, "project")
		targetPath := filepath.Join(outputPath, fileName)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0777)
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
		if err := tmpl.Execute(&buf, project); err != nil {
			return err
		}

		// Remove the ".templ" extension if present.
		newFilePath := strings.TrimSuffix(targetPath, ".templ")
		return os.WriteFile(newFilePath, buf.Bytes(), 0777)
	})
	if err != nil {
		return err
	}

	// Run post-generation commands.
	if err := runPostCommands(project, outputPath); err != nil {
		return err
	}

	// Optionally run Git commands.
	if initGit {
		if err := runGitCommands(outputPath); err != nil {
			return err
		}
	}

	fmt.Println("Done!")
	return nil
}

// runCmd executes a command with the specified working directory, message, command name, and arguments.
func runCmd(workingDir, message, name string, args ...string) error {
	fmt.Println(message)
	cmd := exec.Command(name, args...)
	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runPostCommands executes the Go module and build commands.
func runPostCommands(project *Project, outputPath string) error {
	commands := []struct {
		message string
		name    string
		args    []string
	}{
		{"Running go mod init...", "go", []string{"mod", "init", project.Name}},
		{"Running go mod tidy...", "go", []string{"mod", "tidy"}},
		{"Running gofmt...", "gofmt", []string{"-w", "-s", "."}},
		{"Running make build...", "make", []string{"build"}},
	}

	for _, cmdDef := range commands {
		if err := runCmd(outputPath, cmdDef.message, cmdDef.name, cmdDef.args...); err != nil {
			return err
		}
	}
	return nil
}

// runGitCommands executes the Git commands to initialize a repository and commit the changes.
func runGitCommands(outputPath string) error {
	commands := []struct {
		message string
		name    string
		args    []string
	}{
		{"Running git init...", "git", []string{"init"}},
		{"Running git add...", "git", []string{"add", "-A"}},
		{"Running git commit...", "git", []string{"commit", "-m", "Initialize Snowflake project"}},
	}

	for _, cmdDef := range commands {
		if err := runCmd(outputPath, cmdDef.message, cmdDef.name, cmdDef.args...); err != nil {
			return err
		}
	}
	return nil
}
