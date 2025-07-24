package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init(url string) error {
	var err error
	db, err = sql.Open("postgres", url)
	if err != nil {
		return err
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
