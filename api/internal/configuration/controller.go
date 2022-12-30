package configuration

import (
	"github.com/gin-gonic/gin"
	"log"
)

func GetSettingsEndpoint(c *gin.Context) {
	config, err := GetConfiguration()
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	} else {
		c.JSON(200, config)
	}
}

func SaveSettingsEndpoint(c *gin.Context) {
	var config Config
	decodeErr := c.BindJSON(&config)
	if decodeErr != nil {
		log.Println("decoding error: " + decodeErr.Error())
		c.JSON(400, gin.H{
			"message": decodeErr.Error(),
		})
		return
	} else {
		saveErr := SaveConfiguration(config, true)
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
