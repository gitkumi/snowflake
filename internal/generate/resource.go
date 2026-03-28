package generate

import (
	"fmt"
	"strings"
	"unicode"
)

type Resource struct {
	Name      string
	NameTitle string
	PluralName string

	ModuleName string
	Database   string
	Fields     []Field
}

type Field struct {
	Name      string
	NameTitle string
	Type      string
	SQLType   string
	GoType    string
}

func NewResource(name string, fields []Field, cfg *ProjectConfig) *Resource {
	return &Resource{
		Name:       name,
		NameTitle:  toTitle(name),
		PluralName: pluralize(name),
		ModuleName: cfg.Module,
		Database:   cfg.Database,
		Fields:     fields,
	}
}

func pluralize(s string) string {
	if s == "" {
		return s
	}
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "z") ||
		strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") {
		return s + "es"
	}
	if strings.HasSuffix(s, "y") && len(s) > 1 && !isVowel(rune(s[len(s)-2])) {
		return s[:len(s)-1] + "ies"
	}
	return s + "s"
}

func isVowel(r rune) bool {
	switch r {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}

func ParseFields(rawFields []string, database string) ([]Field, error) {
	fields := make([]Field, 0, len(rawFields))
	for _, raw := range rawFields {
		parts := strings.SplitN(raw, ":", 2)
		name := parts[0]
		if name == "" {
			return nil, fmt.Errorf("empty field name in %q", raw)
		}

		typeName := "string"
		if len(parts) == 2 {
			typeName = parts[1]
		}

		mapping, ok := typeMapping[typeName]
		if !ok {
			return nil, fmt.Errorf("unknown field type %q (valid: string, text, int, bigint, bool, float, timestamp)", typeName)
		}

		sqlType, ok := mapping.SQLTypes[database]
		if !ok {
			return nil, fmt.Errorf("unsupported database %q", database)
		}

		fields = append(fields, Field{
			Name:      name,
			NameTitle: toTitle(name),
			Type:      typeName,
			SQLType:   sqlType,
			GoType:    mapping.GoType,
		})
	}
	return fields, nil
}

func toTitle(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

