package models

type TorrentResponse struct {
	Data []struct {
		Name       string   `json:"name"`
		Size       string   `json:"size"`
		Date       string   `json:"date"`
		Seeders    string   `json:"seeders"`
		Leechers   string   `json:"leechers"`
		Url        string   `json:"url"`
		Uploader   string   `json:"uploader"`
		Screenshot []string `json:"screenshot"`
		Category   string   `json:"category"`
		Files      []string `json:"files"`
		Poster     string   `json:"poster"`
		Magnet     string   `json:"magnet"`
		Hash       string   `json:"hash"`
	} `json:"data"`
	Time  float64 `json:"time"`
	Total int     `json:"total"`
}

type Torrent struct {
	Name       string   `json:"name"`
	Size       string   `json:"size"`
	Date       string   `json:"date"`
	Seeders    string   `json:"seeders"`
	Leechers   string   `json:"leechers"`
	Url        string   `json:"url"`
	Uploader   string   `json:"uploader"`
	Screenshot []string `json:"screenshot"`
	Category   string   `json:"category"`
	Files      []string `json:"files"`
	Poster     string   `json:"poster"`
	Magnet     string   `json:"magnet"`
	Hash       string   `json:"hash"`
}
