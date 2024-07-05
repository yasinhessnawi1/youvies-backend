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

func GetShows(w http.ResponseWriter, r *http.Request) {
	collection := database.Client.Database("youvies").Collection("shows")

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

	var shows []models.Show
	if err = cursor.All(context.Background(), &shows); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove duplicates
	shows = removeDuplicateShows(shows)

	// Encode and send the result
	err = json.NewEncoder(w).Encode(shows)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Helper function to remove duplicate shows
func removeDuplicateShows(shows []models.Show) []models.Show {
	seen := make(map[string]bool)
	var result []models.Show
	for _, show := range shows {
		if _, ok := seen[show.Title+strconv.Itoa(show.Year)]; !ok {
			seen[show.ID.Hex()] = true
			result = append(result, show)
		}
	}
	return result
}

func GetShow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var show models.Show
	collection := database.Client.Database("youvies").Collection("shows")
	if err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&show); err != nil {
		http.Error(w, "Show not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(show)
}
func CreateShow(w http.ResponseWriter, r *http.Request) {
	var show models.Show
	if err := json.NewDecoder(r.Body).Decode(&show); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := database.InsertItem(show, show.Title, "shows")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := map[string]string{
		"message": "Show created successfully",
		"ID":      show.ID.Hex(),
	}
	json.NewEncoder(w).Encode(result)
}

func UpdateShow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var show models.Show
	if err := json.NewDecoder(r.Body).Decode(&show); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.EditItem(bson.M{"_id": objID}, show, "shows")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := map[string]string{
		"message": "Show updated successfully",
	}
	json.NewEncoder(w).Encode(result)
}

func DeleteShow(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err = database.DeleteItem(bson.M{"_id": objID}, "shows"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func SearchShows(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	collection := database.Client.Database("youvies").Collection("shows")
	cursor, err := collection.Find(context.Background(), bson.M{"title": bson.M{"$regex": title, "$options": "i"}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var shows []models.Show
	if err = cursor.All(context.Background(), &shows); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove duplicates
	shows = removeDuplicateShows(shows)

	json.NewEncoder(w).Encode(shows)
}
