package router

import (
	"net/http"
	"traffic-go/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.StaticFS("/img", http.Dir("./static/img"))

	v1 := router.Group("api/v1")
	{
		v1.POST("/records", controllers.GetRecords)
		v1.GET("/pictures", controllers.GetPictures)
		v1.POST("/pictures", controllers.SetPicture)
	}
	return router
}
