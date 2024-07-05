package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type BaseScraper struct {
	Name    string
	BaseURL string
	APIKey  string
}

type Scraper interface {
	Scrape() error
}

func NewBaseScraper(name, baseURL string) *BaseScraper {
	return &BaseScraper{
		Name:    name,
		BaseURL: baseURL,
		APIKey:  os.Getenv(fmt.Sprintf("%s_KEY", strings.ToUpper(name))),
	}
}

func (s *BaseScraper) FetchURL(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	return resp, nil
}

func (s *BaseScraper) FetchJSON(url string, target interface{}) error {
	resp, err := s.FetchURL(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	if err := json.Unmarshal(bodyBytes, target); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	return nil
}
