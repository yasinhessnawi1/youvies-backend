package scraper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"strconv"
	"youvies-backend/database"
	"youvies-backend/models"
	"youvies-backend/utils"
)

type ShowScraper struct {
	BaseScraper
	TMDBKey string
	TVDBKey string
}

func NewShowScraper(tmdbKey, tvdbKey string) *ShowScraper {
	return &ShowScraper{
		TMDBKey: tmdbKey,
		TVDBKey: tvdbKey,
	}
}

func (s *ShowScraper) FetchShowsFromTMDB(showID string) (*models.Show, error) {
	url := fmt.Sprintf("%s/tv/%s?api_key=%s", utils.TMDBBaseURL, showID, s.TMDBKey)
	var response struct {
		Name         string `json:"name"`
		Overview     string `json:"overview"`
		FirstAirDate string `json:"first_air_date"`
		Networks     []struct {
			Name string `json:"name"`
		} `json:"networks"`
		PosterPath        string   `json:"poster_path"`
		VoteAverage       float64  `json:"vote_average"`
		Original_language string   `json:"original_language"`
		Origin_Country    []string `json:"origin_country"`
		Backdrop_Path     string   `json:"backdrop_path"`
	}
	err := utils.FetchJSON(url, "TMDB", &response)
	if err != nil {
		return nil, fmt.Errorf("error fetching TMDB show data: %v", err)
	}

	var networks []string
	for _, network := range response.Networks {
		networks = append(networks, network.Name)
	}

	return &models.Show{
		Title:        response.Name,
		Description:  response.Overview,
		FirstAirDate: response.FirstAirDate,
		Networks:     networks,
		ImageURL:     fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", response.PosterPath),
		Rating:       response.VoteAverage,
		Year:         utils.GetYear(response.FirstAirDate),
		Language:     response.Original_language,
		Country:      response.Origin_Country,
		Backdrop:     fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", response.Backdrop_Path),
	}, nil
}

func (s *ShowScraper) FetchShowIDsFromTMDB() ([]string, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/tv/popular?api_key=%s", s.TMDBKey)
	var response struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}

	err := s.FetchJSON(url, &response)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, result := range response.Results {
		ids = append(ids, fmt.Sprintf("%d", result.ID))
	}
	return ids, nil
}
func (s *ShowScraper) Scrape() error {
	showIDs, err := s.FetchShowIDsFromTMDB()
	if err != nil {
		log.Printf("Error fetching show IDs from TMDB: %v\n", err)
		return err
	}

	for _, id := range showIDs {
		show, err := s.FetchShowsFromTMDB(id)
		if err != nil {
			return fmt.Errorf("error fetching TMDB show data: %v", err)
		}
		exict, _ := database.IfItemExists(bson.M{"title": show.Title}, "shows")
		if exict {
			continue
		}
		torrents, err := utils.FetchTorrents(show.Title, []string{"Show", "TV", "TV Show", "Shows", "Series"})
		if err != nil {
			return fmt.Errorf("error fetching torrents: %v", err)
		}
		show.Episodes = torrents
		err = database.InsertItem(show, show.Title, "shows")
		if err != nil {
			log.Printf("error inserting shows into database: %v", err)
		}
	}
	log.Println("Fetching new shows completed")
	return nil
}

func (s *ShowScraper) FetchOldShows() {
	var ids []string
	file, err := os.Open("utils/tv_series/tv_series_ids_05_15_2024.json")
	if err != nil {
		log.Fatal(fmt.Sprintf("Error while oppning ids file: %s", err))
	}
	defer file.Close()

	// Scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Unmarshal the JSON line into a Movie struct
		type Show struct {
			ID           int     `json:"id"`
			OriginalName string  `json:"original_name"`
			Popularity   float64 `json:"popularity"`
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

	for _, id := range ids {
		show, err := s.FetchShowsFromTMDB(id)
		if err != nil {
			log.Printf("error fetching TMDB show data: %v", err)
		}

		exist, err := database.IfItemExists(bson.M{"title": show.Title}, "shows")
		if err != nil || exist {
			continue
		}
		torrents, err := utils.FetchTorrents(show.Title, []string{"Show", "TV", "TV Show", "Shows", "Series"})
		if err != nil {
			log.Printf("error fetching torrents: %v", err)
		}
		show.Episodes = torrents
		err = database.InsertItem(show, show.Title, "shows")
		if err != nil {
			log.Printf("error inserting shows into database: %v", err)
		}

	}
	log.Println("Fetching old shows is completed")
}
