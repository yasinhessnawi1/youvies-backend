package api

import (
	"net/http"
	"strconv"
	"youvies-backend/database"
	"youvies-backend/models"

	"github.com/gin-gonic/gin"
)

// GetAnimeShows retrieves anime shows with pagination.
func GetAnimeShows(c *gin.Context) {
	version := c.Query("type")
	collection := "anime_shows"
	if version == "sorted" {
		collection = "sorted_anime_shows"
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
	var animeShows []models.AnimeTiny
	err = database.FindMany(collection, &animeShows, pageSize, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, animeShows)
}

// GetAnimeShowByID retrieves an anime show by its ID from the database.
func GetAnimeShowByID(c *gin.Context) {
	id := c.Param("id")

	var animeShow models.AnimeShow
	err := database.FindItem(id, "anime_shows", &animeShow)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime show not found"})
		return
	}

	c.JSON(http.StatusOK, animeShow)
}

// GetAnimeShowsByGenre retrieves anime shows by genre from the database.
func GetAnimeShowsByGenre(c *gin.Context) {
	version := c.Query("type")
	collection := "anime_shows"
	if version == "sorted" {
		collection = "sorted_anime_shows"
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
	var animeShows []models.AnimeTiny
	err = database.FindByGenre(collection, genre, &animeShows, pageSize, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, animeShows)
}

// SearchAnimeShows searches anime shows by title.
func SearchAnimeShows(c *gin.Context) {
	version := c.Query("type")
	collection := "anime_shows"
	if version == "sorted" {
		collection = "sorted_anime_shows"
	}
	title := c.Query("title")

	var animeShows []models.AnimeTiny
	err := database.SearchItems(collection, title, &animeShows, 10, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, animeShows)
}
