package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func ConnectionToDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database/address.db")
	if err != nil {
		log.Printf("error while getting connection to db")
		return nil, err
	}
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS addresses (id INTEGER PRIMARY KEY, address TEXT, counter INTEGER)")
	if err != nil {
		log.Printf("error while creating table to db")
		return nil, err
	}
	statement.Exec()

	return db, nil
}

func DeleteRowInDB() {
	db, err := sql.Open("sqlite3", "./database/address2.db")
	if err != nil {
		log.Printf("error while getting connection to db")
	}
	statement, err := db.Prepare("DELETE FROM addresses2 where id = 10")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec("DELETE FROM addresses2 where id = 10")
}
