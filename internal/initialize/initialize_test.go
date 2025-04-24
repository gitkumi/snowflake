package initialize_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gitkumi/snowflake/internal/initialize"
)

func TestGenerateNoDB(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:          true,
		Name:           "acme",
		Database:       initialize.DatabaseNone,
		BackgroundJob:  initialize.BackgroundJobNone,
		AppType:        initialize.AppTypeAPI,
		OutputDir:      tmpDir,
		Git:            false,
		SMTP:           true,
		Storage:        true,
		Redis:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateSQLite3(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:          true,
		Name:           "acme",
		Database:       initialize.DatabaseSQLite3,
		BackgroundJob:  initialize.BackgroundJobNone,
		AppType:        initialize.AppTypeAPI,
		OutputDir:      tmpDir,
		Git:            false,
		SMTP:           true,
		Storage:        true,
		Redis:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGeneratePostgres(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:          true,
		Name:           "acme",
		Database:       initialize.DatabasePostgres,
		BackgroundJob:  initialize.BackgroundJobNone,
		AppType:        initialize.AppTypeAPI,
		OutputDir:      tmpDir,
		Git:            false,
		SMTP:           true,
		Storage:        true,
		Redis:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateMySQL(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:          true,
		Name:           "acme",
		Database:       initialize.DatabaseMySQL,
		BackgroundJob:  initialize.BackgroundJobNone,
		AppType:        initialize.AppTypeAPI,
		OutputDir:      tmpDir,
		Git:            false,
		SMTP:           true,
		Storage:        true,
		Redis:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateWebApp(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:          true,
		Name:           "acme",
		Database:       initialize.DatabaseSQLite3,
		BackgroundJob:  initialize.BackgroundJobNone,
		AppType:        initialize.AppTypeWeb,
		OutputDir:      tmpDir,
		Git:            false,
		SMTP:           true,
		Storage:        true,
		Redis:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}

	apiMainPath := filepath.Join(projectDir, "cmd", "api")
	webMainPath := filepath.Join(projectDir, "cmd", "web")
	htmlDirPath := filepath.Join(projectDir, "internal", "html")

	if _, err := os.Stat(webMainPath); os.IsNotExist(err) {
		t.Fatal("Web directory was not created")
	}

	if _, err := os.Stat(apiMainPath); err == nil {
		t.Fatal("API directory should not exist for Web app type")
	}

	if _, err := os.Stat(htmlDirPath); os.IsNotExist(err) {
		t.Fatal("HTML directory created for web app type")
	}
}

func TestGenerateNoSMTP(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:          true,
		Name:           "acme",
		Database:       initialize.DatabaseSQLite3,
		BackgroundJob:  initialize.BackgroundJobNone,
		AppType:        initialize.AppTypeAPI,
		OutputDir:      tmpDir,
		Git:            false,
		SMTP:           false,
		Storage:        true,
		Redis:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateNoStorage(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:          true,
		Name:           "acme",
		Database:       initialize.DatabaseSQLite3,
		BackgroundJob:  initialize.BackgroundJobNone,
		AppType:        initialize.AppTypeAPI,
		OutputDir:      tmpDir,
		Git:            false,
		SMTP:           true,
		Storage:        false,
		Redis:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateNoRedis(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:          true,
		Name:           "acme",
		Database:       initialize.DatabaseSQLite3,
		BackgroundJob:  initialize.BackgroundJobNone,
		AppType:        initialize.AppTypeAPI,
		OutputDir:      tmpDir,
		Git:            false,
		SMTP:           true,
		Storage:        true,
		Redis:          false,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateBackgroundJobBasic(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:          true,
		Name:           "acme",
		Database:       initialize.DatabaseSQLite3,
		BackgroundJob:  initialize.BackgroundJobBasic,
		AppType:        initialize.AppTypeAPI,
		OutputDir:      tmpDir,
		Git:            false,
		SMTP:           true,
		Storage:        true,
		Redis:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateBackgroundJobSQS(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:          true,
		Name:           "acme",
		Database:       initialize.DatabaseSQLite3,
		BackgroundJob:  initialize.BackgroundJobSQS,
		AppType:        initialize.AppTypeAPI,
		OutputDir:      tmpDir,
		Git:            false,
		SMTP:           true,
		Storage:        true,
		Redis:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func FuzzGenerate(f *testing.F) {
	f.Add(
		true, true, true,
		0, 0, 0,
	)

	f.Fuzz(func(t *testing.T,
		withSMTP, withStorage, withRedis bool,
		appTypeInt, dbTypeInt, jobTypeInt int,
	) {
		tmpDir := t.TempDir()

		appTypes := []initialize.AppType{
			initialize.AppTypeAPI,
			initialize.AppTypeWeb,
		}
		databases := []initialize.Database{
			initialize.DatabaseSQLite3,
			initialize.DatabasePostgres,
			initialize.DatabaseMySQL,
			initialize.DatabaseNone,
		}
		backgroundJobs := []initialize.BackgroundJob{
			initialize.BackgroundJobBasic,
			initialize.BackgroundJobSQS,
			initialize.BackgroundJobNone,
		}

		appType := appTypes[abs(appTypeInt)%len(appTypes)]
		database := databases[abs(dbTypeInt)%len(databases)]
		backgroundJob := backgroundJobs[abs(jobTypeInt)%len(backgroundJobs)]

		err := initialize.Run(&initialize.Config{
			Git:            false,
			Quiet:          true,
			Name:           "acme",
			OutputDir:      tmpDir,
			Database:       database,
			BackgroundJob:  backgroundJob,
			AppType:        appType,
			SMTP:           withSMTP,
			Storage:        withStorage,
			Redis:          withRedis,
		})
		if err != nil {
			t.Logf("initialize.Run returned error: %v", err)
			return
		}

		projectDir := filepath.Join(tmpDir, "acme")
		if _, err := os.Stat(projectDir); os.IsNotExist(err) {
			t.Fatal("project directory not created")
		}
	})
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
