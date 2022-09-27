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
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS addresses (id INTEGER PRIMARY KEY, address TEXT, privateKeys TEXT, counter INTEGER)")
	if err != nil {
		log.Printf("error while creating table to db")
		return nil, err
	}
	statement.Exec()
	//statement, _ = database.Prepare("INSERT INTO addresses (addres, counter) VALUES (?, ?)")
	//statement.Exec("iriba", "1")
	//rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
	//var id int
	//var firstname string
	//var lastname string
	//for rows.Next() {
	//	rows.Scan(&id, &firstname, &lastname)
	//	fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	//}
	return db, nil
}
