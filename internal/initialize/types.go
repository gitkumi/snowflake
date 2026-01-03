package initialize

import (
	"fmt"
)

type Database string

const (
	DatabaseSQLite3  Database = "sqlite3"
	DatabasePostgres Database = "postgres"
	DatabaseMySQL    Database = "mysql"
	DatabaseMariaDB  Database = "mariadb"
	DatabaseNone     Database = "none"
)

var AllDatabases = []Database{
	DatabaseNone,
	DatabaseSQLite3,
	DatabasePostgres,
	DatabaseMySQL,
	DatabaseMariaDB,
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
	case DatabaseMariaDB:
		return fmt.Sprintf("mariadb:mariadb@tcp(localhost:3306)/%s?parseTime=true", projectName)
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
	case DatabaseMariaDB:
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
	case DatabaseMariaDB:
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
	case DatabaseMariaDB:
		return "github.com/go-sql-driver/mysql"
	default:
		return ""
	}
}

type Queue string

const (
	QueueNone Queue = "none"
	QueueSQS  Queue = "sqs"
)

var AllQueues = []Queue{
	QueueNone,
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

type ContainerRuntime string

const (
	ContainerRuntimePodman ContainerRuntime = "podman"
	ContainerRuntimeDocker ContainerRuntime = "docker"
)

var AllContainerRuntimes = []ContainerRuntime{
	ContainerRuntimePodman,
	ContainerRuntimeDocker,
}

func (c ContainerRuntime) IsValid() bool {
	for _, runtime := range AllContainerRuntimes {
		if runtime == c {
			return true
		}
	}
	return false
}

func (c ContainerRuntime) String() string {
	return string(c)
}
