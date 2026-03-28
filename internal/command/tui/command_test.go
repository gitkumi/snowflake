package tui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseProjectPath(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		desc      string
		input     string
		wantName  string
		wantDir   string
		wantError bool
	}{
		{
			desc:     "simple name",
			input:    "acme",
			wantName: "acme",
			wantDir:  cwd,
		},
		{
			desc:     "relative path with dot",
			input:    "./acme",
			wantName: "acme",
			wantDir:  cwd,
		},
		{
			desc:     "nested relative path",
			input:    "./foo/bar/acme",
			wantName: "acme",
			wantDir:  filepath.Join(cwd, "foo", "bar"),
		},
		{
			desc:     "relative path without dot",
			input:    "foo/acme",
			wantName: "acme",
			wantDir:  filepath.Join(cwd, "foo"),
		},
		{
			desc:     "absolute path",
			input:    "/tmp/projects/acme",
			wantName: "acme",
			wantDir:  "/tmp/projects",
		},
		{
			desc:     "absolute path single component",
			input:    "/acme",
			wantName: "acme",
			wantDir:  "/",
		},
		{
			desc:     "trailing slash is cleaned",
			input:    "./acme/",
			wantName: "acme",
			wantDir:  cwd,
		},
		{
			desc:     "dot-dot in path",
			input:    "./foo/../acme",
			wantName: "acme",
			wantDir:  cwd,
		},
		{
			desc:     "whitespace is trimmed",
			input:    "  acme  ",
			wantName: "acme",
			wantDir:  cwd,
		},
		{
			desc:     "name with hyphen",
			input:    "my-project",
			wantName: "my-project",
			wantDir:  cwd,
		},
		{
			desc:     "name with underscore",
			input:    "my_project",
			wantName: "my_project",
			wantDir:  cwd,
		},
		{
			desc:     "deeply nested path",
			input:    "/a/b/c/d/e/project",
			wantName: "project",
			wantDir:  "/a/b/c/d/e",
		},
		{
			desc:      "empty string",
			input:     "",
			wantError: true,
		},
		{
			desc:      "whitespace only",
			input:     "   ",
			wantError: true,
		},
		{
			desc:      "just a dot",
			input:     ".",
			wantError: true,
		},
		{
			desc:      "just a slash",
			input:     "/",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			name, dir, err := ParseProjectPath(tt.input)

			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error, got name=%q dir=%q", name, dir)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if name != tt.wantName {
				t.Errorf("name = %q, want %q", name, tt.wantName)
			}

			if dir != tt.wantDir {
				t.Errorf("dir = %q, want %q", dir, tt.wantDir)
			}
		})
	}
}
