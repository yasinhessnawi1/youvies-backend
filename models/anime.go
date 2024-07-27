package models

import "time"

type AnimeMovie struct {
	ID       int                  `json:"id"`
	Title    string               `json:"title"`
	Torrents map[string][]Torrent `bson:"torrents" json:"torrents"` // Transient field for processed torrents
	Links    struct {
		Self string `json:"self"`
	} `json:"links"`
	Attributes    Attributes    `json:"attributes"`
	Relationships Relationships `json:"relationships"`
	Genres        []string      `json:"genres"`
}

type AnimeShow struct {
	ID      int            `json:"id"`
	Title   string         `bson:"title" json:"title"`
	Seasons map[int]Season `bson:"seasons" json:"seasons"` // Transient field for processed torrents
	Links   struct {
		Self string `json:"self"`
	} `json:"links"`
	Attributes    Attributes           `bson:"attributes" json:"attributes"`
	Relationships Relationships        `bson:"relationships" json:"relationships"`
	Genres        []string             `bson:"genres" json:"genres"`
	Episodes      []EpisodeInfo        `bson:"episodes" json:"episodes"`
	FullContent   map[string][]Torrent `bson:"fullContent" json:"fullContent"`
}

type Attributes struct {
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	Slug                string    `json:"slug"`
	Synopsis            string    `json:"synopsis"`
	Description         string    `json:"description"`
	CoverImageTopOffset int       `json:"coverImageTopOffset"`
	Titles              struct {
		En   string `json:"en"`
		EnJp string `json:"en_jp"`
		EnUs string `json:"en_us"`
		JaJp string `json:"ja_jp"`
	} `json:"titles"`
	CanonicalTitle    string        `json:"canonicalTitle"`
	AbbreviatedTitles []interface{} `json:"abbreviatedTitles"`
	AverageRating     string        `json:"averageRating"`
	RatingFrequencies struct {
		Field1  string `json:"2"`
		Field2  string `json:"3"`
		Field3  string `json:"4"`
		Field4  string `json:"5"`
		Field5  string `json:"6"`
		Field6  string `json:"7"`
		Field7  string `json:"8"`
		Field8  string `json:"9"`
		Field9  string `json:"10"`
		Field10 string `json:"11"`
		Field11 string `json:"12"`
		Field12 string `json:"13"`
		Field13 string `json:"14"`
		Field14 string `json:"15"`
		Field15 string `json:"16"`
		Field16 string `json:"17"`
		Field17 string `json:"18"`
		Field18 string `json:"19"`
		Field19 string `json:"20"`
	} `json:"ratingFrequencies"`
	UserCount      int         `json:"userCount"`
	FavoritesCount int         `json:"favoritesCount"`
	StartDate      string      `json:"startDate"`
	EndDate        string      `json:"endDate"`
	NextRelease    interface{} `json:"nextRelease"`
	PopularityRank int         `json:"popularityRank"`
	RatingRank     int         `json:"ratingRank"`
	AgeRating      string      `json:"ageRating"`
	AgeRatingGuide string      `json:"ageRatingGuide"`
	Subtype        string      `json:"subtype"`
	Status         string      `json:"status"`
	Tba            interface{} `json:"tba"`
	PosterImage    struct {
		Tiny     string `json:"tiny"`
		Large    string `json:"large"`
		Small    string `json:"small"`
		Medium   string `json:"medium"`
		Original string `json:"original"`
		Meta     struct {
			Dimensions struct {
				Tiny struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"tiny"`
				Large struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"large"`
				Small struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"small"`
				Medium struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"medium"`
			} `json:"dimensions"`
		} `json:"meta"`
	} `json:"posterImage"`
	CoverImage struct {
		Tiny     string `json:"tiny"`
		Large    string `json:"large"`
		Small    string `json:"small"`
		Original string `json:"original"`
		Meta     struct {
			Dimensions struct {
				Tiny struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"tiny"`
				Large struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"large"`
				Small struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"small"`
			} `json:"dimensions"`
		} `json:"meta"`
	} `json:"coverImage"`
	EpisodeCount   int    `json:"episodeCount"`
	EpisodeLength  int    `json:"episodeLength"`
	TotalLength    int    `json:"totalLength"`
	YoutubeVideoId string `json:"youtubeVideoId"`
	ShowType       string `json:"showType"`
	Nsfw           bool   `json:"nsfw"`
}

type Relationships struct {
	Genres struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"links"`
	} `json:"genres"`
	Categories struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"categories"`
	} `json:"categories"`
	Castings struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"castings"`
	} `json:"castings"`
	Installments struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"installments"`
	} `json:"installments"`
	Mappings struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"mappings"`
	} `json:"mappings"`
	Reviews struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"reviews"`
	} `json:"reviews"`
	MediaRelationships struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"mediaRelationships"`
	} `json:"mediaRelationships"`
	Characters struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"characters"`
	} `json:"characters"`
	Staff struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"staff"`
	} `json:"staff"`
	Productions struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"productions"`
	} `json:"productions"`
	Quotes struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"quotes"`
	} `json:"quotes"`
	Episodes struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"episodes"`
	} `json:"episodes"`
	StreamingLinks struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"streamingLinks"`
	} `json:"streamingLinks"`
	AnimeProductions struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"links"`
	} `json:"animeProductions"`
	AnimeCharacters struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"animeCharacters"`
	} `json:"animeCharacters"`
	AnimeStaff struct {
		Links struct {
			Self    string `json:"self"`
			Related string `json:"related"`
		} `json:"animeStaff"`
	} `json:"animeStaff"`
}

type AnimeResponse struct {
	Data []Anime `json:"data"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
	Links struct {
		First string `json:"first"`
		Next  string `json:"next"`
		Last  string `json:"last"`
	} `json:"links"`
}
type Anime struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
	Attributes    Attributes    `json:"attributes"`
	Relationships Relationships `json:"relationships"`
}
type EpisodeResponse struct {
	Data []EpisodeInfo `json:"data"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
	Links struct {
		First string `json:"first"`
		Next  string `json:"next"`
		Last  string `json:"last"`
	} `json:"links"`
}
type EpisodeInfo struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
	Attributes struct {
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
		Synopsis    string    `json:"synopsis"`
		Description string    `json:"description"`
		Titles      struct {
			EnJp string `json:"en_jp"`
			EnUs string `json:"en_us"`
			JaJp string `json:"ja_jp"`
		} `json:"titles"`
		CanonicalTitle string      `json:"canonicalTitle"`
		SeasonNumber   int         `json:"seasonNumber"`
		Number         int         `json:"number"`
		RelativeNumber interface{} `json:"relativeNumber"`
		Airdate        string      `json:"airdate"`
		Length         int         `json:"length"`
		Thumbnail      struct {
			Original string `json:"original"`
			Meta     struct {
				Dimensions struct {
				} `json:"dimensions"`
			} `json:"meta"`
		} `json:"thumbnail"`
	} `json:"attributes"`
	Relationships struct {
		Media struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"media"`
		Videos struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"videos"`
	} `json:"relationships"`
}
