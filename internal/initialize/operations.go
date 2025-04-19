package initialize

import (
	"fmt"
	"io"
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
		{
			Message: "snowflake: go mod init",
			Name:    "go",
			Args:    []string{"mod", "init", project.Name},
		},
		{
			Message: "snowflake: go mod tidy",
			Name:    "go",
			Args:    []string{"mod", "tidy"},
		},
		{
			Message: "snowflake: gofmt",
			Name:    "gofmt",
			Args:    []string{"-w", "-s", "."},
		},
		{
			Message: "snowflake: make build",
			Name:    "make",
			Args:    []string{"build"},
		},
	}

	for _, cmd := range commands {
		if _, err := exec.LookPath(cmd.Name); err != nil {
			return fmt.Errorf("%s is not installed or not found in PATH", cmd.Name)
		}

		if cmd.Message != "" && !quiet {
			fmt.Println(cmd.Message)
		}

		execCmd := exec.Command(cmd.Name, cmd.Args...)
		execCmd.Dir = outputPath

		if quiet {
			execCmd.Stdout = io.Discard
		} else {
			execCmd.Stdout = os.Stdout
		}
		execCmd.Stderr = os.Stderr

		if err := execCmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func runGitCommands(outputPath string, quiet bool) error {
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("git is not installed or not found in PATH. skipping git initialization")
		return nil
	}

	commands := []Command{
		{
			Message: "",
			Name:    "git",
			Args:    []string{"init"},
		},
		{
			Message: "",
			Name:    "git",
			Args:    []string{"add", "-A"},
		},
		{
			Message: "",
			Name:    "git",
			Args:    []string{"commit", "-m", "Initialize Snowflake project"},
		},
	}

	if !quiet {
		fmt.Println("snowflake: initializing git")
	}

	for _, cmd := range commands {
		if cmd.Message != "" && !quiet {
			fmt.Println(cmd.Message)
		}

		execCmd := exec.Command(cmd.Name, cmd.Args...)
		execCmd.Dir = outputPath

		if quiet {
			execCmd.Stdout = io.Discard
		} else {
			execCmd.Stdout = os.Stdout
		}
		execCmd.Stderr = os.Stderr

		if err := execCmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
