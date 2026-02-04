package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type link struct {
	ID        int64
	URL       string
	CreatedAt string
}

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS links (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		url        TEXT    NOT NULL,
		created_at TEXT    NOT NULL DEFAULT (datetime('now'))
	)`); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func insertLink(db *sql.DB, url string) (int64, error) {
	res, err := db.Exec("INSERT INTO links (url) VALUES (?)", url)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func deleteLink(db *sql.DB, id int64) error {
	res, err := db.Exec("DELETE FROM links WHERE id = ?", id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("link not found")
	}
	return nil
}

func getURL(db *sql.DB, id int64) (string, error) {
	var url string
	err := db.QueryRow("SELECT url FROM links WHERE id = ?", id).Scan(&url)
	if err != nil {
		return "", fmt.Errorf("link not found")
	}
	return url, nil
}

func listLinks(db *sql.DB) ([]link, error) {
	rows, err := db.Query("SELECT id, url, created_at FROM links ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var links []link
	for rows.Next() {
		var l link
		if err := rows.Scan(&l.ID, &l.URL, &l.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return links, rows.Err()
}
