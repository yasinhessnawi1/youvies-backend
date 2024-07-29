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

type AnimeShowScraper struct {
	BaseScraper
}

func NewAnimeShowScraper() *AnimeShowScraper {
	return &AnimeShowScraper{
		BaseScraper: *NewBaseScraper("anime_show", utils.KitsuBaseURL),
	}
}

// FetchAnimeDetailsFromKitsu fetches anime details from Kitsu with pagination
func (s *AnimeShowScraper) FetchAnimeDetailsFromKitsu() ([]models.AnimeResponse, error) {
	var allAnimes []models.AnimeResponse
	url := fmt.Sprintf("%s", s.BaseScraper.BaseURL)

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
	fmt.Printf("found this many animes: %d\n", len(allAnimes))

	return allAnimes, nil
}

// Scrape fetches data from various APIs and inserts them into the database
func (s *AnimeShowScraper) Scrape() error {
	animes, err := s.FetchAnimeDetailsFromKitsu()
	if err != nil {
		return fmt.Errorf("error fetching Kitsu anime data: %v", err)
	}

	for _, animeResp := range animes {
		for _, anime := range animeResp.Data {
			// Fetch episodes
			episodes, err := utils.FetchAllEpisodes(anime.Id)
			if err != nil {
				log.Printf("Failed to fetch episodes for anime %s: %v", anime.Attributes.CanonicalTitle, err)
				continue
			}
			anime.Attributes.EpisodeCount = len(episodes)

			genres, err := utils.FetchGenres(anime.Relationships.Genres.Links.Related)
			if err != nil {
				log.Printf("Failed to fetch genres for %s: %v", anime.Attributes.CanonicalTitle, err)
				continue
			}

			animeDoc := s.createAnimeShowDoc(anime, genres)
			animeDoc.Episodes = episodes
			exists, err := database.IfItemExists(map[string]interface{}{"title": animeDoc.Title}, "anime_shows")
			if err != nil {
				log.Fatalf("Error checking if item exists: %v", err)
			}

			if exists && animeDoc.Title == "" {
				log.Printf("Anime %s already exists in database", animeDoc.Title)
				continue
			}
			torrents, err := utils.FetchTorrents(animeDoc.Attributes.CanonicalTitle)
			if err != nil {
				log.Printf("error fetching torrents: %v", err)
			}

			categorizedTorrents, fullContent := utils.CategorizeAnimeTorrentsBySeasonsAndEpisodes(torrents, anime.Attributes.EpisodeCount)
			animeDoc.Seasons = categorizedTorrents
			animeDoc.FullContent = fullContent
			for _, torrent := range categorizedTorrents {
				for _, episode := range torrent.Episodes {
					for _, torrent := range episode.Torrents {
						for _, torrent := range torrent {
							err := utils.SaveMetadata(torrent.Magnet, torrent.Name)
							if err != nil {
								log.Printf("Failed to save torrent metadata for %s: %v", torrent.Name, err)
								continue
							}
						}

					}
				}
			}
			for _, content := range fullContent {
				for _, torrent := range content {
					err := utils.SaveMetadata(torrent.Magnet, torrent.Name)
					if err != nil {
						log.Printf("Failed to save torrent metadata for %s: %v", torrent.Name, err)
						continue
					}
				}
			}
			missingEpisodes := s.checkForMissingEpisodes(animeDoc)
			if len(missingEpisodes) > 0 {
				missingTorrents, err := utils.FetchMissingTorrentsAnime(animeDoc.Title, missingEpisodes)
				if err != nil {
					log.Printf("error fetching missing torrents: %v", err)
				}
				for _, torrent := range missingTorrents {
					err := utils.SaveMetadata(torrent.Magnet, torrent.Name)
					if err != nil {
						log.Printf("Failed to save torrent metadata for %s: %v", torrent.Name, err)
						continue
					}
					if torrent.Name != "" {
						episodeNum := utils.GetEpisodeNumberFromTorrentName(torrent.Name)
						quality := utils.ExtractQuality(torrent.Name)
						if episode, ok := animeDoc.Seasons[1].Episodes[episodeNum]; ok {
							episode.Torrents[quality] = append(episode.Torrents[quality], torrent)
						} else {
							animeDoc.Seasons[1].Episodes[episodeNum] = models.Episode{
								Torrents: map[string][]models.Torrent{
									quality: {torrent},
								},
							}
						}
					}
				}
			}

			if exists {
				var existingAnime models.AnimeShow
				if err := database.FindItem(map[string]interface{}{"title": animeDoc.Title}, "anime_shows", &existingAnime); err != nil {
					log.Printf("Failed to fetch existing anime show: %v", err)
					continue
				}
				if s.hasAnimeShowChanged(existingAnime, animeDoc, categorizedTorrents) {
					if err := database.EditItem(map[string]interface{}{"title": animeDoc.Title}, animeDoc, "anime_shows"); err != nil {
						log.Printf("Failed to update anime show %s in database: %v", animeDoc.Title, err)
					}
				}
			} else {
				if err := database.InsertItem(animeDoc, animeDoc.Title, "anime_shows"); err != nil {
					log.Printf("Failed to save anime show %s to database: %v", animeDoc.Title, err)
				}
			}
		}
	}
	log.Println("Fetching new anime shows completed")
	return nil
}

// checkForMissingEpisodes checks if any episodes are missing torrents and returns a list of missing episodes
func (s *AnimeShowScraper) checkForMissingEpisodes(animeDoc models.AnimeShow) []models.EpisodeInfo {
	var missingEpisodes []models.EpisodeInfo
	for _, episode := range animeDoc.Episodes {
		episodeNum := episode.Attributes.Number
		if _, ok := animeDoc.Seasons[1].Episodes[episodeNum]; !ok {
			missingEpisodes = append(missingEpisodes, episode)
		}
	}
	return missingEpisodes
}

// createAnimeShowDoc constructs an anime show document from Kitsu data.
func (s *AnimeShowScraper) createAnimeShowDoc(anime models.Anime, genres []string) models.AnimeShow {
	title := anime.Attributes.CanonicalTitle
	if title == "" {
		title = anime.Attributes.Titles.EnJp
		if title == "" {
			title = anime.Attributes.Titles.En
		}
	}
	return models.AnimeShow{
		ID:            primitive.NewObjectID(),
		Attributes:    anime.Attributes,
		Relationships: anime.Relationships,
		Genres:        genres,
		Title:         title, // Including the title attribute
	}
}

// hasAnimeShowChanged checks if the anime show details or torrents have changed.
func (s *AnimeShowScraper) hasAnimeShowChanged(existingAnimeShow models.AnimeShow, newDetails models.AnimeShow, newTorrents map[int]models.Season) bool {
	return existingAnimeShow.Attributes.Synopsis != newDetails.Attributes.Synopsis ||
		existingAnimeShow.Attributes.StartDate != newDetails.Attributes.StartDate ||
		existingAnimeShow.Attributes.AverageRating != newDetails.Attributes.AverageRating ||
		existingAnimeShow.Attributes.PopularityRank != newDetails.Attributes.PopularityRank ||
		existingAnimeShow.Attributes.PosterImage.Original != newDetails.Attributes.PosterImage.Original ||
		!compareTorrentsBySeason(existingAnimeShow.Seasons, newTorrents)
}
