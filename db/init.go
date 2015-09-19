package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

func init() {
	dsn, err := getDSN()
	if err != nil {
		panic(err)
	}

	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(20)
	db.Mapper = reflectx.NewMapper("sqlx")
}

func getDSN() (string, error) {
	return "chanxuehong:chanxuehong@tcp(xxxxx:3306)/cxhtest?clientFoundRows=false&parseTime=true&loc=Asia%2FShanghai&timeout=5s&charset=utf8&collation=utf8_general_ci", nil
}
