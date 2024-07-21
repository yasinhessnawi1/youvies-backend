package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

	// Setup the router and register handlers
	r := mux.NewRouter()
	r.Use(enableCors)

	api.RegisterHandlers(r)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("Starting server on port %s", port)

	go func() {
		log.Println(http.ListenAndServe(":"+port, loggedRouter))
                log.Println("started listning")
	}()

	tmdb := os.Getenv("TMDB_KEY")
	omdb := os.Getenv("OMDB_KEY")
	tvdb := os.Getenv("TVDB_KEY")
	log.Println("connecting to database")
	// Connect to the database
	database.ConnectDB()

	// Initialize scrapers with API keys
	animeShowScraper := scraper.NewAnimeShowScraper()
	animeMovieScraper := scraper.NewAnimeMovieScraper()
	movieScraper := scraper.NewMovieScraper(tmdb, omdb)
	showScraper := scraper.NewShowScraper(tmdb, tvdb)

	// Initialize bulk scraper
	bulkScraper := scraper.NewBulkScraper([]scraper.Scraper{
		animeShowScraper,
		animeMovieScraper,
		showScraper,
		movieScraper,
	})

	go func() {
		showScraper.FetchOldShows()
	}()

	// Run the bulk scraper
	err = bulkScraper.ScrapeAll()
	if err != nil {
		log.Printf("Error scraping data: %v", err)
	}

	// Block forever to keep the program running
	select {}
}
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
