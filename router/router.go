package router

import (
	"traffic-go/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("v1/api")
	{
		v1.POST("/records", controllers.GetRecords)
		// v1.GET("/database", controllers.GetRecordList)
		// v1.POST("/database", controllers.GetRecordList)
	}
	return router
}
