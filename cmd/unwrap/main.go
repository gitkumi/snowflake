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

	err = files.Create("devapi", true, path.Join(cwd, "template"))
	if err != nil {
		log.Fatal(err)
	}
}
