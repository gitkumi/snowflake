package generate

import (
	"testing"
	"time"
)

func TestMigrationNumber(t *testing.T) {
	num := MigrationNumber()

	if len(num) != 14 {
		t.Errorf("expected 14-char timestamp, got %q (len %d)", num, len(num))
	}

	_, err := time.Parse("20060102150405", num)
	if err != nil {
		t.Errorf("migration number %q is not a valid timestamp: %v", num, err)
	}
}

func TestMigrationFilePath(t *testing.T) {
	path := MigrationFilePath("/project/cmd/app/sql/migrations", "20260323161113", "posts")
	expected := "/project/cmd/app/sql/migrations/20260323161113_posts.sql"
	if path != expected {
		t.Errorf("got %q, want %q", path, expected)
	}
}
