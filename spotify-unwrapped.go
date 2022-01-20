package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/jemgunay/spotify-unwrapped/api"
	"github.com/jemgunay/spotify-unwrapped/config"
	"github.com/jemgunay/spotify-unwrapped/spotify"
)

func main() {
	conf := config.New()

	spotifyReq := spotify.New(conf.Spotify)
	if err := spotifyReq.Auth(); err != nil {
		log.Printf("failed to auth spotify client: %s", err)
	}

	// define HTTP handlers
	handlers := api.New(spotifyReq)
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/data/{target}/{id}", handlers.DataHandler).Methods(http.MethodGet)

	// start HTTP server
	err := http.ListenAndServe(":"+strconv.Itoa(conf.Port), r)
	log.Printf("HTTP server shut down: %s", err)
}
