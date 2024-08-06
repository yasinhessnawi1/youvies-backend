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

type ShowScraper struct {
	BaseScraper
	tmdbAPIKey string
}

func NewShowScraper(tmdbKey string) *ShowScraper {
	return &ShowScraper{
		BaseScraper: *NewBaseScraper("show", utils.TMDBBaseURL),
		tmdbAPIKey:  tmdbKey,
	}
}

func (ss *ShowScraper) FetchChangedShowIDs() ([]string, error) {
	log.Println("Fetching changed show IDs")
	url := fmt.Sprintf("%s/tv/changes?api_key=%s", utils.TMDBBaseURL, ss.tmdbAPIKey)
	var response struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}
	if err := ss.FetchJSON(url, &response); err != nil {
		return nil, err
	}
	var ids []string
	for _, result := range response.Results {
		ids = append(ids, strconv.Itoa(result.ID))
	}
	log.Printf("%v Changed show IDs fetched", len(ids))
	return ids, nil
}

func (ss *ShowScraper) FetchAiringTodayShowIDs() ([]string, error) {
	log.Println("Fetching airing today show IDs")
	url := fmt.Sprintf("%s/tv/airing_today?api_key=%s", utils.TMDBBaseURL, ss.tmdbAPIKey)
	var response struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}
	if err := ss.FetchJSON(url, &response); err != nil {
		return nil, err
	}
	var ids []string
	for _, result := range response.Results {
		ids = append(ids, strconv.Itoa(result.ID))
	}
	log.Printf("%v Airing show IDs fetched", len(ids))
	return ids, nil
}

// Scrape orchestrates the fetching of show data from multiple sources.
func (ss *ShowScraper) Scrape() error {
	log.Println("Scraping shows")
	changedIDs, err := ss.FetchChangedShowIDs()
	if err != nil {
		return fmt.Errorf("error fetching changed show IDs from TMDB: %v", err)
	}

	airingTodayIDs, err := ss.FetchAiringTodayShowIDs()
	if err != nil {
		return fmt.Errorf("error fetching airing today show IDs from TMDB: %v", err)
	}

	ids := append(changedIDs, airingTodayIDs...)
	ids = utils.RemoveDuplicateStrings(ids) // Remove duplicates

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 20) // Limit the number of concurrent goroutines

	for _, id := range ids {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(id string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			exists, err := database.IfItemExists(id, "show")
			if err != nil {
				log.Fatalf("Error checking if item exists: %v", err)
			}
			if exists {
				showDetails, err := ss.FetchShowDetailsFromTMDB(id)
				if err != nil {
					log.Printf("Failed to fetch TMDB details for %s: %v", id, err)
					return
				}
				seasons, err := ss.FetchSeasonDetailsFromTMDB(id)
				if err != nil {
					log.Printf("Failed to fetch seasons for show %s: %v", showDetails.Title, err)
					return
				}
				showDetails.SeasonsInfo = seasons
				var existingShow models.Show
				if err := database.FindItem(id, "show", &existingShow); err != nil {
					log.Printf("Failed to fetch existing show: %v", err)
					return
				}
				if ss.hasShowChanged(existingShow, showDetails) {
					if err := database.EditItem(showDetails, "show"); err != nil {
						log.Printf("Failed to update show %s in database: %v", showDetails.Title, err)
					}
				}
			} else {
				showDetails, err := ss.FetchShowDetailsFromTMDB(id)
				if err != nil {
					log.Printf("Failed to fetch TMDB details for %s: %v", id, err)
					return
				}
				seasons, err := ss.FetchSeasonDetailsFromTMDB(id)
				if err != nil {
					log.Printf("Failed to fetch seasons for show %s: %v", showDetails.Title, err)
					return
				}
				showDetails.SeasonsInfo = seasons
				if err := database.InsertItem(showDetails, "show"); err != nil {
					log.Printf("Failed to save show %s to database: %v", showDetails.Title, err)
				}
			}
		}(id)
	}

	wg.Wait()
	fmt.Println("Fetching shows completed")
	return nil
}

