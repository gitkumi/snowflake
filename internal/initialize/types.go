package initialize

import (
	"fmt"
)

type AppType string

const (
	AppTypeAPI AppType = "api"
	AppTypeWeb AppType = "web"
)

var AllAppTypes = []AppType{
	AppTypeAPI,
	AppTypeWeb,
}

func (t AppType) IsValid() bool {
	for _, appType := range AllAppTypes {
		if appType == t {
			return true
		}
	}
	return false
}

func (t AppType) String() string {
	return string(t)
}

type Database string

const (
	DatabaseSQLite3  Database = "sqlite3"
	DatabasePostgres Database = "postgres"
	DatabaseMySQL    Database = "mysql"
	DatabaseNone     Database = "none"
)

var AllDatabases = []Database{
	DatabaseNone,
	DatabaseSQLite3,
	DatabasePostgres,
	DatabaseMySQL,
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
	case DatabaseSQLite3:
		return projectName + "_dev.db"
	case DatabasePostgres:
		return fmt.Sprintf("user=postgres password=postgres dbname=%s host=localhost port=5432 sslmode=disable", projectName)
	case DatabaseMySQL:
		return fmt.Sprintf("mysql:mysql@tcp(localhost:3306)/%s?parseTime=true", projectName)
	default:
		return ""
	}
}

func (d Database) Driver() string {
	switch d {
	case DatabaseSQLite3:
		return "sqlite3"
	case DatabasePostgres:
		return "postgres"
	case DatabaseMySQL:
		return "mysql"
	default:
		return ""
	}
}

func (d Database) SQLCEngine() string {
	switch d {
	case DatabaseSQLite3:
		return "sqlite"
	case DatabasePostgres:
		return "postgresql"
	case DatabaseMySQL:
		return "mysql"
	default:
		return ""
	}
}

func (d Database) Import() string {
	switch d {
	case DatabaseSQLite3:
		return "github.com/mattn/go-sqlite3"
	case DatabasePostgres:
		return "github.com/lib/pq"
	case DatabaseMySQL:
		return "github.com/go-sql-driver/mysql"
	default:
		return ""
	}
}

type Queue string

const (
	QueueNone  Queue = "none"
	QueueBasic Queue = "basic"
	QueueSQS   Queue = "sqs"
)

var AllQueues = []Queue{
	QueueNone,
	QueueBasic,
	QueueSQS,
}

func (t Queue) IsValid() bool {
	for _, bg := range AllQueues {
		if bg == t {
			return true
		}
	}
	return false
}

func (t Queue) String() string {
	return string(t)
}
