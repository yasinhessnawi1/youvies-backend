package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
)

type MagnetRequest struct {
	MagnetURI string `json:"magnet"`
}

func StreamTorrentHandler(w http.ResponseWriter, r *http.Request) {
	var req MagnetRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	magnetURI := req.MagnetURI
	if magnetURI == "" {
		http.Error(w, "Magnet link is required", http.StatusBadRequest)
		return
	}

	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.DataDir = "./torrents"
	client, err := torrent.NewClient(clientConfig)
	if err != nil {
		http.Error(w, "Failed to create torrent client", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	t, err := client.AddMagnet(magnetURI)
	if err != nil {
		http.Error(w, "Failed to add magnet link", http.StatusInternalServerError)
		return
	}

	<-t.GotInfo()

	var videoFile *torrent.File
	for _, file := range t.Files() {
		if strings.HasSuffix(file.Path(), ".mp4") || strings.HasSuffix(file.Path(), ".mkv") {
			videoFile = file
			break
		}
	}

	if videoFile == nil {
		http.Error(w, "No video file found in torrent", http.StatusNotFound)
		return
	}

	videoFile.Download()

	time.Sleep(5 * time.Second)

	videoPath := filepath.Join(clientConfig.DataDir, videoFile.Path())

	fmt.Println(videoPath)

	http.ServeFile(w, r, videoPath)
}
