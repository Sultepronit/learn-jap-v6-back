package db

import (
	"japv6/models"
	"log"
)

func InsertWordCard(card models.WordCard) error {
	query := `
		INSERT INTO words (id, card_v, card_sync_v, card_data)
		VALUES (?, ?, ?, ?)
	`
	r, err := conn.Exec(query, card.ID, card.V, card.SyncV, card.Data)
	log.Println(r.RowsAffected())
	return err
}

func FillWordCards(cards []models.WordCard) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO words (id, card_v, card_sync_v, card_data)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range cards {
		_, err := stmt.Exec(c.ID, c.V, c.SyncV, c.Data)
		if err != nil {
			return err
		}
		// log.Println(r.LastInsertId())
	}

	return tx.Commit()
}

// GUESS
func SelectVoices(isMale bool) ([]models.Voice, error) {
	query := `
		SELECT name, code_name, rate, rating, comment
		FROM voices
		WHERE excluded = false
			AND is_male = ?
	`
	rows, err := conn.Query(query, isMale)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	re := make([]models.Voice, 0, 40)
	for rows.Next() {
		var v models.Voice
		err = rows.Scan(&v.Name, &v.CodeName, &v.Rate, &v.Rating, &v.Comment)
		if err != nil {
			return nil, err
		}
		re = append(re, v)
	}

	return re, nil
}
