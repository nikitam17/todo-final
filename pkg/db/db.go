package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// Переменная для использования другими функциями
var TodoDB *sql.DB

// Открывает базу данных и при необходимости создавать таблицу с индексом
func Init(dbFile string) (*sql.DB, error) {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	if install {
		// Создаём файл
		file, err := os.Create(dbFile)
		if err != nil {
			return nil, fmt.Errorf("ошибка при создании файла БД: %w", err)
		}
		defer file.Close()
	}
	// Открываем БД
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к БД: %w", err)
	}
	if install {
		// Создаём таблицу
		sqlQuery := `CREATE TABLE scheduler (
							id INTEGER PRIMARY KEY AUTOINCREMENT,
							date CHAR(8) NOT NULL DEFAULT "",
							title VARCHAR NOT NULL DEFAULT "",
							comment TEXT NOT NULL DEFAULT "",
							repeat VARCHAR(128) NOT NULL DEFAULT ""
							)`
		_, err := db.Exec(sqlQuery)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания таблицы: %w", err)
		}
		// Создаём индекс
		sqlQuery = `CREATE INDEX idx_scheduler_date ON scheduler(date)`
		_, err = db.Exec(sqlQuery)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания индекса: %w", err)
		}
	}
	TodoDB = db
	return db, nil
}
