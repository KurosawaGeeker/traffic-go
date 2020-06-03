package models

import (
	"fmt"
	"time"

	"traffic-go/libs"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//DB 全局数据库
var DB *gorm.DB

func InitDB() *gorm.DB {
	config := libs.LoadConfig()
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		config.Mysql.User, config.Mysql.Passwd,
		config.Mysql.Host, config.Mysql.Port, config.Mysql.Name)
	db, err := gorm.Open("mysql", connStr) //连接数据库
	if err != nil {
		panic(err)
	}
	if config.Server.Mode == "debug" {
		db.LogMode(true)
	} else {
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
