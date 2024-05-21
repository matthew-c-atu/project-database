/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/matthew-c-atu/project-database/internal/pkg/db"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
)

type RootCfg struct{ *cobra.Command }

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "project-database",
	Short: "A server application which provides access to SQLite music database",
	Long: `### project-database ###
A microservice for providing access to a SQLite music database which is populated
based on the values obtained from a connection to project-audio-streamer.
It provides REST API endpoints for retrieving all items from the database and for 
retrieving items which match a search query, matching on name or genre.

### Endpoints ###
GET /search?name={name}&genre={genre}	- Get songs by seaching on name and/or genre. Leaving these parameters empty returns all songs.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		root := &RootCfg{cmd}

		println("about to serve...")
		root.serve()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose logging")
	rootCmd.PersistentFlags().BoolP("debug", "g", false, "Debug logging")
	rootCmd.PersistentFlags().BoolP("docker", "d", false, "Toggle Docker mode - Uses host.docker.internal instead of localhost")
	rootCmd.PersistentFlags().IntP("port", "p", 9002, "The port on which to run the service")
	rootCmd.PersistentFlags().IntP("fileport", "f", 9001, "The port on which to look for the file server service")
}

func (r *RootCfg) serve() {
	port, err := r.Flags().GetInt("port")
	if err != nil {
		log.Fatal("Couldn't get port flag")
	}
	verbose, _ := r.Flags().GetBool("verbose")
	debug, _ := r.Flags().GetBool("debug")

	ctx := context.Background()
	musicDb, err := r.setupDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer musicDb.Close()

	// TODO: FIX CROSS ORIGIN HEADER ON GIN CONTEXT!!!
	mux := http.NewServeMux()

	mux.Handle("/search", addHeaders(searchSongs(ctx, musicDb)))

	if debug {
		r.printSongsInDB(musicDb)
	}

	slog.Info(fmt.Sprintf("Starting datbase service on port %v\n", port))

	if verbose {
		slog.Info(fmt.Sprintf("Database info: %s", musicDb.String()))
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func (r *RootCfg) printSongsInDB(musicDb *bun.DB) {
	songs, err := db.GetSongsFromRows("songs", musicDb)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %v songs:\n", len(songs))
	for song := range songs {
		marshaled, err := json.Marshal(songs[song])
		if err != nil {
			log.Fatal("failed to marshal JSON: ", err)
		}
		resp := fmt.Sprintf("%s", marshaled)
		println(resp)
	}
}

func (r *RootCfg) setupDb(ctx context.Context) (*bun.DB, error) {
	fileServerPort, err := r.Flags().GetInt("fileport")
	if err != nil {
		return nil, err
	}

	dockerFlag, err := r.Flags().GetBool("docker")
	if err != nil {
		return nil, err
	}

	var fileServerUrl string
	if dockerFlag {
		fileServerUrl = fmt.Sprintf("http://host.docker.internal:%v", fileServerPort)
	} else {
		fileServerUrl = fmt.Sprintf("http://localhost:%v", fileServerPort)
	}

	musicDb, err := db.CreateDB()
	if err != nil {
		log.Fatal(err)
	}

	err = db.CreateTable(ctx, musicDb)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Second)
	err = db.PopulateDatabase(ctx, fileServerUrl, musicDb)
	if err != nil {
		println("failed to populate DB.")
		log.Fatal(err)
	}

	return musicDb, nil
}

func searchSongsGin(ctx context.Context, musicDb *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Header("Access-Control-Allow-Origin", "*")
		// c.Header("Cache-Control", "no-cache, no-store")
		nameQuery := c.Query("name")
		// idQuery := c.Query("id")
		genreQuery := c.Query("genre")
		fmt.Printf("nameQuery: %v", nameQuery)

		var songs []db.Song

		err := musicDb.NewSelect().
			Model(&songs).
			Where("? LIKE ?", bun.Ident("name"), fmt.Sprintf("%%%v%%", nameQuery)).
			Where("? LIKE ?", bun.Ident("genre"), fmt.Sprintf("%%%v%%", genreQuery)).
			Scan(ctx)
		if err != nil {
			slog.Info(err.Error())
		}

		marshaled, err := json.Marshal(songs)
		if err != nil {
			slog.Info(err.Error())
		}
		c.Writer.Write(marshaled)
	}
}

func searchSongs(ctx context.Context, musicDb *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// c.Header("Access-Control-Allow-Origin", "*")
		// c.Header("Cache-Control", "no-cache, no-store")
		nameQuery := r.URL.Query().Get("name")
		// idQuery := c.Query("id")
		genreQuery := r.URL.Query().Get("genre")
		fmt.Printf("nameQuery: %v", nameQuery)

		var songs []db.Song

		err := musicDb.NewSelect().
			Model(&songs).
			Where("? LIKE ?", bun.Ident("name"), fmt.Sprintf("%%%v%%", nameQuery)).
			Where("? LIKE ?", bun.Ident("genre"), fmt.Sprintf("%%%v%%", genreQuery)).
			Scan(ctx)
		if err != nil {
			slog.Info(err.Error())
		}

		marshaled, err := json.Marshal(songs)
		if err != nil {
			slog.Info(err.Error())
		}
		w.Write(marshaled)
	}
}

func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Cache-Control", "no-cache, no-store")
		if h != nil {
			h.ServeHTTP(w, r)
		}
	}
}
