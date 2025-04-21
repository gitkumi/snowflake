package parser_test

import (
	"reflect"
	"testing"

	parser "github.com/gitkumi/snowflake/internal/parser/sqlite"
)

func TestParseQueries(t *testing.T) {
	tests := []struct {
		name     string
		sqlInput string
		want     []parser.Query
		wantErr  bool
	}{
		{
			name: "select query",
			sqlInput: `-- name: GetUser :one
SELECT * FROM users WHERE id = ?`,
			want: []parser.Query{
				{
					Name:     "GetUser",
					Type:     parser.One,
					SQL:      "SELECT * FROM users WHERE id = ?",
					Table:    "users",
					Params:   []string{"?"},
					Columns:  []string{"*"},
					IsSelect: true,
				},
			},
			wantErr: false,
		},
		{
			name: "insert query",
			sqlInput: `-- name: CreateUser :exec
INSERT INTO users (name, email) VALUES (?, ?)`,
			want: []parser.Query{
				{
					Name:     "CreateUser",
					Type:     parser.Exec,
					SQL:      "INSERT INTO users (name, email) VALUES (?, ?)",
					Table:    "users",
					Params:   []string{"?", "?"},
					Columns:  []string{"name", "email"},
					IsInsert: true,
				},
			},
			wantErr: false,
		},
		{
			name: "update query",
			sqlInput: `-- name: UpdateUser :exec
UPDATE users SET name = ?, email = ? WHERE id = ?`,
			want: []parser.Query{
				{
					Name:     "UpdateUser",
					Type:     parser.Exec,
					SQL:      "UPDATE users SET name = ?, email = ? WHERE id = ?",
					Table:    "users",
					Params:   []string{"?", "?", "?"},
					IsUpdate: true,
				},
			},
			wantErr: false,
		},
		{
			name: "delete query",
			sqlInput: `-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?`,
			want: []parser.Query{
				{
					Name:     "DeleteUser",
					Type:     parser.Exec,
					SQL:      "DELETE FROM users WHERE id = ?",
					Table:    "users",
					Params:   []string{"?"},
					IsDelete: true,
				},
			},
			wantErr: false,
		},
		{
			name: "multiple queries",
			sqlInput: `-- name: GetUser :one
SELECT id, name, email FROM users WHERE id = ?

-- name: ListUsers :many
SELECT * FROM users ORDER BY name

-- name: CreateUser :exec
INSERT INTO users (name, email) VALUES (?, ?)`,
			want: []parser.Query{
				{
					Name:     "GetUser",
					Type:     parser.One,
					SQL:      "SELECT id, name, email FROM users WHERE id = ?",
					Table:    "users",
					Params:   []string{"?"},
					Columns:  []string{"id", "name", "email"},
					IsSelect: true,
				},
				{
					Name:     "ListUsers",
					Type:     parser.Many,
					SQL:      "SELECT * FROM users ORDER BY name",
					Table:    "users",
					Params:   []string{},
					Columns:  []string{"*"},
					IsSelect: true,
				},
				{
					Name:     "CreateUser",
					Type:     parser.Exec,
					SQL:      "INSERT INTO users (name, email) VALUES (?, ?)",
					Table:    "users",
					Params:   []string{"?", "?"},
					Columns:  []string{"name", "email"},
					IsInsert: true,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseQueries(tt.sqlInput)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(tt.want) != len(got) {
				t.Errorf("different number of queries parsed: want %d, got %d", len(tt.want), len(got))
				return
			}

			for i, wantQuery := range tt.want {
				if i < len(got) {
					gotQuery := got[i]
					if wantQuery.Name != gotQuery.Name {
						t.Errorf("query name mismatch: want %q, got %q", wantQuery.Name, gotQuery.Name)
					}
					if wantQuery.Type != gotQuery.Type {
						t.Errorf("query type mismatch: want %q, got %q", wantQuery.Type, gotQuery.Type)
					}
					if wantQuery.SQL != gotQuery.SQL {
						t.Errorf("SQL mismatch: want %q, got %q", wantQuery.SQL, gotQuery.SQL)
					}
					if wantQuery.Table != gotQuery.Table {
						t.Errorf("table mismatch: want %q, got %q", wantQuery.Table, gotQuery.Table)
					}
					if !reflect.DeepEqual(wantQuery.Params, gotQuery.Params) {
						t.Errorf("params mismatch: want %v, got %v", wantQuery.Params, gotQuery.Params)
					}
					if wantQuery.IsSelect != gotQuery.IsSelect {
						t.Errorf("IsSelect mismatch: want %v, got %v", wantQuery.IsSelect, gotQuery.IsSelect)
					}
					if wantQuery.IsInsert != gotQuery.IsInsert {
						t.Errorf("IsInsert mismatch: want %v, got %v", wantQuery.IsInsert, gotQuery.IsInsert)
					}
					if wantQuery.IsUpdate != gotQuery.IsUpdate {
						t.Errorf("IsUpdate mismatch: want %v, got %v", wantQuery.IsUpdate, gotQuery.IsUpdate)
					}
					if wantQuery.IsDelete != gotQuery.IsDelete {
						t.Errorf("IsDelete mismatch: want %v, got %v", wantQuery.IsDelete, gotQuery.IsDelete)
					}

					// Only compare columns for INSERT and SELECT statements
					if wantQuery.IsInsert || wantQuery.IsSelect {
						if !compareStringSlices(wantQuery.Columns, gotQuery.Columns) {
							t.Errorf("columns mismatch: want %v, got %v", wantQuery.Columns, gotQuery.Columns)
						}
					}
				}
			}
		})
	}
}
