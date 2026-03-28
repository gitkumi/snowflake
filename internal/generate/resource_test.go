package generate

import (
	"testing"
)

func TestPluralize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"post", "posts"},
		{"category", "categories"},
		{"address", "addresses"},
		{"bus", "buses"},
		{"box", "boxes"},
		{"church", "churches"},
		{"dish", "dishes"},
		{"quiz", "quizes"},
		{"key", "keys"},
		{"day", "days"},
		{"toy", "toys"},
		{"user", "users"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := pluralize(tt.input)
			if got != tt.expected {
				t.Errorf("pluralize(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestToTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"post", "Post"},
		{"user", "User"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toTitle(tt.input)
			if got != tt.expected {
				t.Errorf("toTitle(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseFields(t *testing.T) {
	fields, err := ParseFields([]string{"title:string", "body:text", "count:int", "active:bool"}, "postgres")
	if err != nil {
		t.Fatal(err)
	}

	if len(fields) != 4 {
		t.Fatalf("expected 4 fields, got %d", len(fields))
	}

	if fields[0].Name != "title" || fields[0].SQLType != "TEXT" || fields[0].GoType != "string" {
		t.Errorf("unexpected first field: %+v", fields[0])
	}

	if fields[2].Name != "count" || fields[2].SQLType != "INTEGER" || fields[2].GoType != "int64" {
		t.Errorf("unexpected third field: %+v", fields[2])
	}
}

func TestParseFieldsDefaultType(t *testing.T) {
	fields, err := ParseFields([]string{"name"}, "postgres")
	if err != nil {
		t.Fatal(err)
	}

	if fields[0].Type != "string" {
		t.Errorf("expected default type 'string', got %q", fields[0].Type)
	}
}

func TestParseFieldsInvalidType(t *testing.T) {
	_, err := ParseFields([]string{"name:invalid"}, "postgres")
	if err == nil {
		t.Fatal("expected error for invalid type")
	}
}

func TestParseFieldsEmptyName(t *testing.T) {
	_, err := ParseFields([]string{":string"}, "postgres")
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNewResource(t *testing.T) {
	cfg := &ProjectConfig{Module: "acme", Database: "postgres"}
	fields := []Field{{Name: "title", NameTitle: "Title", Type: "string", SQLType: "TEXT", GoType: "string"}}

	r := NewResource("post", fields, cfg)

	if r.Name != "post" {
		t.Errorf("expected Name 'post', got %q", r.Name)
	}
	if r.PluralName != "posts" {
		t.Errorf("expected PluralName 'posts', got %q", r.PluralName)
	}
	if r.NameTitle != "Post" {
		t.Errorf("expected NameTitle 'Post', got %q", r.NameTitle)
	}
}
