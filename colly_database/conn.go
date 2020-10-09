package colly_database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err    error
)

func init() {
	// dbusername+":"+dbpassword+"@tcp("+dbhostip+")/"+dbname+"?charset=utf8"
	ip := "127.0.0.1"
	pwd := "123456789"
	username := "root"

	url := username + ":" + pwd + "@tcp(" + ip + ":3307)/tests?charset=utf8mb4"

	dbConn, err = sql.Open("mysql", url)
	if err != nil {
		panic(err.Error())
	}
}
