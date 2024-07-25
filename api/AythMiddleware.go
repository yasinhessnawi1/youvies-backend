package api

import (
	"net/http"
	"strings"
	"youvies-backend/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden no key provided"})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden, token not valid"})
			c.Abort()
			return
		}

		c.Set("user", claims.Username)
		c.Set("role", claims.Role)

		if role == "admin" && claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden, admin role required"})
			c.Abort()
			return
		}

		if role == "user" && claims.Role != "user" && claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden, user role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