// FetchShowDetailsFromTMDB retrieves detailed show data from TMDB.
func (ss *ShowScraper) FetchShowDetailsFromTMDB(id string) (*models.Show, error) {
	url := fmt.Sprintf("%s/tv/%s?api_key=%s", utils.TMDBBaseURL, id, ss.tmdbAPIKey)
	var response struct {
		ID               int                   `json:"id"`
		Name             string                `json:"name"`
		Overview         string                `json:"overview"`
		FirstAirDate     string                `json:"first_air_date"`
		Genres           []models.GenreMapping `json:"genres"`
		PosterPath       string                `json:"poster_path"`
		VoteAverage      float64               `json:"vote_average"`
		VoteCount        int                   `json:"vote_count"`
		OriginalLanguage string                `json:"original_language"`
		Popularity       float64               `json:"popularity"`
		BackdropPath     string                `json:"backdrop_path"`
	}
	err := utils.FetchJSON(url, "", &response)
	if err != nil {
		return nil, fmt.Errorf("error fetching show details: %v", err)
	}

	return &models.Show{
		ID:               strconv.Itoa(response.ID),
		Title:            response.Name,
		Overview:         response.Overview,
		FirstAirDate:     response.FirstAirDate,
		Genres:           response.Genres,
		PosterPath:       response.PosterPath,
		VoteAverage:      response.VoteAverage,
		VoteCount:        response.VoteCount,
		OriginalLanguage: response.OriginalLanguage,
		Popularity:       response.Popularity,
		BackdropPath:     response.BackdropPath,
		LastUpdated:      time.Now().Format(time.RFC3339),
	}, nil
}

// FetchSeasonDetailsFromTMDB retrieves detailed season data from TMDB.
func (ss *ShowScraper) FetchSeasonDetailsFromTMDB(showID string) ([]models.SeasonInfo, error) {
	var seasons []models.SeasonInfo
	url := fmt.Sprintf("%s/tv/%s?api_key=%s&append_to_response=seasons", utils.TMDBBaseURL, showID, ss.tmdbAPIKey)
	var response struct {
		Seasons []models.SeasonInfo `json:"seasons"`
	}
	err := utils.FetchJSON(url, "", &response)
	if err != nil {
		return nil, fmt.Errorf("error fetching season details: %v", err)
	}
	seasons = response.Seasons
	return seasons, nil
}

// hasShowChanged checks if the show details have changed.
func (ss *ShowScraper) hasShowChanged(existingShow models.Show, newDetails *models.Show) bool {
	// Compare relevant fields to determine if there are changes.
	return existingShow.Overview != newDetails.Overview ||
		existingShow.VoteAverage != newDetails.VoteAverage ||
		existingShow.VoteCount != newDetails.VoteCount ||
		existingShow.Popularity != newDetails.Popularity ||
		existingShow.BackdropPath != newDetails.BackdropPath ||
		existingShow.PosterPath != newDetails.PosterPath
}

/*
// FetchShowIDsFromTMDB retrieves a list of popular show IDs from TMDB.
func (ss *ShowScraper) FetchShowIDsFromTMDB() ([]string, error) {
	today := time.Now().Format("01_02_2006")
	url := fmt.Sprintf("http://files.tmdb.org/p/exports/tv_series_ids_%s.json.gz", today)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error downloading show IDs file: %v", err)
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
		var show struct {
			ID int `json:"id"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &show); err != nil {
			log.Printf("Error decoding JSON: %v\n", err)
			continue
		}
		ids = append(ids, strconv.Itoa(show.ID))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading IDs file: %v", err)
	}

	return ids, nil
}


// Scrape orchestrates the fetching of show data from multiple sources.
func (ss *ShowScraper) Scrape() error {
	ids, err := ss.FetchShowIDsFromTMDB()
	if err != nil {
		return fmt.Errorf("error fetching show IDs from TMDB: %v", err)
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 20) // Limit the number of concurrent goroutines

	for _, id := range ids {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(id string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			exists, err := database.IfItemExists(id, "show")
			if err != nil {
				log.Fatalf("Error checking if item exists: %v", err)
			}
			if exists {
				return
			}
			showDetails, err := ss.FetchShowDetailsFromTMDB(id)
			if err != nil {
				log.Printf("Failed to fetch TMDB details for %s: %v", id, err)
				return
			}

			seasons, err := ss.FetchSeasonDetailsFromTMDB(id)
			if err != nil {
				log.Printf("Failed to fetch seasons for show %s: %v", showDetails.Title, err)
				return
			}

			showDetails.SeasonsInfo = seasons

			if exists {
				var existingShow models.Show
				if err := database.FindItem(id, "show", &existingShow); err != nil {
					log.Printf("Failed to fetch existing show: %v", err)
					return
				}
			} else {
				if err := database.InsertItem(showDetails, "show"); err != nil {
					log.Printf("Failed to save show %s to database: %v", showDetails.Title, err)
				}
			}
		}(id)
	}

	wg.Wait()
	fmt.Println("Fetching shows completed")
	return nil
}

*/
