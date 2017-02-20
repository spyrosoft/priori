package main

import (
	"database/sql"
	"log"

	_ "github.com/spyrosoft/pq"
)

func connectToDatabase() *sql.DB {

	//TODO: This is a possible landmine.
	//What happens if postgres is restarted?
	//Is the connection intact?
	//Or do all queries break?

	//TODO: Verify that the connection actually opens
	//Nonexistent databases still open a connection...

	db, err := sql.Open("postgres", "dbname="+siteData.DatabaseName+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
