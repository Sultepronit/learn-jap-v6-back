package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dgraph-io/badger/v4"
)

type Card0 struct {
	ID    int             `json:"id"`
	V     int             `json:"v"`
	SyncV int             `json:"syncV"`
	Data  json.RawMessage `json:"data"`
}

func save() {

}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func syncHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	var cards []Card0
	if err := json.NewDecoder(r.Body).Decode(&cards); err != nil {
		log.Println(err)
		http.Error(w, "Bad JSON", 400)
		return
	}
	// log.Println(cards)
	// for i, c := range cards {
	// 	log.Println(i, c)
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(cards); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	db, err := badger.Open(badger.DefaultOptions("badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// mux := http.NewServeMux()

	// mux.HandleFunc("POST /sync", func(w http.ResponseWriter, r *http.Request) {
	http.HandleFunc("POST /sync", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		var cards []Card0
		if err := json.NewDecoder(r.Body).Decode(&cards); err != nil {
			log.Println(err)
			http.Error(w, "Bad JSON", 400)
			return
		}
		// log.Println(cards)
		// for i, c := range cards {
		// 	log.Println(i, c)
		// }

		err := db.Update(func(txn *badger.Txn) error {
			for _, c := range cards {
				k := []byte("wc:" + strconv.Itoa(c.ID))
				v, _ := json.Marshal(c)
				err := txn.Set(k, v)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(cards); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/dump", func(w http.ResponseWriter, r *http.Request) {
		db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				item.Value(func(v []byte) error {
					fmt.Fprintf(w, "Key: %s\nValue: %s\n\n", item.Key(), v)
					return nil
				})
			}
			return nil
		})
	})
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println(r.URL.Path, "*")
	// })

	log.Println("Server starts on :8080")
	// log.Fatal(http.ListenAndServe(":8080", mux))
	log.Fatal(http.ListenAndServe(":8080", nil))
	// log.Fatal(http.ListenAndServe(":8080", cors(http.HandlerFunc(syncHandler))))
	// handler := cors(http.HandlerFunc(syncHandler))
	// log.Fatal(http.ListenAndServe(":8080", handler))
}
