package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeviceOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "device" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Hanya mesin yang diizinkan mengakses resource ini"})
			return
		}
		c.Next()
	}
}
