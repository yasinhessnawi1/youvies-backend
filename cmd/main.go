package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"youvies-backend/api"
	"youvies-backend/database"
	"youvies-backend/scraper"
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
	router.Use(gin.ErrorLoggerT(gin.ErrorTypeAny))
	// Register routes
	api.RegisterRoutes(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("Starting server on port %s", port)

	go func() {
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	//tmdb := os.Getenv("TMDB_KEY")

	// Connect to the database
	database.ConnectDB()

	//movieScraper := scraper.NewMovieScraper(tmdb)
	//showScraper := scraper.NewShowScraper(tmdb)
	animeShowScraper := scraper.NewAnimeShowScraper()
	//animeMovieScraper := scraper.NewAnimeMovieScraper()

	// Initialize bulk scraper
	bulkScraper := scraper.NewBulkScraper([]scraper.Scraper{
		//showScraper,
		//movieScraper,
		animeShowScraper,
		//animeMovieScraper,
	})

	// Run the bulk scraper
	if err := bulkScraper.ScrapeAll(); err != nil {
		log.Printf("Error scraping data: %v", err)
	}

	// Block forever to keep the program running
	select {}
}

func enableCors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	c.Next()
}
