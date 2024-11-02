package database

import (
	"log"
	"net/http"
	"os"
)

func Main(dbkey struct {
	Url      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}) {
	var key = dbkey.User + ":" + dbkey.Password + "@tcp(" + dbkey.Url + ")/" + dbkey.Database
	log.Printf("Connecting to %v", dbkey.Database)
	var Db = SqlInit(key)
	defer Db.Close()
	var apiKeysMap = updateAPIMap(Db)
	r := HttpServer(Db, apiKeysMap)
	var localIP = os.Getenv("HTTP_HOST") + ":36020"
	err := http.ListenAndServe(localIP, r)
	if err != nil {
		log.Panicf("Error starting server: %v", err)
	}

}
