package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/jemgunay/spotify-unwrapped/api"
	"github.com/jemgunay/spotify-unwrapped/config"
	"github.com/jemgunay/spotify-unwrapped/spotify"
)

func main() {
	conf := config.New()
	logger := conf.Logger

	// generate Spotify client auth token
	spotifyReq := spotify.New(logger, conf.Spotify)

	// define HTTP handlers
	handlers := api.New(logger, spotifyReq)
	r := mux.NewRouter()
	r.Use(allowCORSMiddleware)
	r.HandleFunc("/api/v1/playlists/{playlistID}", handlers.PlaylistsHandler).Methods(http.MethodGet)

	// start HTTP server
	logger.Info("starting HTTP server", zap.Int("port", conf.Port))
	err := http.ListenAndServe(":"+strconv.Itoa(conf.Port), r)
	logger.Info("HTTP server shut down", zap.Error(err))
}

func allowCORSMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}
