package db

import (
	"database/sql"
	"fmt"
	"japv6/models"
)

// func SelectMetaCardsByIds(table string, group string, ids []int) ([]models.CardMeta, error) {
// 	j, err := json.Marshal(ids)
// 	if err != nil {
// 		return nil, err
// 	}

// 	query := fmt.Sprintf(`
// 		SELECT id, %[1]s_v, %[1]s_sync_v
// 		FROM %[2]ss
// 		WHERE id IN (SELECT value FROM json_each(?))
// 	`, group, table)

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

func SelectCardsSyncRange(table string, group string, from int, to int) ([]models.Card, error) {
	query := fmt.Sprintf(`
		SELECT id, %[1]s_v, %[1]s_sync_v, %[1]s_data
		FROM %[2]ss
        WHERE %[1]s_sync_v BETWEEN ? AND ?
    `, group, table)
	// fmt.Println(query)
	rows, err := conn.Query(query, from, to)
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

// func UpdateCards(cards []models.Card, v int, table string, group string) (re []models.CardMeta, newV int, err error) {
// 	re = make([]models.CardMeta, len(cards))
// 	newV = 0

// 	tx, err := conn.Begin()
// 	if err != nil {
// 		return
// 	}
// 	defer tx.Rollback()

// 	query := fmt.Sprintf(`
// 		UPDATE %ss
// 		SET %[2]s_v = ?, %[2]s_sync_v = ?, %[2]s_data = ?
// 		WHERE id = ?;
// 	`, table, group)
// 	// fmt.Println(query)
// 	stmt, err := tx.Prepare(query)
// 	if err != nil {
// 		return
// 	}
// 	defer stmt.Close()

// 	for i, c := range cards {
// 		v++
// 		c.SyncV = v
// 		re[i] = c.CardMeta
// 		_, err = stmt.Exec(c.V, c.SyncV, c.Data, c.ID)
// 		if err != nil {
// 			return
// 		}
// 	}

// 	tn := fmt.Sprintf("%s_%ss", table, group)
// 	_, err = tx.Exec("UPDATE versions SET val = ? WHERE id = ?", v, tn)
// 	if err != nil {
// 		return
// 	}

// 	return re, v, tx.Commit()
// }

// func selectMetaCardById(tx *sql.Tx, table string, group string, id int) (*models.CardMeta, error) {
func selectMetaCardById(tx *sql.Tx, table string, group string, id any) (*models.CardMeta, error) {
	query := fmt.Sprintf(`
		SELECT %[1]s_v, %[1]s_sync_v
		FROM %[2]s
		WHERE id = ?
	`, group, table)

	var card models.CardMeta
	err := tx.QueryRow(query, id).Scan(&card.V, &card.SyncV)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &card, nil
}

// refactore these thigs with creation of the statement!
func updateCard(tx *sql.Tx, c models.Card, table string, group string) error {
	query := fmt.Sprintf(`
		UPDATE %[1]s
		SET %[2]s_v = ?, %[2]s_sync_v = ?, %[2]s_data = ?
		WHERE id = ?;
	`, table, group)

	_, err := tx.Exec(query, c.V, c.SyncV, c.Data, c.ID)
	return err
}

func createCard(tx *sql.Tx, c models.Card, table string, group string) error {
	query := fmt.Sprintf(`
		INSERT INTO %[1]s (id, %[2]s_v, %[2]s_sync_v, %[2]s_data)
		VALUES (?, ?, ?, ?);
	`, table, group)

	_, err := tx.Exec(query, c.ID, c.V, c.SyncV, c.Data)
	return err
}

func UpsertCards(inputCards []models.Card, v int, tableEntry string, group string) (re []models.CardMeta, newV int, err error) {
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

// func TempFillCards(cards []models.Card, table string, group string) error {
func TempFillCards(cards []models.AnyCard, table string, group string) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := fmt.Sprintf(`
		INSERT INTO %[1]s (id, %[2]s_v, %[2]s_sync_v, %[2]s_data)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
		%[2]s_v = excluded.%[2]s_v,
		%[2]s_sync_v = excluded.%[2]s_sync_v,
		%[2]s_data = excluded.%[2]s_data;
	`, table, group)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range cards {
		_, err := stmt.Exec(c.ID, c.V, c.SyncV, c.Data)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
