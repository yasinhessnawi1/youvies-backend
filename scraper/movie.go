package scraper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"os"
	"strconv"
	"youvies-backend/database"
	"youvies-backend/models"
	"youvies-backend/utils"
)

// Base URL setup for various APIs
const (
	tmdbBaseURL = "https://api.themoviedb.org/3/movie/"
	omdbBaseURL = "http://www.omdbapi.com/?"
	ytsBaseURL  = "https://yts.mx/api/v2/list_movies.json"
)

type MovieScraper struct {
	BaseScraper
	tmdbAPIKey string
	omdbAPIKey string
}

func NewMovieScraper(tmdbKey, omdbKey string) *MovieScraper {
	return &MovieScraper{
		BaseScraper: *NewBaseScraper("movie", "https://api.themoviedb.org/3"),
		tmdbAPIKey:  tmdbKey,
		omdbAPIKey:  omdbKey,
	}
}

// Scrape orchestrates the fetching of movie data from multiple sources.
func (ms *MovieScraper) Scrape() error {

	ids, err := ms.FetchMovieIDsFromTMDB()
	if err != nil {
		return fmt.Errorf("error fetching movie IDs from TMDB: %v", err)

	}
	for _, id := range ids {
		url := fmt.Sprintf("%s/movie/%s?api_key=%s", utils.TMDBBaseURL, id, ms.tmdbAPIKey)
		var response struct {
			ImdbId string `json:"imdb_id"`
			Title  string `json:"title"`
		}
		err := utils.FetchJSON(url, "TMDB", &response)
		if err != nil {
			return fmt.Errorf("error fetching TMDB movie data: %v", err)
		}
		omdbMovie, err := ms.fetchOMDBDetails(response.ImdbId)
		if err != nil || omdbMovie.Response == "false" {
			log.Printf("Failed to fetch OMDB details for %s: %v", response.ImdbId, err)
			continue
		}
		exist, err := database.IfItemExists(bson.M{"title": response.Title}, "movies")
		if err != nil || exist {

			continue
		}
		torrents, err := utils.FetchTorrents(response.Title, []string{"Movies"})
		if err != nil {
			log.Printf("Failed to fetch torrents for %s: %v", response.Title, err)
			continue

		}

		movie := ms.createMovieDoc(omdbMovie, torrents)
		if err := database.InsertItem(movie, movie.Title, "movies"); err != nil {
			log.Printf("Failed to save movie %s to database: %v", response.Title, err)
			continue
		}
	}
	fmt.Println("Fetching old movies are completed")
	return nil
}

// fetchOMDBDetails retrieves detailed movie data from the OMDB API.
func (ms *MovieScraper) fetchOMDBDetails(imdbID string) (*models.OmdbMovie, error) {
	resp, err := http.Get(fmt.Sprintf("%si=%s&plot=full&apikey=%s", omdbBaseURL, imdbID, ms.omdbAPIKey))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var omdbMovie models.OmdbMovie
	if err := json.NewDecoder(resp.Body).Decode(&omdbMovie); err != nil {
		return nil, err
	}
	return &omdbMovie, nil
}

// createMovieDoc constructs a movie document from TMDB and OMDB data.
func (ms *MovieScraper) createMovieDoc(omdbMovie *models.OmdbMovie, torrents []models.Torrent) models.Movie {
	return models.Movie{
		ID:          primitive.NewObjectID(),
		Title:       omdbMovie.Title,
		Description: omdbMovie.Plot,
		Year:        omdbMovie.Year,
		Language:    omdbMovie.Language,
		Director:    omdbMovie.Director,
		Genres:      omdbMovie.Genre,
		Torrents:    torrents,
		Rating:      omdbMovie.ImdbRating,
		PosterURL:   omdbMovie.Poster,
	}
}

func (ms *MovieScraper) FetchMovieIDsFromTMDB() ([]string, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s", ms.tmdbAPIKey)
	var response struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}

	err := ms.FetchJSON(url, &response)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, result := range response.Results {
		ids = append(ids, fmt.Sprintf("%d", result.ID))
	}

	file, err := os.Open("utils/movie_ids_05_15_2024.json/movie_ids_05_15_2024.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Unmarshal the JSON line into a Movie struct
		type Movie struct {
			Adult         bool    `json:"adult"`
			ID            int     `json:"id"`
			OriginalTitle string  `json:"original_title"`
			Popularity    float64 `json:"popularity"`
			Video         bool    `json:"video"`
		}
		var movie Movie
		err := json.Unmarshal([]byte(line), &movie)
		if err != nil {
			log.Printf("Error decoding JSON: %v\n", err)
			continue
		}

		// Append the ID to the array
		ids = append(ids, strconv.Itoa(movie.ID))
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file: %v\n", err)
	}

	return ids, nil
}
