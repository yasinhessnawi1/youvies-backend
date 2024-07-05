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

func FetchTorrents(title string, categories []string) ([]models.Torrent, error) {
	var torrents []models.Torrent

	fmt.Printf("Fetching torrents for %s: %s\n", categories[0], title)
	url := fmt.Sprintf("%ssearch?query=%s", TorrentAPIBaseURL, strings.ReplaceAll(title, " ", "%20"))

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching torrents:  %s ", err)
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

	var result models.TorrentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}
	for _, torrent := range result.Data {
		if !contains(categories, torrent.Category) {
			continue
		}

		torrents = append(torrents, torrent)
	}
	if len(torrents) == 0 {
		return nil, fmt.Errorf("no torrents found for title: %s", title)
	}

	return torrents, nil
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
