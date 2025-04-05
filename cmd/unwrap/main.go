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

	cfg := &generator.GeneratorConfig{
		Name:      "acme",
		Database:  generator.SQLite3,
		AppType:   generator.API,
		InitGit:   true,
		OutputDir: path.Join(cwd, "template"),
	}

	err = generator.Generate(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
