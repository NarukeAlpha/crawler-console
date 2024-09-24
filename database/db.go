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
func updateAPIMap(db *sql.DB) map[string]string {
	apiKeysMap := make(map[string]string)
	rows, err := db.Query("SELECT api_key, service_name FROM APIKEYS")
	if err != nil {
		log.Printf("Error querying APIKEYS table: %v", err)
		return apiKeysMap
	}
	defer rows.Close()
	for rows.Next() {
		var apiKey, serviceName string
		if err = rows.Scan(&apiKey, &serviceName); err != nil {
			log.Printf("Error scanning row: %v", err)
			return apiKeysMap
		}
		apiKeysMap[apiKey] = serviceName
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return apiKeysMap
	}
	return apiKeysMap
}

func allowedIPList(db *sql.DB) []string {
	var ipList []string
	rows, err := db.Query("SELECT ip FROM IP_LIST")
	if err != nil {
		log.Printf("Error querying IP_LIST table: %v", err)
		return ipList
	}
	defer rows.Close()
	for rows.Next() {
		var ip string
		if err = rows.Scan(&ip); err != nil {
			log.Printf("Error scanning row: %v", err)
			return ipList
		}
		ipList = append(ipList, ip)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return ipList
	}
	return ipList
}

func addNewIP(db *sql.DB, ip string) {
	_, err := db.Exec("INSERT INTO IP_LIST(ip) VALUES(?)", ip)
	if err != nil {
		log.Printf("Error inserting new IP: %v", err)
	}

}

func SqlWrite(db *sql.DB, data []any) {
	if len(data) < 3 {
		log.Panicf("Request data shorter than expected: size was %v & expected slice ", len(data))
	}
}
