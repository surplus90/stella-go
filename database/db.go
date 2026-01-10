package database

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	
	DB, err = sql.Open("sqlite", "sqlite/stella.db")
	if err != nil {
		panic(err)
	}
}