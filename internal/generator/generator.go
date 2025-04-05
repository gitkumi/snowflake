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
	AppType  AppType
}

type GeneratorConfig struct {
	Name      string
	Database  Database
	AppType   AppType
	InitGit   bool
	OutputDir string
}

type FileExclusions struct {
	ByAppType map[AppType][]string
}

type FileRenaming struct {
	ByAppType map[AppType]map[string]string
}

func Generate(cfg *GeneratorConfig) error {
	project := &Project{
		Name:     cfg.Name,
		Database: cfg.Database,
		AppType:  cfg.AppType,
	}

	outputPath := filepath.Join(cfg.OutputDir, cfg.Name)
	templateFiles := snowflaketemplate.BaseFiles

	templateFuncs := createTemplateFuncs(cfg)
	exclusions := createFileExclusions()
	renaming := createFileRenaming()

	fmt.Println("Generating files...")
	err := fs.WalkDir(templateFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fileName := strings.TrimPrefix(path, "base")
		targetPath := filepath.Join(outputPath, fileName)

		if shouldExcludeFile(path, project, exclusions) {
			return nil
		}

		targetPath = getTargetPath(path, targetPath, project, renaming)

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

	if cfg.InitGit {
		if err := runGitCommands(outputPath); err != nil {
			return err
		}
	}

	fmt.Println("")
	fmt.Printf(`✅ Snowflake project '%s' generated successfully! 🎉

Run your new project:

  $ cd %s
  $ make dev
`, project.Name, project.Name)

	return nil
}

func createTemplateFuncs(cfg *GeneratorConfig) template.FuncMap {
	return template.FuncMap{
		"DatabaseMigration": func(filename string) (string, error) {
			return LoadDatabaseMigration(cfg.Database, filename)
		},
		"DatabaseQuery": func(filename string) (string, error) {
			return LoadDatabaseQuery(cfg.Database, filename)
		},
	}
}

func createFileExclusions() *FileExclusions {
	return &FileExclusions{
		ByAppType: map[AppType][]string{
			API: {
				"/internal/html",
				".templ.templ",
			},
		},
	}
}

func createFileRenaming() *FileRenaming {
	return &FileRenaming{
		ByAppType: map[AppType]map[string]string{
			Web: {
				"/cmd/api/": "/cmd/web/",
				"/cmd/api/main.go": "/cmd/web/main.go",
			},
		},
	}
}

func shouldExcludeFile(path string, project *Project, exclusions *FileExclusions) bool {
	if excludedPaths, ok := exclusions.ByAppType[project.AppType]; ok {
		for _, excludedPath := range excludedPaths {
			if strings.Contains(path, excludedPath) {
				return true
			}
		}
	}

	return false
}

func getTargetPath(srcPath, originalTargetPath string, project *Project, renaming *FileRenaming) string {
	if renamingRules, ok := renaming.ByAppType[project.AppType]; ok {
		targetPath := originalTargetPath
		for pattern, replacement := range renamingRules {
			if strings.Contains(srcPath, pattern) {
				return strings.Replace(targetPath, pattern, replacement, 1)
			}
		}
	}
	return originalTargetPath
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
