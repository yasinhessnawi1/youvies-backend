package api

import (
	"github.com/anacrolix/torrent"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	client         *torrent.Client
	clientMutex    sync.Mutex
	currentTorrent *torrent.Torrent
	defaultDir     = ".tmp"
)

func initClient() error {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	if client == nil {
		clientConfig := torrent.NewDefaultClientConfig()
		clientConfig.NoDefaultPortForwarding = true
		clientConfig.DataDir = defaultDir
		var err error
		client, err = torrent.NewClient(clientConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanupClient() {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	if client != nil {
		client.Close()
		client = nil
	}
}

func streamHandler(c *gin.Context) {
	magnet := c.Query("magnet")
	if magnet == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Magnet link is required"})
		return
	}

	err := initClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create torrent client"})
		return
	}

	clientMutex.Lock()
	if currentTorrent != nil {
		currentTorrent.Drop()
		currentTorrent = nil
	}
	clientMutex.Unlock()

	t, err := client.AddMagnet(magnet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add magnet link"})
		return
	}

	<-t.GotInfo()
	t.DownloadAll()

	clientMutex.Lock()
	currentTorrent = t
	clientMutex.Unlock()

	var videoFile *torrent.File
	for _, f := range t.Files() {
		if strings.HasSuffix(f.Path(), ".mp4") || strings.HasSuffix(f.Path(), ".mkv") && !strings.Contains(f.Path(), "sample") {
			videoFile = f
			break
		}
	}

	if videoFile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No suitable video file found in the torrent"})
		return
	}

	c.Header("Content-Type", "video/mp4")

	c.Stream(func(w io.Writer) bool {
		r := videoFile.NewReader()
		if err != nil {
			return false
		}
		defer r.Close()

		buf := make([]byte, 32*1024)
		for {
			n, err := r.Read(buf)
			if err != nil && err != io.EOF {
				return false
			}
			if n == 0 {
				return false
			}
			_, err = w.Write(buf[:n])
			if err != nil {
				return false
			}
		}
	})
}

func flushHandler() {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	if currentTorrent != nil {
		currentTorrent.Drop()
		currentTorrent = nil
	}
	cleanupClient()

	err := os.RemoveAll(defaultDir)
	if err != nil {
		return
	}

	err = os.Remove(".torrent.bolt.db")
	if err != nil && !os.IsNotExist(err) {
		return
	}

}
