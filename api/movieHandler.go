package api

import (
	"context"
	"net/http"
	"strconv"
	"youvies-backend/database"
	"youvies-backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetMovies retrieves movies with pagination.
func GetMovies(c *gin.Context) {
	collection := database.Client.Database("youvies").Collection("movies")

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

	skip := (page - 1) * pageSize // Find with pagination
	cursor, err := collection.Find(context.Background(), bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var movies []models.Movie
	if err = cursor.All(context.Background(), &movies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

// GetMovieByID retrieves a movie by its ID from the database.
func GetMovieByID(c *gin.Context) {
	id := c.Param("id")

	var movie models.Movie
	collection := database.Client.Database("youvies").Collection("movies")
	if err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&movie); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// GetMoviesByGenre retrieves movies by genre from the database.
func GetMoviesByGenre(c *gin.Context) {
	genre := c.Param("genre")
	var movies []models.Movie
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
	err = database.FindMany(bson.D{{"genres.name", genre}}, "movies", &movies, options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

// CreateMovie creates a new movie in the database.
func CreateMovie(c *gin.Context) {
	var movie models.Movie
	if err := c.BindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.InsertItem(movie, movie.Title, "movies")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := map[string]string{
		"message": "Movie created successfully",
		"ID":      strconv.Itoa(movie.ID),
	}
	c.JSON(http.StatusOK, result)
}

// UpdateMovie updates an existing movie in the database.
func UpdateMovie(c *gin.Context) {
	id := c.Param("id")

	var movie models.Movie
	if err := c.BindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.EditItem(bson.M{"_id": id}, movie, "movies")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := map[string]string{
		"message": "Movie updated successfully",
	}
	c.JSON(http.StatusOK, result)
}

// DeleteMovie deletes a movie from the database.
func DeleteMovie(c *gin.Context) {
	id := c.Param("id")

	if err := database.DeleteItem(bson.M{"_id": id}, "movies"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchMovies searches movies by title.
func SearchMovies(c *gin.Context) {
	title := c.Query("title")
	collection := database.Client.Database("youvies").Collection("movies")
	cursor, err := collection.Find(context.Background(), bson.M{"title": bson.M{"$regex": title, "$options": "i"}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var movies []models.Movie
	if err = cursor.All(context.Background(), &movies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}
