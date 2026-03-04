package db

import (
	// "log"
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