package generator

// Database represents the type of database to use
type Database string

const (
	SQLite3  Database = "sqlite3"
	Postgres Database = "postgres"
	MySQL    Database = "mysql"
)

// AllDatabases contains all supported database types
var AllDatabases = []Database{
	SQLite3,
	Postgres,
	MySQL,
}

// String returns the string representation of the database type
func (d Database) String() string {
	return string(d)
}

// IsValid checks if the database type is supported
func (d Database) IsValid() bool {
	for _, db := range AllDatabases {
		if db == d {
			return true
		}
	}
	return false
}

// Driver returns the database driver name
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

// Project represents the project configuration
type Project struct {
	Name     string
	Database Database
}

