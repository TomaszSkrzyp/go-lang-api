package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func initDB() *sql.DB {
	connectStr := "user=postgres dbname=todo-rest password=Haslopg123 sslmode=disable"
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
