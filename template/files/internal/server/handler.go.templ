package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitkumi/acme/internal/data"
	"github.com/gitkumi/acme/internal/pages"
)

func (app *application) home(c *gin.Context) {
	q := data.New(app.db)
	authors, err := q.ListAuthors(c)

	if err != nil {
		msg := err.Error()
		c.HTML(http.StatusInternalServerError, "", pages.Error(msg))
		return
	}

	c.HTML(http.StatusOK, "", pages.Home(authors))
}

func (app *application) createAuthor(c *gin.Context) {
	q := data.New(app.db)

	author, _ := q.CreateAuthor(c, data.CreateAuthorParams{
		Name: "Test",
		Bio:  sql.NullString{String: "Test", Valid: true},
	})

	fmt.Println(author)

	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (app *application) health(c *gin.Context) {
	err := app.db.Ping()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func (app *application) noRoute(c *gin.Context) {
	c.HTML(http.StatusNotFound, "", pages.Error("404"))
}
