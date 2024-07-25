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

// GetAnimeMovies retrieves anime movies with pagination.
func GetAnimeMovies(c *gin.Context) {
	collection := database.Client.Database("youvies").Collection("anime_movies")

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
	cursor, err := collection.Find(context.Background(), bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var animeMovies []models.AnimeMovie
	if err = cursor.All(context.Background(), &animeMovies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Encode and send the result
	c.JSON(http.StatusOK, animeMovies)
}

// GetAnimeMovieByID retrieves an anime movie by its ID from the database.
func GetAnimeMovieByID(c *gin.Context) {
	id := c.Param("id")

	var animeMovie models.AnimeMovie
	collection := database.Client.Database("youvies").Collection("anime_movies")
	if err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&animeMovie); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime movie not found"})
		return
	}

	c.JSON(http.StatusOK, animeMovie)
}

// GetAnimeMoviesByGenre retrieves anime movies by genre from the database.
func GetAnimeMoviesByGenre(c *gin.Context) {
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
	err = database.FindMany(bson.D{{"genres.name", genre}}, "anime_movies", &animeMovies, options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, animeMovies)
}

// CreateAnimeMovie creates a new anime movie in the database.
func CreateAnimeMovie(c *gin.Context) {
	var animeMovie models.AnimeMovie
	if err := c.BindJSON(&animeMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.InsertItem(animeMovie, animeMovie.Attributes.Titles.En, "anime_movies")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := map[string]string{
		"message": "Anime movie created successfully",
		"ID":      strconv.Itoa(animeMovie.ID),
	}
	c.JSON(http.StatusOK, result)
}

// UpdateAnimeMovie updates an existing anime movie in the database.
func UpdateAnimeMovie(c *gin.Context) {
	id := c.Param("id")

	var animeMovie models.AnimeMovie
	if err := c.BindJSON(&animeMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.EditItem(bson.M{"_id": id}, animeMovie, "anime_movies")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := map[string]string{
		"message": "Anime movie updated successfully",
	}
	c.JSON(http.StatusOK, result)
}

// DeleteAnimeMovie deletes an anime movie from the database.
func DeleteAnimeMovie(c *gin.Context) {
	id := c.Param("id")

	if err := database.DeleteItem(bson.M{"_id": id}, "anime_movies"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchAnimeMovies searches anime movies by title.
func SearchAnimeMovies(c *gin.Context) {
	title := c.Query("title")
	collection := database.Client.Database("youvies").Collection("anime_movies")
	cursor, err := collection.Find(context.Background(), bson.M{"title": bson.M{"$regex": title, "$options": "i"}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var animeMovies []models.AnimeMovie
	if err = cursor.All(context.Background(), &animeMovies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Encode and send the result
	c.JSON(http.StatusOK, animeMovies)
}
