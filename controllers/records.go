package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
	"traffic-go/libs"
	"traffic-go/models"
)

type GetRecordsService struct {
	TimeLast int64  `json:"time_last"`
	TimeNew  int64  `json:"time_new"`
	Location string `json:"location"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	RuleType string `json:"rule_type"`
	Waste 	bool 	`json:"waste"` //为1 废片 0 正片 默认正片
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
	RollBack 	bool `json:"roll_back"` //正片/废片是否经过回滚
	IsValid 	bool `json:"is_valid"` //废片是否有效

}

func GetRecords(c *gin.Context){
	var service GetRecordsService
	var total int
	var picsQuery *gorm.DB
	if err := c.ShouldBindJSON(&service);err !=nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 400, "error": err.Error()})
		return
	}
	timeLayout := "2006-01-02"
	timeLast := time.Unix(service.TimeLast, 0).Format(timeLayout)
	timeNew := time.Unix(service.TimeNew, 0).Format(timeLayout)
	page := service.Page
	pageSize := service.PageSize
	if service.RuleType == "all" {
		if service.Location != "" 	{ //location存在则查 否则不查
			var location models.Location
			models.DB.Select("id").Where("name = ?", service.Location).First(&location)
			locationID := location.ID //查询location_id
			picsQuery = models.DB.Where("location_id = ? AND shoot_time BETWEEN ? AND ?", locationID, timeNew, timeLast)
		} else { //location存在则查 否则不查
			picsQuery = models.DB.Where("shoot_time BETWEEN ? AND ?", timeNew, timeLast)
		}
	}else{
		rule_type_key := ""
		for key := range libs.RuletypeMap{
			if libs.RuletypeMap[key] == service.RuleType{
				rule_type_key = key
				break
			}else{
				rule_type_key = ""
			}
		}
		if service.Location != "" { //location存在则查 否则不查
			var location models.Location
			models.DB.Select("id").Where("name = ?", service.Location).First(&location)
			locationID := location.ID //查询location_id


			picsQuery = models.DB.Where("location_id = ? AND rule_type = ? AND shoot_time BETWEEN ? AND ? ", locationID, rule_type_key,timeNew, timeLast)
		} else { //location存在则查 否则不查
			picsQuery = models.DB.Where("rule_type = ? AND shoot_time BETWEEN ? AND ?", rule_type_key ,timeNew, timeLast)
		}
	}
	data := &ret{Pics: []pic{}}
	if service.Waste==true {
		var badpics []models.Badpic
		picsQuery.Limit(pageSize).Offset((page - 1) * pageSize).Preload("Location").Find(&badpics)
		picsQuery.Table("badpics").Count(&total)
		for _, badpic := range badpics {
			data.Pics = append(data.Pics, pic{
				Location:  badpic.Location.Name,
				ShootTime: badpic.ShootTime,
				RuleType:  libs.RuletypeMap[badpic.RuleType],
				LicPlate:  badpic.LicPlate,
				Direct:    badpic.Direct,
				PicPath:   "img/" + badpic.PicPath,
				RollBack:  badpic.Rollback,
				IsValid:   badpic.Status,
			})
		}
	}else{
		var goodpics []models.Goodpic
		picsQuery.Limit(pageSize).Offset((page - 1) * pageSize).Preload("Location").Find(&goodpics)
		picsQuery.Table("goodpics").Count(&total)
		for _, goodpic := range goodpics {
			data.Pics = append(data.Pics, pic{
				Location:  goodpic.Location.Name,
				ShootTime: goodpic.ShootTime,
				RuleType:  libs.RuletypeMap[goodpic.RuleType],
				LicPlate:  goodpic.LicPlate,
				Direct:    goodpic.Direct,
				PicPath:   "img/" + goodpic.PicPath,
				RollBack:  goodpic.Rollback,
				IsValid:   goodpic.Status,
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "status": 200, "total": total})
}
