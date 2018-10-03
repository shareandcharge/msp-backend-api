package tools

// you should not rely on mysql for checking the history. try fatdb api instead!

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var MDB *sqlx.DB

func MySQLConnect(dbName string) {
	MDB = sqlx.MustConnect("mysql", "andy:hardpassword1@(35.157.14.177:3306)/blockchain")

	//some benchmark should be done here
	MDB.SetMaxOpenConns(300)
	MDB.SetMaxIdleConns(10)
	MDB.SetConnMaxLifetime(10 * time.Second)
}
