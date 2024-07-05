package api

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"youvies-backend/database"
	"youvies-backend/models"
)

func GetAnimeMovies(w http.ResponseWriter, r *http.Request) {
	collection := database.Client.Database("youvies").Collection("anime_movies")

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

	var movies []models.Anime
	if err = cursor.All(context.Background(), &movies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter and remove duplicates
	uniqueMovies := make(map[string]models.Anime)
	for _, movie := range movies {
		if _, exists := uniqueMovies[movie.Id]; !exists {
			uniqueMovies[movie.Id] = movie
		}
	}

	result := make([]models.Anime, 0, len(uniqueMovies))
	for _, movie := range uniqueMovies {
		result = append(result, movie)
	}

	// Encode and send the result
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetAnimeMovie(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var movie models.Anime
	collection := database.Client.Database("youvies").Collection("anime_movies")
	if err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&movie); err != nil {
		http.Error(w, "Anime movie not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(movie)
}

func CreateAnimeMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Anime
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := database.InsertItem(movie, movie.Attributes.Titles.En, "anime_movies")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := map[string]string{
		"message":  "Anime movie created successfully",
		"movie_id": movie.Id,
	}
	json.NewEncoder(w).Encode(result)
}

func UpdateAnimeMovie(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var movie models.Anime
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.EditItem(bson.M{"_id": objID}, movie, "anime_movies")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := map[string]string{
		"message": "Anime movie updated successfully",
	}
	json.NewEncoder(w).Encode(result)
}

func DeleteAnimeMovie(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	if err := database.DeleteItem(bson.M{"_id": objID}, "anime_movies"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func SearchAnimeMovies(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	collection := database.Client.Database("youvies").Collection("anime_movies")
	cursor, err := collection.Find(context.Background(), bson.M{"title": bson.M{"$regex": title, "$options": "i"}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var movies []models.Anime
	if err = cursor.All(context.Background(), &movies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter and remove duplicates
	uniqueMovies := make(map[string]models.Anime)
	for _, movie := range movies {
		if _, exists := uniqueMovies[movie.Attributes.Titles.En]; !exists {
			uniqueMovies[movie.Attributes.Titles.En] = movie
		}
	}

	result := make([]models.Anime, 0, len(uniqueMovies))
	for _, movie := range uniqueMovies {
		result = append(result, movie)
	}

	// Encode and send the result
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
