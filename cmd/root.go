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
	"os"

	"github.com/gin-gonic/gin"
	"github.com/matthew-c-atu/project-database/internal/pkg/db"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
)

type RootCfg struct{ *cobra.Command }

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "project-database",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.project-database.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().IntP("port", "p", 9002, "The port on which to run the service")
}

func (r *RootCfg) serve() {

	port, err := r.Flags().GetInt("port")
	if err != nil {
		log.Fatal("Couldn't get port flag")
	}
	ctx := context.Background()

	// next step: implement REST API for front end to interact with
	musicDb, err := r.setupDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer musicDb.Close()

	router := gin.Default()
	// router.GET("/songs", getSongs(ctx, musicDb))
	router.GET("/search", searchSongs(ctx, musicDb))
	router.Run(fmt.Sprintf("localhost:%v", port))

	// err = db.PrintTable("songs", musicDb)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// endpoint: GET /songs
	songs, err := db.GetSongsFromRows("songs", musicDb)
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

func (r *RootCfg) setupDb(ctx context.Context) (*bun.DB, error) {
	fileServerPort := 9001
	fileServerUrl := fmt.Sprintf("http://localhost:%v", fileServerPort)

	musicDb, err := db.CreateDB()
	if err != nil {
		log.Fatal(err)
	}

	err = db.CreateTable(ctx, musicDb)
	if err != nil {
		log.Fatal(err)
	}
	err = db.PopulateDatabase(ctx, fileServerUrl, musicDb)
	if err != nil {
		log.Fatal(err)
	}

	return musicDb, nil
}

// func getSongs(ctx context.Context, musicDb *bun.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// nameQuery := c.Query("name")
// 		// idQuery := c.Query("id")
// 		// genreQuery := c.Query("genre")
//
// 		var songs []db.Song
// 		err := musicDb.NewSelect().Model(&songs).Scan(ctx)
// 		if err != nil {
// 			slog.Info(err.Error())
// 		}
//
// 		marshaled, err := json.Marshal(songs)
//
// 		if err != nil {
// 			slog.Info(err.Error())
// 		}
// 		c.Writer.Write(marshaled)
// 	}
// }

func searchSongs(ctx context.Context, musicDb *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
