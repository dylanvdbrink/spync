package syncer

import (
	"api/internal/applemusic"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func StatusSocket(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	id := fmt.Sprintf("%s-%s", c.ClientIP(), randSeq(10))
	if sockets[id] == nil {
		connection := new(Connection)
		connection.Socket = conn
		sockets[id] = connection
	}

	for {
		t, msg, readErr := conn.ReadMessage()
		if readErr != nil {
			break
		}
		writeErr := conn.WriteMessage(t, msg)
		if writeErr != nil {
			return
		}
	}
}

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
