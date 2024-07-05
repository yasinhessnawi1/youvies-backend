package api

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"youvies-backend/database"
	"youvies-backend/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetMovies retrieves all movies from the database.
func GetMovies(w http.ResponseWriter, r *http.Request) {
	collection := database.Client.Database("youvies").Collection("movies")

	// Read pagination parameters from URL query
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var movies []models.Movie
	if err = cursor.All(context.Background(), &movies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter and remove duplicates
	uniqueMovies := removeDuplicateMovies(movies)

	// Encode and send the result
	err = json.NewEncoder(w).Encode(uniqueMovies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetMovie retrieves a single movie by ID.
func GetMovie(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var movie models.Movie
	collection := database.Client.Database("youvies").Collection("movies")
	if err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&movie); err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(movie)
}

// CreateMovie adds a new movie to the database.
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := database.InsertItem(movie, movie.Title, "movies")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := map[string]string{
		"message": "Movie created successfully",
		"ID":      movie.ID.Hex(),
	}
	json.NewEncoder(w).Encode(result)
}

// UpdateMovie updates an existing movie in the database.
func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.EditItem(bson.M{"_id": objID}, movie, "movies")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := map[string]string{
		"message": "Movie updated successfully",
	}
	json.NewEncoder(w).Encode(result)
}

// DeleteMovie removes a movie from the database.
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	if err = database.DeleteItem(bson.M{"_id": objID}, "movies"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchMovies finds movies by title.
func SearchMovies(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	collection := database.Client.Database("youvies").Collection("movies")
	cursor, err := collection.Find(context.Background(), bson.M{"title": bson.M{"$regex": title, "$options": "i"}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var movies []models.Movie
	if err = cursor.All(context.Background(), &movies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove duplicates
	movies = removeDuplicateMovies(movies)

	json.NewEncoder(w).Encode(movies)
}

// Helper function to remove duplicate movies
func removeDuplicateMovies(movies []models.Movie) []models.Movie {
	seen := make(map[string]bool)
	var result []models.Movie
	for _, movie := range movies {
		if _, ok := seen[movie.Title+movie.Year]; !ok {
			seen[movie.Title+movie.Year] = true
			result = append(result, movie)
		}
	}
	return result
}
