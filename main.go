package main

import (
	"strconv"
	_ "time"
	"traffic-go/libs"
	"traffic-go/models"
	"traffic-go/router"

	"github.com/gin-gonic/gin"
)

func main() {
	config := libs.LoadConfig()
	gin.SetMode(config.Server.Mode)
	config = libs.LoadConfig()

	db := models.InitDB()
	defer db.Close()

	r := router.InitRouter()
	r.Run(":" + strconv.Itoa(config.Server.Port))
}
