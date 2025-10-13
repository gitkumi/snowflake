package initialize_test

import (
	"os"
	"path/filepath"
	"strings"
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

func TestGenerateMariaDB(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseMariaDB,
		Queue:     initialize.QueueNone,
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

func TestGenerateWithHTML(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseSQLite3,
		Queue:     initialize.QueueNone,
		OutputDir: tmpDir,
		Git:       false,
		SMTP:      true,
		Storage:   true,
		Redis:     true,
		ServeHTML: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmpDir, "acme")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatal("project directory not created")
	}
}

func TestGenerateNoSMTP(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseSQLite3,
		Queue:     initialize.QueueNone,
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

func TestEnvFilesGenerated(t *testing.T) {
	tmpDir := t.TempDir()

	err := initialize.Run(&initialize.Config{
		Quiet:     true,
		Name:      "acme",
		Database:  initialize.DatabaseSQLite3,
		Queue:     initialize.QueueNone,
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

	// Check for .env file
	envPath := filepath.Join(projectDir, "cmd", "app", ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		t.Fatalf(".env file not created at %s", envPath)
	}

	// Check for .env.test file
	envTestPath := filepath.Join(projectDir, "cmd", "app", ".env.test")
	if _, err := os.Stat(envTestPath); os.IsNotExist(err) {
		t.Fatalf(".env.test file not created at %s", envTestPath)
	}

	// Verify .env file has content
	envContent, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("failed to read .env file: %v", err)
	}
	if len(envContent) == 0 {
		t.Fatal(".env file is empty")
	}

	// Verify .env.test file has content
	envTestContent, err := os.ReadFile(envTestPath)
	if err != nil {
		t.Fatalf("failed to read .env.test file: %v", err)
	}
	if len(envTestContent) == 0 {
		t.Fatal(".env.test file is empty")
	}
}

func FuzzGenerate(f *testing.F) {
	f.Add(
		true, true, true, true,
		0, 0,
		// OAuth providers
		false, false, false, false, false, false, false, false, false, false, false, false, false,
		// OIDC providers
		false, false, false, false, false, false,
	)

	f.Fuzz(func(t *testing.T,
		withSMTP, withStorage, withRedis, withServeHTML bool,
		dbTypeInt, jobTypeInt int,
		// OAuth providers
		withOAuthGoogle, withOAuthDiscord, withOAuthGitHub, withOAuthInstagram, withOAuthMicrosoft,
		withOAuthReddit, withOAuthSpotify, withOAuthTwitch, withOAuthFacebook, withOAuthLinkedIn,
		withOAuthSlack, withOAuthStripe, withOAuthX bool,
		// OIDC providers
		withOIDCFacebook, withOIDCGoogle, withOIDCLinkedIn, withOIDCMicrosoft, withOIDCTwitch, withOIDCDiscord bool,
	) {
		tmpDir := t.TempDir()

		databases := []initialize.Database{
			initialize.DatabaseSQLite3,
			initialize.DatabasePostgres,
			initialize.DatabaseMySQL,
			initialize.DatabaseMariaDB,
			initialize.DatabaseNone,
		}
		queues := []initialize.Queue{
			initialize.QueueBasic,
			initialize.QueueSQS,
			initialize.QueueNone,
		}

		database := databases[abs(dbTypeInt)%len(databases)]
		queue := queues[abs(jobTypeInt)%len(queues)]

		err := initialize.Run(&initialize.Config{
			Git:            false,
			Quiet:          true,
			Name:           "acme",
			OutputDir:      tmpDir,
			Database:       database,
			Queue:          queue,
			ServeHTML:      withServeHTML,
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

func containsStr(s, substr string) bool {
	return strings.Contains(s, substr)
}
