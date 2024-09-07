package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
	"youvies-backend/api"
	"youvies-backend/database"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Create a new Gin router
	router := gin.Default()

	// Middleware
	router.Use(corsMiddleware()) // Use custom CORS middleware
	router.HandleMethodNotAllowed = true
	router.Use(gin.Logger())

	// Register routes
	api.RegisterRoutes(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	database.ConnectDB()

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// CORS middleware to restrict allowed origins
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := []string{
			"https://<your-github-username>.github.io", // Replace with your GitHub Pages URL
			"http://localhost:3000",
			"https://localhost:3000", // In the case of HTTPS locally
		}

		origin := c.Request.Header.Get("Origin")
		if origin != "" && isValidOrigin(origin, allowedOrigins) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// Helper function to check if the origin is valid
func isValidOrigin(origin string, allowedOrigins []string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if strings.EqualFold(origin, allowedOrigin) {
			return true
		}
	}
	return false
}
