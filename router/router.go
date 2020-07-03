package router

import (
	"traffic-go/controllers"
	"traffic-go/middleware/jwt"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("api/v1")
	v1.Use(jwt.JWTAuth())  //中间件 需要经过登录验证的
	{
		v1.POST("/records", controllers.GetRecords)
	}
	v2 := router.Group("api/v2")
	{
		v2.GET("/pictures", controllers.GetPictures)
		v2.POST("/pictures", controllers.SetPicture)
		v2.POST("/auth",controllers.LoginUser)
		v2.POST("/register",controllers.RegisterUser)
	}
	return router
}
