package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func getItems(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting items")
	rows, err := db.Query(`
		Select id, name from test
	`)

	if err != nil {
		log.Printf("Could not get items: %v", err)
		writeJSON(w, http.StatusInternalServerError, fmt.Sprintf("Could not get items: %v", err))
		return
	}

	defer rows.Close()
	items := []*item{}
	for rows.Next() {
		i := &item{}
		err = rows.Scan(&i.ID, &i.Name)
		if err != nil {
			log.Printf("Could not get item: %v", err)
			writeJSON(w, http.StatusInternalServerError, fmt.Sprintf("Could not get item: %v", err))
			return
		}
		items = append(items, i)
	}
	writeJSON(w, http.StatusOK, items)
}

func addItem(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding item...")
	body, err := getBodyByteArray(r)
	if err != nil {
		log.Printf("Could not parse body: %v", err)
		writeJSON(w, http.StatusInternalServerError, fmt.Sprintf("Could not parse body: %v", err))
		return
	}
	i := &item{}
	err = json.Unmarshal(body, i)
	if err != nil {
		log.Printf("Could not Unmarshal body: %v", err)
		writeJSON(w, http.StatusInternalServerError, fmt.Sprintf("Could not Unmarshal body: %v", err))
		return
	}

	i.ID = uuid.New().String()
	_, err = db.Exec(`
		insert into test(id, name) values(?,?)
	`, i.ID, i.Name)

	if err != nil {
		log.Printf("Could not add items: %v", err)
		writeJSON(w, http.StatusInternalServerError, fmt.Sprintf("Could not add items: %v", err))
		return
	}
}
