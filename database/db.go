package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
    CREATE TABLE IF NOT EXISTS attributes (
        id SERIAL PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        slug TEXT NOT NULL,
        synopsis TEXT NOT NULL,
        description TEXT NOT NULL,
        canonical_title TEXT NOT NULL,
        average_rating TEXT,
        favorites_count INTEGER,
        start_date TEXT,
        end_date TEXT,
        next_release TEXT,
        popularity_rank INTEGER,
        rating_rank INTEGER,
        age_rating TEXT,
        age_rating_guide TEXT,
        subtype TEXT,
        status TEXT,
        episode_count INTEGER,
        episode_length INTEGER,
        youtube_video_id TEXT,
        show_type TEXT,
        poster_image JSONB NOT NULL,
        cover_image JSONB NOT NULL,
        titles JSONB NOT NULL
    );

    CREATE TABLE IF NOT EXISTS relationships (
        id SERIAL PRIMARY KEY,
        self TEXT,
        related TEXT
    );

    CREATE TABLE IF NOT EXISTS anime_movies (
        id TEXT PRIMARY KEY,
        title TEXT NOT NULL,
        attributes_id INTEGER REFERENCES attributes(id),
        relationships_id INTEGER REFERENCES relationships(id),
        genres JSONB
    );

    CREATE TABLE IF NOT EXISTS anime_shows (
        id TEXT PRIMARY KEY,
        title TEXT NOT NULL,
        attributes_id INTEGER REFERENCES attributes(id),
        relationships_id INTEGER REFERENCES relationships(id),
        genres JSONB
    );

    CREATE TABLE IF NOT EXISTS movie (
        id TEXT PRIMARY KEY,
        original_language TEXT,
        original_title TEXT NOT NULL,
        overview TEXT NOT NULL,
        popularity FLOAT,
        poster_path TEXT,
        release_date TEXT,
        title TEXT NOT NULL,
        vote_average FLOAT,
        vote_count INTEGER,
        backdrop_path TEXT,
        genres JSONB,
        last_updated TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS show (
        id TEXT PRIMARY KEY,
        title TEXT NOT NULL,
        overview TEXT NOT NULL,
        poster_path TEXT,
        first_air_date TEXT,
        genres JSONB,
        vote_average FLOAT,
        vote_count INTEGER,
        original_language TEXT,
        popularity FLOAT,
        backdrop_path TEXT,
        last_updated TIMESTAMP,
        seasons_info JSONB
    );

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
    );

    CREATE TABLE IF NOT EXISTS episodes (
        id TEXT PRIMARY KEY,
        anime_show_id TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        synopsis TEXT NOT NULL,
        description TEXT NOT NULL,
        titles JSONB NOT NULL,
        canonical_title TEXT NOT NULL,
        season_number INTEGER,
        number INTEGER,
        relative_number INTEGER,
        airdate TEXT,
        length INTEGER,
        thumbnail JSONB NOT NULL,
        FOREIGN KEY (anime_show_id) REFERENCES anime_shows(id)
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
	case *models.Movie:
		return insertMovie(*v)
	case *models.Show:
		return insertShow(*v)
	case *models.AnimeShow:
		_, err := insertAnimeShow(*v)
		return err
	case *models.AnimeMovie:
		_, err := insertAnimeMovie(*v)
		return err
	case *models.User:
		return insertUser(*v)
	default:
		return fmt.Errorf("unsupported item type")
	}
}

