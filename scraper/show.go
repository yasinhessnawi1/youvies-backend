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
	"sort"
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
		BaseScraper: *NewBaseScraper("show", tmdbBaseURL),
		tmdbAPIKey:  tmdbKey,
	}
}

// Scrape orchestrates the fetching of show data from multiple sources.
func (ss *ShowScraper) Scrape() error {
	ids, err := ss.FetchShowIDsFromTMDB()
	if err != nil {
		return fmt.Errorf("error fetching show IDs from TMDB: %v", err)
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit the number of concurrent goroutines
	sort.Strings(ids)
	for _, id := range ids {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(id string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			showDetails, err := ss.FetchShowDetailsFromTMDB(id)
			if err != nil {
				log.Printf("Failed to fetch TMDB details for %s: %v", id, err)
				return
			}

			exists, err := database.IfItemExists(bson.M{"title": showDetails.Title}, "shows")
			if err != nil {
				log.Printf("Database error: %v", err)
				return
			}
			if exists && showDetails.Title == "" {
				log.Printf("Show %s already exists in database", showDetails.Title)
				return
			}

			torrents, err := utils.FetchTorrents(showDetails.Title)
			if err != nil {
				log.Printf("Failed to fetch torrents for %s: %v", showDetails.Title, err)
				return
			}

			// Categorize torrents by seasons and episodes
			categorizedTorrents, extra := utils.CategorizeTorrentsBySeasonsAndEpisodes(torrents)
			showDetails.OtherTorrents = extra
			// Check for missing episodes and fetch them if necessary
			missingTorrents, err := utils.FetchMissingTorrents(showDetails.Title, torrents, showDetails.SeasonsInfo)
			if err != nil {
				log.Printf("Failed to fetch missing torrents for %s: %v", showDetails.Title, err)
				return
			}
			if len(missingTorrents) > 0 {
				missingCategorizedTorrents, extra := utils.CategorizeTorrentsBySeasonsAndEpisodes(missingTorrents)
				showDetails.OtherTorrents = append(showDetails.OtherTorrents, extra...)
				for seasonNum, season := range missingCategorizedTorrents {
					if _, exists := categorizedTorrents[seasonNum]; !exists {
						categorizedTorrents[seasonNum] = season
					} else {
						for episodeNum, episode := range season.Episodes {
							if _, exists := categorizedTorrents[seasonNum].Episodes[episodeNum]; !exists {
								categorizedTorrents[seasonNum].Episodes[episodeNum] = episode
							} else {
								for quality, torrents := range episode.Torrents {
									categorizedTorrents[seasonNum].Episodes[episodeNum].Torrents[quality] = append(
										categorizedTorrents[seasonNum].Episodes[episodeNum].Torrents[quality],
										torrents...,
									)
								}
							}
						}
					}
				}
			}

			// Update existing show if changes are found
			if exists {
				// Fetch existing show
				var existingShow models.Show
				if err := database.FindItem(bson.M{"title": showDetails.Title}, "shows", &existingShow); err != nil {
					log.Printf("Failed to fetch existing show: %v", err)
					return
				}
				// Check for updates
				if ss.hasShowChanged(existingShow, showDetails, categorizedTorrents) {
					show := ss.createShowDoc(showDetails, categorizedTorrents, existingShow.ID)
					if err := database.EditItem(bson.M{"title": show.Title}, show, "shows"); err != nil {
						log.Printf("Failed to update show %s in database: %v", show.Title, err)
					}
				}
			} else {
				show := ss.createShowDoc(showDetails, categorizedTorrents, -1)
				show.ID = utils.GetNextShowID()
				if err := database.InsertItem(show, show.Title, "shows"); err != nil {
					log.Printf("Failed to save show %s to database: %v", show.Title, err)
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
	url := fmt.Sprintf("%s/tv/%s?api_key=%s", tmdbBaseURL, id, ss.tmdbAPIKey)
	var response struct {
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
		Adult            bool                  `json:"adult"`
		Rating           string                `json:"rating"`
		NumberOfSeasons  int                   `json:"number_of_seasons"`
		NumberOfEpisodes int                   `json:"number_of_episodes"`
		Seasons          []models.SeasonInfo   `json:"seasons"`
		Networks         []struct {
			ID            int    `json:"id"`
			Name          string `json:"name"`
			OriginCountry string `json:"origin_country"`
		} `json:"networks"`
		ProductionCompanies []struct {
			ID            int    `json:"id"`
			Name          string `json:"name"`
			OriginCountry string `json:"origin_country"`
		} `json:"production_companies"`
		ProductionCountries []struct {
			Iso3166_1 string `json:"iso_3166_1"`
			Name      string `json:"name"`
		} `json:"production_countries"`
		SpokenLanguages []struct {
			EnglishName string `json:"english_name"`
			Iso639_1    string `json:"iso_639_1"`
			Name        string `json:"name"`
		} `json:"spoken_languages"`
		ExternalIDs struct {
			ImdbID string `json:"imdb_id"`
		} `json:"external_ids"`
	}
	err := utils.FetchJSON(url, "", &response)
	if err != nil {
		return nil, fmt.Errorf("error fetching show details: %v", err)
	}

	externalIDs, err := ss.FetchExternalIDs(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching external ids: %v", err)
	}

	return &models.Show{
		Title:               response.Name,
		Overview:            response.Overview,
		FirstAirDate:        response.FirstAirDate,
		Genres:              response.Genres,
		PosterPath:          response.PosterPath,
		VoteAverage:         response.VoteAverage,
		VoteCount:           response.VoteCount,
		OriginalLanguage:    response.OriginalLanguage,
		Popularity:          response.Popularity,
		BackdropPath:        response.BackdropPath,
		SeasonsInfo:         response.Seasons,
		Networks:            ss.parseNetworks(response.Networks),
		ProductionCompanies: ss.parseProductionCompanies(response.ProductionCompanies),
		ProductionCountries: ss.parseProductionCountries(response.ProductionCountries),
		SpokenLanguages:     ss.parseSpokenLanguages(response.SpokenLanguages),
		ExternalIDs:         externalIDs,
	}, nil
}

// FetchExternalIDs fetches external IDs for a show from TMDB.
func (ss *ShowScraper) FetchExternalIDs(id string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/tv/%s/external_ids?api_key=%s", tmdbBaseURL, id, ss.tmdbAPIKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ss.tmdbAPIKey))

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
		return nil, fmt.Errorf("error decoding external ids %v: %v", url, err)
	}

	return externalIDs, nil
}

// FetchShowIDsFromTMDB retrieves a list of popular show IDs from TMDB.
func (ss *ShowScraper) FetchShowIDsFromTMDB() ([]string, error) {
	url := fmt.Sprintf("%s/tv/popular?api_key=%s", tmdbBaseURL, ss.tmdbAPIKey)
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

	file, err := os.Open("utils/tv_series_ids_05_15_2024.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Unmarshal the JSON line into a Show struct
		type Show struct {
			Adult         bool    `json:"adult"`
			ID            int     `json:"id"`
			OriginalTitle string  `json:"original_name"`
			Popularity    float64 `json:"popularity"`
			Video         bool    `json:"video"`
		}
		var show Show
		err := json.Unmarshal([]byte(line), &show)
		if err != nil {
			log.Printf("Error decoding JSON: %v\n", err)
			continue
		}

		// Append the ID to the array
		ids = append(ids, strconv.Itoa(show.ID))
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file: %v\n", err)
	}

	return ids, nil
}

// createShowDoc constructs a show document from TMDB data.
func (ss *ShowScraper) createShowDoc(showDetails *models.Show, torrents map[int]models.Season, id int) models.Show {
	if id == -1 {
		id = utils.GetNextShowID()
	}
	return models.Show{
		ID:                  id,
		Title:               showDetails.Title,
		Overview:            showDetails.Overview,
		FirstAirDate:        showDetails.FirstAirDate,
		Genres:              showDetails.Genres,
		PosterPath:          showDetails.PosterPath,
		VoteAverage:         showDetails.VoteAverage,
		VoteCount:           showDetails.VoteCount,
		OriginalLanguage:    showDetails.OriginalLanguage,
		Popularity:          showDetails.Popularity,
		BackdropPath:        showDetails.BackdropPath,
		Seasons:             torrents,
		Networks:            showDetails.Networks,
		ProductionCompanies: showDetails.ProductionCompanies,
		ProductionCountries: showDetails.ProductionCountries,
		SpokenLanguages:     showDetails.SpokenLanguages,
		ExternalIDs:         showDetails.ExternalIDs,
		SeasonsInfo:         showDetails.SeasonsInfo,
		LastUpdated:         time.Now().Format(time.RFC3339),
		OtherTorrents:       showDetails.OtherTorrents,
		Country:             showDetails.Country,
	}
}

// hasShowChanged checks if the show details or torrents have changed.
func (ss *ShowScraper) hasShowChanged(existingShow models.Show, newDetails *models.Show, newTorrents map[int]models.Season) bool {
	// Compare relevant fields to determine if there are changes.
	return existingShow.Overview != newDetails.Overview ||
		existingShow.FirstAirDate != newDetails.FirstAirDate ||
		existingShow.VoteAverage != newDetails.VoteAverage ||
		existingShow.VoteCount != newDetails.VoteCount ||
		existingShow.Popularity != newDetails.Popularity ||
		existingShow.BackdropPath != newDetails.BackdropPath ||
		existingShow.PosterPath != newDetails.PosterPath ||
		!compareTorrentsBySeason(existingShow.Seasons, newTorrents)
}

// compareTorrentsBySeason compares two maps of torrents to see if they are different.
func compareTorrentsBySeason(oldTorrents, newTorrents map[int]models.Season) bool {
	if len(oldTorrents) != len(newTorrents) {
		return false
	}
	for season, oldSeason := range oldTorrents {
		newSeason, exists := newTorrents[season]
		if !exists || len(oldSeason.Episodes) != len(newSeason.Episodes) {
			return false
		}
		for episode, oldEpisode := range oldSeason.Episodes {
			newEpisode, exists := newSeason.Episodes[episode]
			if !exists || len(oldEpisode.Torrents) != len(newEpisode.Torrents) {
				return false
			}
			for quality, oldList := range oldEpisode.Torrents {
				newList, exists := newEpisode.Torrents[quality]
				if !exists || len(oldList) != len(newList) {
					return false
				}
				for i, oldTorrent := range oldList {
					if oldTorrent.Name != newList[i].Name || oldTorrent.Seeders != newList[i].Seeders {
						return false
					}
				}
			}
		}
	}
	return true
}

// parseNetworks parses the networks data from TMDB response.
func (ss *ShowScraper) parseNetworks(networks []struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}) []string {
	result := make([]string, len(networks))
	for i, network := range networks {
		result[i] = network.Name
	}
	return result
}

// parseProductionCompanies parses the production companies data from TMDB response.
func (ss *ShowScraper) parseProductionCompanies(companies []struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}) []string {
	result := make([]string, len(companies))
	for i, company := range companies {
		result[i] = company.Name
	}
	return result
}

// parseProductionCountries parses the production countries data from TMDB response.
func (ss *ShowScraper) parseProductionCountries(countries []struct {
	Iso3166_1 string `json:"iso_3166_1"`
	Name      string `json:"name"`
}) []string {
	result := make([]string, len(countries))
	for i, country := range countries {
		result[i] = country.Name
	}
	return result
}

// parseSpokenLanguages parses the spoken languages data from TMDB response.
func (ss *ShowScraper) parseSpokenLanguages(languages []struct {
	EnglishName string `json:"english_name"`
	Iso639_1    string `json:"iso_639_1"`
	Name        string `json:"name"`
}) []string {
	result := make([]string, len(languages))
	for i, language := range languages {
		result[i] = language.EnglishName
	}
	return result
}
