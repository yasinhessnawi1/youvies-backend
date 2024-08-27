package api

import (
	"net/http"
	"time"
	"youvies-backend/database"
	"youvies-backend/models"
	"youvies-backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles the user registration process.
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := utils.CheckUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user.Password = string(hashedPassword)
	user.Created = time.Now()
	user.Updated = time.Now()
	user.ID = utils.GenerateUUID()
	user.Role = "user"

	if err := database.InsertItem(&user, "users"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// LoginUser handles the user login process.
func LoginUser(c *gin.Context) {
	var creds *models.User
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	id, err := database.FindUser(creds.Username, "users")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := database.FindItem(id, "users", &user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if creds.Role == "" {
		creds.Role = "user"
	} else if user.Role != creds.Role {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token, err := utils.GenerateJWT(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(100 * 365 * 24 * time.Hour),
	})
	response := map[string]interface{}{
		"token": token,
		"user":  user,
	}
	c.JSON(http.StatusOK, response)
}

// EditUser handles updating user details.
func EditUser(c *gin.Context) {
	var user models.User
	userID := c.GetHeader("user")
	err := database.FindItem(userID, "users", &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var userUpdate models.User
	if err := c.BindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.EditItem(&userUpdate, "users"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// LogoutUser handles the user logout process.
func LogoutUser(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
	})
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
