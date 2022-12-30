package syncer

import (
	"api/internal/applemusic"
	"github.com/gin-gonic/gin"
)

func SyncPlaylistEndpoint(c *gin.Context) {
	playlistId := c.Param("playlistId")

	err := applemusic.CheckAuth()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = SyncPlaylist(playlistId)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Sync done",
	})
}

func SyncPlaylistsEndpoint(c *gin.Context) {
	SyncAllPlaylists()

	c.JSON(200, gin.H{
		"message": "Sync done",
	})
}
