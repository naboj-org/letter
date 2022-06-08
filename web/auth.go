package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthHandler(authToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Token")
		if token != authToken {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied."})
			c.Abort()
			return
		}
	}
}
