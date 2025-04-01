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

	err = generator.Generate("devproject", true, path.Join(cwd, "template"))
	if err != nil {
		log.Fatal(err)
	}
}
