package main

import (
	// "github.com/gin-gonic/gin"

	"encoding/json"
	"fmt"
	"log"

	"github.com/matthew-c-atu/project-database/internal/pkg/db"
)

func main() {

	data := db.CreateDB()
	s := []db.Song{
		{
			Id:    1,
			Name:  "Binary - Symptome",
			Url:   "http://localhost:8080/hls/Binary%20-%20Symptome.m3u8",
			Genre: "Drum and Bass",
		},
		{
			Id:    2,
			Name:  "Foo - Bar",
			Url:   "http://localhost:8080/hls/Foo%20-%Bar.m3u8",
			Genre: "Drum and Bass",
		},
	}
	for song := range s {
		db.InsertData(s[song], data)
	}
	db.PrintTable("songs", data)
	songs, err := db.GetSongsFromRows("songs", data)
	if err != nil {
		log.Fatal(err)
	}
	for song := range songs {
		// fmt.Printf("%v", songs[song])
		marshaled, err := json.Marshal(songs[song])
		if err != nil {
			log.Fatal("failed to marshal JSON: ", err)
		}
		resp := fmt.Sprintf("%s", marshaled)
		println(resp)
	}
}
