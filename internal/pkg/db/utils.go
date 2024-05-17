package db

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

func MapSongNamesToSongFiles(songNames, songFiles []string) (map[string]string, error) {
	if len(songNames) != len(songFiles) {
		return nil, errors.New("Could not create mapping: len(songNames) != len(SongFiles)")
	}

	// sort slices in ascending order first
	slices.SortFunc(songNames, func(a, b string) int {
		return cmp.Compare(strings.ToLower(a), strings.ToLower(b))
	})

	slices.SortFunc(songFiles, func(a, b string) int {
		return cmp.Compare(strings.ToLower(a), strings.ToLower(b))
	})

	songsMap := make(map[string]string)

	for k, v := range songNames {
		// if key doesnt exist in map
		if _, ok := songsMap[v]; ok {
			break
		} else {
			if strings.Compare(v, strings.TrimSuffix(songFiles[k], ".m3u8")) != 0 {
				return nil, errors.New("Could not create mapping: songName != songFile")
			}
			songsMap[v] = songFiles[k]
		}
	}
	return songsMap, nil
}

func MapSongNamesToGenre(songNames []string, genre string) (map[string]string, error) {
	// sort slices in ascending order first
	slices.SortFunc(songNames, func(a, b string) int {
		return cmp.Compare(strings.ToLower(a), strings.ToLower(b))
	})

	songsMap := make(map[string]string)

	for _, v := range songNames {
		// if key doesnt exist in map
		if _, ok := songsMap[v]; ok {
			break
		} else {
			songsMap[v] = genre
		}
	}
	return songsMap, nil
}

func GetSongStringsFromJsonResponse(url, endpoint string) ([]string, error) {
	songFilesUrl := fmt.Sprintf("%v%v", url, endpoint)
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

	for _, item := range items {
		println(item)
	}
	return items, nil
}
