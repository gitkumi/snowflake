package generator

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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

type FileRenames struct {
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
	renames := createFileRenames()

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

		targetDir := filepath.Dir(newFilePath)
		if err := os.MkdirAll(targetDir, 0777); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", targetDir, err)
		}

		return os.WriteFile(newFilePath, buf.Bytes(), 0777)
	})
	if err != nil {
		return err
	}

	if err := processFileRenames(project, outputPath, renames); err != nil {
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

func createFileRenames() *FileRenames {
	return &FileRenames{
		ByAppType: map[AppType]map[string]string{
			Web: {
				"/cmd/api/main.go": "/cmd/web/main.go",
			},
		},
	}
}

func processFileRenames(project *Project, outputPath string, renames *FileRenames) error {
	if renameMappings, ok := renames.ByAppType[project.AppType]; ok {
		// Track directories that might need cleanup
		dirsToCheck := make(map[string]bool)

		// Process all file renames
		for oldPath, newPath := range renameMappings {
			if err := renameFile(outputPath, oldPath, newPath); err != nil {
				return err
			}

			// Add source directory to cleanup list
			sourceDir := filepath.Dir(filepath.Join(outputPath, oldPath))
			dirsToCheck[sourceDir] = true
		}

		// Clean up empty directories
		if err := cleanupEmptyDirs(dirsToCheck); err != nil {
			return err
		}
	}
	return nil
}

func renameFile(basePath, oldRelPath, newRelPath string) error {
	fullOldPath := filepath.Join(basePath, oldRelPath)
	fullNewPath := filepath.Join(basePath, newRelPath)

	if _, err := os.Stat(fullOldPath); os.IsNotExist(err) {
		return nil
	}

	targetDir := filepath.Dir(fullNewPath)
	if err := os.MkdirAll(targetDir, 0777); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", targetDir, err)
	}

	data, err := os.ReadFile(fullOldPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", fullOldPath, err)
	}

	if err := os.WriteFile(fullNewPath, data, 0666); err != nil {
		return fmt.Errorf("failed to write file %s: %v", fullNewPath, err)
	}

	if err := os.Remove(fullOldPath); err != nil {
		return fmt.Errorf("failed to remove file %s: %v", fullOldPath, err)
	}

	return nil
}

func cleanupEmptyDirs(dirs map[string]bool) error {
	var dirList []string
	for dir := range dirs {
		dirList = append(dirList, dir)
	}

	sort.Slice(dirList, func(i, j int) bool {
		return len(dirList[i]) > len(dirList[j])
	})

	for _, dir := range dirList {
		// Check if directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		// Check if directory is empty
		entries, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("failed to read directory %s: %v", dir, err)
		}

		// Remove if empty
		if len(entries) == 0 {
			if err := os.Remove(dir); err != nil {
				return fmt.Errorf("failed to remove empty directory %s: %v", dir, err)
			}
		}
	}

	return nil
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
