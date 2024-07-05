package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Show struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title        string             `bson:"title" json:"title"`
	Description  string             `bson:"description" json:"description"`
	Year         int                `bson:"year" json:"year"`
	Rating       float64            `bson:"rating" json:"rating"`
	ImageURL     string             `bson:"image_url" json:"image_url"`
	Language     string             `bson:"language" json:"language"`
	Networks     []string           `bson:"networks" json:"networks"`
	FirstAirDate string             `bson:"first_air_date" json:"first_air_date"`
	Episodes     []Torrent          `bson:"episodes" json:"episodes"`
	Country      []string           `json:"country"`
	Backdrop     string             `json:"backdrop"`
}

type Season struct {
	SeasonNumber int       `bson:"season_number" json:"season_number"`
	Episodes     []Episode `bson:"episodes" json:"episodes"`
}

type Episode struct {
	EpisodeNumber int       `bson:"episode_number" json:"episode_number"`
	Title         string    `bson:"title" json:"title"`
	Torrents      []Torrent `bson:"torrents" json:"torrents"`
}

type TMDBShow struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Overview     string   `json:"overview"`
	FirstAirDate string   `json:"first_air_date"`
	Genres       []string `json:"genres"`
}

type TVDBShow struct {
	SeriesID   string   `json:"id"`
	SeriesName string   `json:"seriesName"`
	Network    string   `json:"network"`
	Overview   string   `json:"overview"`
	FirstAired string   `json:"firstAired"`
	Genres     []string `json:"genres"`
}
