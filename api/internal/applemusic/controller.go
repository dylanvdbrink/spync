package applemusic

import (
	"api/internal/spotify"
	"github.com/gin-gonic/gin"
)

type TokenObject struct {
	Token string `json:"token"`
}

func FindTrackEndpoint(c *gin.Context) {
	trackId := c.Param("trackId")
	track, err := spotify.GetTrack(trackId)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	} else {
		applemusicTrack, findErr := FindTrack(track)
		if findErr != nil {
			return
		}

		c.JSON(200, applemusicTrack)
	}
}

func GetSpotifyPlaylistsEndpoint(c *gin.Context) {
	playlists, err := GetSpotifyPlaylists()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(200, playlists)
		return
	}
}

func SaveAuthEndpoint(c *gin.Context) {
	var tokenObj TokenObject
	decodeErr := c.BindJSON(&tokenObj)
	if decodeErr != nil {
		c.JSON(400, gin.H{
			"message": decodeErr.Error(),
		})
		return
	} else {
		err := SaveAuth(tokenObj.Token)
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.Status(204)
	}
}
