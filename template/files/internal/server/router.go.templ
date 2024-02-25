package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gitkumi/acme/internal/gintemplrenderer"
	"github.com/gitkumi/acme/static"
)

func initRouter(app *application) *gin.Engine {
	r := gin.Default()
	r.HTMLRender = gintemplrenderer.Default
	r.SetTrustedProxies(nil)

	r.Use(static.ServePublic())

	r.GET("/", app.home)
	r.GET("/create", app.createAuthor)
	r.GET("/health", app.health)
	r.NoRoute(app.noRoute)

	return r
}
