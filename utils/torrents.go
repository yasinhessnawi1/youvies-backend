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

	var result models.TorrentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	adultKeywords := []string{
		"adult", "explicit", "porn", "pornographic", "xxx", "sex", "sexual", "nudity", "nude",
		"erotic", "hardcore", "softcore", "fetish", "bdsm", "incest", "bestiality", "rape",
		"molestation", "pedophilia", "child porn", "child pornography", "underage", "prostitution",
		"escort", "hooker", "whore", "slut", "strip", "stripper", "peep show", "voyeur",
		"masturbation", "masturbate", "orgasm", "orgy", "gangbang", "cum", "anal", "oral", "blowjob",
		"handjob", "penetration", "climax", "ejaculation", "bondage", "submission", "dominance",
		"kink", "taboo", "swinger", "swinging", "adult toys", "sex toys", "lube", "lubricant",
		"voyeurism", "exhibitionism", "smut", "brothel", "red light district", "x-rated", "raunchy",
		"strip club", "naked", "bare", "risque", "provocative", "suggestive",
	}

	for _, torrent := range result.Data {
		if !containsAny(torrent.Category, adultKeywords) && !containsAny(torrent.Name, adultKeywords) {
			torrents = append(torrents, torrent)
		}
	}

	fmt.Printf("Found %d torrents for %s: %s \n", len(torrents), title, url)
	if len(torrents) == 0 {
		return nil, fmt.Errorf("no torrents found for title: %s", title)
	}

	return torrents, nil
}

func containsAny(text string, items []string) bool {
	text = strings.TrimSpace(strings.ToLower(text))
	for _, item := range items {
		item = strings.TrimSpace(strings.ToLower(item))
		if strings.Contains(text, item) {
			return true
		}
	}
	return false
}
