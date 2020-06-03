package libs

import (
	"log"
	"sync"

	"github.com/BurntSushi/toml"
)

//AppConfig 项目配置
type AppConfig struct {
	Server serverConfig
	Mysql  mysqlConfig
}

type serverConfig struct {
	Mode         string
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

type mysqlConfig struct {
	User   string
	Passwd string
	Host   string
	Port   int
	Name   string
}

var (
	configPath string = "./config/app.toml"
	config     AppConfig
	once       sync.Once
)

//LoadConfig 获取项目配置
func LoadConfig() AppConfig {
	once.Do(func() {
		if _, err := toml.DecodeFile(configPath, &config); err != nil {
			log.Fatal(err)
		}
	})
	return config
}
