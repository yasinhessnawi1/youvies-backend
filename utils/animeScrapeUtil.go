package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"youvies-backend/models"
)

// FetchAllEpisodes fetches all episodes for a given anime ID from Kitsu API
func FetchAllEpisodes(animeID string) ([]models.EpisodeInfo, error) {
	var allEpisodes []models.EpisodeInfo
	url := fmt.Sprintf("https://kitsu.io/api/edge/anime/%s/episodes", animeID)

	for url != "" {
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("got an error while fetching url %s: %v", url, err)
		}

		var episodesResponse models.EpisodeResponse
		if err := json.NewDecoder(resp.Body).Decode(&episodesResponse); err != nil {
			return nil, err
		}

		allEpisodes = append(allEpisodes, episodesResponse.Data...)
		url = episodesResponse.Links.Next
		err = resp.Body.Close()
		if err != nil {
			return nil, err
		}
	}
	return allEpisodes, nil
}

// FetchMissingTorrentsAnime fetches missing torrents for a list of episodes
func FetchMissingTorrentsAnime(title string, episodes []models.EpisodeInfo) ([]models.Torrent, error) {
	var missingTorrents []models.Torrent

	for _, episode := range episodes {
		query := fmt.Sprintf("%s S%02dE%02d", title, episode.Attributes.SeasonNumber, episode.Attributes.Number)
		torrents, err := FetchTorrents(query)
		if err != nil || len(torrents) == 0 {
			query = fmt.Sprintf("%s %02d", title, episode.Attributes.Number)
			torrents, err = FetchTorrents(query)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch torrents for %s: %v", query, err)
			}
		}
		missingTorrents = append(missingTorrents, torrents...)
	}
	return missingTorrents, nil
}
