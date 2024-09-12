package database

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func HttpServer(db *sql.DB) *mux.Router {
	rt := mux.NewRouter()
	api := rt.PathPrefix("/v1").Subrouter()
	api.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Command to exit recieved:%v", api.Get("exit"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Shutting down Service"))
		time.Sleep(1000)
		os.Exit(2)
	}).Methods("POST")
	api.HandleFunc("/ping-db", func(w http.ResponseWriter, r *http.Request) {
		if SqlPing(db) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ping successful"))

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Ping failed"))
		}
	}).Methods("GET")
	api.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy"))
	}).Methods("GET")
	log.Printf("Server started on %v:36020", os.Getenv("HTTP_HOST"))
	return api

}
