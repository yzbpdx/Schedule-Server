package main

import (
	config "schedule/conf"
	"schedule/gorm"
	"schedule/logs"
	"schedule/redis"
)

func main() {
	logs.GetInstance().Logger.Infof("logger start!")
	config.InitServerConfig("conf/server.yaml")
	config := config.GetServerConfig()
	logs.GetInstance().Logger.Infof("config %+v", config)
	redis.RedisInit(&config.Redis)
	gorm.MysqlInit(config.MySQL)

}
