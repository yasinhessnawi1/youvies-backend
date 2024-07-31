package scraper

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"youvies-backend/database"
	"youvies-backend/models"
	"youvies-backend/utils"
)

type AnimeMovieScraper struct {
	BaseScraper
}

func NewAnimeMovieScraper() *AnimeMovieScraper {
	return &AnimeMovieScraper{
		BaseScraper: *NewBaseScraper("anime_movie", utils.KitsuBaseURL),
	}
}

// FetchAnimeMoviesFromKitsu fetches anime movie details from Kitsu with pagination
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

// Scrape fetches data from various APIs and inserts them into the database
func (s *AnimeMovieScraper) Scrape() error {
	animes, err := s.FetchAnimeMoviesFromKitsu()
	if err != nil {
		return fmt.Errorf("error fetching Kitsu anime movie data: %v", err)
	}
	for _, animeResp := range animes {
		for _, anime := range animeResp.Data {
			genres, err := utils.FetchGenres(anime.Relationships.Genres.Links.Related)
			if err != nil {
				log.Printf("Failed to fetch genres for %s: %v", anime.Attributes.CanonicalTitle, err)
				continue
			}

			animeDoc := s.createAnimeMovieDoc(anime, genres)

			exists, err := database.IfItemExists(map[string]interface{}{"title": animeDoc.Title}, "anime_movies")
			if err != nil {
				log.Fatalf("Error checking if item exists: %v", err)
			}
			if exists {
				continue
			}
			torrents, err := utils.FetchTorrents(animeDoc.Title, "anime movie")
			if err != nil || len(torrents) == 0 {
				log.Printf("error fetching torrents: %v", err)
				continue
			}
			categorizedTorrents := utils.CategorizeTorrentsByQuality(torrents)
			animeDoc.Torrents = categorizedTorrents

			if exists {
				var existingAnime models.AnimeMovie
				if err := database.FindItem(map[string]interface{}{"title": animeDoc.Title}, "anime_movies", &existingAnime); err != nil {
					log.Printf("Failed to fetch existing anime movie: %v", err)
					continue
				}
				if s.hasAnimeMovieChanged(existingAnime, animeDoc, categorizedTorrents) {
					if err := database.EditItem(map[string]interface{}{"title": animeDoc.Title}, animeDoc, "anime_movies"); err != nil {
						log.Printf("Failed to update anime movie %s in database: %v", animeDoc.Title, err)
					}
				}
			} else {
				if err := database.InsertItem(animeDoc, animeDoc.Title, "anime_movies"); err != nil {
					log.Printf("Failed to save anime movie %s to database: %v", animeDoc.Title, err)
				}
			}
		}
	}
	log.Println("Fetching new anime movies completed")
	return nil
}

// createAnimeMovieDoc constructs an anime movie document from Kitsu data.
func (s *AnimeMovieScraper) createAnimeMovieDoc(anime models.Anime, genres []string) models.AnimeMovie {
	title := anime.Attributes.Titles.En
	if title == "" {
		title = anime.Attributes.Titles.EnJp
		if title == "" {
			title = anime.Attributes.CanonicalTitle
		} else {
			title = anime.Attributes.Titles.EnUs
		}
	}
	return models.AnimeMovie{
		ID:            primitive.NewObjectID(),
		Attributes:    anime.Attributes,
		Relationships: anime.Relationships,
		Genres:        genres,
		Title:         title,
	}
}

// hasAnimeMovieChanged checks if the anime movie details or torrents have changed.
func (s *AnimeMovieScraper) hasAnimeMovieChanged(existingAnimeMovie models.AnimeMovie, newDetails models.AnimeMovie, newTorrents map[string][]models.Torrent) bool {
	return existingAnimeMovie.Attributes.Synopsis != newDetails.Attributes.Synopsis ||
		existingAnimeMovie.Attributes.StartDate != newDetails.Attributes.StartDate ||
		existingAnimeMovie.Attributes.AverageRating != newDetails.Attributes.AverageRating ||
		existingAnimeMovie.Attributes.PopularityRank != newDetails.Attributes.PopularityRank ||
		existingAnimeMovie.Attributes.PosterImage.Original != newDetails.Attributes.PosterImage.Original ||
		!compareTorrents(existingAnimeMovie.Torrents, newTorrents)
}
