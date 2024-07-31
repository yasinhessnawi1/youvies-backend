package scraper

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"sort"
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

// FetchAnimeDetailsFromKitsu fetches anime details from Kitsu with pagination
func (s *AnimeShowScraper) FetchAnimeDetailsFromKitsu() ([]models.AnimeResponse, error) {
	var allAnime []models.AnimeResponse
	var wg sync.WaitGroup
	var mu sync.Mutex
	semaphore := make(chan struct{}, 10) // Limit the number of concurrent goroutines

	pageNum := 1 // Start from page 1150 for debug
	for {
		wg.Add(1)
		semaphore <- struct{}{}
		url := fmt.Sprintf("%s?filter[subtype]=TV&page[limit]=20&page[offset]=%d", s.BaseScraper.BaseURL, pageNum*20)

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

			mu.Lock()
			fmt.Printf("fetched %d animes from %v\n", len(animes.Data), url)
			allAnime = append(allAnime, animes)
			mu.Unlock()
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
	}
	wg.Wait()
	fmt.Printf("found this many animes: %d\n", len(allAnime))

	sort.Slice(allAnime, func(i, j int) bool {
		return allAnime[i].Data[0].Attributes.CanonicalTitle < allAnime[j].Data[0].Attributes.CanonicalTitle
	})

	return allAnime, nil
}

func (s *AnimeShowScraper) Scrape() error {
	animes, err := s.FetchAnimeDetailsFromKitsu()
	if err != nil {
		return fmt.Errorf("error fetching Kitsu anime data: %v", err)
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit the number of concurrent goroutines

	for _, animeResp := range animes {
		for _, anime := range animeResp.Data {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(anime models.Anime) {
				defer wg.Done()
				defer func() { <-semaphore }()

				// Fetch episodes
				episodes, err := utils.FetchAllEpisodes(anime.Id)
				if err != nil {
					log.Printf("Failed to fetch episodes for anime %s: %v", anime.Attributes.CanonicalTitle, err)
					return
				}
				anime.Attributes.EpisodeCount = len(episodes)

				genres, err := utils.FetchGenres(anime.Relationships.Genres.Links.Related)
				if err != nil {
					log.Printf("Failed to fetch genres for %s: %v", anime.Attributes.CanonicalTitle, err)
					return
				}

				animeDoc := s.createAnimeShowDoc(anime, genres)
				animeDoc.Episodes = episodes
				exists, err := database.IfItemExists(map[string]interface{}{"title": animeDoc.Title}, "anime_shows")
				if err != nil {
					log.Fatalf("Error checking if item exists: %v", err)
				}

				if exists && animeDoc.Title == "" {
					log.Printf("Anime %s already exists in database", animeDoc.Title)
					return
				}
				fmt.Printf("Fetching torrents for %s\n", animeDoc.Title)
				torrents, err := utils.FetchTorrents(animeDoc.Title, "anime show")
				if err != nil || len(torrents) == 0 {
					log.Printf("error fetching torrents for Anime %s: %v", animeDoc.Title, err)
					return
				}

				categorizedTorrents, fullContent := utils.CategorizeTorrentsBySeasonsAndEpisodes(torrents)
				animeDoc.Seasons = categorizedTorrents
				animeDoc.FullContent = fullContent
				missingEpisodes := s.checkForMissingEpisodes(animeDoc)
				if len(missingEpisodes) > 0 {
					missingTorrents, err := s.fetchMissingTorrentsForAnime(animeDoc.Title, missingEpisodes)
					if err != nil {
						log.Printf("error fetching missing torrents: %v", err)
					}
					s.addMissingTorrentsToAnimeDoc(animeDoc, missingTorrents)
				}

				if exists {
					var existingAnime models.AnimeShow
					if err := database.FindItem(map[string]interface{}{"title": animeDoc.Title}, "anime_shows", &existingAnime); err != nil {
						log.Printf("Failed to fetch existing anime show: %v", err)
						return
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
			}(anime)
		}
	}

	wg.Wait()
	log.Println("Fetching new anime shows completed")
	return nil
}

// checkForMissingEpisodes checks if any episodes are missing torrents and returns a list of missing episodes
func (s *AnimeShowScraper) checkForMissingEpisodes(animeDoc models.AnimeShow) []models.EpisodeInfo {
	var missingEpisodes []models.EpisodeInfo
	for _, episode := range animeDoc.Episodes {
		seasonNum := episode.Attributes.SeasonNumber
		episodeNum := episode.Attributes.Number
		if _, ok := animeDoc.Seasons[seasonNum].Episodes[episodeNum]; !ok {
			missingEpisodes = append(missingEpisodes, episode)
		}
	}
	return missingEpisodes
}

// fetchMissingTorrentsForAnime fetches missing torrents for the given episodes of an anime
func (s *AnimeShowScraper) fetchMissingTorrentsForAnime(title string, episodes []models.EpisodeInfo) ([]models.Torrent, error) {
	var missingTorrents []models.Torrent

	for _, episode := range episodes {
		query := fmt.Sprintf("%s S%02dE%02d", title, episode.Attributes.SeasonNumber, episode.Attributes.Number)
		torrents, err := utils.FetchTorrents(query, "anime show")
		if err != nil || len(torrents) == 0 {
			query = fmt.Sprintf("%s %02d", title, episode.Attributes.Number)
			torrents, err = utils.FetchTorrents(query, "anime show")
			if err != nil {
				log.Printf("Failed to fetch torrents for %s: %v", query, err)
				continue
			}
		}
		missingTorrents = append(missingTorrents, torrents...)
	}
	return missingTorrents, nil
}

// addMissingTorrentsToAnimeDoc adds the missing torrents to the anime document
func (s *AnimeShowScraper) addMissingTorrentsToAnimeDoc(animeDoc models.AnimeShow, missingTorrents []models.Torrent) {
	for _, torrent := range missingTorrents {
		if torrent.Name != "" {
			seasonNum, episodeNum, err := utils.ExtractSeasonAndEpisode(torrent.Name)
			if err != nil {
				log.Printf("Failed to extract season and episode from torrent name: %v", err)
				continue
			}
			quality := utils.ExtractQuality(torrent.Name)
			if season, ok := animeDoc.Seasons[seasonNum]; ok {
				if episode, ok := season.Episodes[episodeNum]; ok {
					episode.Torrents[quality] = append(episode.Torrents[quality], torrent)
				} else {
					season.Episodes[episodeNum] = models.Episode{
						Torrents: map[string][]models.Torrent{
							quality: {torrent},
						},
					}
				}
			} else {
				animeDoc.Seasons[seasonNum] = models.Season{
					Episodes: map[int]models.Episode{
						episodeNum: {
							Torrents: map[string][]models.Torrent{
								quality: {torrent},
							},
						},
					},
				}
			}
		}
	}
}

// createAnimeShowDoc constructs an anime show document from Kitsu data.
func (s *AnimeShowScraper) createAnimeShowDoc(anime models.Anime, genres []string) models.AnimeShow {
	title := anime.Attributes.CanonicalTitle
	if title == "" {
		title = anime.Attributes.Titles.En
		if title == "" {
			title = anime.Attributes.Titles.EnUs
		} else {
			title = anime.Attributes.Titles.EnJp
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
