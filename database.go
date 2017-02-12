package main

import (
	"database/sql"
	"log"

	_ "github.com/spyrosoft/pq"
)

func connectToDatabase() *sql.DB {
	//TODO: Verify that the connection actually opens
	//Nonexistent databases still open a connection...
	db, err := sql.Open("postgres", "dbname="+siteData.DatabaseName+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
