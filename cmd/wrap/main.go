package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	testPath := filepath.Join(cwd, "template", "devapi")

	buildFiles := []string{
		".git",
		"bin/main",
		"go.mod",
		"go.sum",
	}

	for _, fileName := range buildFiles {
		filePath := filepath.Join(testPath, fileName)
		os.Remove(filePath)
	}

	err = filepath.WalkDir(testPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		newContent := strings.ReplaceAll(string(data), "devapi", "{{ .Name }}")

		newFilePath := filePath + ".templ"

		fmt.Println("Created: " + newFilePath)
		err = os.WriteFile(newFilePath, []byte(newContent), 0644)
		if err != nil {
			return err
		}

		err = os.Remove(filePath)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
