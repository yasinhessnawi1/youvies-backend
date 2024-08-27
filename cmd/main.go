package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"youvies-backend/api"
	"youvies-backend/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	// Create a new Gin router
	router := gin.Default()
	router.Use(enableCors)
	router.HandleMethodNotAllowed = true
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(logClientIP)
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

func enableCors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "https://youvies.online")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	c.Next()
}

// getClientIP tries to get the client's real IP address
func getClientIP(c *gin.Context) string {
	// First try to get the IP from X-Real-IP header (used by some proxies)
	ip := c.GetHeader("X-Real-IP")
	if ip == "" {
		// Next try to get the IP from X-Forwarded-For header (used by load balancers)
		ip = c.GetHeader("X-Forwarded-For")
	}
	if ip == "" {
		// Fall back to using the direct remote address
		ip = c.ClientIP()
	}
	return ip
}

// Middleware to log the client's IP address
func logClientIP(c *gin.Context) {
	ip := getClientIP(c)
	log.Printf("Client IP: %s", ip)
	c.Next() // continue to the next handler
}
