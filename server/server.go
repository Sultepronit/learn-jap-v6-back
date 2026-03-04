package server

import (
	"encoding/json"
	"japv6/db"
	"japv6/models"
	"japv6/sync"
	"log"
	"net/http"
)

func Start() {
	http.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		log.Println(r.URL.Path)
		var cards []models.Card
		if err := json.NewDecoder(r.Body).Decode(&cards); err != nil {
			log.Println("Bad JSON:", err)
			http.Error(w, "Bad JSON", 400)
			return
		}

		err := db.FillWordCards(cards)
		if err != nil {
			log.Fatal(err)
		}

		// w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(cards); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("POST /sync", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		log.Println(r.URL.Path)

		var cards []models.Card
		if err := json.NewDecoder(r.Body).Decode(&cards); err != nil {
			log.Println("Bad JSON:", err)
			http.Error(w, "Bad JSON", 400)
			return
		}

		// sync.SyncWordCards(cards)
		// re, err := db.UpdateWordCards(cards)
		re, err := sync.SyncWordCards(cards)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(re); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		log.Println(r.URL.Path)

		cards, err := db.SelectWordCards()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// log.Println(cards)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(cards); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server starts on: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
