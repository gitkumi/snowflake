package files

import "fmt"

type Database string

func (d *Database) Set(value string) error {
	switch value {
	case string(SQLite3), string(Postgres), string(MySQL):
		*d = Database(value)
		return nil
	default:
		return fmt.Errorf("invalid database %q: must be one of 'sqlite3', 'postgres', or 'mysql'", value)
	}
}

func (d *Database) String() string {
	return string(*d)
}

func (d *Database) Type() string {
	return "database"
}
