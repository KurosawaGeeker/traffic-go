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
	migration()
	return DB
}
type Location struct {
	ID int64
	Name string
	Badpics []Badpic	`gorm:"ForeignKey:LocationId"`
	Goodpics []Goodpic `gorm:"ForeignKey:LocationId"`
}
type Badpic struct {
	ID int64 `gorm:"PRIMARY_KEY;auto_increment"` //id
	PicPath string	//图片路径 pic_path
	LocationId int 	`gorm:"column:location_id;"`//违法地点外键 location_id
	LicKind string
	ShootInstitution string
	Backup string
	ShootTime string  //违法时间 shoot_time
	RuleType string  	//违规类型 rule_type
	Rollback bool 	// 是否经过回滚 rollback
	LicPlate string 	//车牌号码 lic_plate
	Direct string  //方向 direct
	Status bool //校对1 未校对 0 status
}

type Goodpic struct {
	ID int64 `gorm:"PRIMARY_KEY;auto_increment"` //id
	PicPath string	//图片路径 pic_path
	LocationId int 	`gorm:"column:location_id;"`//违法地点外键 location_id
	LicKind string
	ShootInstitution string
	Backup string
	ShootTime string  //违法时间 shoot_time
	RuleType string  	//违规类型 rule_type
	Rollback bool 	// 是否经过回滚 rollback
	LicPlate string 	//车牌号码 lic_plate
	Direct string  //方向 direct
	Status bool //校对1 未校对 0 status
}

