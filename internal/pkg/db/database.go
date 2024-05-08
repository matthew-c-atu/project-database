package db

import (
	"log"
	"os"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Song struct {
	Id    int
	Name  string
	Url   string
	Genre string
}

func CreateDB() *sql.DB {
	os.Remove("./music.db")
	db, err := sql.Open("sqlite3", "./music.db")
	if err != nil {
		log.Fatal("Could not create music database")
	}

	CreateTable(db)
	return db
}

func CreateTable(db *sql.DB) {
	// sql := `CREATE TABLE songs (
	// 	id INTEGER PRIMARY KEY,
	// 	name TEXT NOT NULL,
	// 	url TEXT NOT NULL,
	// 	genre TEXT NOT NULL,
	// );`
	sql := `CREATE TABLE songs (id INTEGER PRIMARY KEY, name TEXT NOT NULL, url TEXT NOT NULL, genre TEXT NOT NULL);`
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal("Failed to create tables in database ", err)
	}
}

func InsertData(song Song, db *sql.DB) {
	statement, err := db.Prepare("INSERT INTO songs (name, url, genre) values (?, ?, ?)")
	if err != nil {
		log.Fatal("Failed to create database InsertData statement ", err)
	}
	_, err = statement.Exec(song.Name, song.Url, song.Genre)
	if err != nil {
		log.Fatal("Failed to exec database InsertData statement ", err)
	}
}

func PrintTable(table string, db *sql.DB) {
	var url string
	if err := db.QueryRow("SELECT (url) FROM songs WHERE id = 1").Scan(&url); err != nil {
		if err == sql.ErrNoRows {
			log.Fatal("ErrNoRows: ", err)
		}
		log.Fatal("other err: ", err)
	}
	println(url)
	return
	// statement, err := db.Prepare("SELECT * FROM songs")
	// if err != nil {
	// 	log.Fatal("Could not prepare print table: ", err)
	// }
	// defer statement.Close()
	// _, err = statement.Exec()
	// if err != nil {
	// 	log.Fatal("Could not exec print table: ", err)
	// }

}

func GetSongsFromRows(table string, db *sql.DB) ([]Song, error) {
	rows, err := db.Query("SELECT * FROM songs", table)
	if err != nil {
		log.Fatal("GetSongsFromRows failed: ", err)
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.Id, &song.Name, &song.Url, &song.Genre); err != nil {
			return songs, err
		}
		songs = append(songs, song)
	}
	if err = rows.Err(); err != nil {
		return songs, err
	}
	return songs, nil
}
