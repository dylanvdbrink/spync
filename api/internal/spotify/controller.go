package spotify

import (
	"github.com/gin-gonic/gin"
)

func GetAuthURLEndpoint(c *gin.Context) {
	c.Redirect(301, GetAuthUrl())
}

func SaveAuthEndpoint(c *gin.Context) {
	state := c.Query("state")
	err := SaveAuth(c.Request, state)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
	}
}

func GetPlaylistsEndpoint(c *gin.Context) {
	playlists, err := GetUserPlaylists()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, playlists)
}

func GetPlaylistTracksEndpoint(c *gin.Context) {
	playlistId := c.Param("playlistId")
	tracks, err := GetPlaylistTracks(playlistId)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, tracks)
}

func GetMeEndpoint(c *gin.Context) {
	me, err := GetMe()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, me)
}
