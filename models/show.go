package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Show struct {
	ID                  primitive.ObjectID     `bson:"_id,omitempty"  json:"id,omitempty"`
	Title               string                 `bson:"title" json:"title"`
	Overview            string                 `bson:"overview" json:"overview"`
	PosterPath          string                 `bson:"image_url" json:"image_url"`
	Networks            []string               `bson:"networks" json:"networks"`
	FirstAirDate        string                 `bson:"first_air_date" json:"first_air_date"`
	Country             []string               `bson:"country" json:"country"`
	Seasons             map[int]Season         `bson:"seasons" json:"seasons"`
	Genres              []GenreMapping         `bson:"genres" json:"genres"`
	VoteAverage         float64                `json:"vote_average"`
	VoteCount           int                    `json:"vote_count"`
	OriginalLanguage    string                 `json:"original_language"`
	Popularity          float64                `json:"popularity"`
	BackdropPath        string                 `json:"backdrop_path"`
	ExternalIDs         map[string]interface{} `bson:"external_ids" json:"external_ids"`
	ProductionCompanies []string               `json:"production_companies"`
	ProductionCountries []string               `json:"production_countries"`
	SpokenLanguages     []string               `json:"spoken_languages"`
	LastUpdated         string                 `bson:"last_updated" json:"last_updated"` // Unix timestamp of last update
	SeasonsInfo         []SeasonInfo           `bson:"seasons_info" json:"seasons_info"`
	OtherTorrents       map[string][]Torrent   `bson:"other_torrents" json:"other_torrents"`
}
type SeasonInfo struct {
	SeasonNumber int    `json:"season_number"`
	EpisodeCount int    `json:"episode_count"`
	AirDate      string `json:"air_date"`
	PosterPath   string `json:"poster_path"`
}
type GenreMapping struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Episode struct {
	Torrents map[string][]Torrent `bson:"torrents" json:"torrents"`
}

type Season struct {
	Episodes map[int]Episode `bson:"episodes" json:"episodes"`
}
