package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	var err error
	// Create the database handle, confirm driver is present
	db, err = sql.Open("mysql", "hiram:hiram@/hiram")
	if err != nil {
		log.Fatalf("Could not connected to database: %v", err)
	}
	defer db.Close()

	// Connect and check the server version
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)

	r := mux.NewRouter()
	r.HandleFunc("/items", getItems).Methods("GET")
	r.HandleFunc("/items", addItem).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    ":7070",
	}

	log.Println("Server started")
	log.Fatal(srv.ListenAndServe())
}

func writeJSON(w http.ResponseWriter, httpStatusCode int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	result, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Could not marshal result: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not marshal result: " + err.Error()))
		return
	}
	w.WriteHeader(httpStatusCode)
	w.Write(result)
}

func getBodyByteArray(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1000000))
	if err != nil {
		log.Printf("Could not parse body: %v", err)
		return nil, err
	}

	err = r.Body.Close()
	if err != nil {
		log.Printf("Could not close body: %v", err)
		return nil, err
	}

	return body, nil
}
