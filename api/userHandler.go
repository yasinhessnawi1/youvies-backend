package api

import (
	"encoding/json"
	"net/http"
	"time"
	"youvies-backend/database"
	"youvies-backend/models"
	"youvies-backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := utils.CheckUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)
	user.Created = time.Now()
	user.Updated = time.Now()
	user.ID = primitive.NewObjectID()
	user.Role = "user"

	if err := database.InsertItem(user, user.Username, "users"); err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.FindItem(bson.M{"username": creds.Username}, "users", &user); err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if creds.Role == "" {
		creds.Role = "user"
	} else if user.Role != creds.Role {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := utils.GenerateJWT(user.Username, user.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(100 * 365 * 24 * time.Hour),
	})
	response := map[string]interface{}{
		"token": token,
		"user":  user,
	}
	json.NewEncoder(w).Encode(response)

}

func EditUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserKey).(string)

	var userUpdate models.User
	if err := json.NewDecoder(r.Body).Decode(&userUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userUpdate.Updated = time.Now()
	if userUpdate.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userUpdate.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		userUpdate.Password = string(hashedPassword)
	}

	filter := bson.M{"username": userID}
	update := bson.M{"$set": userUpdate}

	if err := database.EditItem(filter, update, "users"); err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
	})
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}
