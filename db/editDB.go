package db

import "log"

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Edit() {
	// createTable := `
	// DROP TABLE IF EXISTS words;
	// CREATE TABLE words (
	// 	id INTEGER PRIMARY KEY,
	// 	card_v INTEGER NOT NULL,
	// 	card_sync_v INTEGER NOT NULL,
	// 	card_data BLOB NOT NULL
	// )
	// `

	query := `
	DROP TABLE IF EXISTS versions;
	CREATE TABLE versions (
		id TEXT PRIMARY KEY,
		val INTEGER NOT NULL
	)
	`
	r, err := conn.Exec(query)

	query = `INSERT INTO versions (id, val) VALUES (?, ?)`
	r, err = conn.Exec(query, "word_cards", 0)

	handleErr(err)
	log.Println(r.RowsAffected())
}
