package scraper

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"youvies-backend/database"
	"youvies-backend/models"
	"youvies-backend/utils"
)

type AnimeMovieScraper struct {
	BaseScraper
}

func NewAnimeMovieScraper() *AnimeMovieScraper {
	return &AnimeMovieScraper{
		BaseScraper: *NewBaseScraper("anime_movie", "https://kitsu.io/api/edge/anime"),
	}
}

// fetchAnimeMovies fetches popular anime IDs from Kitsu
func (s *AnimeMovieScraper) fetchAnimeMovies() ([]models.AnimeResponse, error) {
	var allAnimes []models.AnimeResponse
	url := fmt.Sprintf("%s/anime?filter[subtype]=movie", utils.KitsuBaseURL)
	for url != "" {
		resp, err := s.FetchURL(url)
		if err != nil {
			return nil, err
		}

		var animes models.AnimeResponse
		if err := json.NewDecoder(resp.Body).Decode(&animes); err != nil {
			return nil, err
		}

		allAnimes = append(allAnimes, animes)
		url = animes.Links.Next
		err = resp.Body.Close()
		if err != nil {
			return nil, err
		}

	}
	fmt.Printf("found this much animes movies: %d", len(allAnimes))

	return allAnimes, nil
}

// Scrape fetches data from various APIs and inserts them into the database
func (s *AnimeMovieScraper) Scrape() error {
	animes, err := s.fetchAnimeMovies()
	if err != nil {
		return fmt.Errorf("error fetching anime IDs: %v", err)
	}

	for _, data := range animes {
		for _, anime := range data.Data {
			name := "attributes.titles.en"

			exict, _ := database.IfItemExists(bson.M{name: anime.Attributes.Titles.En}, "anime_movies")
			if exict {
				continue
			}
			// Fetch torrents from AniDex
			torrents, err := utils.FetchTorrents(anime.Attributes.Titles.En, []string{"Anime", "Movies"})
			if err != nil {
				log.Printf("error fetching torrents: %v", err)
			}
			anime.Torrents = torrents
			exists, err := database.IfItemExists(bson.M{"attributes.titles.en": anime.Attributes.Titles.En}, "anime_movies")
			if err != nil {
				log.Fatalf("Error checking if item exists: %v", err)
			}
			if exists {
				err = database.EditItem(bson.M{"attributes.titles.en": anime.Attributes.Titles.En}, anime, "anime_movies")
				if err != nil {
					return fmt.Errorf("error editing anime movies: %v", err)
				}

			} else {
				err = database.InsertItem(anime, anime.Attributes.Titles.En, "anime_movies")
				if err != nil {
					return fmt.Errorf("error inserting anime movies into database: %v", err)
				}

			}
		}
	}
	fmt.Println("Fetching new anime movies completed")
	return nil
}
