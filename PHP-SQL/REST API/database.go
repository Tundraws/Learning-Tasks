package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "products.db")
	if err != nil {
		return err
	}

	// Создаем таблицу товаров
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        price REAL NOT NULL,
        stock INTEGER NOT NULL
    );`

	_, err = db.Exec(createTableSQL)
	return err
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
