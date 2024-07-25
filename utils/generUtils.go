package utils

import (
	"encoding/json"
	"net/http"
)

type GenreAnime struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name string `json:"name"`
	} `json:"attributes"`
}

type GenreResponse struct {
	Data []GenreAnime `json:"data"`
}

func FetchGenres(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var genreResponse GenreResponse
	if err := json.NewDecoder(resp.Body).Decode(&genreResponse); err != nil {
		return nil, err
	}

	var genres []string
	for _, genre := range genreResponse.Data {
		genres = append(genres, genre.Attributes.Name)
	}
	return genres, nil
}
