package scraper

import (
	"log"
	"strings"
	"youvies-backend/database"
	"youvies-backend/models"
	"youvies-backend/utils"
)

type AnimeMovieScraper struct {
	BaseScraper
}

func NewAnimeMovieScraper() *AnimeMovieScraper {
	log.Printf("Creating new anime movie scraper")
	return &AnimeMovieScraper{
		BaseScraper: *NewBaseScraper("anime_movie", utils.KitsuBaseURL),
	}
}

func (s *AnimeMovieScraper) Scrape(animeList []models.Anime) error {
	for _, anime := range animeList {
		if strings.Contains(anime.Attributes.Slug, "delete") {
			continue
		}

		// Fetch the last updated timestamp from your database
		existingAnime := models.AnimeMovie{}
		if err := database.FindItem(anime.ID, "anime_movies", &existingAnime); err != nil {
			log.Printf("Failed to fetch existing anime movie: %v", err)
			continue
		}
		exists, err := database.IfItemExists(anime.ID, "anime_movies")
		if err != nil {
			log.Printf("Error checking if item exists: %v", err)
			continue
		}

		// Compare the updated_at timestamp
		if anime.Attributes.UpdatedAt.After(existingAnime.Attributes.UpdatedAt) || !exists {
			// Fetch genres
			genres, err := utils.FetchGenres(anime.Relationships.Genres.Links.Related)
			if err != nil {
				log.Printf("Failed to fetch genres for %s: %v", anime.Attributes.CanonicalTitle, err)
				continue
			}

			animeDoc := s.createAnimeMovieDoc(anime, genres)

			if existingAnime.ID != "" {
				// Update the existing item
				if err := database.EditItem(animeDoc, "anime_movies"); err != nil {
					log.Printf("Failed to update anime movie %s in database: %v", animeDoc.Title, err)
				}
			} else {
				// Insert the new item
				if err := database.InsertItem(animeDoc, "anime_movies"); err != nil {
					log.Printf("Failed to save anime movie %s to database: %v", animeDoc.Title, err)
				}
			}
		}
	}
	log.Println("Fetching updated anime movies completed")
	return nil
}

func (s *AnimeMovieScraper) createAnimeMovieDoc(anime models.Anime, genres []models.GenreMapping) *models.AnimeMovie {
	title := anime.Attributes.CanonicalTitle
	if title == "" {
		title = anime.Attributes.Titles.En
		if title == "" {
			title = anime.Attributes.Titles.EnUs
		} else {
			title = anime.Attributes.Titles.EnJp
		}
	}
	return &models.AnimeMovie{
		ID:            anime.ID,
		Attributes:    anime.Attributes,
		Relationships: anime.Relationships,
		Genres:        genres,
		Title:         title,
	}
}

/*

func (s *AnimeMovieScraper) FetchAnimeMoviesFromKitsu() ([]models.AnimeResponse, error) {
	var allAnimes []models.AnimeResponse
	url := fmt.Sprintf("%s?filter[subtype]=movie", s.BaseScraper.BaseURL)

	for url != "" {
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("got an error while fetching url %s: %v", url, err)
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
	fmt.Printf("found this many anime movies: %d\n", len(allAnimes))
	return allAnimes, nil
}

func (s *AnimeMovieScraper) Scrape() error {
	animes, err := s.FetchAnimeMoviesFromKitsu()
	if err != nil {
		return fmt.Errorf("error fetching Kitsu anime movie data: %v", err)
	}
	for _, animeResp := range animes {
		for _, anime := range animeResp.Data {
			if strings.Contains(anime.Attributes.Slug, "delete") {
				continue
			}
			exists, err := database.IfItemExists(anime.ID, "anime_movies")
			if err != nil {
				log.Printf("Error checking if item exists: %v", err)
				continue
			}
			if exists {
				continue
			}
			genres, err := utils.FetchGenres(anime.Relationships.Genres.Links.Related)
			if err != nil {
				log.Printf("Failed to fetch genres for %s: %v", anime.Attributes.CanonicalTitle, err)
				continue
			}

			animeDoc := s.createAnimeMovieDoc(anime, genres)
			if err := database.InsertItem(animeDoc, "anime_movies"); err != nil {
				log.Printf("Failed to save anime movie %s to database: %v", animeDoc.Title, err)
			}
		}
	}

	log.Println("Fetching new anime movies completed")
	return nil
}

func (s *AnimeMovieScraper) createAnimeMovieDoc(anime models.Anime, genres []models.GenreMapping) *models.AnimeMovie {
	title := anime.Attributes.CanonicalTitle
	if title == "" {
		title = anime.Attributes.Titles.En
		if title == "" {
			title = anime.Attributes.Titles.EnUs
		} else {
			title = anime.Attributes.Titles.EnJp
		}
	}
	return &models.AnimeMovie{
		ID:            anime.ID,
		Attributes:    anime.Attributes,
		Relationships: anime.Relationships,
		Genres:        genres,
		Title:         title,
	}
}

*/
