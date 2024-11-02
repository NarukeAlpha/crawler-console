package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type LogMessage struct {
	time     string `json:time`
	app      string `json:app`
	loglevel string `json:loglevel`
	msg      string `json:msg`
}

func HttpServer(db *sql.DB, apiKeysMap *sync.Map) *mux.Router {

	rt := mux.NewRouter()
	health := rt.PathPrefix("/v1").Subrouter()
	health.HandleFunc("/ping-db", func(w http.ResponseWriter, r *http.Request) {
		if SqlPing(db) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ping successful"))

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Ping failed"))
		}
	}).Methods("GET")
	health.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy"))
	}).Methods("GET")

	auth := rt.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/verifyKey", func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("ada-api-Key")
		service, exist := apiKeysMap.Load(apiKey)
		if exist {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Authorized for service: " + service.(string)))
		}

	}).Methods("POST")

	adminapi := rt.PathPrefix("/adminAPI").Subrouter()
	adminapi.HandleFunc("/exit", ApiKeyAuth(apiKeysMap, func(w http.ResponseWriter, r *http.Request) {
		key, _ := apiKeysMap.Load(r.Header.Get("ada-api-Key"))
		log.Printf("Command to exit recieved from:%v", key.(string))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Shutting down Service"))
		time.Sleep(1000)
		os.Exit(2)
	})).Methods("POST")

	api := rt.PathPrefix("/v2").Subrouter()
	api.HandleFunc("/add-entry", ApiKeyAuth(apiKeysMap, func(w http.ResponseWriter, r *http.Request) {
		var msg LogMessage
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "Could not decode json request: ", http.StatusBadRequest)
			return
		}
		SqlWrite(db, msg)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Entry added"))
	})).Methods("POST")

	log.Printf("Server started on %v:36020", os.Getenv("HTTP_HOST"))
	return rt

}

func ApiKeyAuth(AKM *sync.Map, indexedFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("ada-api-Key")
		_, exist := AKM.Load(apiKey)
		if exist {
			indexedFunction(w, r)
		} else {
			http.Error(w, "Unauthorized Key", http.StatusUnauthorized)
			return
		}
	}
}
