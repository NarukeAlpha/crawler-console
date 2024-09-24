package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/gorilla/mux"
)

func HttpServer(db *sql.DB, apiKeysMap map[string]string, allowedIPList []string) *mux.Router {

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
		service, exist := apiKeysMap[apiKey]
		if exist {
			if !slices.Contains(allowedIPList, r.RemoteAddr) {
				allowedIPList = append(allowedIPList, r.RemoteAddr)
				addNewIP(db, r.RemoteAddr)
				w.Write([]byte("New IP authorized"))
				w.WriteHeader(http.StatusOK)
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Authorized for service: " + service))
	}).Methods("POST")

	adminapi := rt.PathPrefix("/adminAPI").Subrouter()
	adminapi.HandleFunc("/exit", ApiKeyAuth(apiKeysMap, allowedIPList, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Command to exit recieved from:%v", apiKeysMap[r.Header.Get("ada-api-Key")])
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Shutting down Service"))
		time.Sleep(1000)
		os.Exit(2)
	})).Methods("POST")

	api := rt.PathPrefix("/v2").Subrouter()
	api.HandleFunc("/add-entry", ApiKeyAuth(apiKeysMap, allowedIPList, func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode; err != nil {
			http.Error(w, "Could not decode json request: ", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Entry added"))
	})).Methods("POST")

	log.Printf("Server started on %v:36020", os.Getenv("HTTP_HOST"))
	return rt

}

func ApiKeyAuth(AKM map[string]string, AIL []string, indexedFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("ada-api-Key")
		_, exist := AKM[apiKey]
		if !exist || !slices.Contains(AIL, r.RemoteAddr) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		indexedFunction(w, r)
	}
}
