package api

import (
	"net/http"
	"strconv"
	"youvies-backend/database"
	"youvies-backend/models"

	"github.com/gin-gonic/gin"
)

// GetShows retrieves shows with pagination.
func GetShows(c *gin.Context) {
	version := c.Query("type")
	collection := "show"
	if version == "sorted" {
		collection = "sorted_shows"
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
	var shows []models.Show
	err = database.FindMany(collection, &shows, pageSize, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shows)
}

// GetShowByID retrieves a show by its ID from the database.
func GetShowByID(c *gin.Context) {
	id := c.Param("id")

	var show models.Show
	err := database.FindItem(id, "show", &show)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shows not found"})
		return
	}

	c.JSON(http.StatusOK, show)
}

// GetShowsByGenre retrieves shows by genre from the database.
func GetShowsByGenre(c *gin.Context) {
	version := c.Query("type")
	collection := "show"
	if version == "sorted" {
		collection = "sorted_shows"
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
	var shows []models.Show
	err = database.FindByGenre(collection, genre, &shows, pageSize, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shows)
}

// SearchShows searches shows by title.
func SearchShows(c *gin.Context) {
	version := c.Query("type")
	collection := "show"
	if version == "sorted" {
		collection = "sorted_shows"
	}
	title := c.Query("title")

	var shows []models.Show
	err := database.SearchItems(collection, title, &shows, 10, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shows)
}
