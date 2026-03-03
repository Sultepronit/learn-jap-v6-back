package db

import "log"

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Edit() {
	createTable := `
	DROP TABLE IF EXISTS words;
	CREATE TABLE words (
		id INTEGER PRIMARY KEY,
		card_v INTEGER NOT NULL,
		card_sync_v INTEGER NOT NULL,
		card_data BLOB NOT NULL
	)
	`

	r, err := conn.Exec(createTable)
	handleErr(err)
	log.Println(r.RowsAffected())
}
