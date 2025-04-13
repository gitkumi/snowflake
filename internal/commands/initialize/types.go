package initialize

import (
	"fmt"
)

type AppType string

const (
	Web AppType = "web"
	API AppType = "api"
)

var AllAppTypes = []AppType{
	Web,
	API,
}

func (t AppType) IsValid() bool {
	for _, appType := range AllAppTypes {
		if appType == t {
			return true
		}
	}
	return false
}

type Database string

const (
	SQLite3  Database = "sqlite3"
	Postgres Database = "postgres"
	MySQL    Database = "mysql"
	// None     Database = "none"
)

var AllDatabases = []Database{
	SQLite3,
	Postgres,
	MySQL,
	// None,
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

func (d Database) ConnString(projectName string) string {
	switch d {
	case SQLite3:
		return projectName + "_dev.db"
	case Postgres:
		return fmt.Sprintf("user=postgres password=postgres dbname=%s host=localhost port=5432 sslmode=disable", projectName)
	case MySQL:
		return fmt.Sprintf("mysql:mysql@tcp(localhost:3306)/%s?parseTime=true", projectName)
	default:
		return ""
	}
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
