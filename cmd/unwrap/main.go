package main

import (
	"log"
	"os"
	"path"

	"github.com/gitkumi/snowflake/internal/files"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = files.Create(&files.Config{
		Name:      "devapi",
		Git:       true,
		OutputDir: path.Join(cwd, "template"),
		Database:  "sqlite3",
	})
	if err != nil {
		log.Fatal(err)
	}
}
