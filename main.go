package main

import (
	"Institution/config"
	"Institution/logs"
	"Institution/mysql"
	"Institution/redis"
	"Institution/router"
	"fmt"
	"os"
)

func main() {
	var rootPath string
	if len(os.Args) > 1 {
		rootPath = os.Args[1]
	} else {
		rootPath, _ = os.Getwd()
	}

	logs.SetRootPath(rootPath)
	logs.GetInstance().Logger.Infof("logger start!")
	config.InitServerConfig(rootPath + "/config/config.yaml")
	config := config.GetServerConfig()
	logs.GetInstance().Logger.Infof("config %+v", config)
	redis.RedisInit(&config.Redis)
	mysql.MysqlInit(config.MySQL)
	ginRouter := router.RouterInit(config, rootPath)

	ginRouter.Run(fmt.Sprintf(":%d", config.Server.Port))
}
