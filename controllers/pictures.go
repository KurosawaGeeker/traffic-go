package controllers

import (
	"io/ioutil"
	"net/http"
	"traffic-go/models"
	"traffic-go/protoc"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

//GetPicturesServices 获取图片接口参数
type GetPicturesServices struct {
	Type   string `form:"type"`
	Number int    `form:"number"`
}

//SetPicturesServices 设置图片接口参数
type SetPicturesServices struct {
	IsValid bool `json:"is_valid"`
	ID      int  `json:"id"`
}

//GetPictures 获取图片接口
func GetPictures(c *gin.Context) {
	var service GetPicturesServices
	var badpics []models.Badpic
	var protocResp protoc.Pics

	if err := c.ShouldBindQuery(&service); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "error": err.Error()})
		return
	}

	models.DB.Where("status = 0 and rule_type = ?", service.Type).
		Limit(service.Number).Preload("Location").Find(&badpics)
	for _, pic := range badpics {
		buffer, err := ioutil.ReadFile("static/" + pic.PicPath)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": 500, "error": err.Error()})
			return
		}
		protocResp.Pic = append(protocResp.Pic, &protoc.Pics_Picture{
			Id: int32(pic.ID), Location: pic.Location.Name, PicData: buffer})
	}
	c.ProtoBuf(http.StatusOK, &protocResp)

}

//SetPicture 设置图片接口
func SetPicture(c *gin.Context) {
	var service SetPicturesServices
	var badpic models.Badpic
	var goodpic models.Goodpic

	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "error": err.Error()})
		return
	}

	if err := models.DB.Where("id = ?", service.ID).First(&badpic).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "error": "Not Found"})
		return
	}
	models.DB.Model(&badpic).Update("status", true)
	if service.IsValid {
		copier.Copy(&goodpic, &badpic)
		goodpic.ID = 0
		if err := models.DB.Delete(&badpic).Create(&goodpic).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"status": 500, "error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": 200, "is_ok": true})
}
