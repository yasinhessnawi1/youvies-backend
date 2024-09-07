package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"youvies-backend/models"

	"github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	connStr := os.Getenv("POSTGRES_URI")
	if connStr == "" {
		log.Println("POSTGRES_URI not found in environment, using default URI")
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL!")
	createTables()
}

func createTables() {
	tableCreationQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        username TEXT NOT NULL,
        password TEXT NOT NULL,
        email TEXT NOT NULL,
        role TEXT,
        active BOOLEAN,
        created TIMESTAMP,
        updated TIMESTAMP,
        avatar TEXT,
        favorites TEXT[],
        friends TEXT[],
        rooms TEXT[],
        watched TEXT[]
    );`
	_, err := DB.Exec(tableCreationQuery)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Tables created successfully!")
}

func InsertItem(item interface{}, tableName string) error {
	log.Printf("Inserting Item in table %s\n ", tableName)
	switch v := item.(type) {
	case *models.User:
		return insertUser(*v)
	default:
		return fmt.Errorf("unsupported item type")
	}
}

func GetItem(id string, itemType string) (interface{}, error) {
	switch itemType {
	case "user":
		return getUser(id)
	default:
		return nil, fmt.Errorf("unsupported item type")
	}
}

func EditItem(update interface{}, tableName string) error {
	var query string
	var args []interface{}

	switch v := update.(type) {
	case *models.User:
		query = fmt.Sprintf(`UPDATE %s SET 
			username=$1, password=$2, email=$3, role=$4, active=$5, created=$6, updated=$7, avatar=$8, favorites=$9, friends=$10, rooms=$11, watched=$12 
			WHERE id=$13`, tableName)
		args = append(args, v.Username, v.Password, v.Email, v.Role, v.Active, v.Created, v.Updated, v.Avatar, pq.Array(v.Favorites), pq.Array(v.Friends), pq.Array(v.Rooms), pq.Array(v.Watched), v.ID)
	default:
		return fmt.Errorf("unsupported item type")
	}

	_, err := DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error updating item: %v", err)
	}
	log.Printf("Item updated successfully in table %s \n", tableName)

	return nil
}

func IfItemExists(id string, tableName string) (bool, error) {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = $1", tableName)
	var count int
	err := DB.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func FindItem(id string, tableName string, result interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tableName)
	row := DB.QueryRow(query, id)

	switch v := result.(type) {
	case *models.User:
		err := row.Scan(&v.ID, &v.Username, &v.Password, &v.Email, &v.Role, &v.Active, &v.Created, &v.Updated, &v.Avatar, pq.Array(&v.Favorites), pq.Array(&v.Friends), pq.Array(&v.Rooms), pq.Array(&v.Watched))
		if err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}
	default:
		return fmt.Errorf("unsupported result type")
	}

	return nil
}

func FindUser(itemName, tableName string) (string, error) {
	// Prepare SQL query to check both username and email
	query := fmt.Sprintf("SELECT id FROM %s WHERE username = $1 OR email = $1", tableName)

	var id string
	err := DB.QueryRow(query, itemName).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
