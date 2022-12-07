package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
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

	// caching middleware
	cacheMiddleware, err := newCacheMiddleware()
	if err != nil {
		logger.Fatal("failed to create HTTP cache", zap.Error(err))
		return
	}

	// define HTTP handlers
	handlers := api.New(logger, spotifyReq)
	r := mux.NewRouter()
	r.Use(allowCORSMiddleware, cacheMiddleware)
	r.HandleFunc("/api/v1/playlists/{playlistID}", handlers.PlaylistsHandler).Methods(http.MethodGet)

	// start HTTP server
	logger.Info("starting HTTP server", zap.Int("port", conf.Port))
	err = http.ListenAndServe(":"+strconv.Itoa(conf.Port), r)
	logger.Info("HTTP server shut down", zap.Error(err))
}

func allowCORSMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}

func newCacheMiddleware() (func(h http.Handler) http.Handler, error) {
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(1000),
	)
	if err != nil {
		return nil, err
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(10*time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		return nil, err
	}
	return cacheClient.Middleware, nil
}
