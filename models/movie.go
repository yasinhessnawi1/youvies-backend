package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Movie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Year        string             `bson:"year" json:"year"`
	Director    string             `bson:"director" json:"director"`
	Genres      string             `bson:"genres" json:"genres"`
	Torrents    []Torrent          `bson:"torrents" json:"torrents"`
	Rating      string             `bson:"rating,omitempty" json:"rating,omitempty"`
	PosterURL   string             `bson:"poster_url" json:"poster_url"`
	Language    string             `bson:"language" json:"language"`
}

type ImdbMovie struct {
	Page    int `json:"page"`
	Results []struct {
		Adult            bool    `json:"adult"`
		BackdropPath     string  `json:"backdrop_path"`
		GenreIds         []int   `json:"genre_ids"`
		Id               int     `json:"id"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		Overview         string  `json:"overview"`
		Popularity       float64 `json:"popularity"`
		PosterPath       string  `json:"poster_path"`
		ReleaseDate      string  `json:"release_date"`
		Title            string  `json:"title"`
		Video            bool    `json:"video"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int     `json:"vote_count"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
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

type YTSResponse struct {
	Status        string `json:"status"`
	StatusMessage string `json:"status_message"`
	Data          struct {
		MovieCount int `json:"movie_count"`
		Limit      int `json:"limit"`
		PageNumber int `json:"page_number"`
		Movies     []struct {
			Id                      int       `json:"id"`
			Url                     string    `json:"url"`
			ImdbCode                string    `json:"imdb_code"`
			Title                   string    `json:"title"`
			TitleEnglish            string    `json:"title_english"`
			TitleLong               string    `json:"title_long"`
			Slug                    string    `json:"slug"`
			Year                    int       `json:"year"`
			Rating                  float64   `json:"rating"`
			Runtime                 int       `json:"runtime"`
			Genres                  []string  `json:"genres"`
			Summary                 string    `json:"summary"`
			DescriptionFull         string    `json:"description_full"`
			Synopsis                string    `json:"synopsis"`
			YtTrailerCode           string    `json:"yt_trailer_code"`
			Language                string    `json:"language"`
			MpaRating               string    `json:"mpa_rating"`
			BackgroundImage         string    `json:"background_image"`
			BackgroundImageOriginal string    `json:"background_image_original"`
			SmallCoverImage         string    `json:"small_cover_image"`
			MediumCoverImage        string    `json:"medium_cover_image"`
			LargeCoverImage         string    `json:"large_cover_image"`
			State                   string    `json:"state"`
			Torrents                []Torrent `json:"torrents"`
			DateUploaded            string    `json:"date_uploaded"`
			DateUploadedUnix        int       `json:"date_uploaded_unix"`
		} `json:"movies"`
	} `json:"data"`
	Meta struct {
		ServerTime     int    `json:"server_time"`
		ServerTimezone string `json:"server_timezone"`
		ApiVersion     int    `json:"api_version"`
		ExecutionTime  string `json:"execution_time"`
	} `json:"@meta"`
}
