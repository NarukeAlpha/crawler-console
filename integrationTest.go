package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/narukealpha/crawler-console/database"
)

func main() {
	url, isPresent := os.LookupEnv("LS_URL")
	var password, dbname string
	if !isPresent {
		url = "127.0.0.1"
		password = os.Getenv("LS_PASSWORD")
		dbname = os.Getenv("LS_DB")
	} else {
		password = "testing"
		dbname = "main"
	}
	var host = "http://" + url + ":36020/v1/healthcheck"
	var pingerr error
	var req *http.Response
	for i := 0; i < 10; i++ {
		req, pingerr = http.Get(host)
		if pingerr != nil {
			log.Println("Failed to ping service: %v", pingerr)
		}
		if req.StatusCode == 200 {
			err := req.Body.Close()
			if req.Body.Close(); err != nil {
				log.Panicf("Failed to close response body: %v", err)
			}
			break
		}
		req.Body.Close()
		time.Sleep(5 * time.Second)
	}
	if pingerr == nil {
		//Test #1 testing the API table schema2
		APIkeySchemaTest(password, url, dbname)

	} else {
		log.Println("Failed to connect to log service in 30 seconds, closing with error \nNo tests were ran")
		os.Exit(1)
	}

}

func APIkeySchemaTest(password string, url string, dbname string) {
	key := "integrationTest:" + password + "@tcp(" + url + ")/" + dbname
	db := database.SqlInit(key)
	var aktable = "CREATE TABLE APIKEYS (apiKey VARCHAR(50) PRIMARY KEY ,serviceName VARCHAR(50) NOT NULL)"

	defer db.Close()
	var sqlerr *mysql.MySQLError
	log.Println("Staring API table test")
	time.Sleep(10 * time.Second)

	rows, err := db.Query("SELECT api_keys, service_name FROM APIKEYS")
	if err != nil {
		if errors.As(err, &sqlerr) && sqlerr.Number == 1146 {
			_, err = db.Exec(aktable)
			if err != nil {
				log.Fatalf("No table found and unable to create the table: " + err.Error())
			}
		} else if errors.As(err, &sqlerr) && sqlerr.Number == 1054 {
			_, err = db.Exec("DROP TABLE APIKEYS")
			if err != nil {
				log.Fatalf("Failed to drop table" + err.Error())
			}
			_, err = db.Exec(aktable)
			if err != nil {
				log.Fatalf("Unable to create the table after dropping" + err.Error())
			}
		} else {
			log.Fatalf("Failed API Keys:" + err.Error())
		}
	}

	var mp sync.Map
	for rows.Next() {
		var apiKey, serviceName string
		if err = rows.Scan(&apiKey, &serviceName); err != nil {
			log.Fatalf("Error scanning row: %v \n\n Something went wrong after schema matching", err)
		}
		mp.Store(apiKey, serviceName)
	}
	rows.Close()
	lng := 0
	mp.Range(func(key, value any) bool { lng++; return true })
	if lng <= 0 {
		//add a test entry and re read rows, then load into mp
		query := "INSERT INTO APIKEYS (apiKey, serviceName) VALUES ('testKey', 'testService')"
		_, err = db.Exec(query)
		if err != nil {
			log.Fatalf("Failed to add test entry: %v", err)
		}
		rows, err = db.Query("SELECT api_keys, service_name FROM APIKEYS")
		if err != nil {
			log.Fatalf("Failed to re-read rows: %v", err)
		}
		for rows.Next() {
			var apiKey, serviceName string
			if err = rows.Scan(&apiKey, &serviceName); err != nil {
				log.Fatalf("Error scanning row: %v \n\n something went wrong after map comparison", err)
			}
			mp.Store(apiKey, serviceName)
		}
		rows.Close()
		lng = 0
		mp.Range(func(key, value any) bool { lng++; return true })
		if lng <= 0 {
			log.Fatalf("Failed to add test entry, something went wrong with the database")
		}
	}
	mp.Clear()
	log.Println("Database API table test passed")
}
