package controllers

import (
	"net/http"
	"fmt"
	"time"
	"traffic-go/models"
	"traffic-go/libs"

	"github.com/gin-gonic/gin"
)

type query struct {
	TimeLast int64  `json:"time_last"`
	TimeNew  int64  `json:"time_new"`
	Location string `json:"location"`
}

type ret struct {
	Pics[] pic `json:"pics"`
}

type pic struct {
	Location string `json:"location"`
	ShootTime string `'json:"shoot_time"`
	RuleType string `json:"rule_type"`
	LicPlate string `json:"lic_plate"`
	Direct string `json:"direct"`
}

//GetRecords 获取记录列表
func GetRecords(c *gin.Context) {
	var json query
	var badpics []models.Badpic
	if err := c.ShouldBindJSON(&json); err != nil { // 如果绑定的字段为空的话 就会返回error 否则不
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		timeLayout := "2006-01-02 15:04:05"
		timeLast := time.Unix(json.TimeLast, 0).Format(timeLayout)
		timeNew := time.Unix(json.TimeNew, 0).Format(timeLayout)
		if json.Location != ""{ //location存在则查 否则不查
			var location models.Location
			models.DB.Select("id").Where("name = ?", json.Location).First(&location)
			locationID := location.ID //查询location_id
			fmt.Print("location_id:",locationID)
			models.DB.Where("location_id = ? AND shoot_time BETWEEN ? AND ?",locationID,timeNew,timeLast).Find(&badpics)
		}else{//location存在则查 否则不查
			//models.DB.Where("shoot_time BETWEEN ? AND ?",time_new,time_last).Find(&badpics)
			models.DB.Find(&badpics)
		}
		data := &ret{Pics: []pic{}}
		
		for i:=0;i<5;i++ {
			var loctemp models.Location
			models.DB.Select("name").Where("id = ?", badpics[i].LocationID).First(&loctemp)
			data.Pics = append(data.Pics, pic{
				Location:loctemp.Name,
				ShootTime: badpics[i].ShootTime,
				RuleType: libs.RuletypeMap[badpics[i].RuleType],
				LicPlate:badpics[i].LicPlate,
				Direct:badpics[i].Direct,
			})
		}
		fmt.Print("data:",data)
		c.JSON(http.StatusOK,gin.H{"data":data,"status":200})
		return
	}
}
