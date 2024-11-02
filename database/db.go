package database

import (
	"database/sql"
	"log"
	"sync"

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
func updateAPIMap(db *sql.DB) *sync.Map {
	apiKeysMap := sync.Map{}
	rows, err := db.Query("SELECT api_keys, service_name FROM APIKEYS")
	if err != nil {
		log.Printf("Error querying APIKEYS table: %v", err)
		return &apiKeysMap
	}
	defer rows.Close()
	for rows.Next() {
		var apiKey, serviceName string
		if err = rows.Scan(&apiKey, &serviceName); err != nil {
			log.Printf("Error scanning row: %v", err)
			return &apiKeysMap
		}
		apiKeysMap.Store(apiKey, serviceName)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return &apiKeysMap
	}
	return &apiKeysMap
}

func SqlWrite(db *sql.DB, data LogMessage) {
	/*
			To define the way logs are received and passed on to the DB
		1- time
		2- App Name
		3- Log Level
		4- message
		lenght of slice should be data[3]
	*/
	query := "INSERT INTO "

}
