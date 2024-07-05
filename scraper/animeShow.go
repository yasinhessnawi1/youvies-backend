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

type AnimeShowScraper struct {
	BaseScraper
}

func NewAnimeShowScraper() *AnimeShowScraper {
	return &AnimeShowScraper{
		BaseScraper: *NewBaseScraper("anime_show", "https://kitsu.io/api/edge/anime"),
	}
}

// FetchAnimeDetailsFromKitsu fetches anime details from Kitsu with pagination
func (s *AnimeShowScraper) FetchAnimeDetailsFromKitsu() ([]models.AnimeResponse, error) {
	var allAnimes []models.AnimeResponse
	url := fmt.Sprintf("%s/anime", utils.KitsuBaseURL)

	for url != "" {
		resp, err := s.FetchURL(url)
		if err != nil {
			return nil, fmt.Errorf("Got an error while fetching url %s: %v", url, err)
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
	fmt.Printf("found this much animes: %d", len(allAnimes))
	return allAnimes, nil
}

// Scrape fetches data from various APIs and inserts them into the database
func (s *AnimeShowScraper) Scrape() error {

	animes, err := s.FetchAnimeDetailsFromKitsu()
	if err != nil {
		return fmt.Errorf("error fetching Kitsu anime data: %v", err)
	}
	for _, data := range animes {
		for _, anime := range data.Data {
			name := "attributes.titles.en"
			exict, _ := database.IfItemExists(bson.M{name: anime.Attributes.Titles.En}, "anime_shows")
			if exict {
				continue
			}
			// Fetch torrents from AniDex
			torrents, err := utils.FetchTorrents(anime.Attributes.Titles.En, []string{"Anime", "TV", "Show"})
			if err != nil {
				log.Printf("error fetching torrents: %v", err)
			}
			anime.Torrents = torrents
			exists, err := database.IfItemExists(bson.M{"attributes.titles.en": anime.Attributes.Titles.En}, "anime_show")
			if err != nil {
				log.Fatalf("Error checking if item exists: %v", err)
			}
			if exists {
				err = database.EditItem(bson.M{"attributes.titles.en": anime.Attributes.Titles.En}, anime, "anime_show")
				if err != nil {
					return fmt.Errorf("error editing show anime: %v", err)
				}
			} else {
				err = database.InsertItem(anime, anime.Attributes.Titles.En, "anime_shows")
				if err != nil {
					return fmt.Errorf("error inserting show anime: %v", err)
				}
			}
		}
	}
	log.Println("Fetching new anime shows completed")
	return nil
}
