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
		Authentication: initialize.AuthenticationEmail,
		OAuthGoogle:    false,
		OAuthDiscord:   false,
		OAuthFacebook:  false,
		OAuthGitHub:    false,
		OAuthInstagram: false,
		OAuthLinkedIn:  false,
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
		Authentication: initialize.AuthenticationEmail,
		OAuthGoogle:    false,
		OAuthDiscord:   false,
		OAuthFacebook:  false,
		OAuthGitHub:    false,
		OAuthInstagram: false,
		OAuthLinkedIn:  false,
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
		Authentication: initialize.AuthenticationEmail,
		OAuthGoogle:    false,
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
		Authentication: initialize.AuthenticationEmail,
		OAuthGoogle:    false,
		OAuthDiscord:   false,
		OAuthFacebook:  false,
		OAuthGitHub:    false,
		OAuthInstagram: false,
		OAuthLinkedIn:  false,
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
		Authentication: initialize.AuthenticationEmail,
		OAuthGoogle:    false,
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
		Authentication: initialize.AuthenticationEmail,
		OAuthGoogle:    false,
		OAuthDiscord:   false,
		OAuthFacebook:  false,
		OAuthGitHub:    false,
		OAuthInstagram: false,
		OAuthLinkedIn:  false,
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
		Authentication: initialize.AuthenticationEmail,
		OAuthGoogle:    false,
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

func TestGenerateAuthenticationNone(t *testing.T) {
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
		Authentication: initialize.AuthenticationNone,
		OAuthGoogle:    false,
		OAuthDiscord:   false,
		OAuthFacebook:  false,
		OAuthGitHub:    false,
		OAuthInstagram: false,
		OAuthLinkedIn:  false,
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

func TestGenerateAuthenticationEmail(t *testing.T) {
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
		Authentication: initialize.AuthenticationEmail,
		OAuthGoogle:    false,
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

func TestGenerateAuthenticationEmailWithUsername(t *testing.T) {
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
		Authentication: initialize.AuthenticationEmailWithUsername,
		OAuthGoogle:    false,
		OAuthDiscord:   false,
		OAuthFacebook:  false,
		OAuthGitHub:    false,
		OAuthInstagram: false,
		OAuthLinkedIn:  false,
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
		Authentication: initialize.AuthenticationEmail,
		OAuthGoogle:    false,
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
		Authentication: initialize.AuthenticationEmail,
		OAuthGoogle:    false,
		OAuthDiscord:   false,
		OAuthFacebook:  false,
		OAuthGitHub:    false,
		OAuthInstagram: false,
		OAuthLinkedIn:  false,
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
		Authentication: initialize.AuthenticationEmail,
		OAuthGoogle:    false,
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

func TestGenerateBackgroundJobAsynq(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:          true,
		Name:           "acme",
		Database:       initialize.DatabaseSQLite3,
		BackgroundJob:  initialize.BackgroundJobAsynq,
		AppType:        initialize.AppTypeAPI,
		OutputDir:      tmpDir,
		Git:            false,
		SMTP:           true,
		Storage:        true,
		Authentication: initialize.AuthenticationEmail,
		OAuthGoogle:    false,
		OAuthDiscord:   false,
		OAuthFacebook:  false,
		OAuthGitHub:    false,
		OAuthInstagram: false,
		OAuthLinkedIn:  false,
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

func TestGenerateOAuthGoogle(t *testing.T) {
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
		Authentication: initialize.AuthenticationEmail,
		OAuthGoogle:    true,
		OAuthDiscord:   false,
		OAuthFacebook:  false,
		OAuthGitHub:    false,
		OAuthInstagram: false,
		OAuthLinkedIn:  false,
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
		true, true, true, true,
		true, true, true, true, true, true,
		0, 0, 0, 0,
	)

	f.Fuzz(func(t *testing.T,
		withSMTP, withStorage, withAuth, withRedis bool,
		oauthGoogle, oauthDiscord, oauthFacebook, oauthGitHub, oauthInstagram, oauthLinkedIn bool,
		appTypeInt, dbTypeInt, jobTypeInt, authenticationTypeInt int,
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
			initialize.BackgroundJobAsynq,
			initialize.BackgroundJobNone,
		}
		authentications := []initialize.Authentication{
			initialize.AuthenticationNone,
			initialize.AuthenticationEmail,
			initialize.AuthenticationEmailWithUsername,
		}

		appType := appTypes[abs(appTypeInt)%len(appTypes)]
		database := databases[abs(dbTypeInt)%len(databases)]
		backgroundJob := backgroundJobs[abs(jobTypeInt)%len(backgroundJobs)]
		authentication := authentications[abs(authenticationTypeInt)%len(authentications)]

		err := initialize.Run(&initialize.Config{
			Quiet:          true,
			Name:           "acme",
			Database:       database,
			BackgroundJob:  backgroundJob,
			Authentication: authentication,
			AppType:        appType,
			OutputDir:      tmpDir,
			Git:            true,
			SMTP:           withSMTP,
			Storage:        withStorage,
			Redis:          withRedis,
			OAuthGoogle:    oauthGoogle,
			OAuthDiscord:   oauthDiscord,
			OAuthFacebook:  oauthFacebook,
			OAuthGitHub:    oauthGitHub,
			OAuthInstagram: oauthInstagram,
			OAuthLinkedIn:  oauthLinkedIn,
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
