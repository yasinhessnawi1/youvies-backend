package scraper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"youvies-backend/database"
	"youvies-backend/models"
	"youvies-backend/utils"
)

const (
	tmdbBaseURL = utils.TMDBBaseURL
)

type MovieScraper struct {
	BaseScraper
	tmdbAPIKey string
}

func NewMovieScraper(tmdbKey string) *MovieScraper {
	return &MovieScraper{
		BaseScraper: *NewBaseScraper("movie", tmdbBaseURL),
		tmdbAPIKey:  tmdbKey,
	}
}

// Scrape orchestrates the fetching of movie data from multiple sources.
func (ms *MovieScraper) Scrape() error {
	ids, err := ms.FetchMovieIDsFromTMDB()
	if err != nil {
		return fmt.Errorf("error fetching movie IDs from TMDB: %v", err)
	}
	for _, id := range ids {
		movieDetails, err := ms.FetchMovieDetailsFromTMDB(id)
		if err != nil {
			log.Printf("Failed to fetch TMDB details for %s: %v", id, err)
			continue
		}

		exists, err := database.IfItemExists(bson.M{"title": movieDetails.Title}, "movies")
		if err != nil {
			log.Printf("Database error: %v", err)
			continue
		}
		if exists {
			continue
		}
		torrents, err := utils.FetchTorrents(movieDetails.Title)
		if err != nil {
			log.Printf("Failed to fetch torrents for %s: %v", movieDetails.Title, err)
			continue
		}
		categorizedTorrents := utils.CategorizeTorrentsByQuality(torrents)

		// Update existing movie if changes are found
		if exists {
			// Fetch existing movie
			var existingMovie models.Movie
			if err := database.FindItem(bson.M{"title": movieDetails.Title}, "movies", &existingMovie); err != nil {
				log.Printf("Failed to fetch existing movie: %v", err)
				continue
			}
			// Check for updates
			if ms.hasMovieChanged(existingMovie, movieDetails, categorizedTorrents) {
				movie := ms.createMovieDoc(movieDetails, categorizedTorrents, existingMovie.ID)
				if err := database.EditItem(bson.M{"title": movie.Title}, movie, "movies"); err != nil {
					log.Printf("Failed to update movie %s in database: %v", movie.Title, err)
				}
			}
		} else {
			movie := ms.createMovieDoc(movieDetails, categorizedTorrents, -1)
			if err := database.InsertItem(movie, movie.Title, "movies"); err != nil {
				log.Printf("Failed to save movie %s to database: %v", movie.Title, err)
				continue
			}
		}
	}
	fmt.Println("Fetching movies completed")
	return nil
}

// FetchMovieDetailsFromTMDB retrieves detailed movie data from TMDB.
func (ms *MovieScraper) FetchMovieDetailsFromTMDB(id string) (*models.Movie, error) {
	url := fmt.Sprintf("%s/movie/%s?api_key=%s", tmdbBaseURL, id, ms.tmdbAPIKey)
	var response struct {
		Title            string                `json:"title"`
		OriginalLanguage string                `json:"original_language"`
		OriginalTitle    string                `json:"original_title"`
		Overview         string                `json:"overview"`
		Popularity       float64               `json:"popularity"`
		PosterPath       string                `json:"poster_path"`
		ReleaseDate      string                `json:"release_date"`
		VoteAverage      float64               `json:"vote_average"`
		VoteCount        int                   `json:"vote_count"`
		BackdropPath     string                `json:"backdrop_path"`
		Adult            bool                  `json:"adult"`
		Genres           []models.GenreMapping `json:"genres"`
	}
	err := utils.FetchJSON(url, "", &response)
	if err != nil {
		return nil, err
	}

	externalIDs, err := ms.FetchExternalIDs(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching external ids : %v", err)
	}

	return &models.Movie{
		Title:            response.Title,
		OriginalLanguage: response.OriginalLanguage,
		OriginalTitle:    response.OriginalTitle,
		Overview:         response.Overview,
		Popularity:       response.Popularity,
		PosterPath:       response.PosterPath,
		ReleaseDate:      response.ReleaseDate,
		VoteAverage:      response.VoteAverage,
		VoteCount:        response.VoteCount,
		BackdropPath:     response.BackdropPath,
		Adult:            response.Adult,
		Genres:           response.Genres,
		ExternalIDs:      externalIDs,
	}, nil
}

