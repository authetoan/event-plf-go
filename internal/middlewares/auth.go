package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		userRole := c.GetHeader("X-User-Role")

		if userID == "" || userRole == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - missing authentication headers",
			})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("userRole", userRole)

		c.Next()
	}
}
