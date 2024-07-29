package utils

import (
	"encoding/json"
	"fmt"
	"github.com/anacrolix/torrent"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"youvies-backend/database"
	"youvies-backend/models"
)

const (
	TorrentStorageDir = "./torrents" // Directory to store torrent files
	DefaultClientPort = 42069
)

var mutex sync.Mutex

func FetchTorrents(title string) ([]models.Torrent, error) {
	var torrents []models.Torrent

	url := fmt.Sprintf("%ssearch?query=%s", TorrentAPIBaseURL, strings.ReplaceAll(title, " ", "%20"))

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching torrents: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(fmt.Errorf("error closing request body: %s", err))
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var result models.TorrentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	var movieCategories = []string{
		"movie", "film", "cinema", "blockbuster", "feature film",
		"motion picture", "flick", "biopic", "documentary", "short film",
		"thriller", "comedy", "drama", "action", "adventure",
		"animation", "crime", "fantasy", "historical", "horror",
		"musical", "mystery", "romance", "sci-fi", "science fiction",
		"war", "western", "independent film", "indie film", "art house",
		"silent film", "noir", "cult film", "Video > Movies", "Video > HD - TV shows",
	}
	var showCategories = []string{
		"show", "shows", "tv show", "tv shows", "television show", "series", "tv series",
		"sitcom", "reality show", "talk show", "drama series", "comedy series",
		"mini-series", "soap opera", "docuseries", "children's show",
		"news show", "variety show", "game show", "late-night show",
		"cooking show", "competition show", "talent show", "true crime",
		"crime drama", "fantasy series", "sci-fi series", "science fiction series",
		"historical drama", "superhero series", "animated series", "anime series",
		"documentary series", "medical drama", "legal drama", "reality competition", "Video > HD - TV shows",
	}
	var animeCategories = []string{
		"anime", "manga", "ova", "ona", "anime series", "anime movie",
		"light novel", "hentai", "josei", "seinen", "shonen", "shojo",
		"yaoi", "yuri", "anime film", "isekai", "mecha", "slice of life",
		"shoujo-ai", "shounen-ai", "magical girl", "sports anime", "supernatural",
		"fantasy anime", "sci-fi anime", "science fiction anime", "romance anime",
		"action anime", "adventure anime", "comedy anime", "drama anime",
		"historical anime", "horror anime", "music anime", "psychological anime",
		"school anime", "space anime", "thriller anime", "military anime",
	}

	for _, torrent := range result.Data {
		category := strings.ToLower(torrent.Category)

		if containsAny(category, movieCategories) || containsAny(category, showCategories) || containsAny(category, animeCategories) {
			torrents = append(torrents, torrent)
		}
	}

	fmt.Printf("Found %d torrents for %s: %s \n", len(torrents), title, url)
	if len(torrents) == 0 {
		return nil, fmt.Errorf("no torrents found for title: %s", title)
	}

	return torrents, nil
}

func containsAny(text string, items []string) bool {
	text = strings.TrimSpace(strings.ToLower(text))
	for _, item := range items {
		item = strings.TrimSpace(strings.ToLower(item))
		if strings.Contains(text, item) {
			return true
		}
	}
	return false
}

func SaveMetadata(magnetURI, torrentName string) error {
	if !strings.HasPrefix(magnetURI, "magnet:") {
		return fmt.Errorf("invalid magnet URI: %s", magnetURI)
	}

	// Create a unique directory for this torrent
	dirPath := filepath.Join(TorrentStorageDir, generateUniqueID())
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory for torrent: %w", err)
	}

	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.DataDir = dirPath
	clientConfig.ListenPort = getUniquePort()

	torrentClient, err := torrent.NewClient(clientConfig)
	if err != nil {
		return fmt.Errorf("failed to create torrent client: %w", err)
	}

	t, err := torrentClient.AddMagnet(magnetURI)
	if err != nil {
		torrentClient.Close()
		return fmt.Errorf("failed to add magnet: %w", err)
	}

	<-t.GotInfo()
	mi := t.Metainfo()

	infoHash := t.InfoHash().HexString()
	filePath := filepath.Join(dirPath, infoHash+".torrent")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f, err := os.Create(filePath)
		if err != nil {
			torrentClient.Close()
			return fmt.Errorf("failed to create metadata file: %w", err)
		}
		defer f.Close()

		err = mi.Write(f)
		if err != nil {
			torrentClient.Close()
			return fmt.Errorf("failed to write metadata to file: %w", err)
		}
	}

	exists, err := database.IfItemExists(bson.M{"_id": infoHash}, "torrent_files")
	if exists {
		torrentClient.Close()
		return nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		torrentClient.Close()
		return fmt.Errorf("failed to read metadata file: %w", err)
	}

	err = storeTorrentFileInDB(infoHash, content)
	if err != nil {
		torrentClient.Close()
		return fmt.Errorf("failed to store torrent file in database: %w", err)
	}

	torrentClient.Close()

	// Retry file deletion
	err = retry(3, 100*time.Millisecond, func() error {
		return os.RemoveAll(dirPath)
	})
	if err != nil {
		log.Printf("Failed to delete directory %s after retries: %v", dirPath, err)
	}

	return nil
}

func storeTorrentFileInDB(infoHash string, content []byte) error {
	err := database.InsertItem(bson.M{
		"_id":     infoHash,
		"content": content,
	}, infoHash, "torrent_files")
	if err != nil {
		return fmt.Errorf("failed to insert torrent file into database: %w", err)
	}

	return nil
}

func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

var portMutex sync.Mutex
var currentPort = DefaultClientPort

func getUniquePort() int {
	portMutex.Lock()
	defer portMutex.Unlock()
	currentPort++
	return currentPort
}

func retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return retry(attempts, sleep, f)
		}
		return err
	}
	return nil
}
