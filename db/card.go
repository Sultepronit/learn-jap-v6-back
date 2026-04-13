package db

import (
	"database/sql"
	"fmt"
	"japv6/models"
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

func createStmts(tx *sql.Tx, table string, group string) (*sql.Stmt, *sql.Stmt, error) {
	updateStmt, err := tx.Prepare(fmt.Sprintf(`
		UPDATE %[1]s
		SET %[2]s_v = ?, %[2]s_sync_v = ?, %[2]s_data = ?
		WHERE id = ?;
	`, table, group))
	if err != nil {
		return nil, nil, err
	}

	createStmt, err := tx.Prepare(fmt.Sprintf(`
		INSERT INTO %[1]s (id, %[2]s_v, %[2]s_sync_v, %[2]s_data)
		VALUES (?, ?, ?, ?);
	`, table, group))
	if err != nil {
		return nil, nil, err
	}

	return updateStmt, createStmt, err
}

// refactor these thigs with creation of the statement!
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

func UpsertCards(inputCards []models.Card, v int, isOutdated bool, tableEntry string, group string) (re []models.CardMeta, newV int, err error) {
	re = make([]models.CardMeta, 0, len(inputCards))
	newV = v

	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	table := tableEntry + "s"

	updateStmt, createStmt, err := createStmts(tx, table, group)
	if err != nil {
		return
	}
	defer updateStmt.Close()
	defer createStmt.Close()

	for _, ic := range inputCards {
		sc, err := selectMetaCardById(tx, table, group, ic.ID)
		if err != nil {
			return nil, 0, err
		}
		fmt.Println("sc:", sc)

		// var action func(*sql.Tx, models.Card, string, string) error
		var stmt *sql.Stmt

		if sc == nil {
			fmt.Println("new card!")
			// action = createCard
			stmt = createStmt
		} else if ic.SyncV == sc.SyncV || (ic.V > sc.V && (!isOutdated || ic.SyncV == -1)) {
			// 1) the client has up to date card
			// 2) this client is more sedulous than the previous one & is not outdated
			// 3) the client is more sedulous & it is the origin of the card
			// action = updateCard
			stmt = updateStmt
		}

		// if action != nil {
		if stmt != nil {
			newV = v + 1
			ic.SyncV = newV
			re = append(re, ic.CardMeta)
			fmt.Println(ic.CardMeta)
			// err = action(tx, ic, table, group)
			_, err := stmt.Exec(ic.ID, ic.V, ic.SyncV, ic.Data)
			if err != nil {
				return nil, 0, err
			}
		}
	}

	if newV > v {
		if err := updateVersion(tx, newV, tableEntry, group); err != nil {
			return nil, 0, err
		}
	}

	return re, newV, tx.Commit()
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
