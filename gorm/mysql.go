package gorm

import (
	"database/sql"
	"fmt"
	config "schedule/conf"
	"schedule/logs"

	_ "github.com/go-sql-driver/mysql"
)

var client map[string]*sql.DB

func MysqlInit(mysqlConfig map[string]config.MySQLConfig) {
	for clientName, config := range mysqlConfig {
		dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Name, config.PassWord, config.Addr, config.DB)
		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			logs.GetInstance().Logger.Errorf("init mysql error %s", err)
		}
		client[clientName] = db
	}
}

func GetClient(clientName string) *sql.DB {
	return client[clientName]
}
