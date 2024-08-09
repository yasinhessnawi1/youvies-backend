package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"youvies-backend/models"
)

func FetchURLWithKey(url, apiKey string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}
func FetchURL(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func FetchJSON(url, apiKey string, target interface{}) error {
	var resp *http.Response
	var err error
	if apiKey != "" {
		resp, err = FetchURLWithKey(url, apiKey)
	} else {
		resp, err = FetchURL(url)
	}
	if err != nil {
		log.Printf("Error fetching URL: %v %v\n", err, url)
		return nil
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("error decoding JSON for %v: %v", url, err)
	}

	return nil
}

// FetchAllEpisodes fetches all episodes for a given anime ID from Kitsu API and converts them to Episode objects
func FetchAllEpisodes(animeID string) ([]models.Episode, error) {
	var allEpisodes []models.Episode
	url := fmt.Sprintf("https://kitsu.io/api/edge/anime/%s/episodes", animeID)

	for url != "" {
		var episodesResponse models.EpisodeResponse
		err := FetchJSON(url, "", &episodesResponse)
		if err != nil {
			return nil, fmt.Errorf("got an error while fetching url %s: %v", url, err)
		}
		if len(episodesResponse.Data) == 0 {
			continue
		}
		for _, episodeInfo := range episodesResponse.Data {
			episode := models.Episode{
				ID:             episodeInfo.ID,
				AnimeShowID:    animeID,
				CreatedAt:      episodeInfo.Attributes.CreatedAt,
				UpdatedAt:      episodeInfo.Attributes.UpdatedAt,
				Synopsis:       episodeInfo.Attributes.Synopsis,
				Description:    episodeInfo.Attributes.Description,
				Titles:         episodeInfo.Attributes.Titles,
				CanonicalTitle: episodeInfo.Attributes.CanonicalTitle,
				SeasonNumber:   episodeInfo.Attributes.SeasonNumber,
				Number:         episodeInfo.Attributes.Number,
				RelativeNumber: episodeInfo.Attributes.RelativeNumber,
				Airdate:        episodeInfo.Attributes.Airdate,
				Length:         episodeInfo.Attributes.Length,
				Thumbnail:      episodeInfo.Attributes.Thumbnail,
			}
			allEpisodes = append(allEpisodes, episode)
		}
		url = episodesResponse.Links.Next
	}

	return allEpisodes, nil
}

// FetchSortedAnimeByUpdatedAt fetches anime by updated_at timestamp using concurrent URL fetching without total items
func FetchSortedAnimeByUpdatedAt(baseURL string) ([]models.Anime, error) {
	var allAnime []models.Anime
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // Limit the number of concurrent goroutines

	offset := 1
	for {
		wg.Add(1)
		semaphore <- struct{}{}
		url := fmt.Sprintf("%s?filter[subtype]=TV,movie,OVA,ONA&page[limit]=20&page[offset]=%d", baseURL, offset*20)

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

			var animeResponse models.AnimeResponse
			if err := json.NewDecoder(resp.Body).Decode(&animeResponse); err != nil {
				log.Printf("error decoding response body: %v", err)
				return
			}
			resp.Body.Close()

			// If no data is returned, stop further processing
			if len(animeResponse.Data) == 0 {
				return
			}
			allAnime = append(allAnime, animeResponse.Data...)
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
		offset++
		if len(allAnime)%200 == 0 {
			fmt.Print(len(allAnime), "=>")
		}
	}

	wg.Wait()
	fmt.Printf("found this many animes: %d\n", len(allAnime))
	return allAnime, nil

}

func RemoveDuplicateStrings(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	result := make([]string, 0, len(input))
	for _, val := range input {
		if _, ok := seen[val]; !ok {
			seen[val] = struct{}{}
			result = append(result, val)
		}
	}
	return result
}
