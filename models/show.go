package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Show struct {
	ID                  primitive.ObjectID     `bson:"_id,omitempty"  json:"id,omitempty"`
	Title               string                 `bson:"title" json:"title"`
	Overview            string                 `bson:"overview" json:"overview"`
	PosterPath          string                 `bson:"image_url" json:"image_url"`
	Networks            []string               `bson:"networks,omitempty" json:"networks,omitempty"`
	FirstAirDate        string                 `bson:"first_air_date" json:"first_air_date"`
	Country             []string               `bson:"country,omitempty" json:"country,omitempty"`
	Seasons             map[int]Season         `bson:"seasons,omitempty" json:"seasons,omitempty"`
	Genres              []GenreMapping         `bson:"genres" json:"genres"`
	VoteAverage         float64                `bson:"vote_average" json:"vote_average"`
	VoteCount           int                    `bson:"vote_count" json:"vote_count"`
	OriginalLanguage    string                 `bson:"original_language,omitempty " json:"original_language,omitempty"`
	Popularity          float64                `bson:"popularity" json:"popularity"`
	BackdropPath        string                 `bson:"backdrop_path" json:"backdrop_path"`
	ExternalIDs         map[string]interface{} `bson:"external_i_ds,omitempty" json:"external_ids,omitempty"`
	ProductionCompanies []string               `bson:"production_companies,omitempty" json:"production_companies,omitempty"`
	ProductionCountries []string               `bson:"production_countries,omitempty" json:"production_countries,omitempty"`
	SpokenLanguages     []string               `bson:"spoken_languages,omitempty" json:"spoken_languages,omitempty"`
	LastUpdated         string                 `bson:"last_updated,omitempty" json:"last_updated,omitempty"` // Unix timestamp of last update
	SeasonsInfo         []SeasonInfo           `bson:"seasons_info,omitempty" json:"seasons_info,omitempty"`
	OtherTorrents       map[string][]Torrent   `bson:"other_torrents,omitempty" json:"other_torrents,omitempty"`
}
type SeasonInfo struct {
	SeasonNumber int    `json:"season_number"`
	EpisodeCount int    `json:"episode_count"`
	AirDate      string `json:"air_date,omitempty"`
	PosterPath   string `json:"poster_path,omitempty"`
}
type GenreMapping struct {
	ID   int    `bson:",omitempty" json:"id,omitempty"`
	Name string `bson:"name" json:"name"`
}

type Episode struct {
	Torrents map[string][]Torrent `bson:"torrents" json:"torrents"`
}

type Season struct {
	Episodes map[int]Episode `bson:"episodes" json:"episodes"`
}
