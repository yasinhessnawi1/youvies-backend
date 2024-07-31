package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AnimeMovie struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty"  json:"id,omitempty"`
	Title    string               `bson:"title" json:"title"`
	Torrents map[string][]Torrent `bson:"torrents,omitempty" json:"torrents,omitempty"` // Transient field for processed torrents
	Links    struct {
		Self string `json:"self"`
	} `bson:"links,omitempty" json:"links,omitempty"`
	Attributes    Attributes    `json:"attributes"`
	Relationships Relationships `bson:"relationships,omitempty" json:"relationships,omitempty"`
	Genres        []string      `bson:"genres" json:"genres"`
}

type AnimeShow struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"  json:"id,omitempty"`
	Title   string             `bson:"title" json:"title"`
	Seasons map[int]Season     `bson:"seasons,omitempty" json:"seasons,omitempty"` // Transient field for processed torrents
	Links   struct {
		Self string `json:"self"`
	} `bson:"links,omitempty" json:"links,omitempty"`
	Attributes    Attributes           `bson:"attributes" json:"attributes"`
	Relationships Relationships        `bson:"relationships,omitempty" json:"relationships,omitempty"`
	Genres        []string             `bson:"genres" json:"genres"`
	Episodes      []EpisodeInfo        `bson:"episodes,omitempty" json:"episodes,omitempty"`
	FullContent   map[string][]Torrent `bson:"fullContent,omitempty" json:"fullContent,omitempty"`
}

type Attributes struct {
	CreatedAt           time.Time `bson:"created-at" json:"createdAt"`
	UpdatedAt           time.Time `bson:"updated-at" json:"updatedAt"`
	Slug                string    `bson:"slug" json:"slug"`
	Synopsis            string    `bson:"synopsis" json:"synopsis"`
	Description         string    `bson:"description" json:"description"`
	CoverImageTopOffset int       `bson:"cover-image-top-offset" json:"coverImageTopOffset"`
	Titles              struct {
		En   string `json:"en"`
		EnJp string `json:"en_jp"`
		EnUs string `json:"en_us"`
		JaJp string `json:"ja_jp"`
	} `bson:"titles" json:"titles"`
	CanonicalTitle    string        `bson:"canonical-title" json:"canonicalTitle"`
	AbbreviatedTitles []interface{} `bson:"abbreviated-titles,omitempty" json:"abbreviatedTitles,omitempty"`
	AverageRating     string        `bson:"average-rating" json:"averageRating"`
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
	} `bson:"rating-frequencies,omitempty" json:"ratingFrequencies,omitempty"`
	UserCount      int         `bson:"user-count,omitempty" json:"userCount,omitempty"`
	FavoritesCount int         `bson:"favorites-count" json:"favoritesCount"`
	StartDate      string      `bson:"start-date" json:"startDate"`
	EndDate        string      `bson:"end-date" json:"endDate"`
	NextRelease    interface{} `bson:"next-release" json:"nextRelease"`
	PopularityRank int         `bson:"popularity-rank" json:"popularityRank"`
	RatingRank     int         `bson:"rating-rank" json:"ratingRank"`
	AgeRating      string      `bson:"age-rating" json:"ageRating"`
	AgeRatingGuide string      `bson:"age-rating-guide" json:"ageRatingGuide"`
	Subtype        string      `bson:"subtype" json:"subtype"`
	Status         string      `bson:"status" json:"status"`
	Tba            interface{} `bson:"tba" json:"tba"`
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
		} `bson:"meta,omitempty" json:"meta,omitempty"`
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
	} `bson:"cover-image" json:"coverImage"`
	EpisodeCount   int    `bson:"episode-count" json:"episodeCount"`
	EpisodeLength  int    `bson:"episode-length" json:"episodeLength"`
	TotalLength    int    `bson:"total-length" json:"totalLength"`
	YoutubeVideoId string `bson:"youtube-video-id" json:"youtubeVideoId"`
	ShowType       string `bson:"show-type" json:"showType"`
	Nsfw           bool   `bson:"nsfw,omitempty" json:"nsfw,omitempty"`
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
			Self    string `json:"self,omitempty"`
			Related string `json:"related,omitempty"`
		} `json:"categories,omitempty"`
	} `json:"categories,omitempty"`
	Castings struct {
		Links struct {
			Self    string `json:"self,omitempty"`
			Related string `json:"related,omitempty"`
		} `json:"castings,omitempty"`
	} `json:"castings,omitempty"`
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
	Data []Anime `json:"data,omitempty"`
	Meta struct {
		Count int `json:"count,omitempty"`
	} `json:"meta,omitempty"`
	Links struct {
		First string `json:"first,omitempty"`
		Next  string `json:"next"`
		Last  string `json:"last,omitempty"`
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
	Type  string `bson:"type,omitempty" json:"type,omitempty"`
	Links struct {
		Self string `json:"self"`
	} `bson:"links,omitempty" json:"links,omitempty"`
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
		} `bson:"media,omitempty" json:"media,omitempty"`
		Videos struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"videos"`
	} `bson:"relationships,omitempty" json:"relationships,omitempty"`
}
