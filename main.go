package main

import (
	// "github.com/gin-gonic/gin"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/matthew-c-atu/project-database/internal/pkg/db"
)

func main() {

	// next step: implement REST API for front end to interact with

	data := db.CreateDB()
	// songFilesUrl := "http://localhost:9001/songfiles"
	hlsUrlBase := "http://localhost:9001/hls/"
	// songFilesReq, err := http.NewRequest(http.MethodGet, songFilesUrl, nil)
	songNames, err := getFromJsonResponse("songnames")
	if err != nil {
		log.Fatal(err)
	}

	songFiles, err := getFromJsonResponse("songfiles")
	if err != nil {
		log.Fatal(err)
	}

	songsMap := make(map[string]string)

	for _, v := range songNames {
		println(v)
	}

	for _, v := range songFiles {
		println(v)
	}

	for k, v := range songNames {
		// if key doesnt exist in map
		if _, ok := songsMap[v]; ok {
			break
			// println(songNames[n])
			// println(songFiles[f])db
		} else {
			songsMap[v] = songFiles[k]
		}
	}
	var songsColletion []db.Song
	var i int
	for k, v := range songsMap {
		i++
		// fmt.Println(k, "\n", v)
		songsColletion = append(songsColletion, db.Song{
			Id:    i,
			Name:  k,
			Url:   fmt.Sprintf("%s%s", hlsUrlBase, strings.ReplaceAll(v, " ", "%20")),
			Genre: "Drum and Bass",
		})
	}

	// sfJson := json.Unmarshal(songFilesResp.Body, songFiles)
	// sfJson := songFilesResp.Body()

	// s := []db.Song{
	// 	{
	// 		Id:    1,
	// 		Name:  "Binary - Symptome",
	// 		Url:   "http://localhost:8080/hls/Binary%20-%20Symptome.m3u8",
	// 		Genre: "Drum and Bass",
	// 	},
	// 	{
	// 		Id:    2,
	// 		Name:  "Foo - Bar",
	// 		Url:   "http://localhost:8080/hls/Foo%20-%Bar.m3u8",
	// 		Genre: "Drum and Bass",
	// 	},
	// }
	for _, v := range songsColletion {
		db.InsertData(v, data)
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

func getFromJsonResponse(endpoint string) ([]string, error) {
	songFilesUrl := fmt.Sprintf("http://localhost:9001/%v", endpoint)
	// songNamesUrl := "localhost:9001/songnames"
	// songFilesReq, err := http.NewRequest(http.MethodGet, songFilesUrl, nil)
	resp, err := http.Get(songFilesUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	var items []string
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&items)

	return items, nil
}
