package db

import (
	"japv6/models"
    "fmt"
)

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

func UpdateCards(cards []models.Card, v int, table string, group string) (re []models.CardMeta, newV int, err error) {
    re = make([]models.CardMeta, len(cards))
    newV = 0

	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()


    query := fmt.Sprintf(`
		UPDATE %ss
		SET %[2]s_v = ?, %[2]s_sync_v = ?, %[2]s_data = ?
		WHERE id = ?;
	`, table, group)
    // fmt.Println(query)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	for i, c := range cards {
		v++
		c.SyncV = v
		re[i] = c.CardMeta
		_, err = stmt.Exec(c.V, c.SyncV, c.Data, c.ID)
		if err != nil {
			return
		}
	}

    tn := fmt.Sprintf("%s_%ss", table, group)
	_, err = tx.Exec("UPDATE versions SET val = ? WHERE id = ?", v, tn)
	if err != nil {
		return
	}

	return re, v, tx.Commit()
}

func TempFillCards(cards []models.Card, table string, group string) error {
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
