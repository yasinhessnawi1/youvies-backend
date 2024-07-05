package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func FetchURL(url, apiKey string) (*http.Response, error) {
	client := &http.Client{Timeout: DefaultTimeout * time.Second}
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

func FetchJSON(url, apiKey string, target interface{}) error {
	resp, err := FetchURL(url, apiKey)
	if err != nil {
		log.Printf("Error fetching URL: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	return nil
}

func GetYear(dateStr string) int {
	if len(dateStr) >= 4 {
		var year int
		fmt.Sscanf(dateStr, "%4d", &year)
		return year
	}
	return 0
}
