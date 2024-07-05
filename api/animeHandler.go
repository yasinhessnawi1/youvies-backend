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

func GetAnimeShows(w http.ResponseWriter, r *http.Request) {
	collection := database.Client.Database("youvies").Collection("anime_shows")

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

	var shows []models.Anime
	if err = cursor.All(context.Background(), &shows); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter and remove duplicates
	uniqueShows := make(map[string]models.Anime)
	for _, show := range shows {
		if _, exists := uniqueShows[show.Attributes.Titles.En]; !exists {
			uniqueShows[show.Attributes.Titles.En] = show
		}
	}

	result := make([]models.Anime, 0, len(uniqueShows))
	for _, show := range uniqueShows {
		result = append(result, show)
	}

	// Encode and send the result
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetAnimeShow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var show models.Anime
	collection := database.Client.Database("youvies").Collection("anime_shows")
	if err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&show); err != nil {
		http.Error(w, "Anime show not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(show)
}

func CreateAnimeShow(w http.ResponseWriter, r *http.Request) {
	var show models.Anime
	if err := json.NewDecoder(r.Body).Decode(&show); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := database.InsertItem(show, show.Attributes.Titles.En, "anime_shows")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := map[string]string{
		"message":    "Show created successfully",
		"insertedID": show.Id,
	}
	json.NewEncoder(w).Encode(result)
}

func UpdateAnimeShow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var show models.Anime
	if err := json.NewDecoder(r.Body).Decode(&show); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.EditItem(bson.M{"_id": objID}, show, "anime_shows")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	results := map[string]string{
		"message": "Show updated successfully",
	}
	json.NewEncoder(w).Encode(results)
}

func DeleteAnimeShow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	if err := database.DeleteItem(bson.M{"_id": objID}, "anime_shows"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func SearchAnimeShows(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	collection := database.Client.Database("youvies").Collection("anime_shows")
	cursor, err := collection.Find(context.Background(), bson.M{"title": bson.M{"$regex": title, "$options": "i"}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var shows []models.Anime
	if err = cursor.All(context.Background(), &shows); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter and remove duplicates
	uniqueShows := make(map[string]models.Anime)
	for _, show := range shows {
		if _, exists := uniqueShows[show.Attributes.Titles.En]; !exists {
			uniqueShows[show.Attributes.Titles.En] = show
		}
	}

	result := make([]models.Anime, 0, len(uniqueShows))
	for _, show := range uniqueShows {
		result = append(result, show)
	}

	// Encode and send the result
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
