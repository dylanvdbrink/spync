package controllers

import (
	"api/internal/configuration"
	"github.com/gin-gonic/gin"
	"log"
)

func GetSettings(c *gin.Context) {
	config, err := configuration.GetConfiguration()
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	} else {
		c.JSON(200, config)
	}
}

func SaveSettings(c *gin.Context) {
	var config configuration.Config
	decodeErr := c.BindJSON(&config)
	if decodeErr != nil {
		log.Println("decoding error: " + decodeErr.Error())
		c.JSON(400, gin.H{
			"message": decodeErr.Error(),
		})
		return
	} else {
		saveErr := configuration.SaveConfiguration(config)
		if saveErr != nil {
			log.Println("save error: " + saveErr.Error())
			c.JSON(400, gin.H{
				"message": saveErr.Error(),
			})
			return
		}
		c.Status(204)
	}
}
