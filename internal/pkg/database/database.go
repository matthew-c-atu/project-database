package db

import (
	"log"
	"os"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Song struct {
	Id    int
	Url   string
	Genre string
}

func CreateDB() {
	os.Remove("./music.db")
	db, err := sql.Open("sqlite3", "./music.db")
	if err != nil {
		log.Fatal("Could not create music database")
	}
	defer db.Close()

	CreateTable(db)
}

func CreateTable(db *sql.DB) {
	sql := `CREATE TABLE songs (
		id INTEGER PRIMARY KEY,
		url TEXT NOT NULL,
		genre TEXT NOT NULL,
		
	);`
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal("Failed to create tables in database")
	}
}

func InsertData(song Song, db *sql.DB) {
	statement, err := db.Prepare("INSERT INTO songs (url, genre) values (?, ?)")
	if err != nil {
		log.Fatal("Failed to create database InsertData statement")
	}
	_, err = statement.Exec(song.Url, song.Genre)
	if err != nil {
		log.Fatal("Failed to exec database InsertData statement")
	}
}
