package generate

import (
	"fmt"
	"path/filepath"
	"time"
)

func MigrationNumber() string {
	return time.Now().Format("20060102150405")
}

func MigrationFilePath(migrationsDir string, number string, resourcePlural string) string {
	return filepath.Join(migrationsDir, fmt.Sprintf("%s_%s.sql", number, resourcePlural))
}
