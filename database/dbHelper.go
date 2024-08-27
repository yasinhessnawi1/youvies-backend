package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"youvies-backend/models"
)

func unmarshalJSON(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("error unmarshaling data: %v", err)
	}
	return nil
}

func insertUser(user models.User) error {
	query := `INSERT INTO users (
		id, username, password, email, role, active, created, updated, avatar, favorites, friends, rooms, watched
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	)`

	_, err := DB.Exec(query, user.ID, user.Username, user.Password, user.Email, user.Role, user.Active, user.Created, user.Updated, user.Avatar, pq.Array(user.Favorites), pq.Array(user.Friends), pq.Array(user.Rooms), pq.Array(user.Watched))
	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}
	return nil
}

func getUser(userID string) (models.User, error) {
	var user models.User
	query := `SELECT id, username, password, email, role, active, created, updated, avatar, favorites, friends, rooms, watched 
              FROM users WHERE id=$1`
	err := DB.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.Role, &user.Active, &user.Created, &user.Updated, &user.Avatar,
		pq.Array(&user.Favorites), pq.Array(&user.Friends), pq.Array(&user.Rooms), pq.Array(&user.Watched),
	)
	if err != nil {
		return user, fmt.Errorf("error getting user: %v", err)
	}

	return user, nil
}

func scanRows(rows *sql.Rows, results interface{}) error {
	switch v := results.(type) {

	case *[]models.User:
		var users []models.User
		for rows.Next() {
			var user models.User
			err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role, &user.Active, &user.Created, &user.Updated, &user.Avatar, pq.Array(&user.Favorites), pq.Array(&user.Friends), pq.Array(&user.Rooms), pq.Array(&user.Watched))
			if err != nil {
				return fmt.Errorf("error scanning row: %v", err)
			}
			users = append(users, user)
		}
		*v = users

	default:
		return fmt.Errorf("unsupported result type")
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %v", err)
	}

	return nil
}
