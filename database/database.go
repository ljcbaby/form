package database

import (
	"fmt"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ljcbaby/select/config"
)

var (
	MySQL *sql.DB
)

func ConnectMySQL() error {
	conf := config.Conf.MySQL
	var err error
	MySQL, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
	))
	if err != nil {
		return err
	}
	err = MySQL.Ping()
	return err
}
