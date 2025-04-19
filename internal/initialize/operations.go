package initialize

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

func runPostCommands(project *Project, outputPath string, quiet bool) error {
	commands := []Command{
		{"snowflake: go mod init", "go", []string{"mod", "init", project.Name}},
		{"snowflake: go mod tidy", "go", []string{"mod", "tidy"}},
		{"snowflake: gofmt", "gofmt", []string{"-w", "-s", "."}},
		{"snowflake: make build", "make", []string{"build"}},
	}

	for _, cmdDef := range commands {
		if err := runCommand(outputPath, cmdDef, quiet); err != nil {
			return err
		}
	}
	return nil
}

func runGitCommands(outputPath string, quiet bool) error {
	commands := []Command{
		{"", "git", []string{"init"}},
		{"", "git", []string{"add", "-A"}},
		{"", "git", []string{"commit", "-m", "Initialize Snowflake project"}},
	}

	if quiet {
		for i, cmd := range commands {
			commands[i].Args = append([]string{cmd.Args[0], "-q"}, cmd.Args[1:]...)
		}
	}

	if !quiet {
		fmt.Println("snowflake: initializing git")
	}
	for _, cmdDef := range commands {
		if err := runCommand(outputPath, cmdDef, quiet); err != nil {
			return err
		}
	}
	return nil
}

func runCommand(workingDir string, command Command, quiet bool) error {
	if command.Message != "" && !quiet {
		fmt.Println(command.Message)
	}

	cmd := exec.Command(command.Name, command.Args...)
	cmd.Dir = workingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
