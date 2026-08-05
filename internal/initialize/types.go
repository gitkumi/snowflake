package initialize

import (
	"fmt"
	"strings"
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

// SQLDriver returns the database/sql driver name passed to sql.Open, matching
// the driver registered by the imported package. This differs from Driver
// (the goose dialect): PostgreSQL uses the "pgx" driver but the "postgres"
// goose dialect.
func (d Database) SQLDriver() string {
	switch d {
	case DatabaseSQLite3:
		return "sqlite3"
	case DatabasePostgres:
		return "pgx"
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
		return "github.com/jackc/pgx/v5/stdlib"
	case DatabaseMySQL:
		return "github.com/go-sql-driver/mysql"
	case DatabaseMariaDB:
		return "github.com/go-sql-driver/mysql"
	default:
		return ""
	}
}

type KeyValueStore string

const (
	KeyValueStoreNone   KeyValueStore = "none"
	KeyValueStoreRedis  KeyValueStore = "redis"
	KeyValueStoreValkey KeyValueStore = "valkey"
)

var AllKeyValueStores = []KeyValueStore{
	KeyValueStoreNone,
	KeyValueStoreRedis,
	KeyValueStoreValkey,
}

func (k KeyValueStore) IsValid() bool {
	for _, kvs := range AllKeyValueStores {
		if kvs == k {
			return true
		}
	}
	return false
}

func (k KeyValueStore) String() string {
	return string(k)
}

type JobProcessor string

const (
	JobProcessorNone   JobProcessor = "none"
	JobProcessorAbsurd JobProcessor = "absurd"
)

var AllJobProcessors = []JobProcessor{
	JobProcessorNone,
	JobProcessorAbsurd,
}

func (j JobProcessor) IsValid() bool {
	for _, jp := range AllJobProcessors {
		if jp == j {
			return true
		}
	}
	return false
}

func (j JobProcessor) String() string {
	return string(j)
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

type enum interface {
	~string
	IsValid() bool
}

func parseEnum[T enum](value string, defaultValue T, typeName string, allValues any) (T, error) {
	v := T(strings.TrimSpace(value))
	if v == "" {
		v = defaultValue
	}
	if !v.IsValid() {
		return "", fmt.Errorf("invalid %s: %s. Must be one of: %v", typeName, value, allValues)
	}
	return v, nil
}

func ParseDatabase(value string) (Database, error) {
	return parseEnum[Database](value, DatabaseNone, "database type", AllDatabases)
}

func ParseKeyValueStore(value string) (KeyValueStore, error) {
	return parseEnum[KeyValueStore](value, KeyValueStoreNone, "key-value store", AllKeyValueStores)
}

func ParseContainerRuntime(value string) (ContainerRuntime, error) {
	return parseEnum[ContainerRuntime](value, ContainerRuntimePodman, "container runtime", AllContainerRuntimes)
}

func ParseJobProcessor(value string) (JobProcessor, error) {
	return parseEnum[JobProcessor](value, JobProcessorNone, "job processor", AllJobProcessors)
}
