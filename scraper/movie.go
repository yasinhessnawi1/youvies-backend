package scraper

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
	"youvies-backend/database"
	"youvies-backend/models"
	"youvies-backend/utils"
)

type MovieScraper struct {
	BaseScraper
	tmdbAPIKey string
}

func NewMovieScraper(tmdbKey string) *MovieScraper {
	return &MovieScraper{
		BaseScraper: *NewBaseScraper("movie", utils.TMDBBaseURL),
		tmdbAPIKey:  tmdbKey,
	}
}

func (ms *MovieScraper) FetchChangedMovieIDs() ([]string, error) {
	log.Println("Fetching changed movie IDs")
	url := fmt.Sprintf("%s/movie/changes?api_key=%s", utils.TMDBBaseURL, ms.tmdbAPIKey)
	var response struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}
	if err := ms.FetchJSON(url, &response); err != nil {
		return nil, err
	}
	var ids []string
	for _, result := range response.Results {
		ids = append(ids, strconv.Itoa(result.ID))
	}
	log.Printf("%v Changed movie IDs fetched\n", len(ids))
	return ids, nil
}

func (ms *MovieScraper) FetchNowPlayingMovieIDs() ([]string, error) {
	log.Println("Fetching now playing movie IDs")
	url := fmt.Sprintf("%s/movie/now_playing?api_key=%s", utils.TMDBBaseURL, ms.tmdbAPIKey)
	var response struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}
	if err := ms.FetchJSON(url, &response); err != nil {
		return nil, err
	}
	var ids []string
	for _, result := range response.Results {
		ids = append(ids, strconv.Itoa(result.ID))
	}
	log.Printf("%v Now playing movie IDs fetched\n", len(ids))
	return ids, nil
}

// Scrape orchestrates the fetching of movie data from multiple sources.
func (ms *MovieScraper) Scrape() error {
	log.Println("Starting movie scraping...")
	changedIDs, err := ms.FetchChangedMovieIDs()
	if err != nil {
		return fmt.Errorf("error fetching changed movie IDs from TMDB: %v", err)
	}

	nowPlayingIDs, err := ms.FetchNowPlayingMovieIDs()
	if err != nil {
		return fmt.Errorf("error fetching now playing movie IDs from TMDB: %v", err)
	}

	ids := append(changedIDs, nowPlayingIDs...)
	ids = utils.RemoveDuplicateStrings(ids) // Remove duplicates

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 3) // Limit the number of concurrent goroutines

	for _, id := range ids {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(id string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			exists, err := database.IfItemExists(id, "movie")
			if err != nil {
				log.Fatalf("Error checking if item exists: %v", err)
			}
			if exists {
				movieDetails, err := ms.FetchMovieDetailsFromTMDB(id)
				if err != nil {
					log.Printf("Failed to fetch TMDB details for %s: %v", id, err)
					return
				}
				var existingMovie models.Movie
				if err := database.FindItem(id, "movie", &existingMovie); err != nil {
					log.Printf("Failed to fetch existing movie: %v", err)
					return
				}
				if ms.hasMovieChanged(existingMovie, movieDetails) {
					if err := database.EditItem(movieDetails, "movie"); err != nil {
						log.Printf("Failed to update movie %s in database: %v", movieDetails.Title, err)
					}
				}
			} else {
				movieDetails, err := ms.FetchMovieDetailsFromTMDB(id)
				if err != nil {
					log.Printf("Failed to fetch TMDB details for %s: %v", id, err)
					return
				}
				if err := database.InsertItem(movieDetails, "movie"); err != nil {
					log.Printf("Failed to save movie %s to database: %v", movieDetails.Title, err)
				}
			}
		}(id)
	}

	wg.Wait()
	fmt.Println("Fetching movies completed")
	return nil
}

