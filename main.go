package main

import (
	"fmt"
	config "schedule/conf"
	"schedule/logs"
	"schedule/mysql"
	"schedule/redis"
	"schedule/router"
)

func main() {
	logs.GetInstance().Logger.Infof("logger start!")
	config.InitServerConfig("conf/config.yaml")
	config := config.GetServerConfig()
	logs.GetInstance().Logger.Infof("config %+v", config)
	redis.RedisInit(&config.Redis)
	mysql.MysqlInit(config.MySQL)
	ginRouter := router.RouterInit()

	ginRouter.Run(fmt.Sprintf(":%s", config.Server.Port))
}
