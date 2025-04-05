package generator

import (
	"text/template"
)

func CreateTemplateFuncs(cfg *GeneratorConfig) template.FuncMap {
	return template.FuncMap{
		"DatabaseMigration": func(filename string) (string, error) {
			return LoadDatabaseMigration(cfg.Database, filename)
		},
		"DatabaseQuery": func(filename string) (string, error) {
			return LoadDatabaseQuery(cfg.Database, filename)
		},
	}
}
