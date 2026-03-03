package db

import (
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

var conn *sql.DB

func Open() error {
	var err error
	conn, err = sql.Open("sqlite", "db.sqlite?_journal=WAL&_sync=NORMAL")
	if err != nil {
		return err
	}

	// conn.Exec("PRAGMA journal_mode=WAL;")
	// conn.Exec("PRAGMA synchronous=NORMAL;")
	conn.SetMaxOpenConns(1)
	conn.SetConnMaxLifetime(time.Hour)

	var v string
	conn.QueryRow("SELECT sqlite_version()").Scan(&v)
	log.Println("DB opened! SQLite:", v)

	return nil
}
