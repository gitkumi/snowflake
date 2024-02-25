package static

import (
	"embed"
	"mime"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

//go:embed public/*
var public embed.FS

//go:embed sql/migrations/*.sql
var Migration embed.FS

func ServePublic() gin.HandlerFunc {
	return func(c *gin.Context) {
		publicPath := path.Join("public/", c.Request.URL.Path)
		file, err := public.ReadFile(publicPath)

		// If file does not exist, go next.
		if err != nil {
			c.Next()
			return
		}

		contentType := mime.TypeByExtension(publicPath)
		c.Header("Content-Type", contentType)
		c.Data(http.StatusOK, contentType, file)
	}
}
