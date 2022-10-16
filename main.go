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
	for {
		if err := spotifyReq.Auth(); err != nil {
			log.Printf("failed to auth spotify client: %s", err)
			continue
		}
		break
	}

	// define HTTP handlers
	handlers := api.New(spotifyReq)
	r := mux.NewRouter()
	r.Use(allowCORSMiddleware)
	r.HandleFunc("/api/v1/data/playlists/{id}", handlers.PlaylistsHandler).Methods(http.MethodGet)

	// start HTTP server
	err := http.ListenAndServe(":"+strconv.Itoa(conf.Port), r)
	log.Printf("HTTP server shut down: %s", err)
}

func allowCORSMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}
