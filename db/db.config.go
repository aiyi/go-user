package db

import (
	"fmt"

	"github.com/aiyi/go-user/config"
)

func getDSN() (string, error) {
	mysqlConfig := config.ConfigData.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysqlConfig.UserName, mysqlConfig.Password,
		mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Database) +
		"?clientFoundRows=false&parseTime=true&loc=Asia%2FShanghai&timeout=5s&charset=utf8&collation=utf8_general_ci"
	return dsn, nil
}
