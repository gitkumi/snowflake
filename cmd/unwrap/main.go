package main

import (
	"log"
	"os"
	"path"

	"github.com/gitkumi/snowflake/internal/generator"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = generator.Generate("acme", true, path.Join(cwd, "template"), generator.SQLite3)
	if err != nil {
		log.Fatal(err)
	}
}
