package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	// _ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

type Song struct {
	bun.BaseModel `bun:"table:songs,alias:s"`
	Id            int `bun:"id,pk,autoincrement"`
	Name          string
	Url           string
	Genre         string
}

func CreateDB() (*bun.DB, error) {
	os.Remove("./music.db")
	sqldb, err := sql.Open(sqliteshim.ShimName, "./music.db")
	sqldb.SetMaxIdleConns(1000)
	sqldb.SetConnMaxLifetime(0)
	if err != nil {
		return nil, errors.New("Could not create music database")
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))
	return db, nil
}

func CreateTable(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().
		Model((*Song)(nil)).
		Exec(ctx)
	if err != nil {
		panic(err)
	}

	// sql := `CREATE TABLE IF NOT EXISTS songs(id INTEGER PRIMARY KEY, name TEXT NOT NULL, url TEXT NOT NULL, genre TEXT NOT NULL);`
	// ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancelfunc()
	// _, err := db.ExecContext(ctx, sql)
	// if err != nil {
	// 	return errors.New("Failed to create tables in database ")
	return nil
}

func InsertData(ctx context.Context, song Song, db *bun.DB) error {
	s := &song
	_, err := db.NewInsert().Model(s).Exec(ctx)
	statement, err := db.Prepare("INSERT INTO songs (name, url, genre) values (?, ?, ?)")
	if err != nil {
		return errors.New("Failed to create database InsertData statement")
	}
	_, err = statement.Exec(song.Name, song.Url, song.Genre)
	if err != nil {
		return errors.New("Failed to exec database InsertData statement")
	}
	return nil
}

func PrintTable(table string, db *bun.DB) error {
	println("printing table")
	var url string
	if err := db.QueryRow("SELECT (url) FROM songs WHERE id = 1").Scan(&url); err != nil {
		return err
	}
	println(url)
	println("printing table")
	return nil
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

func GetSongsFromRows(table string, db *bun.DB) ([]Song, error) {
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

func PopulateDatabase(ctx context.Context, fileServerUrl string, db *bun.DB) error {
	songNames, err := GetSongStringsFromJsonResponse(fileServerUrl, "/songnames")
	if err != nil {
		return err
	}
	println("got songnames")

	songFiles, err := GetSongStringsFromJsonResponse(fileServerUrl, "/songfiles")
	if err != nil {
		return err
	}

	// need to wait for these to come back from server?

	songsMap, err := MapSongNamesToSongFiles(songNames, songFiles)
	if err != nil {
		return err
	}

	genresMap, err := MapSongNamesToGenre(songNames, "Drum and Bass")
	if err != nil {
		return err
	}

	// hlsUrlBase := fmt.Sprintf("%v/music/hls", fileServerUrl)
	var songsCollection []Song
	var i int
	for k, v := range songsMap {
		i++
		fmt.Println(k, "\n", v)
		song := Song{
			// Id:    i,
			Name:  k,
			Url:   fmt.Sprintf("%s/%s", fileServerUrl, strings.ReplaceAll(v, " ", "%20")),
			Genre: genresMap[k],
		}
		println(song.Url)
		songsCollection = append(songsCollection, song)
	}

	_, err = db.NewInsert().Model(&songsCollection).Exec(ctx)
	if err != nil {
		return err
	}
	for _, v := range songsCollection {
		fmt.Println(v.Id)
	}
	// for _, v := range songsColletion {
	// 	InsertData(v, db)
	// }
	//
	return nil
}
