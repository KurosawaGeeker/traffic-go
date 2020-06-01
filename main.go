package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
	_ "time"
	"traffic-go/models"
)
var DB *gorm.DB
type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}
type Query struct {
	TimeLast 	int64 `form:"time_last" json:"time_last" xml:"time_last"`
	TimeNew 	int64 `form:"time_new" json:"time_new" xml:"time_new"`
	Location 	string `form:"location" json:"location" xml:"location"`
}
func main(){
	DB = models.Database("root:114514@tcp(127.0.0.1:3306)/traffic?charset=utf8")
	r := gin.Default()
	r.GET("/", func(ctx * gin.Context){
		ctx.String(200,"添加成功")
	})
	r.POST("/loginJSON", func(c *gin.Context) {
		var json Login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.User != "manu" || json.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		fmt.Print("User:",json.User)
		fmt.Print("psw:",json.Password)
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})
	r.POST("/query", func(c *gin.Context) {
		var json Query
		var badpics []models.Badpic
		var locations []models.Location
		if err := c.ShouldBindBodyWith(&json,binding.JSON); err != nil { // 如果绑定的字段为空的话 就会返回error 否则不
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}else{
			fmt.Print("json:",json)
			timeLayout := "2006-01-02 15:04:05"
			time_last := time.Unix(json.TimeLast, 0).Format(timeLayout)
			time_new :=time.Unix(json.TimeNew, 0).Format(timeLayout)
			if json.Location != ""{ //location存在则查 否则不查
				location_id  := DB.Where("name = ?", json.Location).First(&locations) //查询location_id]
				DB.Where("location_id = ? AND shoot_time BETWEEN ? AND ?",location_id,time_new,time_last).Find(&badpics)
			}else{//location存在则查 否则不查
				DB.Where("shoot_time BETWEEN ? AND ?",time_new,time_last).Find(&badpics)
			}
			fmt.Print(badpics)
			return
		}
	})
	r.Run(":8081")
}

