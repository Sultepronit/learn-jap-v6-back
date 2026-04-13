package db

import (
	"database/sql"
	"fmt"
)

func GetVersion(id string) (int, error) {
    query := "SELECT val FROM versions WHERE id = ?"
    var re int
    
    err := conn.QueryRow(query, id).Scan(&re)
    if err != nil {
        return 0, err
    }

    return re, nil
}

func getVersionTx(tx *sql.Tx, id string) (int, error) {
	query := "SELECT val FROM versions WHERE id = ?"
	var re int

	err := tx.QueryRow(query, id).Scan(&re)
	if err != nil {
		return 0, err
	}

	return re, nil
}

func updateVersion(tx *sql.Tx, newV int, tableEntry string, group string) error {
	tn := fmt.Sprintf("%s_%ss", tableEntry, group)
	_, err := tx.Exec("UPDATE versions SET val = ? WHERE id = ?", newV, tn)
	return err
}