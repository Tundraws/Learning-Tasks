package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	ID      int
	Date    string
	Name    string
	Content string
}

var db *sql.DB

func InitDB() error {
	var err error
	// Открываем соединение с базой данных с WAL режимом
	db, err = sql.Open("sqlite3", "file:guestbook.db?_journal=WAL&_fk=true")
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных: %v", err)
	}

	// Устанавливаем ограничения соединений
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Проверяем соединение
	if err = db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("ошибка подключения к базе: %v", err)
	}

	// Создаем таблицу в транзакции
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		name TEXT NOT NULL DEFAULT 'Анонимно',
		content TEXT NOT NULL CHECK(content != ''),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_created_at ON messages (created_at);
	`)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы: %v", err)
	}

	return tx.Commit()
}

func AddMessage(name, content string) error {
	if content == "" {
		return fmt.Errorf("сообщение не может быть пустым")
	}
	if name == "" {
		name = "Анонимно"
	}

	currentTime := time.Now().Format("02.01.2006 15:04")

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"INSERT INTO messages (date, name, content) VALUES (?, ?, ?)",
		currentTime, name, content,
	)
	if err != nil {
		return fmt.Errorf("ошибка вставки сообщения: %v", err)
	}

	return tx.Commit()
}

func GetMessages() ([]Message, error) {
	rows, err := db.Query(`
		SELECT id, date, name, content 
		FROM messages 
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса сообщений: %v", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Date, &m.Name, &m.Content); err != nil {
			return nil, fmt.Errorf("ошибка сканирования сообщения: %v", err)
		}
		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %v", err)
	}

	return messages, nil
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
