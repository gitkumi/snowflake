package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileExclusions struct {
	ByAppType map[AppType][]string
}

type FileRenames struct {
	ByAppType map[AppType]map[string]string
}

func CreateFileExclusions() *FileExclusions {
	return &FileExclusions{
		ByAppType: map[AppType][]string{
			API: {
				"/internal/html",
				"/internal/application/handler/html_handler.go",
				".templ.templ",
			},
		},
	}
}

func CreateFileRenames() *FileRenames {
	return &FileRenames{
		ByAppType: map[AppType]map[string]string{
			Web: {
				"/cmd/api/main.go": "/cmd/web/main.go",
			},
		},
	}
}

func ShouldExcludeFile(path string, project *Project, exclusions *FileExclusions) bool {
	if excludedPaths, ok := exclusions.ByAppType[project.AppType]; ok {
		for _, excludedPath := range excludedPaths {
			if strings.Contains(path, excludedPath) {
				return true
			}
		}
	}
	return false
}

func ProcessFileRenames(project *Project, outputPath string, renames *FileRenames) error {
	renameMappings, ok := renames.ByAppType[project.AppType]
	if !ok {
		return nil
	}
	
	sourceDirs := make(map[string]bool)
	
	for oldPath, newPath := range renameMappings {
		fullOldPath := filepath.Join(outputPath, oldPath)
		fullNewPath := filepath.Join(outputPath, newPath)
		
		if _, err := os.Stat(fullOldPath); os.IsNotExist(err) {
			continue
		}
		
		targetDir := filepath.Dir(fullNewPath)
		if err := os.MkdirAll(targetDir, 0777); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", targetDir, err)
		}
		
		if err := os.Rename(fullOldPath, fullNewPath); err != nil {
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
		}
		
		sourceDir := filepath.Dir(fullOldPath)
		sourceDirs[sourceDir] = true
	}
	
	return cleanupSourceDirs(sourceDirs)
}

func cleanupSourceDirs(dirs map[string]bool) error {
	var dirList []string
	for dir := range dirs {
		dirList = append(dirList, dir)
	}
	
	// Sort by length in descending order to ensure child directories
	// are processed before their parents
	sort.Slice(dirList, func(i, j int) bool {
		return len(dirList[i]) > len(dirList[j])
	})
	
	for _, dir := range dirList {
		if err := cleanupEmptyDir(dir); err != nil {
			return err
		}
	}
	
	return nil
}

func cleanupEmptyDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", dir, err)
	}

	// Process subdirectories first
	for _, entry := range entries {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			if err := cleanupEmptyDir(subdir); err != nil {
				return err
			}
		}
	}

	entries, err = os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", dir, err)
	}

	if len(entries) == 0 {
		if err := os.Remove(dir); err != nil {
			return fmt.Errorf("failed to remove empty directory %s: %v", dir, err)
		}
	}

	return nil
}
