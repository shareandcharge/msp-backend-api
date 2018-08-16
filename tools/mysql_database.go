package tools

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

var MDB *sqlx.DB

func MySQLConnect(dbName string) {
	MDB, err := sqlx.Connect("mysql", "andy:hardpassword1@(18.197.172.83:3306)/"+dbName)
	ErrorCheck(err, "mysql_database.go", true)

	//some benchmark should be done here
	MDB.SetMaxOpenConns(300)
	MDB.SetMaxIdleConns(10)
	MDB.SetConnMaxLifetime(10 * time.Second)
}
