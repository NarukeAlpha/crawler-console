package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	url, isPresent := os.LookupEnv("LS_URL")
	if !isPresent {
		url = "127.0.0.1"
	}
	var host = "http://" + url + ":36020/v1/healthcheck"
	for i := 0; i < 10; i++ {
		req, err := http.Get(host)
		if err != nil {
			log.Println("Failed to ping service: %v", err)
		}
		defer req.Body.Close()
		if req.StatusCode == 200 {
			log.Println("Log Service is healthy")
			os.Exit(0)
		}
		time.Sleep(5 * time.Second)
	}
	log.Println("Failed to connect to log service in 30 seconds, closing with error")
	os.Exit(1)
}
