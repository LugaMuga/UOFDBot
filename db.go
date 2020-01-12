package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func InitDb() {
	var err error
	DB, err = sql.Open(`sqlite3`, config.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	Migrate()
}

func Migrate() {
	migrationsStmt := `
	CREATE TABLE IF NOT EXISTS chat_user (
	    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER,
		user_id INTEGER,
		username TEXT,
		user_first_name TEXT,
		user_last_name TEXT,
		enabled boolean,
		pidor_score INTEGER DEFAULT 0,
		pidor_last_timestamp INTEGER DEFAULT 0,
		hero_score INTEGER DEFAULT 0,
		hero_last_timestamp INTEGER DEFAULT 0                         
	);
	`
	_, err := DB.Exec(migrationsStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, migrationsStmt)
		return
	}
}

func CloseDb() {
	DB.Close()
}
