package generate

import (
	"fmt"
	"strings"
	"unicode"
)

type Resource struct {
	Name       string
	NameTitle  string
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

func NewResource(name string, plural string, fields []Field, cfg *ProjectConfig) *Resource {
	return &Resource{
		Name:       name,
		NameTitle:  toTitle(name),
		PluralName: plural,
		ModuleName: cfg.Module,
		Database:   cfg.Database,
		Fields:     fields,
	}
}


func ParseFields(rawFields []string, database string) ([]Field, error) {
	fields := make([]Field, 0, len(rawFields))
	for _, raw := range rawFields {
		parts := strings.SplitN(raw, ":", 2)
		if len(parts) != 2 || parts[1] == "" {
			return nil, fmt.Errorf("invalid field %q, expected name:type format (e.g. title:string)\nValid types: string, text, int, bigint, bool, float, timestamp", raw)
		}

		name := parts[0]
		if name == "" {
			return nil, fmt.Errorf("empty field name in %q", raw)
		}

		typeName := parts[1]

		mapping, ok := typeMapping[typeName]
		if !ok {
			return nil, fmt.Errorf("unknown field type %q in %q\nValid types: string, text, int, bigint, bool, float, timestamp", typeName, raw)
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

	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-'
	})
	if len(parts) == 0 {
		return s
	}

	var builder strings.Builder
	for _, part := range parts {
		runes := []rune(part)
		if len(runes) == 0 {
			continue
		}
		runes[0] = unicode.ToUpper(runes[0])
		builder.WriteString(string(runes))
	}

	return builder.String()
}
