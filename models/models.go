package models

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)
var DB *gorm.DB

//database 在中间件中初始化mysql连接
func migration() {
	// 自动迁移模式
	DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&Location{},&Badpic{},&Goodpic{})
}
func Database(connString string) *gorm.DB{
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
	//migration()
	return DB
}
type Location struct {
	ID int64
	Name string
}
type Badpic struct {
	ID int64 `gorm:"PRIMARY_KEY;auto_increment"`
	PicPath string
	LocationId int
	ShootTime string
	RuleType int
	Rollback bool
}

type Goodpic struct {
	ID int64 `gorm:"PRIMARY_KEY;auto_increment"`
	PicPath string
	LocationId int
	ShootTime string
	RuleType int
	Rollback bool
}

