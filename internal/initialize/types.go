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

func ParseDatabase(value string) (Database, error) {
	database := Database(strings.TrimSpace(value))
	if database == "" {
		database = DatabaseNone
	}
	if !database.IsValid() {
		return "", fmt.Errorf("invalid database type: %s. Must be one of: %v", value, AllDatabases)
	}
	return database, nil
}

func ParseKeyValueStore(value string) (KeyValueStore, error) {
	store := KeyValueStore(strings.TrimSpace(value))
	if store == "" {
		store = KeyValueStoreNone
	}
	if !store.IsValid() {
		return "", fmt.Errorf("invalid key-value store: %s. Must be one of: %v", value, AllKeyValueStores)
	}
	return store, nil
}

func ParseContainerRuntime(value string) (ContainerRuntime, error) {
	runtime := ContainerRuntime(strings.TrimSpace(value))
	if runtime == "" {
		runtime = ContainerRuntimePodman
	}
	if !runtime.IsValid() {
		return "", fmt.Errorf("invalid container runtime: %s. Must be one of: %v", value, AllContainerRuntimes)
	}
	return runtime, nil
}
