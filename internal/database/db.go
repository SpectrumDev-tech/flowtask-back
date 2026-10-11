// Package database
package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		log.Fatalf("Error open database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Database is not responding")
	}

	createTableQuery := `
			CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			energy TEXT DEFAULT 'medium',
			category TEXT DEFAULT 'code',
			completed BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Error creating table")
	}

	createCategoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		position INTEGER DEFAULT 0
	);
	`
	_, err = db.Exec(createCategoriesTable)
	if err != nil {
		log.Fatalf("Error creating categories table: %v", err)
	}

	// Добавим дефолтные категории, если их еще нет
	var count int
	db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
	if count == 0 {
		defaults := []string{"code", "study", "life", "work"}
		for _, name := range defaults {
			db.Exec("INSERT OR IGNORE INTO categories (name) VALUES (?)", name)
		}
	}

	log.Println("The database has been successfully initialized.")
	return db
}