func GetItem(id string, itemType string) (interface{}, error) {
	switch itemType {
	case "movie":
		return getMovie(id)
	case "show":
		return getShow(id)
	case "animeShow":
		return getAnimeShow(id)
	case "animeMovie":
		return getAnimeMovie(id)
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
	case *models.Movie:
		genres, err := json.Marshal(v.Genres)
		if err != nil {
			return fmt.Errorf("error marshaling genres: %v", err)
		}
		query = fmt.Sprintf(`UPDATE %s SET 
			original_language=$1, original_title=$2, overview=$3, popularity=$4, poster_path=$5, release_date=$6, title=$7, vote_average=$8, vote_count=$9, backdrop_path=$10,  genres=$11, last_updated=$12 
			WHERE id=$13`, tableName)
		args = append(args, v.OriginalLanguage, v.OriginalTitle, v.Overview, v.Popularity, v.PosterPath, v.ReleaseDate, v.Title, v.VoteAverage, v.VoteCount, v.BackdropPath, genres, v.LastUpdated, v.ID)
	case *models.Show:
		genres, err := json.Marshal(v.Genres)
		if err != nil {
			return fmt.Errorf("error marshaling genres: %v", err)
		}
		seasonsInfo, err := json.Marshal(v.SeasonsInfo)
		if err != nil {
			return fmt.Errorf("error marshaling seasons info: %v", err)
		}
		query = fmt.Sprintf(`UPDATE %s SET 
			title=$1, overview=$2, poster_path=$3, first_air_date=$4, genres=$5, vote_average=$6, vote_count=$7, original_language=$8, popularity=$9, backdrop_path=$10, last_updated=$11, seasons_info=$12 
			WHERE id=$13`, tableName)
		args = append(args, v.Title, v.Overview, v.PosterPath, v.FirstAirDate, genres, v.VoteAverage, v.VoteCount, v.OriginalLanguage, v.Popularity, v.BackdropPath, v.LastUpdated, seasonsInfo, v.ID)
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

func EditAnime(oldItem, newItem interface{}, tableName string) error {
	var query string
	var args []interface{}

	err := deleteItem(oldItem, tableName)
	if err != nil {
		return fmt.Errorf("error deleting item: %v", err)
	}

	switch v := newItem.(type) {
	case *models.AnimeShow:
		genres, err := json.Marshal(v.Genres)
		if err != nil {
			return fmt.Errorf("error marshaling genres: %v", err)
		}
		episodes, err := getEpisodesByAnimeShowID(v.ID)
		if err != nil {
			return fmt.Errorf("error marshaling episodes: %v", err)
		}
		query = fmt.Sprintf(`UPDATE %s SET 
			title=$1, attributes_id=$2, relationships_id=$3, genres=$4, episodes=$5 
			WHERE id=$6`, tableName)
		args = append(args, v.Title, v.Attributes.ID, v.Relationships.ID, genres, episodes, v.ID)
	case *models.AnimeMovie:
		genres, err := json.Marshal(v.Genres)
		if err != nil {
			return fmt.Errorf("error marshaling genres: %v", err)
		}
		query = fmt.Sprintf(`UPDATE %s SET 
			title=$1, attributes_id=$2, relationships_id=$3, genres=$4 
			WHERE id=$5`, tableName)
		args = append(args, v.Title, v.Attributes.ID, v.Relationships.ID, genres, v.ID)
	default:
		return fmt.Errorf("unsupported item type")
	}

	_, err = DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error updating item: %v", err)
	}

	return nil
}

func deleteItem(item interface{}, name string) error {
	var query string
	var args []interface{}

	switch v := item.(type) {
	case *models.AnimeShow:
		query = fmt.Sprintf("DELETE FROM %s WHERE id=$1", name)
		args = append(args, v.ID)
	case *models.AnimeMovie:
		query = fmt.Sprintf("DELETE FROM %s WHERE id=$1", name)
		args = append(args, v.ID)
	default:
		return fmt.Errorf("unsupported item type")
	}

	_, err := DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error deleting item: %v", err)
	}

	return nil
}

func IfItemExists(id string, tableName string) (bool, error) {
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
	case *models.Movie:
		var genres []byte
		err := row.Scan(&v.ID, &v.OriginalLanguage, &v.OriginalTitle, &v.Overview, &v.Popularity, &v.PosterPath, &v.ReleaseDate, &v.Title, &v.VoteAverage, &v.VoteCount, &v.BackdropPath, &genres, &v.LastUpdated)
		if err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}
		unmarshalJSON(genres, &v.Genres)
	case *models.Show:
		var genres, seasonsInfo []byte
		err := row.Scan(&v.ID, &v.Title, &v.Overview, &v.PosterPath, &v.FirstAirDate, &genres, &v.VoteAverage, &v.VoteCount, &v.OriginalLanguage, &v.Popularity, &v.BackdropPath, &v.LastUpdated, &seasonsInfo)
		if err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}
		unmarshalJSON(genres, &v.Genres)
		unmarshalJSON(seasonsInfo, &v.SeasonsInfo)
	case *models.AnimeShow:
		_, err := getAnimeShow(v.ID)
		if err != nil {
			return err
		}
	case *models.AnimeMovie:
		var genres []byte
		var attributesID, relationshipsID int64
		err := row.Scan(&v.ID, &v.Title, &attributesID, &relationshipsID, &genres)
		if err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}
		v.Attributes, err = getAttributes(attributesID)
		if err != nil {
			return err
		}
		v.Relationships, err = getRelationships(relationshipsID)
		if err != nil {
			return err
		}
		unmarshalJSON(genres, &v.Genres)
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
	query := fmt.Sprintf("SELECT id FROM %s WHERE username = $1", tableName)
	var id string
	err := DB.QueryRow(query, itemName).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func FindMany(tableName string, results interface{}, limit, offset int) error {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT %d OFFSET %d", tableName, limit, offset)
	rows, err := DB.Query(query)
	if err != nil {
		return fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	return scanRows(rows, results)
}

func SearchItems(tableName, title string, results interface{}, limit, offset int) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE title ILIKE $1 LIMIT $2 OFFSET $3", tableName)
	args := []interface{}{fmt.Sprintf("%%%s%%", title), limit, offset}
	rows, err := DB.Query(query, args...)
	if err != nil {
		return fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	return scanRows(rows, results)
}

func FindByGenre(tableName, genre string, results interface{}, limit, offset int) error {
	var query string
	args := []interface{}{limit, offset}

	// Adjust the query based on the table name
	query = fmt.Sprintf(`SELECT * FROM %s WHERE genres @> '[{"name": "%s"}]'::jsonb LIMIT $1 OFFSET $2`, tableName, genre)
	// Execute the query
	rows, err := DB.Query(query, args...)
	if err != nil {
		return fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	// Process the results based on the type
	return scanRows(rows, results)
}
