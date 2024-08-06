package api

import (
	"net/http"
	"strconv"
	"youvies-backend/database"
	"youvies-backend/models"

	"github.com/gin-gonic/gin"
)

// GetAnimeMovies retrieves anime movies with pagination.
func GetAnimeMovies(c *gin.Context) {
	version := c.Query("type")
	collection := "anime_movies"
	if version == "tiny" {
		collection = "tiny_anime_movies"
	}

	// Read pagination parameters from URL query
	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")

	// Set default values if parameters are not provided
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize

	// Find with pagination
	var animeMovies []models.AnimeMovie
	err = database.FindMany(collection, &animeMovies, pageSize, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, animeMovies)
}

// GetAnimeMovieByID retrieves an anime movie by its ID from the database.
func GetAnimeMovieByID(c *gin.Context) {
	id := c.Param("id")

	var animeMovie models.AnimeMovie
	err := database.FindItem(id, "anime_movies", &animeMovie)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime movie not found"})
		return
	}

	c.JSON(http.StatusOK, animeMovie)
}

// GetAnimeMoviesByGenre retrieves anime movies by genre from the database.
func GetAnimeMoviesByGenre(c *gin.Context) {
	version := c.Query("type")
	collection := "anime_movies"
	if version == "tiny" {
		collection = "tiny_anime_movies"
	}
	genre := c.Param("genre")

	// Read pagination parameters from URL query
	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")

	// Set default values if parameters are not provided
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize
	var animeMovies []models.AnimeMovie
	err = database.FindByGenre(collection, genre, &animeMovies, pageSize, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, animeMovies)
}

// SearchAnimeMovies searches anime movies by title.
func SearchAnimeMovies(c *gin.Context) {
	version := c.Query("type")
	collection := "anime_movies"
	if version == "tiny" {
		collection = "tiny_anime_movies"
	}
	title := c.Query("title")

	var animeMovies []models.AnimeMovie
	err := database.SearchItems(collection, title, &animeMovies, 10, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, animeMovies)
}
