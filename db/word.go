package db

import (
	"japv6/models"
	"log"
)

func DeleteWords(ids []int) (re []models.CardMeta, newV int, err error) {
	re = make([]models.CardMeta, 0, len(inputCards))
	newV = 0

	isFresh := true

	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	table := tableEntry + "s"
	for _, ic := range inputCards {
		sc, err := selectMetaCardById(tx, table, group, ic.ID)
		if err != nil {
			return nil, 0, err
		}
		fmt.Println("sc:", sc)

		var action func(*sql.Tx, models.Card, string, string) error

		if sc == nil {
			fmt.Println("new card!")
			action = createCard
		} else if ic.SyncV == sc.SyncV || (ic.V > sc.V && isFresh) || sc.SyncV < 1 {
			action = updateCard
		}

		if action != nil {
			v++
			ic.SyncV = v
			re = append(re, ic.CardMeta)
			fmt.Println(ic.CardMeta)
			err = action(tx, ic, table, group)
			if err != nil {
				return nil, 0, err
			}
		}
	}

	tn := fmt.Sprintf("%s_%ss", tableEntry, group)
	_, err = tx.Exec("UPDATE versions SET val = ? WHERE id = ?", v, tn)
	if err != nil {
		return
	}

	return re, v, tx.Commit()
}

func InsertWordCard(card models.Card) error {
	query := `
		INSERT OR REPLACE INTO words (id, card_v, card_sync_v, card_data)
		VALUES (?, ?, ?, ?)
	`
	r, err := conn.Exec(query, card.ID, card.V, card.SyncV, card.Data)
	log.Println(r.RowsAffected())
	return err
}

func FillWordCards(cards []models.Card) error {
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

// func SelectMetaCardsByIds(ids []int) ([]models.CardMeta, error) {
// 	j, err := json.Marshal(ids)
// 	if err != nil {
// 		return nil, err
// 	}

// 	query := `
// 		SELECT id, card_v, card_sync_v
// 		FROM words
// 		WHERE id IN (SELECT value FROM json_each(?))`

// 	rows, err := conn.Query(query, j)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	re := make([]models.CardMeta, 0, 10)
// 	for rows.Next() {
// 		var c models.CardMeta
// 		err = rows.Scan(&c.ID, &c.V, &c.SyncV)
// 		if err != nil {
// 			return nil, err
// 		}
// 		re = append(re, c)
// 	}

// 	return re, nil
// }

func UpdateWordCards(cards []models.Card) ([]models.CardMeta, error) {
	v, err := GetVersion("word_cards")
	if err != nil {
		return nil, err
	}

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
