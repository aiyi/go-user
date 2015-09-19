package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

func init() {
	var err error

	db, err = sqlx.Open("mysql", getDSN())
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(20)
	db.Mapper = reflectx.NewMapper("json")
}

func getDSN() string {
	return "chanxuehong:chanxuehong@tcp(xxxxx:3306)/cxhtest?clientFoundRows=false&parseTime=true&loc=Asia%2FShanghai&timeout=5s&charset=utf8&collation=utf8_general_ci"
}
