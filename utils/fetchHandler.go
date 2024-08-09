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
	const (
		workerCount = 10 // Adjust based on your server and network capabilities
		pageLimit   = 20 // Number of items per page
	)
	var allAnime []models.Anime
	var wg sync.WaitGroup
	urls := make(chan int)
	results := make(chan []models.Anime, workerCount)
	errors := make(chan error, workerCount)

	// Worker function
	worker := func() {
		defer wg.Done()
		for offset := range urls {
			url := fmt.Sprintf("%s?sort=updated_at&page[limit]=%d&page[offset]=%d", baseURL, pageLimit, offset)
			resp, err := http.Get(url)
			if err != nil {
				errors <- fmt.Errorf("error fetching URL %s: %v", url, err)
				continue
			}

			var animeResponse models.AnimeResponse
			if err := json.NewDecoder(resp.Body).Decode(&animeResponse); err != nil {
				errors <- fmt.Errorf("error decoding response body: %v", err)
				resp.Body.Close()
				continue
			}
			resp.Body.Close()

			// If no data is returned, stop further processing for this worker
			if len(animeResponse.Data) == 0 {
				return
			}

			var filteredAnime []models.Anime
			for _, anime := range animeResponse.Data {
				if anime.Attributes.Subtype == "movie" || anime.Attributes.Subtype == "TV" ||
					anime.Attributes.Subtype == "ONA" || anime.Attributes.Subtype == "OVA" {
					filteredAnime = append(filteredAnime, anime)
				}
			}

			results <- filteredAnime

			// Generate the next offset for this worker
			urls <- offset + pageLimit
		}
	}

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker()
	}

	// Send the initial offset to start the process
	go func() {
		urls <- 1 // Start with the first offset
	}()

	// Close the results and errors channels once all workers are done
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results and handle errors concurrently
	go func() {
		for result := range results {
			allAnime = append(allAnime, result...)
			if len(allAnime)%10 == 0 {
				fmt.Printf("%d=>", len(allAnime))
			}
		}
	}()

	for err := range errors {
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("Fetched %d anime by updated_at\n", len(allAnime))
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
