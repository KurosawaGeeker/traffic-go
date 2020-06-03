package router

import (
	"traffic-go/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("api/v1")
	{
		v1.POST("/records", controllers.GetRecords)
		v1.GET("/pictures", controllers.GetPictures)
		v1.POST("/pictures", controllers.SetPicture)
	}
	return router
}
