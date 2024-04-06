package main

import (
	"Institution/config"
	"Institution/logs"
	"Institution/mysql"
	"Institution/redis"
	"Institution/router"
	"fmt"
)

func main() {
	logs.GetInstance().Logger.Infof("logger start!")
	config.InitServerConfig("config/config.yaml")
	config := config.GetServerConfig()
	logs.GetInstance().Logger.Infof("config %+v", config)
	redis.RedisInit(&config.Redis)
	mysql.MysqlInit(config.MySQL)
	ginRouter := router.RouterInit(config)

	ginRouter.Run(fmt.Sprintf(":%d", config.Server.Port))
}
