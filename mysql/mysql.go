package mysql

import (
	"Institution/config"
	"Institution/logs"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// 初始化gorm
func MysqlInit(config config.MySQLConfig) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Name, config.PassWord, config.Addr, config.DB)
	DB, err := gorm.Open(mysql.Open(dataSourceName))
	if err != nil {
		logs.GetInstance().Logger.Errorf("init mysql error %s", err)
	}
	db = DB
}

// 返回gorm对象
func GetClient() *gorm.DB {
	return db
}
