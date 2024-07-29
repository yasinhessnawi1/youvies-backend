package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"youvies-backend/models"
)

func FetchTorrents(title string) ([]models.Torrent, error) {
	var torrents []models.Torrent

	url := fmt.Sprintf("%ssearch?query=%s", TorrentAPIBaseURL, strings.ReplaceAll(title, " ", "%20"))

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching torrents: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(fmt.Errorf("error closing request body: %s", err))
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	var result models.TorrentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	var movieCategories = []string{
		"movie", "movies", "film", "cinema", "blockbuster", "feature film",
		"motion picture", "flick", "biopic", "documentary", "short film",
		"thriller", "comedy", "drama", "action", "adventure",
		"animation", "crime", "fantasy", "historical", "horror",
		"musical", "mystery", "romance", "sci-fi", "science fiction",
		"war", "western", "independent film", "indie film", "art house",
		"silent film", "noir", "cult film", "Video > Movies", "Video > HD - movies",
	}
	var showCategories = []string{
		"show", "shows", "tv show", "tv shows", "television show", "series", "tv series",
		"sitcom", "reality show", "talk show", "drama series", "comedy series",
		"mini-series", "soap opera", "docuseries", "children's show",
		"news show", "variety show", "game show", "late-night show",
		"cooking show", "competition show", "talent show", "true crime",
		"crime drama", "fantasy series", "sci-fi series", "science fiction series",
		"historical drama", "superhero series", "animated series", "anime series",
		"documentary series", "medical drama", "legal drama", "reality competition", "Video > HD - TV shows",
	}
	var animeCategories = []string{
		"anime", "manga", "ova", "ona", "anime series", "anime movie",
		"light novel", "hentai", "josei", "seinen", "shonen", "shojo",
		"yaoi", "yuri", "anime film", "isekai", "mecha", "slice of life",
		"shoujo-ai", "shounen-ai", "magical girl", "sports anime", "supernatural",
		"fantasy anime", "sci-fi anime", "science fiction anime", "romance anime",
		"action anime", "adventure anime", "comedy anime", "drama anime",
		"historical anime", "horror anime", "music anime", "psychological anime",
		"school anime", "space anime", "thriller anime", "military anime",
	}

	for _, torrent := range result.Data {
		category := strings.ToLower(torrent.Category)

		if containsAny(category, movieCategories) || containsAny(category, showCategories) || containsAny(category, animeCategories) {
			if torrent.Magnet != "" || torrent.Hash != "" {
				torrents = append(torrents, torrent)
			}
		}
	}
	if len(torrents) == 0 {
		return nil, fmt.Errorf("no torrents found for title: %s", title)
	}
	fmt.Printf("Found %d torrents for %s: %s \n", len(torrents), torrents[0].Category+title, url)

	return torrents, nil
}

func containsAny(text string, items []string) bool {
	text = strings.TrimSpace(strings.ToLower(text))
	for _, item := range items {
		item = strings.TrimSpace(strings.ToLower(item))
		if strings.Contains(text, item) && !strings.Contains(text, "porn") {
			return true
		}
	}
	return false
}
