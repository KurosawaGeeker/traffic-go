package models

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"time"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
var DB *gorm.DB

//database 在中间件中初始化mysql连接
func migration() {
	// 自动迁移模式
	DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&User{}).
		AutoMigrate(&Video{})
}
func Database(connString string) {
	db, err := gorm.Open("mysql", connString) //连接数据库
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	if gin.Mode() == "release" {
		db.LogMode(false)

	}
	db.DB().SetMaxIdleConns(20)
	//打开
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Second * 30)
	DB = db
	migration()
}
type User struct{
	gorm.Model// 包含id 创建 修改 删除时间
	username string
}// 用户模型