package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(192.168.137.3:3306)/yiitest")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	rows, err := db.Query("select * from alpha")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id               int64
			stringIdentifier string
		)
		err := rows.Scan(&id, &stringIdentifier)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("id: %d, string_identifier: %s\n", id, stringIdentifier)
	}

}
