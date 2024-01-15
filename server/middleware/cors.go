package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests only from 'http://localhost:3000'
		c.Header("Access-Control-Allow-Origin", "*")

		// Allow necessary HTTP methods
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")

		// Allow necessary headers
		c.Header("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")

		// Allow credentials
		c.Header("Access-Control-Allow-Credentials", "true")

		c.Status(http.StatusOK)
		c.Next()
	}
}
