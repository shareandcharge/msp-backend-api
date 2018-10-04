package tools

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

var MDB *sqlx.DB

func MySQLConnect(dbName string) {
	MDB = sqlx.MustConnect("mysql", "andy:hardpassword1@(54.93.205.218:3306)/blockchain")

	//some benchmark should be done here
	MDB.SetMaxOpenConns(300)
	MDB.SetMaxIdleConns(10)
	MDB.SetConnMaxLifetime(10 * time.Second)
}
