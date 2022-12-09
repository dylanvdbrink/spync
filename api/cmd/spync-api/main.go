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

	err = router.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
