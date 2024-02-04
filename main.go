package main

import (
	"schedule/algorithm"
	config "schedule/conf"
	"schedule/logs"
	"schedule/mysql"
	"schedule/redis"
)

func main() {
	logs.GetInstance().Logger.Infof("logger start!")
	config.InitServerConfig("conf/config.yaml")
	config := config.GetServerConfig()
	logs.GetInstance().Logger.Infof("config %+v", config)
	redis.RedisInit(&config.Redis)
	mysql.MysqlInit(config.MySQL)

	algorithm.StartSchedule()
}
