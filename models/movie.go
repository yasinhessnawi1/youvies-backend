package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Movie struct {
	ID               primitive.ObjectID     `bson:"_id,omitempty"  json:"id,omitempty"`
	OriginalLanguage string                 `bson:"original_language" json:"original_language"`
	OriginalTitle    string                 `bson:"original_title" json:"original_title"`
	Overview         string                 `bson:"overview" json:"overview"`
	Popularity       float64                `bson:"popularity" json:"popularity"`
	PosterPath       string                 `bson:"poster_path" json:"poster_path"`
	ReleaseDate      string                 `bson:"release_date" json:"release_date"`
	Title            string                 `bson:"title" json:"title"`
	VoteAverage      float64                `bson:"vote_average" json:"vote_average"`
	VoteCount        int                    `bson:"vote_count" json:"vote_count"`
	BackdropPath     string                 `bson:"backdrop_path" json:"backdrop_path"`
	Adult            bool                   `bson:"adult" json:"adult"`
	Genres           []GenreMapping         `bson:"genres" json:"genres"`
	Torrents         map[string][]Torrent   `bson:"torrents" json:"torrents"`
	ExternalIDs      map[string]interface{} `bson:"external_ids" json:"external_ids"`
	LastUpdated      string                 `bson:"last_updated" json:"last_updated"` // Unix timestamp of last update
}

type OmdbMovie struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Runtime  string `json:"Runtime"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Language string `json:"Language"`
	Country  string `json:"Country"`
	Awards   string `json:"Awards"`
	Poster   string `json:"Poster"`
	Ratings  []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	ImdbID     string `json:"imdbID"`
	Type       string `json:"Type"`
	DVD        string `json:"DVD"`
	BoxOffice  string `json:"BoxOffice"`
	Production string `json:"Production"`
	Website    string `json:"Website"`
	Response   string `json:"Response"`
}
