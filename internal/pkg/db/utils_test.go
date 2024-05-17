package db_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	dbutils "github.com/matthew-c-atu/project-database/internal/pkg/db"
)

func TestMapSongNamesToSongFiles(t *testing.T) {
	type testCase struct {
		a []string
		b []string
	}
	validCases := []testCase{
		{
			a: []string{"foo", "bar"},
			b: []string{"foo.m3u8", "bar.m3u8"},
		},
		{
			a: []string{"foo", "bar"},
			b: []string{"foo", "bar"},
		},
		{
			a: []string{"", ""},
			b: []string{"", ""},
		},
	}

	invalidCases := []testCase{
		{
			a: []string{"foo", "bar"},
			b: []string{"baz.m3u8", "bar.m3u8"},
		},
		{
			a: []string{"foo", "bar"},
			b: []string{"baz", "bar"},
		},
	}

	t.Run("test valid cases", func(t *testing.T) {
		for _, v := range validCases {
			_, err := dbutils.MapSongNamesToSongFiles(v.a, v.b)
			if err != nil {
				t.Errorf("Invalid test case %s : %s - error: %s", v.a, v.b, err.Error())
			}
		}
	})
	t.Run("test invalid cases", func(t *testing.T) {
		for _, v := range invalidCases {

			songsMap, err := dbutils.MapSongNamesToSongFiles(v.a, v.b)
			if err != nil {
				// detected error - PASS
				fmt.Printf("Successfully caught error: %v {a:%v, b:%v}\n", err.Error(), v.a, v.b)
			}
			// if songsMap is nil we detected an error successfully
			if songsMap != nil {
				t.Errorf("Failed on test case %s : %s", v.a, v.b)
			}
		}
	})
}

func TestMapSongNamesToGenre(t *testing.T) {
	type testCase struct {
		a []string
		b string
	}

	validCase := testCase{
		a: []string{"foo", "bar"},
		b: "foo and bass",
	}

	t.Run("test valid case", func(t *testing.T) {
		mapping, err := dbutils.MapSongNamesToGenre(validCase.a, validCase.b)
		if len(mapping) != len(validCase.a) {
			t.Errorf("Mapping failed - err: %v", err.Error())
		}
	})
}

func TestGetSongStringsFromJsonResponse(t *testing.T) {
	t.Run("valid response", func(t *testing.T) {
		handlerFunc := func(w http.ResponseWriter, r *http.Request) {
			resp := `["Foo - The Test", "Bar - The Unit"]`
			w.Write([]byte(resp))
		}
		server := httptest.NewServer(http.HandlerFunc(handlerFunc))

		songs, err := dbutils.GetSongStringsFromJsonResponse(server.URL, "/")
		if err != nil {
			t.Error(err.Error())
		}
		println("Results of call to GetSongStringsFromJsonResponse:")
		for _, v := range songs {
			println(v)
		}

		if songs[0] != "Foo - The Test" && songs[1] != "Bar - The Unit" && len(songs) != 2 {
			t.Error("Failed to parse list of strings from json response")
		}
	})

	t.Run("empty response", func(t *testing.T) {
		handlerFunc := func(w http.ResponseWriter, r *http.Request) {
			resp := `[]`
			w.Write([]byte(resp))
		}
		server := httptest.NewServer(http.HandlerFunc(handlerFunc))

		songs, err := dbutils.GetSongStringsFromJsonResponse(server.URL, "/")
		if err != nil {
			t.Error(err.Error())
		}
		println("Results of call to GetSongStringsFromJsonResponse:")
		for _, v := range songs {
			println(v)
		}

		if len(songs) != 0 {
			t.Error("Failed to parse empty json response")
		}
	})

	t.Run("invalid url and endpoint", func(t *testing.T) {
		_, err := dbutils.GetSongStringsFromJsonResponse("asdf", "hjkl")
		if err != nil {
			return
		}
		t.Error("Failed to test invalid url and endpoint - err was nil")
	})
}
