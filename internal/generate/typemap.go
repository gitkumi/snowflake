package generate

type TypeInfo struct {
	SQLTypes map[string]string
	GoType   string
}

var typeMapping = map[string]TypeInfo{
	"string": {
		SQLTypes: map[string]string{
			"postgres": "TEXT",
			"mysql":    "VARCHAR(255)",
			"mariadb":  "VARCHAR(255)",
			"sqlite3":  "TEXT",
		},
		GoType: "string",
	},
	"text": {
		SQLTypes: map[string]string{
			"postgres": "TEXT",
			"mysql":    "TEXT",
			"mariadb":  "TEXT",
			"sqlite3":  "TEXT",
		},
		GoType: "string",
	},
	"int": {
		SQLTypes: map[string]string{
			"postgres": "INTEGER",
			"mysql":    "INT",
			"mariadb":  "INT",
			"sqlite3":  "INTEGER",
		},
		GoType: "int64",
	},
	"bigint": {
		SQLTypes: map[string]string{
			"postgres": "BIGINT",
			"mysql":    "BIGINT",
			"mariadb":  "BIGINT",
			"sqlite3":  "INTEGER",
		},
		GoType: "int64",
	},
	"bool": {
		SQLTypes: map[string]string{
			"postgres": "BOOLEAN",
			"mysql":    "BOOLEAN",
			"mariadb":  "BOOLEAN",
			"sqlite3":  "BOOLEAN",
		},
		GoType: "bool",
	},
	"float": {
		SQLTypes: map[string]string{
			"postgres": "DOUBLE PRECISION",
			"mysql":    "FLOAT",
			"mariadb":  "FLOAT",
			"sqlite3":  "REAL",
		},
		GoType: "float64",
	},
	"timestamp": {
		SQLTypes: map[string]string{
			"postgres": "TIMESTAMP",
			"mysql":    "DATETIME",
			"mariadb":  "DATETIME",
			"sqlite3":  "DATETIME",
		},
		GoType: "time.Time",
	},
}
