package controllers

import (
	"api/internal/clients"
	"github.com/gin-gonic/gin"
)

func GetAuthURL(c *gin.Context) {
	c.Redirect(301, clients.GetAuthUrl())
}

func SaveAuth(c *gin.Context) {
	state := c.Query("state")
	err := clients.SaveAuth(c.Request, state)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
	}
}

func GetPlaylists(c *gin.Context) {
	playlists, err := clients.GetUserPlaylists()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, playlists)
}

func GetMe(c *gin.Context) {
	me, err := clients.GetMe()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, me)
}
