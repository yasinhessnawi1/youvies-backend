package api

import (
	"net/http"
	"strconv"
	"youvies-backend/database"
	"youvies-backend/models"

	"github.com/gin-gonic/gin"
)

// GetMovies retrieves movies with pagination.
func GetMovies(c *gin.Context) {
	version := c.Query("type")
	collection := "movie"
	if version == "sorted" {
		collection = "sorted_movies"
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
	var movies []models.Movie
	err = database.FindMany(collection, &movies, pageSize, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

// GetMovieByID retrieves a movie by its ID from the database.
func GetMovieByID(c *gin.Context) {
	id := c.Param("id")

	var movie models.Movie
	err := database.FindItem(id, "movie", &movie)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movies not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// GetMoviesByGenre retrieves movies by genre from the database.
func GetMoviesByGenre(c *gin.Context) {
	version := c.Query("type")
	collection := "movie"
	if version == "sorted" {
		collection = "sorted_movies"
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
	var movies []models.Movie
	err = database.FindByGenre(collection, genre, &movies, pageSize, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

// SearchMovies searches movies by title.
func SearchMovies(c *gin.Context) {
	version := c.Query("type")
	collection := "movie"
	if version == "sorted" {
		collection = "sorted_movies"
	}
	title := c.Query("title")

	var movies []models.Movie
	err := database.SearchItems(collection, title, &movies, 10, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}
