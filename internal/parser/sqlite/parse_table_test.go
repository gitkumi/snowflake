package parser_test

import (
	"testing"

	parser "github.com/gitkumi/snowflake/internal/parser/sqlite"
)

func TestParseTable(t *testing.T) {
	tests := []struct {
		name     string
		sqlInput string
		want     []parser.TableSchema
		wantErr  bool
	}{
		{
			name: "simple table",
			sqlInput: `CREATE TABLE users (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT NOT NULL,
                email TEXT NOT NULL UNIQUE,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )`,
			want: []parser.TableSchema{
				{
					TableName: "users",
					Columns: []parser.ColumnDef{
						{
							Name:            "id",
							Type:            "INTEGER",
							Constraints:     []string{"PRIMARY KEY", "AUTOINCREMENT"},
							IsPrimaryKey:    true,
							IsAutoIncrement: true,
							IsNotNull:       false,
						},
						{
							Name:            "name",
							Type:            "TEXT",
							Constraints:     []string{"NOT NULL"},
							IsPrimaryKey:    false,
							IsAutoIncrement: false,
							IsNotNull:       true,
						},
						{
							Name:            "email",
							Type:            "TEXT",
							Constraints:     []string{"NOT NULL", "UNIQUE"},
							IsPrimaryKey:    false,
							IsAutoIncrement: false,
							IsNotNull:       true,
						},
						{
							Name:            "created_at",
							Type:            "TIMESTAMP",
							Constraints:     []string{"DEFAULT CURRENT_TIMESTAMP"},
							IsPrimaryKey:    false,
							IsAutoIncrement: false,
							IsNotNull:       false,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple tables",
			sqlInput: `CREATE TABLE users (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT NOT NULL
            );
            
            CREATE TABLE posts (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                user_id INTEGER NOT NULL,
                title TEXT NOT NULL,
                content TEXT
            )`,
			want: []parser.TableSchema{
				{
					TableName: "users",
					Columns: []parser.ColumnDef{
						{
							Name:            "id",
							Type:            "INTEGER",
							Constraints:     []string{"PRIMARY KEY", "AUTOINCREMENT"},
							IsPrimaryKey:    true,
							IsAutoIncrement: true,
							IsNotNull:       false,
						},
						{
							Name:            "name",
							Type:            "TEXT",
							Constraints:     []string{"NOT NULL"},
							IsPrimaryKey:    false,
							IsAutoIncrement: false,
							IsNotNull:       true,
						},
					},
				},
				{
					TableName: "posts",
					Columns: []parser.ColumnDef{
						{
							Name:            "id",
							Type:            "INTEGER",
							Constraints:     []string{"PRIMARY KEY", "AUTOINCREMENT"},
							IsPrimaryKey:    true,
							IsAutoIncrement: true,
							IsNotNull:       false,
						},
						{
							Name:            "user_id",
							Type:            "INTEGER",
							Constraints:     []string{"NOT NULL"},
							IsPrimaryKey:    false,
							IsAutoIncrement: false,
							IsNotNull:       true,
						},
						{
							Name:            "title",
							Type:            "TEXT",
							Constraints:     []string{"NOT NULL"},
							IsPrimaryKey:    false,
							IsAutoIncrement: false,
							IsNotNull:       true,
						},
						{
							Name:            "content",
							Type:            "TEXT",
							Constraints:     []string{},
							IsPrimaryKey:    false,
							IsAutoIncrement: false,
							IsNotNull:       false,
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseTable(tt.sqlInput)
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
				t.Errorf("different number of tables parsed: want %d, got %d", len(tt.want), len(got))
				return
			}

			for i, wantTable := range tt.want {
				if i < len(got) {
					gotTable := got[i]
					if wantTable.TableName != gotTable.TableName {
						t.Errorf("table name mismatch: want %q, got %q", wantTable.TableName, gotTable.TableName)
					}

					if len(wantTable.Columns) != len(gotTable.Columns) {
						t.Errorf("different number of columns: want %d, got %d", len(wantTable.Columns), len(gotTable.Columns))
						continue
					}

					for j, wantColumn := range wantTable.Columns {
						if j < len(gotTable.Columns) {
							gotColumn := gotTable.Columns[j]
							if wantColumn.Name != gotColumn.Name {
								t.Errorf("column name mismatch: want %q, got %q", wantColumn.Name, gotColumn.Name)
							}
							if wantColumn.Type != gotColumn.Type {
								t.Errorf("column type mismatch: want %q, got %q", wantColumn.Type, gotColumn.Type)
							}
							if wantColumn.IsPrimaryKey != gotColumn.IsPrimaryKey {
								t.Errorf("IsPrimaryKey mismatch: want %v, got %v", wantColumn.IsPrimaryKey, gotColumn.IsPrimaryKey)
							}
							if wantColumn.IsAutoIncrement != gotColumn.IsAutoIncrement {
								t.Errorf("IsAutoIncrement mismatch: want %v, got %v", wantColumn.IsAutoIncrement, gotColumn.IsAutoIncrement)
							}
							if wantColumn.IsNotNull != gotColumn.IsNotNull {
								t.Errorf("IsNotNull mismatch: want %v, got %v", wantColumn.IsNotNull, gotColumn.IsNotNull)
							}

							// Compare constraints with relaxed ordering
							if !compareStringSlices(wantColumn.Constraints, gotColumn.Constraints) {
								t.Errorf("constraints mismatch: want %v, got %v", wantColumn.Constraints, gotColumn.Constraints)
							}
						}
					}
				}
			}
		})
	}
}