// FetchExternalIDs fetches external IDs for a movie from TMDB.
func (ms *MovieScraper) FetchExternalIDs(id string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/movie/%s/external_ids?api_key=%s", tmdbBaseURL, id, ms.tmdbAPIKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ms.tmdbAPIKey))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var externalIDs map[string]interface{}
	if err := json.Unmarshal(body, &externalIDs); err != nil {
		return nil, fmt.Errorf("error marchling external ids : %v", err)
	}

	return externalIDs, nil
}

// FetchMovieIDsFromTMDB retrieves a list of popular movie IDs from TMDB.
func (ms *MovieScraper) FetchMovieIDsFromTMDB() ([]string, error) {
	url := fmt.Sprintf("%s/movie/popular?api_key=%s", tmdbBaseURL, ms.tmdbAPIKey)
	var response struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}

	err := utils.FetchJSON(url, "", &response)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, result := range response.Results {
		ids = append(ids, fmt.Sprintf("%d", result.ID))
	}

	file, err := os.Open("utils/movie_ids_05_15_2024.json")
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

// createMovieDoc constructs a movie document from TMDB and OMDB data.
func (ms *MovieScraper) createMovieDoc(movieDetails *models.Movie, torrents map[string][]models.Torrent, id int) models.Movie {
	if id == -1 {
		id = utils.GetNextMovieID()
	}
	return models.Movie{
		ID:               id,
		OriginalLanguage: movieDetails.OriginalLanguage,
		OriginalTitle:    movieDetails.OriginalTitle,
		Overview:         movieDetails.Overview,
		Popularity:       movieDetails.Popularity,
		PosterPath:       movieDetails.PosterPath,
		ReleaseDate:      movieDetails.ReleaseDate,
		Title:            movieDetails.Title,
		VoteAverage:      movieDetails.VoteAverage,
		VoteCount:        movieDetails.VoteCount,
		BackdropPath:     movieDetails.BackdropPath,
		Adult:            movieDetails.Adult,
		Genres:           movieDetails.Genres,
		Torrents:         torrents,
		ExternalIDs:      movieDetails.ExternalIDs,
		LastUpdated:      time.Now().Format(time.RFC3339),
	}
}

// hasMovieChanged checks if the movie details or torrents have changed.
func (ms *MovieScraper) hasMovieChanged(existingMovie models.Movie, newDetails *models.Movie, newTorrents map[string][]models.Torrent) bool {
	// Compare relevant fields to determine if there are changes.
	return existingMovie.Overview != newDetails.Overview ||
		existingMovie.ReleaseDate != newDetails.ReleaseDate ||
		existingMovie.VoteAverage != newDetails.VoteAverage ||
		existingMovie.VoteCount != newDetails.VoteCount ||
		existingMovie.Popularity != newDetails.Popularity ||
		existingMovie.BackdropPath != newDetails.BackdropPath ||
		existingMovie.PosterPath != newDetails.PosterPath ||
		!compareTorrents(existingMovie.Torrents, newTorrents)
}

// compareTorrents compares two maps of torrents to see if they are different.
func compareTorrents(oldTorrents, newTorrents map[string][]models.Torrent) bool {
	if len(oldTorrents) != len(newTorrents) {
		return false
	}
	for quality, oldList := range oldTorrents {
		newList, exists := newTorrents[quality]
		if !exists || len(oldList) != len(newList) {
			return false
		}
		for i, oldTorrent := range oldList {
			if oldTorrent.Name != newList[i].Name || oldTorrent.Seeders != newList[i].Seeders {
				return false
			}
		}
	}
	return true
}
