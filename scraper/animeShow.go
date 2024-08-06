package scraper

import (
	"log"
	"strings"
	"sync"
	"youvies-backend/database"
	"youvies-backend/models"
	"youvies-backend/utils"
)

type AnimeShowScraper struct {
	BaseScraper
}

func NewAnimeShowScraper() *AnimeShowScraper {
	return &AnimeShowScraper{
		BaseScraper: *NewBaseScraper("anime_show", utils.KitsuBaseURL),
	}
}

func (s *AnimeShowScraper) Scrape(animeList []models.Anime) error {
	log.Printf("Scraping %d anime shows", len(animeList))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 3) // Limit the number of concurrent goroutines

	for _, anime := range animeList {
		if strings.Contains(anime.Attributes.Slug, "delete") {
			continue
		}
		wg.Add(1)
		semaphore <- struct{}{}

		go func(anime models.Anime) {
			defer wg.Done()
			defer func() { <-semaphore }()
			exists, err := database.IfItemExists(anime.ID, "anime_shows")
			if err != nil {
				log.Printf("Error checking if item exists: %v", err)
				return
			}
			// Fetch the last updated timestamp from your database
			existingAnime := models.AnimeShow{}
			if err := database.FindItem(anime.ID, "anime_shows", &existingAnime); err != nil {
				log.Printf("Failed to fetch existing anime show: %v", err)
				return
			}

			// Compare the updated_at timestamp
			if anime.Attributes.UpdatedAt.After(existingAnime.Attributes.UpdatedAt) || !exists {
				// Fetch episodes
				episodes, err := utils.FetchAllEpisodes(anime.ID)
				if err != nil {
					log.Printf("Failed to fetch episodes for anime %s: %v", anime.Attributes.CanonicalTitle, err)
					return
				}
				if len(episodes) > anime.Attributes.EpisodeCount {
					anime.Attributes.EpisodeCount = len(episodes)
				}

				genres, err := utils.FetchGenres(anime.Relationships.Genres.Links.Related)
				if err != nil {
					log.Printf("Failed to fetch genres for %s: %v", anime.Attributes.CanonicalTitle, err)
					return
				}

				animeDoc := s.createAnimeShowDoc(anime, genres)
				if existingAnime.ID != "" {
					// Update the existing item
					if err := database.EditItem(animeDoc, "anime_shows"); err != nil {
						log.Printf("Failed to update anime show %s in database: %v", animeDoc.Title, err)
					}
				} else {
					// Insert the new item
					if err := database.InsertItem(animeDoc, "anime_shows"); err != nil {
						log.Printf("Failed to save anime show %s to database: %v", animeDoc.Title, err)
					}
				}
			}
		}(anime)
	}

	wg.Wait()
	log.Println("Fetching updated anime shows completed")
	return nil
}

// createAnimeShowDoc constructs an anime show document from Kitsu data.
func (s *AnimeShowScraper) createAnimeShowDoc(anime models.Anime, genres []models.GenreMapping) *models.AnimeShow {
	title := anime.Attributes.CanonicalTitle
	if title == "" {
		title = anime.Attributes.Titles.En
		if title == "" {
			title = anime.Attributes.Titles.EnUs
		} else {
			title = anime.Attributes.Titles.EnJp
		}
	}
	return &models.AnimeShow{
		ID:            anime.ID,
		Attributes:    anime.Attributes,
		Relationships: anime.Relationships,
		Genres:        genres,
		Title:         title,
	}
}