// FetchMovieDetailsFromTMDB retrieves detailed movie data from TMDB.
func (ms *MovieScraper) FetchMovieDetailsFromTMDB(id string) (*models.Movie, error) {
	url := fmt.Sprintf("%s/movie/%s?api_key=%s", utils.TMDBBaseURL, id, ms.tmdbAPIKey)
	var response struct {
		ID               int                   `json:"id"`
		Title            string                `json:"title"`
		OriginalTitle    string                `json:"original_title"`
		Overview         string                `json:"overview"`
		Popularity       float64               `json:"popularity"`
		PosterPath       string                `json:"poster_path"`
		ReleaseDate      string                `json:"release_date"`
		Genres           []models.GenreMapping `json:"genres"`
		VoteAverage      float64               `json:"vote_average"`
		VoteCount        int                   `json:"vote_count"`
		OriginalLanguage string                `json:"original_language"`
		BackdropPath     string                `json:"backdrop_path"`
	}
	err := utils.FetchJSON(url, "", &response)
	if err != nil {
		return nil, fmt.Errorf("error fetching movie details: %v", err)
	}

	return &models.Movie{
		ID:               strconv.Itoa(response.ID),
		Title:            response.Title,
		OriginalTitle:    response.OriginalTitle,
		Overview:         response.Overview,
		Popularity:       response.Popularity,
		PosterPath:       response.PosterPath,
		ReleaseDate:      response.ReleaseDate,
		Genres:           response.Genres,
		VoteAverage:      response.VoteAverage,
		VoteCount:        response.VoteCount,
		OriginalLanguage: response.OriginalLanguage,
		BackdropPath:     response.BackdropPath,
		LastUpdated:      time.Now().Format(time.RFC3339),
	}, nil
}

// hasMovieChanged checks if the movie details have changed.
func (ms *MovieScraper) hasMovieChanged(existingMovie models.Movie, newDetails *models.Movie) bool {
	// Compare relevant fields to determine if there are changes.
	return existingMovie.Overview != newDetails.Overview ||
		existingMovie.ReleaseDate != newDetails.ReleaseDate ||
		existingMovie.VoteAverage != newDetails.VoteAverage ||
		existingMovie.VoteCount != newDetails.VoteCount ||
		existingMovie.Popularity != newDetails.Popularity ||
		existingMovie.BackdropPath != newDetails.BackdropPath ||
		existingMovie.PosterPath != newDetails.PosterPath
}

/*
// FetchMovieIDsFromTMDB retrieves a list of popular movie IDs from TMDB.
func (ms *MovieScraper) FetchMovieIDsFromTMDB() ([]string, error) {
	today := time.Now().Format("01_02_2006")
	url := fmt.Sprintf("http://files.tmdb.org/p/exports/movie_ids_%s.json.gz", today)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error downloading movie IDs file: %v", err)
	}
	defer resp.Body.Close()

	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error creating gzip reader: %v", err)
	}
	defer gzipReader.Close()

	var ids []string
	scanner := bufio.NewScanner(gzipReader)
	for scanner.Scan() {
		var movie struct {
			ID int `json:"id"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &movie); err != nil {
			log.Printf("Error decoding JSON: %v\n", err)
			continue
		}
		ids = append(ids, strconv.Itoa(movie.ID))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading IDs file: %v", err)
	}

	return ids, nil
}


// Scrape orchestrates the fetching of movie data from multiple sources.
func (ms *MovieScraper) Scrape() error {
	ids, err := ms.FetchMovieIDsFromTMDB()
	if err != nil {
		return fmt.Errorf("error fetching movie IDs from TMDB: %v", err)
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 20) // Limit the number of concurrent goroutines

	for _, id := range ids {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(id string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			exists, err := database.IfItemExists(id, "movie")
			if err != nil {
				log.Fatalf("Error checking if item exists: %v", err)
			}
			if exists {
				return
			}
			movieDetails, err := ms.FetchMovieDetailsFromTMDB(id)
			if err != nil {
				log.Printf("Failed to fetch TMDB details for %s: %v", id, err)
				return
			}

			if exists {
				var existingMovie models.Movie
				if err := database.FindItem(id, "movie", &existingMovie); err != nil {
					log.Printf("Failed to fetch existing movie: %v", err)
					return
				}
			} else {
				if err := database.InsertItem(movieDetails, "movie"); err != nil {
					log.Printf("Failed to save movie %s to database: %v", movieDetails.Title, err)
				}
			}
		}(id)
	}

	wg.Wait()
	fmt.Println("Fetching movies completed")
	return nil
}

*/
