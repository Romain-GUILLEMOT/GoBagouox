package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Authorization header not provided")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerSchema)

		if tokenString != os.Getenv("WEBSERVER_API_KEY") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid token provided")
		}
		c.Next()

	}
}
