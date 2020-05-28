package main

import (
	"github.com/gin-gonic/gin"
	"traffic-go/models"
)

func main(){
	models.Database("root:114514@tcp(127.0.0.1:3306)/test?charset=utf8")
	r := gin.Default()
	r.GET("/", func(ctx * gin.Context){
		ctx.String(200,"index")
	})
	r.POST("/query", func(cxt *gin.Context) {
		 // 这里编写通过时间地点进行查询接口
	})
	r.Run()
}