package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/secure" // Security headers middleware
	"golang.org/x/time/rate"        // Rate limiting package
	"youvies-backend/api"
	"youvies-backend/database"
)

var rateLimiters = make(map[string]*rate.Limiter)
var rateLimit = rate.Every(1 * time.Second)
var burstLimit = 5 // Adjust this as needed

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Create a new Gin router
	router := gin.Default()

	// Middleware
	router.Use(corsMiddleware())            // Custom CORS middleware
	router.Use(rateLimitingMiddleware())    // Rate limiting middleware
	router.Use(securityHeadersMiddleware()) // Security headers middleware
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
			"https://yasinhessnawi1.github.io", // Replace with your GitHub Pages URL
			"http://localhost:3000",
			"https://localhost:3000", // In the case of HTTPS locally
			"https://youvies.online/",
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

// Rate limiting middleware to protect against bot traffic
func rateLimitingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if _, exists := rateLimiters[ip]; !exists {
			rateLimiters[ip] = rate.NewLimiter(rateLimit, burstLimit)
		}

		limiter := rateLimiters[ip]
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			return
		}

		c.Next()
	}
}

// Security headers middleware for additional protection
func securityHeadersMiddleware() gin.HandlerFunc {
	return secure.New(secure.Config{
		FrameDeny:          true,          // Prevents clickjacking
		ContentTypeNosniff: true,          // Prevents MIME type sniffing
		BrowserXssFilter:   true,          // Enables XSS protection in browsers
		ReferrerPolicy:     "no-referrer", // Controls the information sent in the Referer header
	})
}
