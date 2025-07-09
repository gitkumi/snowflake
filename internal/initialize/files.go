package initialize

import (
	"fmt"
	"os"
	"slices"
)

func RemoveEmptyDirs(paths map[string]bool) error {
	// Sort paths by length in descending order to process deeper directories first
	var sortedPaths []string
	for dir := range paths {
		sortedPaths = append(sortedPaths, dir)
	}
	slices.SortFunc(sortedPaths, func(a, b string) int {
		if len(a) > len(b) {
			return -1
		} else if len(a) < len(b) {
			return 1
		}
		return 0
	})

	for _, dir := range sortedPaths {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// Directory doesn't exist, skip it
			continue
		} else if err != nil {
			return fmt.Errorf("failed to check if directory %s exists: %v", dir, err)
		}

		isEmpty, err := IsDirectoryEmpty(dir)
		if err != nil {
			return fmt.Errorf("failed to check if directory %s is empty: %v", dir, err)
		}
		if isEmpty {
			if err := os.Remove(dir); err != nil {
				return fmt.Errorf("failed to remove empty directory %s: %v", dir, err)
			}
		}
	}
	return nil
}

func IsDirectoryEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err != nil {
		return true, nil
	}

	return false, nil
}
