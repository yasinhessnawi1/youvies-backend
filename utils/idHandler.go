package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

var (
	movieIDCounter      int
	showIDCounter       int
	animeMovieIDCounter int
	animeShowIDCounter  int
	counterMutex        sync.Mutex
)

// InitializeIDCounters initializes the ID counters from the database.
func InitializeIDCounters(client *mongo.Client) {
	counterMutex.Lock()
	defer counterMutex.Unlock()

	movieIDCounter = getMaxID(client, "movies")
	showIDCounter = getMaxID(client, "shows")
	animeMovieIDCounter = getMaxID(client, "anime_movies")
	animeShowIDCounter = getMaxID(client, "anime_shows")
}

// getMaxID retrieves the maximum ID from a collection.
func getMaxID(client *mongo.Client, collectionName string) int {
	collection := client.Database("youvies").Collection(collectionName)
	var result struct {
		MaxID int `bson:"max_id"`
	}
	err := collection.FindOne(context.Background(), bson.M{}, options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})).Decode(&result)
	if err != nil {
		log.Printf("Failed to get max ID for %s: %v", collectionName, err)
		return 0
	}
	return result.MaxID
}

// GetNextMovieID increments and returns the next movie ID.
func GetNextMovieID() int {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	movieIDCounter++
	return movieIDCounter
}

// GetNextShowID increments and returns the next show ID.
func GetNextShowID() int {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	showIDCounter++
	return showIDCounter
}

// GetNextAnimeMovieID increments and returns the next anime movie ID.
func GetNextAnimeMovieID() int {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	animeMovieIDCounter++
	return animeMovieIDCounter
}

// GetNextAnimeShowID increments and returns the next anime show ID.
func GetNextAnimeShowID() int {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	animeShowIDCounter++
	return animeShowIDCounter
}
