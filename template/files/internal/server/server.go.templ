package server

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitkumi/acme/static"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type config struct {
	port         int
	runMigration bool
}

type application struct {
	*config
	db *sql.DB
}

type server struct {
	*application
	router *gin.Engine
}

func New() *server {
	app := &application{
		config: initConfig(),
		db:     initDb(),
	}

	app.migrate()

	return &server{
		application: app,
		router:      initRouter(app),
	}
}

func (s *server) Start() {
	port := strconv.Itoa(s.application.config.port)

	if err := s.router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func (app *application) migrate() {
	if !app.config.runMigration {
		return
	}

	goose.SetBaseFS(static.Migration)

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(app.db, "sql/migrations"); err != nil {
		log.Fatal(err)
	}
}

func initConfig() *config {
	config := &config{}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}
	config.port = port

	runMigration, err := strconv.ParseBool(os.Getenv("MIGRATE_DB"))
	if err != nil {
		log.Fatal(err)
	}
	config.runMigration = runMigration

	return config
}

func initDb() *sql.DB {
	db, err := sql.Open("sqlite3", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	return db
}
