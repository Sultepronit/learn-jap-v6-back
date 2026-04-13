package db

import (
	"fmt"
	"japv6/models"
)

func DeleteWords(ids []int) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	vIds := []string{"word_cards", "word_progs"}
	vs := make([]int, 0, 2)
	for _, id := range vIds {
		v, err := getVersionTx(tx, id)
		if err != nil {
			return err
		}
		vs = append(vs, v)
	}

	fmt.Println(vs)

	stmt, err := tx.Prepare(`
		UPDATE words 
		SET card_v = -100, card_sync_v = ?, card_data = X'7b7d',
			prog_v = -100, prog_sync_v = ?, prog_data = X'7b7d'
		WHERE id = ?;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, id := range ids {
		vs[0]++
		vs[1]++
		_, err := stmt.Exec(vs[0], vs[1], id)
		if err != nil {
			return err
		}
	}

	for i, id := range vIds {
		_, err = tx.Exec("UPDATE versions SET val = ? WHERE id = ?", vs[i], id)
		if err != nil {
			return err
		}
	}

	// return nil
	return tx.Commit()
}

// func InsertWordCard(card models.Card) error {
// 	query := `
// 		INSERT OR REPLACE INTO words (id, card_v, card_sync_v, card_data)
// 		VALUES (?, ?, ?, ?)
// 	`
// 	r, err := conn.Exec(query, card.ID, card.V, card.SyncV, card.Data)
// 	log.Println(r.RowsAffected())
// 	return err
// }

// func FillWordCards(cards []models.Card) error {
// 	tx, err := conn.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	stmt, err := tx.Prepare(`
// 		INSERT OR REPLACE INTO words (id, card_v, card_sync_v, card_data)
// 		VALUES (?, ?, ?, ?)
// 	`)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	for _, c := range cards {
// 		_, err := stmt.Exec(c.ID, c.V, c.SyncV, c.Data)
// 		if err != nil {
// 			return err
// 		}
// 		// log.Println(r.LastInsertId())
// 	}

// 	return tx.Commit()
// }

// still may be used for the /test request
func SelectWordCards() ([]models.Card, error) {
	query := `
		SELECT id, card_v, card_sync_v, card_data
		FROM words
	`
	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	re := make([]models.Card, 0, 40)
	for rows.Next() {
		var c models.Card
		err = rows.Scan(&c.ID, &c.V, &c.SyncV, &c.Data)
		if err != nil {
			return nil, err
		}
		re = append(re, c)
	}

	return re, nil
}

// func UpdateWordCards(cards []models.Card) ([]models.CardMeta, error) {
// 	v, err := GetVersion("word_cards")
// 	if err != nil {
// 		return nil, err
// 	}

// 	tx, err := conn.Begin()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback()

// 	stmt, err := tx.Prepare(`
// 		UPDATE words
// 		SET card_v = ?, card_sync_v = ?, card_data = ?
// 		WHERE id = ?;
// 	`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer stmt.Close()

// 	re := make([]models.CardMeta, len(cards))
// 	for i, c := range cards {
// 		v++
// 		c.SyncV = v
// 		re[i] = c.CardMeta
// 		_, err := stmt.Exec(c.V, c.SyncV, c.Data, c.ID)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	_, err = tx.Exec("UPDATE versions SET val = ? WHERE id = ?", v, "word_cards")
// 	if err != nil {
// 		return nil, err
// 	}

// 	return re, tx.Commit()
// }
