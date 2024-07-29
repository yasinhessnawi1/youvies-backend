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

// GetShows retrieves shows with pagination.
func GetShows(c *gin.Context) {
	version := c.Query("type")
	var collection *mongo.Collection
	if version == "tiny" {
		collection = database.Client.Database("youvies").Collection("tiny_shows")
	} else {
		collection = database.Client.Database("youvies").Collection("shows")
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

	var shows []models.Show
	if err = cursor.All(context.Background(), &shows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shows)
}

// GetShowByID retrieves a show by its ID from the database.
func GetShowByID(c *gin.Context) {
	id := c.Param("id")

	// Convert the string ID to a MongoDB ObjectId
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var show models.Show
	collection := database.Client.Database("youvies").Collection("shows")
	if err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&show); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Show not found"})
		return
	}

	c.JSON(http.StatusOK, show)
}

// GetShowsByGenre retrieves shows by genre from the database.
func GetShowsByGenre(c *gin.Context) {
	version := c.Query("type")
	var collection string
	if version == "tiny" {
		collection = version + "_shows"
	} else {
		collection = "shows"
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
	err = database.FindMany(bson.D{{"genres.name", genre}}, collection, &shows, options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shows)
}

// CreateShow creates a new show in the database.
func CreateShow(c *gin.Context) {
	var show models.Show
	if err := c.BindJSON(&show); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.InsertItem(show, show.Title, "shows")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := map[string]string{
		"message": "Show created successfully",
		"ID":      primitive.NewObjectID().Hex(),
	}
	c.JSON(http.StatusOK, result)
}

// UpdateShow updates an existing show in the database.
func UpdateShow(c *gin.Context) {
	id := c.Param("id")

	// Convert the string ID to a MongoDB ObjectId
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var show models.Show
	if err := c.BindJSON(&show); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := database.Client.Database("youvies").Collection("shows")
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": show})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := map[string]string{
		"message": "Show updated successfully",
	}
	c.JSON(http.StatusOK, result)
}

// DeleteShow deletes a show from the database.
func DeleteShow(c *gin.Context) {
	id := c.Param("id")

	// Convert the string ID to a MongoDB ObjectId
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := database.Client.Database("youvies").Collection("shows")
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchShows searches shows by title.
func SearchShows(c *gin.Context) {
	version := c.Query("type")
	var collection *mongo.Collection
	if version == "tiny" {
		collection = database.Client.Database("youvies").Collection("tiny_shows")
	} else {
		collection = database.Client.Database("youvies").Collection("shows")
	}
	title := c.Query("title")
	cursor, err := collection.Find(context.Background(), bson.M{"title": bson.M{"$regex": title, "$options": "i"}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var shows []models.Show
	if err = cursor.All(context.Background(), &shows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shows)
}

func GetShowByVoteAverage(c *gin.Context) {
	version := c.Query("type")
	var collection *mongo.Collection
	if version == "tiny" {
		collection = database.Client.Database("youvies").Collection("tiny_shows")
	} else {
		collection = database.Client.Database("youvies").Collection("shows")
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

	// Calculate skip value
	skip := (page - 1) * pageSize

	// Define the sorting options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSort(bson.D{{"vote_count", -1}})

	// Find the movies sorted by vote count
	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var shows []models.Show
	if err = cursor.All(context.Background(), &shows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shows)
}
