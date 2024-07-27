package api

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"youvies-backend/database"
	"youvies-backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAnimeShows retrieves anime shows with pagination.
func GetAnimeShows(c *gin.Context) {
	version := c.Query("type")
	var collection *mongo.Collection
	if version == "tiny" {
		collection = database.Client.Database("youvies").Collection("tiny_anime_shows")
	} else {
		collection = database.Client.Database("youvies").Collection("anime_shows")
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
	cursor, err := collection.Find(context.Background(), bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var animeShows []models.AnimeShow
	if err = cursor.All(context.Background(), &animeShows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, animeShows)
}

// GetAnimeShowByID retrieves an anime show by its ID from the database.
func GetAnimeShowByID(c *gin.Context) {
	id := c.Param("id")

	// Convert the string ID to a MongoDB ObjectId
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var animeShow models.AnimeShow
	collection := database.Client.Database("youvies").Collection("anime_shows")
	if err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&animeShow); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime show not found"})
		return
	}

	c.JSON(http.StatusOK, animeShow)
}

// GetAnimeShowsByGenre retrieves anime shows by genre from the database.
func GetAnimeShowsByGenre(c *gin.Context) {
	version := c.Query("type")
	var collection string
	if version == "tiny" {
		collection = "tiny_anime_shows"
	} else {
		collection = "anime_shows"
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
	var animeShows []models.AnimeShow
	err = database.FindMany(bson.D{{"genres.name", genre}}, collection, &animeShows, options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, animeShows)
}

// CreateAnimeShow creates a new anime show in the database.
func CreateAnimeShow(c *gin.Context) {
	var animeShow models.AnimeShow
	if err := c.BindJSON(&animeShow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.InsertItem(animeShow, animeShow.Attributes.Titles.En, "anime_shows")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := map[string]string{
		"message": "Anime show created successfully",
		"ID":      primitive.NilObjectID.Hex(),
	}
	c.JSON(http.StatusOK, result)
}

// UpdateAnimeShow updates an existing anime show in the database.
func UpdateAnimeShow(c *gin.Context) {
	id := c.Param("id")

	// Convert the string ID to a MongoDB ObjectId
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var animeShow models.AnimeShow
	if err := c.BindJSON(&animeShow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = database.EditItem(bson.M{"_id": objID}, animeShow, "anime_shows")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := map[string]string{
		"message": "Anime show updated successfully",
	}
	c.JSON(http.StatusOK, result)
}

// DeleteAnimeShow deletes an anime show from the database.
func DeleteAnimeShow(c *gin.Context) {
	id := c.Param("id")

	// Convert the string ID to a MongoDB ObjectId
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := database.DeleteItem(bson.M{"_id": objID}, "anime_shows"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchAnimeShows searches anime shows by title.
func SearchAnimeShows(c *gin.Context) {
	version := c.Query("type")
	var collection *mongo.Collection
	if version == "tiny" {
		collection = database.Client.Database("youvies").Collection("tiny_anime_shows")
	} else {
		collection = database.Client.Database("youvies").Collection("anime_shows")
	}
	title := c.Query("title")
	cursor, err := collection.Find(context.Background(), bson.M{"title": bson.M{"$regex": title, "$options": "i"}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var animeShows []models.AnimeShow
	if err = cursor.All(context.Background(), &animeShows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, animeShows)
}
