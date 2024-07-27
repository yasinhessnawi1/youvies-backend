package models

type TorrentResponse struct {
	Data  []Torrent `json:"data"`
	Time  float64   `json:"time"`
	Total int       `json:"total"`
}

type Torrent struct {
	Name     string `bson:"name" json:"name"`
	Size     string `bson:"size" json:"size"`
	Date     string `bson:"date" json:"date"`
	Seeders  string `bson:"seeders" json:"seeders"`
	Leechers string `bson:"leechers" json:"leechers"`
	Url      string `bson:"url" json:"url"`
	Uploader string `bson:"uploader" json:"uploader"`
	Category string `bson:"category" json:"category"`
	Poster   string `bson:"poster" json:"poster"`
	Magnet   string `bson:"magnet" json:"magnet"`
	Hash     string `bson:"hash" json:"hash"`
}
