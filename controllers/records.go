package controllers

import (
	"fmt"
	"net/http"
	"time"
	"traffic-go/libs"
	"traffic-go/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type GetRecordsService struct {
	TimeLast int64  `json:"time_last"`
	TimeNew  int64  `json:"time_new"`
	Location string `json:"location"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type ret struct {
	Pics []pic `json:"pics"`
}

type pic struct {
	Location  string `json:"location"`
	ShootTime string `json:"shoot_time"`
	RuleType  string `json:"rule_type"`
	LicPlate  string `json:"lic_plate"`
	Direct    string `json:"direct"`
	PicPath   string `json:"pic_path"`
}

//GetRecords 获取记录列表
func GetRecords(c *gin.Context) {
	var service GetRecordsService
	var badpics []models.Badpic
	var total int //获取总条数
	var badpicsQuery *gorm.DB

	if err := c.ShouldBindJSON(&service); err != nil { // 如果绑定的字段为空的话 就会返回error 否则不
		c.JSON(http.StatusBadRequest, gin.H{"status": 400, "error": err.Error()})
		return
	}
	page := service.Page
	pageSize := service.PageSize
	timeLayout := "2006-01-02"
	timeLast := time.Unix(service.TimeLast, 0).Format(timeLayout)
	timeNew := time.Unix(service.TimeNew, 0).Format(timeLayout)
	if service.Location != "" { //location存在则查 否则不查
		var location models.Location
		models.DB.Select("id").Where("name = ?", service.Location).First(&location)
		locationID := location.ID //查询location_id
		fmt.Print("location_id:", locationID)
		badpicsQuery = models.DB.Where("location_id = ? AND shoot_time BETWEEN ? AND ?", locationID, timeNew, timeLast)
	} else { //location存在则查 否则不查
		//models.DB.Where("shoot_time BETWEEN ? AND ?",time_new,time_last).Find(&badpics).Count(&total)
		badpicsQuery = models.DB
	}
	badpicsQuery.Limit(pageSize).Offset((page - 1) * pageSize).Preload("Location").Find(&badpics)
	badpicsQuery.Table("badpics").Count(&total)

	data := &ret{Pics: []pic{}}

	for _, badpic := range badpics {
		data.Pics = append(data.Pics, pic{
			Location:  badpic.Location.Name,
			ShootTime: badpic.ShootTime,
			RuleType:  libs.RuletypeMap[badpic.RuleType],
			LicPlate:  badpic.LicPlate,
			Direct:    badpic.Direct,
			PicPath:   "img/" + badpic.PicPath,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "status": 200, "total": total})
}
