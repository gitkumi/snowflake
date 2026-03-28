package generate

import (
	"fmt"
	"strings"
	"unicode"
)

type Resource struct {
	Name            string
	NamePlural      string
	NameTitle       string
	NameTitlePlural string
	ModuleName      string
	Database        string
	Fields          []Field
}

type Field struct {
	Name      string
	NameTitle string
	Type      string
	SQLType   string
	GoType    string
	NullType  string
}

func NewResource(name string, fields []Field, cfg *SnowflakeConfig) *Resource {
	title := toTitle(name)
	plural := pluralize(name)
	titlePlural := toTitle(plural)

	return &Resource{
		Name:            name,
		NamePlural:      plural,
		NameTitle:       title,
		NameTitlePlural: titlePlural,
		ModuleName:      cfg.Module,
		Database:        cfg.Database,
		Fields:          fields,
	}
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
			NullType:  mapping.NullType,
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

func pluralize(s string) string {
	if s == "" {
		return s
	}

	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "z") ||
		strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") {
		return s + "es"
	}

	if strings.HasSuffix(s, "y") && len(s) > 1 {
		prev := rune(s[len(s)-2])
		if !isVowel(prev) {
			return s[:len(s)-1] + "ies"
		}
	}

	return s + "s"
}

func isVowel(r rune) bool {
	switch unicode.ToLower(r) {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}
