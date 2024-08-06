package models

import "time"

type AnimeMovie struct {
	ID            string         `json:"id,omitempty"`
	Title         string         `json:"title"`
	Attributes    Attributes     `json:"attributes"`
	Relationships Relationships  `json:"relationships,omitempty"`
	Genres        []GenreMapping `json:"genres"`
}

type AnimeShow struct {
	ID            string         `json:"id,omitempty"`
	Title         string         `json:"title"`
	Attributes    Attributes     `json:"attributes"`
	Relationships Relationships  `json:"relationships,omitempty"`
	Genres        []GenreMapping `json:"genres"`
}

type Attributes struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Slug           string    `json:"slug"`
	Synopsis       string    `json:"synopsis"`
	Description    string    `json:"description"`
	Titles         Titles    `json:"titles"`
	CanonicalTitle string    `json:"canonicalTitle"`
	AverageRating  string    `json:"averageRating"`
	FavoritesCount int       `json:"favoritesCount"`
	StartDate      string    `json:"startDate"`
	EndDate        string    `json:"endDate"`
	NextRelease    string    `json:"nextRelease"`
	PopularityRank int       `json:"popularityRank"`
	RatingRank     int       `json:"ratingRank"`
	AgeRating      string    `json:"ageRating"`
	AgeRatingGuide string    `json:"ageRatingGuide"`
	Subtype        string    `json:"subtype"`
	Status         string    `json:"status"`
	PosterImage    Image     `json:"posterImage"`
	CoverImage     Image     `json:"coverImage"`
	EpisodeCount   int       `json:"episodeCount"`
	EpisodeLength  int       `json:"episodeLength"`
	YoutubeVideoId string    `json:"youtubeVideoId"`
	ShowType       string    `json:"showType"`
}

type Titles struct {
	En   string `json:"en"`
	EnJp string `json:"en_jp"`
	EnUs string `json:"en_us"`
	JaJp string `json:"ja_jp"`
}

type Image struct {
	Tiny     string `json:"tiny"`
	Large    string `json:"large"`
	Small    string `json:"small"`
	Medium   string `json:"medium"`
	Original string `json:"original"`
}

type Relationships struct {
	ID     int64  `json:"id"`
	Genres Genres `json:"genres"`
}

type Genres struct {
	Links Link `json:"links"`
}
type Link struct {
	Self    string `json:"self"`
	Related string `json:"related"`
}

type Anime struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Links         Link          `json:"links"`
	Attributes    Attributes    `json:"attributes"`
	Relationships Relationships `json:"relationships"`
}

type AnimeResponse struct {
	Data  []Anime `json:"data,omitempty"`
	Meta  Meta    `json:"meta,omitempty"`
	Links Links   `json:"links"`
}

type Meta struct {
	Count int `json:"count,omitempty"`
}

type Links struct {
	First string `json:"first,omitempty"`
	Next  string `json:"next"`
	Last  string `json:"last,omitempty"`
}

type EpisodeResponse struct {
	Data  []EpisodeInfo `json:"data"`
	Meta  Meta          `json:"meta"`
	Links Links         `json:"links"`
}

type EpisodeInfo struct {
	ID         string            `json:"id"`
	Attributes EpisodeAttributes `json:"attributes"`
}

type EpisodeAttributes struct {
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Synopsis       string    `json:"synopsis"`
	Description    string    `json:"description"`
	Titles         Titles    `json:"titles"`
	CanonicalTitle string    `json:"canonicalTitle"`
	SeasonNumber   int       `json:"seasonNumber"`
	Number         int       `json:"number"`
	RelativeNumber int       `json:"relativeNumber"`
	Airdate        string    `json:"airdate"`
	Length         int       `json:"length"`
	Thumbnail      Image     `json:"thumbnail"`
}

type Movie struct {
	ID               string         `json:"id,omitempty"`
	OriginalLanguage string         `json:"original_language,omitempty"`
	OriginalTitle    string         `json:"original_title"`
	Overview         string         `json:"overview"`
	Popularity       float64        `json:"popularity"`
	PosterPath       string         `json:"poster_path"`
	ReleaseDate      string         `json:"release_date"`
	Title            string         `json:"title"`
	VoteAverage      float64        `json:"vote_average"`
	VoteCount        int            `json:"vote_count"`
	BackdropPath     string         `json:"backdrop_path"`
	Genres           []GenreMapping `json:"genres"`
	LastUpdated      string         `json:"last_updated,omitempty"`
}

type Show struct {
	ID               string         `json:"id,omitempty"`
	Title            string         `json:"title"`
	Overview         string         `json:"overview"`
	PosterPath       string         `json:"image_url"`
	FirstAirDate     string         `json:"first_air_date"`
	Genres           []GenreMapping `json:"genres"`
	VoteAverage      float64        `json:"vote_average"`
	VoteCount        int            `json:"vote_count"`
	OriginalLanguage string         `json:"original_language,omitempty"`
	Popularity       float64        `json:"popularity"`
	BackdropPath     string         `json:"backdrop_path"`
	LastUpdated      string         `json:"last_updated,omitempty"`
	SeasonsInfo      []SeasonInfo   `json:"seasons_info,omitempty"`
}

type SeasonInfo struct {
	SeasonNumber int    `json:"season_number"`
	EpisodeCount int    `json:"episode_count"`
	AirDate      string `json:"air_date,omitempty"`
	PosterPath   string `json:"poster_path,omitempty"`
}

type GenreMapping struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

type User struct {
	ID        string    `json:"id,omitempty"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Avatar    string    `json:"avatar"`
	Favorites []string  `json:"favorites"`
	Friends   []string  `json:"friends"`
	Rooms     []string  `json:"rooms"`
	Watched   []string  `json:"watched"`
}
type Episode struct {
	ID             string    `json:"id"`
	AnimeShowID    string    `json:"animeShowId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Synopsis       string    `json:"synopsis"`
	Description    string    `json:"description"`
	Titles         Titles    `json:"titles"` // JSON field
	CanonicalTitle string    `json:"canonicalTitle"`
	SeasonNumber   int       `json:"seasonNumber"`
	Number         int       `json:"number"`
	RelativeNumber int       `json:"relativeNumber"`
	Airdate        string    `json:"airdate"`
	Length         int       `json:"length"`
	Thumbnail      Image     `json:"thumbnail"` // JSON field
}
