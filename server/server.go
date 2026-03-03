package server

import (
	"encoding/json"
	"japv6/db"
	"japv6/models"
	"log"
	"net/http"
)

func Start() {
	http.HandleFunc("POST /sync", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		var cards []models.WordCard
		if err := json.NewDecoder(r.Body).Decode(&cards); err != nil {
			log.Println(err)
			http.Error(w, "Bad JSON", 400)
			return
		}
		// log.Println(cards)
		// for i, c := range cards {
		// 	log.Println(i, c)
		// }
		err := db.FillWordCards(cards)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(cards); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server starts on: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
