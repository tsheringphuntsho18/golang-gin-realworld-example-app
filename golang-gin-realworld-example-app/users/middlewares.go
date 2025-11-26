package users

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" || !isValidToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// isValidToken checks if the provided token is valid
func isValidToken(token string) bool {
	// Implement token validation logic here
	return strings.HasPrefix(token, "Bearer ")
}

// AdminMiddleware checks if the user has admin privileges
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if !isAdminToken(token) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// isAdminToken checks if the provided token belongs to an admin user
func isAdminToken(token string) bool {
	// Implement admin token validation logic here
	return strings.HasPrefix(token, "Bearer admin_")
}