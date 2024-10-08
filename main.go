package main

import (
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/narukealpha/crawler-console/database"
)

func AssertErrorToNil(message string, err error) {
	if err != nil {
		log.Panicf(message, err)
	}
}

type SetUp struct {
	Completed bool `json:"completed"`
}

type DbKey struct {
	Url      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Data struct {
	SetUp SetUp `json:"setUp"`
	DbKey DbKey `json:"dbKey"`
}

var data = Data{
	SetUp: SetUp{
		Completed: false,
	},
	DbKey: DbKey{
		Url:      "127.0.0.1:3306",
		User:     "",
		Password: "",
		Database: "",
	},
}

var mw io.Writer

func init() {
	log.Println("Initializing")
	user, isPresent := os.LookupEnv("DB_USER")
	if isPresent {
		data.DbKey.User = user
		data.DbKey.Password = os.Getenv("DB_PASSWORD")
		data.DbKey.Database = os.Getenv("DB_DATABASE")
	} else {
		log.Fatalf("Env variables not set properly")
	}

	_, err := os.Stat("log.txt")
	if os.IsNotExist(err) {
		_, err = os.Create("log.txt")
		if err != nil {
			log.Fatal(err)
		}
	}
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	mw = io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.Println("Started successfully")

}

func main() {

	go database.Main(data.DbKey)
	log.Println("Started DB process")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

}
