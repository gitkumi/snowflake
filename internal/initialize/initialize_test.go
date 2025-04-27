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
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseNone,
		Queue:     initialize.QueueNone,
		AppType:   initialize.AppTypeAPI,
		OutputDir: tmpDir,
		Git:       false,
		SMTP:      true,
		Storage:   true,
		Redis:     true,
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
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseSQLite3,
		Queue:     initialize.QueueNone,
		AppType:   initialize.AppTypeAPI,
		OutputDir: tmpDir,
		Git:       false,
		SMTP:      true,
		Storage:   true,
		Redis:     true,
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
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabasePostgres,
		Queue:     initialize.QueueNone,
		AppType:   initialize.AppTypeAPI,
		OutputDir: tmpDir,
		Git:       false,
		SMTP:      true,
		Storage:   true,
		Redis:     true,
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
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseMySQL,
		Queue:     initialize.QueueNone,
		AppType:   initialize.AppTypeAPI,
		OutputDir: tmpDir,
		Git:       false,
		SMTP:      true,
		Storage:   true,
		Redis:     true,
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
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseSQLite3,
		Queue:     initialize.QueueNone,
		AppType:   initialize.AppTypeWeb,
		OutputDir: tmpDir,
		Git:       false,
		SMTP:      true,
		Storage:   true,
		Redis:     true,
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
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseSQLite3,
		Queue:     initialize.QueueNone,
		AppType:   initialize.AppTypeAPI,
		OutputDir: tmpDir,
		Git:       false,
		SMTP:      false,
		Storage:   true,
		Redis:     true,
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
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseSQLite3,
		Queue:     initialize.QueueNone,
		AppType:   initialize.AppTypeAPI,
		OutputDir: tmpDir,
		Git:       false,
		SMTP:      true,
		Storage:   false,
		Redis:     true,
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
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseSQLite3,
		Queue:     initialize.QueueNone,
		AppType:   initialize.AppTypeAPI,
		OutputDir: tmpDir,
		Git:       false,
		SMTP:      true,
		Storage:   true,
		Redis:     false,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateQueueBasic(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseSQLite3,
		Queue:     initialize.QueueBasic,
		AppType:   initialize.AppTypeAPI,
		OutputDir: tmpDir,
		Git:       false,
		SMTP:      true,
		Storage:   true,
		Redis:     true,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateQueueSQS(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseSQLite3,
		Queue:     initialize.QueueSQS,
		AppType:   initialize.AppTypeAPI,
		OutputDir: tmpDir,
		Git:       false,
		SMTP:      true,
		Storage:   true,
		Redis:     true,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateWithAuthProviders(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:         true,
		Name:          "acme",
		Database:      initialize.DatabaseSQLite3,
		Queue:         initialize.QueueNone,
		AppType:       initialize.AppTypeWeb,
		OutputDir:     tmpDir,
		Git:           false,
		SMTP:          true,
		Storage:       true,
		Redis:         true,
		OAuthGoogle:   true,
		OAuthGitHub:   true,
		OAuthFacebook: true,
		OIDCGoogle:    true,
		OIDCMicrosoft: true,
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
		// OAuth providers
		false, false, false, false, false, false, false, false, false, false, false, false, false,
		// OIDC providers
		false, false, false, false, false, false,
	)

	f.Fuzz(func(t *testing.T,
		withSMTP, withStorage, withRedis bool,
		appTypeInt, dbTypeInt, jobTypeInt int,
		// OAuth providers
		withOAuthGoogle, withOAuthDiscord, withOAuthGitHub, withOAuthInstagram, withOAuthMicrosoft,
		withOAuthReddit, withOAuthSpotify, withOAuthTwitch, withOAuthFacebook, withOAuthLinkedIn,
		withOAuthSlack, withOAuthStripe, withOAuthX bool,
		// OIDC providers
		withOIDCFacebook, withOIDCGoogle, withOIDCLinkedIn, withOIDCMicrosoft, withOIDCTwitch, withOIDCDiscord bool,
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
		queues := []initialize.Queue{
			initialize.QueueBasic,
			initialize.QueueSQS,
			initialize.QueueNone,
		}

		appType := appTypes[abs(appTypeInt)%len(appTypes)]
		database := databases[abs(dbTypeInt)%len(databases)]
		queue := queues[abs(jobTypeInt)%len(queues)]

		err := initialize.Run(&initialize.Config{
			Git:            false,
			Quiet:          true,
			Name:           "acme",
			OutputDir:      tmpDir,
			Database:       database,
			Queue:          queue,
			AppType:        appType,
			SMTP:           withSMTP,
			Storage:        withStorage,
			Redis:          withRedis,
			OAuthGoogle:    withOAuthGoogle,
			OAuthDiscord:   withOAuthDiscord,
			OAuthGitHub:    withOAuthGitHub,
			OAuthInstagram: withOAuthInstagram,
			OAuthMicrosoft: withOAuthMicrosoft,
			OAuthReddit:    withOAuthReddit,
			OAuthSpotify:   withOAuthSpotify,
			OAuthTwitch:    withOAuthTwitch,
			OAuthFacebook:  withOAuthFacebook,
			OAuthLinkedIn:  withOAuthLinkedIn,
			OAuthSlack:     withOAuthSlack,
			OAuthStripe:    withOAuthStripe,
			OAuthX:         withOAuthX,
			OIDCFacebook:   withOIDCFacebook,
			OIDCGoogle:     withOIDCGoogle,
			OIDCLinkedIn:   withOIDCLinkedIn,
			OIDCMicrosoft:  withOIDCMicrosoft,
			OIDCTwitch:     withOIDCTwitch,
			OIDCDiscord:    withOIDCDiscord,
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
