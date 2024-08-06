package utils

import (
	"time"
	"youvies-backend/models"
)

type GenreAnime struct {
	ID         string          `json:"id"`
	Attributes GenreAttributes `json:"attributes"`
}

type GenreAttributes struct {
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
}
type GenreResponse struct {
	Data []GenreAnime `json:"data"`
}

func FetchGenres(url string) ([]models.GenreMapping, error) {
	var genreResponse GenreResponse
	err := FetchJSON(url, "", &genreResponse)
	if err != nil {
		return nil, err
	}

	var genres []models.GenreMapping
	if len(genreResponse.Data) == 0 || genreResponse.Data == nil {
		return genres, nil
	}
	// Initialize the genres slice with the correct length
	genres = make([]models.GenreMapping, len(genreResponse.Data))
	for i, genre := range genreResponse.Data {
		genres[i].Name = genre.Attributes.Name
		genres[i].ID = i
	}
	return genres, nil
}
