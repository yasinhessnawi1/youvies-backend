package utils

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"regexp"
	"youvies-backend/database"
	"youvies-backend/models"
)

// CheckUser validates the user fields.
func CheckUser(user models.User) error {
	if user.Username == "" {
		return errors.New("Username is required")
	}

	if err := checkUsernameAvailability(user.Username); err != nil {
		return err
	}

	if user.Password == "" {
		return errors.New("Password is required")
	} else if len(user.Password) < 8 || len(user.Password) > 50 {
		return errors.New("Password must be between 8 and 50 characters")
	} else if !isValidPassword(user.Password) {
		return errors.New("Password must contain at least one uppercase letter and one number")
	}

	if user.Email == "" {
		return errors.New("Email is required")
	} else if !isValidEmail(user.Email) {
		return errors.New("Invalid email format")
	}

	return nil
}

// isValidEmail checks if the email is valid.
func isValidEmail(email string) bool {
	// Use a regex to validate the email format.
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

// isValidPassword checks if the password contains at least one uppercase letter and one number.
func isValidPassword(password string) bool {
	// Check for at least one uppercase letter.
	reUpper := regexp.MustCompile(`[A-Z]`)
	// Check for at least one number.
	reDigit := regexp.MustCompile(`\d`)

	return reUpper.MatchString(password) && reDigit.MatchString(password)
}

// checkUsernameAvailability checks if the username is already taken.
func checkUsernameAvailability(username string) error {
	collection := database.Client.Database("youvies").Collection("users")
	var existingUser models.User
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&existingUser)
	if err == nil {
		return errors.New("Username is already taken")
	}
	return nil
}
