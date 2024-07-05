package models

import (
	"time"
)

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
	Torrents []Torrent `json:"torrents,omitempty"`
	Id       string    `json:"id"`
	Type     string    `json:"type"`
	Links    struct {
		Self string `json:"self"`
	} `json:"links"`
	Attributes struct {
		CreatedAt           time.Time `json:"createdAt"`
		UpdatedAt           time.Time `json:"updatedAt"`
		Slug                string    `json:"slug"`
		Synopsis            string    `json:"synopsis"`
		Description         string    `json:"description"`
		CoverImageTopOffset int       `json:"coverImageTopOffset"`
		Titles              struct {
			En   string `json:"en,omitempty"`
			EnJp string `json:"en_jp"`
			JaJp string `json:"ja_jp"`
			EnUs string `json:"en_us,omitempty"`
		} `json:"titles"`
		CanonicalTitle    string   `json:"canonicalTitle"`
		AbbreviatedTitles []string `json:"abbreviatedTitles"`
		AverageRating     string   `json:"averageRating"`
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
		Tba            *string     `json:"tba"`
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
		CoverImage *struct {
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
		EpisodeLength  *int   `json:"episodeLength"`
		TotalLength    int    `json:"totalLength"`
		YoutubeVideoId string `json:"youtubeVideoId"`
		ShowType       string `json:"showType"`
		Nsfw           bool   `json:"nsfw"`
	} `json:"attributes"`
	Relationships struct {
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
			} `json:"links"`
		} `json:"categories"`
		Castings struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"castings"`
		Installments struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"installments"`
		Mappings struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"mappings"`
		Reviews struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"reviews"`
		MediaRelationships struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"mediaRelationships"`
		Characters struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"characters"`
		Staff struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"staff"`
		Productions struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"productions"`
		Quotes struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"quotes"`
		Episodes struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"episodes"`
		StreamingLinks struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
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
			} `json:"links"`
		} `json:"animeCharacters"`
		AnimeStaff struct {
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"animeStaff"`
	} `json:"relationships"`
}
