package mysql

import (
	"fmt"
	config "schedule/conf"
	"schedule/logs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func MysqlInit(config config.MySQLConfig) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Name, config.PassWord, config.Addr, config.DB)
	DB, err := gorm.Open(mysql.Open(dataSourceName))
	if err != nil {
		logs.GetInstance().Logger.Errorf("init mysql error %s", err)
	}
	db = DB
}

func GetClient() *gorm.DB {
	return db
}
