package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func SqlInit(key string) *sql.DB {
	db, err := sql.Open("mysql", key)
	if err != nil {
		log.Panicf("Error connecting to database: %v", err)
	}
	log.Println("Connected to database")
	db.Ping()
	return db

}

func SqlPing(db *sql.DB) bool {
	err := db.Ping()
	if err != nil {
		log.Printf("Error pinging DB: %v", err)
		return false
	}
	log.Println("Ping successful")
	return true
}
