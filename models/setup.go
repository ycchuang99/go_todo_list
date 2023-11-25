package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDatabase() {
	database, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/todo_list")

	if err != nil {
		panic(err.Error())
	}

	DB = database
}
