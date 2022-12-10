package main

import (
	"api/internal/controllers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	err := router.SetTrustedProxies([]string{})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	router.GET("/ping", controllers.Ping)

	router.GET("/settings", controllers.GetSettings)
	router.POST("/settings", controllers.SaveSettings)

	router.GET("/spotify/auth", controllers.GetAuthURL)
	router.GET("/spotify/save-auth", controllers.SaveAuth)
	router.GET("/spotify/me", controllers.GetMe)
	router.GET("/spotify/me/playlists", controllers.GetPlaylists)

	err = router.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
