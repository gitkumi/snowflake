package generator

import (
	"fmt"
	"path/filepath"

	snowflaketemplate "github.com/gitkumi/snowflake/template"
)

type Database string

const (
	SQLite3  Database = "sqlite3"
	Postgres Database = "postgres"
	MySQL    Database = "mysql"
)

var AllDatabases = []Database{
	SQLite3,
	Postgres,
	MySQL,
}

func (d Database) String() string {
	return string(d)
}

func (d Database) IsValid() bool {
	for _, db := range AllDatabases {
		if db == d {
			return true
		}
	}
	return false
}

func (d Database) Driver() string {
	switch d {
	case SQLite3:
		return "sqlite3"
	case Postgres:
		return "postgres"
	case MySQL:
		return "mysql"
	default:
		return ""
	}
}

func (d Database) SQLCEngine() string {
	switch d {
	case SQLite3:
		return "sqlite"
	case Postgres:
		return "postgresql"
	case MySQL:
		return "mysql"
	default:
		return ""
	}
}

func (d Database) Import() string {
	switch d {
	case SQLite3:
		return "github.com/mattn/go-sqlite3"
	case Postgres:
		return "github.com/lib/pq"
	case MySQL:
		return "github.com/go-sql-driver/mysql"
	default:
		return ""
	}
}

func LoadDatabaseMigration(db Database, filename string) (string, error) {
	fragmentPath := filepath.Join("fragments/database", string(db), "migrations", filename)
	content, err := snowflaketemplate.DatabaseFragments.ReadFile(fragmentPath)
	if err != nil {
		return "", fmt.Errorf("failed to read database fragment: %w", err)
	}
	return string(content), nil
}
