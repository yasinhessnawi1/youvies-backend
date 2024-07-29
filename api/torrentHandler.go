package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"youvies-backend/database"
)

// GetTorrentFile serves the .torrent file from the database
func GetTorrentFile(c *gin.Context) {
	id := c.Param("id")

	collection := database.Client.Database("youvies").Collection("torrent_files")
	var torrentFile struct {
		ID      string `bson:"_id"`
		Content []byte `bson:"content"`
	}
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&torrentFile)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Torrent file not found"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+id+".torrent")
	c.Data(http.StatusOK, "application/x-bittorrent", torrentFile.Content)
}
