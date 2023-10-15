package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

func CheckMethodAndContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost && c.Request.Method != http.MethodGet {
			c.AbortWithError(http.StatusMethodNotAllowed, entity.ErrMethodNotAllowed)
			return
		}
		if c.GetHeader("Content-Type") != "application/json" && c.GetHeader("Content-Type") != "text/plain" {
			c.AbortWithError(http.StatusUnsupportedMediaType, entity.ErrUnsupportedMediaType)
			return
		}
		c.Next()
	}
}
