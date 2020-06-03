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
	Page int `json:"page"`
	PageSize int `json:"page_size"'`
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
		var total int  //获取总条数
		page := json.Page
		page_size := json.PageSize
		timeLayout := "2006-01-02 15:04:05"
		timeLast := time.Unix(json.TimeLast, 0).Format(timeLayout)
		timeNew := time.Unix(json.TimeNew, 0).Format(timeLayout)
		if json.Location != ""{ //location存在则查 否则不查
			var location models.Location
			models.DB.Select("id").Where("name = ?", json.Location).First(&location)
			locationID := location.ID //查询location_id
			fmt.Print("location_id:",locationID)

			models.DB.Limit(page_size).Offset((page-1)*page_size).Where("location_id = ? AND shoot_time BETWEEN ? AND ?",locationID,timeNew,timeLast).Find(&badpics).Count(&total)
		}else{//location存在则查 否则不查
			//models.DB.Where("shoot_time BETWEEN ? AND ?",time_new,time_last).Find(&badpics).Count(&total)
			models.DB.Limit(page_size).Offset((page-1)*page_size).Find(&badpics).Count(&total)
		}
		data := &ret{Pics: []pic{}}

		for i:=0;i<len(badpics);i++ {
			var loctemp models.Location
			models.DB.Select("name").Where("id = ?", badpics[i].LocationID).First(&loctemp)
			fmt.Print("len:",len(badpics))
			data.Pics = append(data.Pics, pic{
				Location:loctemp.Name,
				ShootTime: badpics[i].ShootTime,
				RuleType: libs.RuletypeMap[badpics[i].RuleType],
				LicPlate:badpics[i].LicPlate,
				Direct:badpics[i].Direct,
			})
		}
		c.JSON(http.StatusOK,gin.H{"data":data,"status":200,"total":total})
		return
	}
}