/*

// FetchAnimeDetailsFromKitsu fetches anime details from Kitsu with pagination
func (s *AnimeShowScraper) FetchAnimeDetailsFromKitsu() ([]models.AnimeResponse, error) {
	var allAnime []models.AnimeResponse
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit the number of concurrent goroutines

	pageNum := 1
	for {
		wg.Add(1)
		semaphore <- struct{}{}
		url := fmt.Sprintf("%s?page[limit]=20&page[offset]=%d", s.BaseScraper.BaseURL, pageNum*20)

		go func(url string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			resp, err := http.Get(url)
			if err != nil {
				log.Printf("got an error while fetching url %s: %v", url, err)
				return
			}
			if resp.StatusCode != http.StatusOK {
				log.Printf("got a non-200 status code while getting animes: %d", resp.StatusCode)
				log.Printf("link was: %s", url)
				return
			}

			var animes models.AnimeResponse
			if err := json.NewDecoder(resp.Body).Decode(&animes); err != nil {
				log.Printf("error decoding response body: %v", err)
				return
			}
			var filtredAnimes []models.Anime
			for _, anime := range animes.Data {
				if anime.Attributes.Subtype == "movie" {
					continue
				}
				filtredAnimes = append(filtredAnimes, anime)
			}
			animes.Data = filtredAnimes
			allAnime = append(allAnime, animes)
			err = resp.Body.Close()
			if err != nil {
				log.Printf("error closing response body: %v", err)
			}
		}(url)

		// Check if there's a next page
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("got an error while fetching url %s: %v", url, err)
			break
		}
		var animes models.AnimeResponse
		if err := json.NewDecoder(resp.Body).Decode(&animes); err != nil {
			log.Printf("error decoding response body: %v", err)
			break
		}
		if animes.Links.Next == "" {
			break
		}
		pageNum++
		if pageNum%10 == 0 {
			fmt.Print(pageNum, " ")
		}
	}
	wg.Wait()
	fmt.Printf("found this many animes: %d\n", len(allAnime))
	return allAnime, nil
}

func (s *AnimeShowScraper) Scrape() error {
	animes, err := s.FetchAnimeDetailsFromKitsu()
	if err != nil {
		return fmt.Errorf("error fetching Kitsu anime data: %v", err)
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 20) // Limit the number of concurrent goroutines

	for _, animeResp := range animes {
		for _, anime := range animeResp.Data {
			if strings.Contains(anime.Attributes.Slug, "delete") {
				continue
			}
			wg.Add(1)
			semaphore <- struct{}{}

			go func(anime models.Anime) {
				defer wg.Done()
				defer func() { <-semaphore }()
				exists, err := database.IfItemExists(anime.ID, "anime_shows")
				if err != nil {
					log.Printf("Error checking if item exists: %v", err)
					return
				}
				if exists {
					return
				}
				// Fetch episodes
				episodes, err := utils.FetchAllEpisodes(anime.ID)
				if err != nil {
					log.Printf("Failed to fetch episodes for anime %s: %v", anime.Attributes.CanonicalTitle, err)
					return
				}
				if len(episodes) > anime.Attributes.EpisodeCount {
					anime.Attributes.EpisodeCount = len(episodes)
				}

				genres, err := utils.FetchGenres(anime.Relationships.Genres.Links.Related)
				if err != nil {
					log.Printf("Failed to fetch genres for %s: %v", anime.Attributes.CanonicalTitle, err)
					return
				}

				animeDoc := s.createAnimeShowDoc(anime, genres)

				if err := database.InsertItem(animeDoc, "anime_shows"); err != nil {
					log.Printf("Failed to save anime show %s to database: %v", animeDoc.Title, err)
					return
				}

				for _, episode := range episodes {
					if _, err := database.InsertEpisode(episode); err != nil {
						log.Printf("Failed to insert episode %d for anime show %s: %v", episode.Number, animeDoc.Title, err)
					}
				}
			}(anime)
		}
	}

	wg.Wait()
	log.Println("Fetching new anime shows completed")
	return nil
}

// createAnimeShowDoc constructs an anime show document from Kitsu data.
func (s *AnimeShowScraper) createAnimeShowDoc(anime models.Anime, genres []models.GenreMapping) *models.AnimeShow {
	title := anime.Attributes.CanonicalTitle
	if title == "" {
		title = anime.Attributes.Titles.En
		if title == "" {
			title = anime.Attributes.Titles.EnUs
		} else {
			title = anime.Attributes.Titles.EnJp
		}
	}
	return &models.AnimeShow{
		ID:            anime.ID,
		Attributes:    anime.Attributes,
		Relationships: anime.Relationships,
		Genres:        genres,
		Title:         title, // Including the title attribute
	}
}

*/
