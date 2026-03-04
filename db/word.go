package db

import (
	"japv6/models"
	"log"
)

func InsertWordCard(card models.WordCard) error {
	query := `
		INSERT OR REPLACE INTO words (id, card_v, card_sync_v, card_data)
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
		INSERT OR REPLACE INTO words (id, card_v, card_sync_v, card_data)
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

func SelectWordCards() ([]models.WordCard, error) {
	query := `
		SELECT id, card_v, card_sync_v, card_data
		FROM words
	`
	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	re := make([]models.WordCard, 0, 40)
	for rows.Next() {
		var c models.WordCard
		err = rows.Scan(&c.ID, &c.V, &c.SyncV, &c.Data)
		if err != nil {
			return nil, err
		}
		re = append(re, c)
	}

	return re, nil
}

func UpdateWordCards(cards []models.WordCard) ([]models.CardMeta, error) {
	v, err := GetVersion("word_cards")

	tx, err := conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		UPDATE words
		SET card_v = ?, card_sync_v = ?, card_data = ?
		WHERE id = ?;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if err != nil {
		return nil, err
	} 

	re := make([]models.CardMeta, len(cards))
	for i, c := range cards {
		v++
		c.SyncV = v
		re[i] = c.CardMeta
		_, err := stmt.Exec(c.V, c.SyncV, c.Data, c.ID)
		if err != nil {
			return nil, err
		}
	}

	_, err = tx.Exec("UPDATE versions SET val = ? WHERE id = ?", v, "word_cards")
	if err != nil {
		return nil, err
	}

	return re, tx.Commit()
}