package generator

import (
	"fmt"
	"os"
	"os/exec"
)

type Command struct {
	Message string
	Name    string
	Args    []string
}

func RunPostCommands(project *Project, outputPath string) error {
	commands := []Command{
		{"snowflake: go mod init", "go", []string{"mod", "init", project.Name}},
		{"snowflake: go mod tidy", "go", []string{"mod", "tidy"}},
		{"snowflake: gofmt", "gofmt", []string{"-w", "-s", "."}},
		{"snowflake: make build", "make", []string{"build"}},
	}

	for _, cmdDef := range commands {
		if err := RunCommand(outputPath, cmdDef); err != nil {
			return err
		}
	}
	return nil
}

func RunGitCommands(outputPath string) error {
	commands := []Command{
		{"", "git", []string{"init"}},
		{"", "git", []string{"add", "-A"}},
		{"", "git", []string{"commit", "-m", "Initialize Snowflake project"}},
	}

	fmt.Println("snowflake: initializing git")
	for _, cmdDef := range commands {
		if err := RunCommand(outputPath, cmdDef); err != nil {
			return err
		}
	}
	return nil
}

func RunCommand(workingDir string, command Command) error {
	if command.Message != "" {
		fmt.Println(command.Message)
	}

	cmd := exec.Command(command.Name, command.Args...)
	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
