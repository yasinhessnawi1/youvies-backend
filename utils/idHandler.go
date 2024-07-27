package utils

import (
	"go.mongodb.org/mongo-driver/mongo"
	rand2 "math/rand"
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
}

// GetNextMovieID increments and returns the next movie ID.
func GetNextMovieID() int {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	movieIDCounter = rand2.Int()
	return movieIDCounter
}

// GetNextShowID increments and returns the next show ID.
func GetNextShowID() int {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	showIDCounter = rand2.Int()
	return showIDCounter
}

// GetNextAnimeMovieID increments and returns the next anime movie ID.
func GetNextAnimeMovieID() int {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	animeMovieIDCounter = rand2.Int()
	return animeMovieIDCounter
}

// GetNextAnimeShowID increments and returns the next anime show ID.
func GetNextAnimeShowID() int {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	animeShowIDCounter = rand2.Int()
	return animeShowIDCounter
}
