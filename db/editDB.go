package db

import "log"

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Edit() {
	// query := `
	// DROP TABLE IF EXISTS words;
	// CREATE TABLE words (
	// 	id INTEGER PRIMARY KEY,
	// 	card_v INTEGER NOT NULL DEFAULT 0,
	// 	card_sync_v INTEGER NOT NULL DEFAULT 0,
	// 	card_data BLOB,
	// 	prog_v INTEGER NOT NULL DEFAULT 0,
	// 	prog_sync_v INTEGER NOT NULL DEFAULT 0,
	// 	prog_data BLOB
	// ) STRICT;
	// `
	
	query := `
	DROP TABLE IF EXISTS kanjis;
	CREATE TABLE kanjis (
		id TEXT PRIMARY KEY,
		card_v INTEGER NOT NULL DEFAULT 0,
		card_sync_v INTEGER NOT NULL DEFAULT 0,
		card_data BLOB,
		prog_v INTEGER NOT NULL DEFAULT 0,
		prog_sync_v INTEGER NOT NULL DEFAULT 0,
		prog_data BLOB
	) STRICT;
	`

	// query := `
	// DROP TABLE IF EXISTS versions;
	// CREATE TABLE versions (
	// 	id TEXT PRIMARY KEY,
	// 	val INTEGER NOT NULL
	// )
	// `
	r, err := conn.Exec(query)

	// query := `INSERT INTO versions (id, val) VALUES (?, ?)`
	// r, err = conn.Exec(query, "word_cards", 0)
	// r, err := conn.Exec(query, "word_progs", 0)
	// r, err := conn.Exec(query, "kanji_progs", 0)
	// r, err := conn.Exec(query, "kanji_cards", 0)

	handleErr(err)
	log.Println(r.RowsAffected())
}
