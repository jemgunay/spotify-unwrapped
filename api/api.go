package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jemgunay/spotify-unwrapped/spotify"
)

type API struct {
	spotifyReq spotify.Requester
}

func New(spotifyReq spotify.Requester) API {
	return API{
		spotifyReq: spotifyReq,
	}
}

type statsGroup struct {
	Min     float64 `json:"min"`
	minID   string
	MinName string  `json:"min_name"`
	Max     float64 `json:"max"`
	maxID   string
	MaxName string `json:"max_name"`
	sum     float64
	count   float64
	Mean    float64 `json:"avg"`
}

func (s *statsGroup) push(id string, val float64) {
	switch {
	case s.minID == "":
		s.Min = val
		s.minID = id
		s.Max = val
		s.maxID = id
	case val > s.Max:
		s.Max = val
		s.maxID = id
	case val < s.Min:
		s.Min = val
		s.minID = id
	}

	s.sum += val
	s.count++
}

func (s *statsGroup) finish(lookup map[string]spotify.TrackDetail) {
	s.MinName = lookup[s.minID].Name
	s.MaxName = lookup[s.maxID].Name
	s.Mean = s.sum / s.count
}

const (
	playlistsTarget = "playlists"
)

// Data handler provides data used to drive visualisations.
func (a API) DataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	switch vars["target"] {
	case playlistsTarget:
		playlistData, err := a.spotifyReq.GetPlaylist(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to fetch playlist data: %s", err)
			return
		}

		var popularity statsGroup
		trackIDsList := make([]string, 0, len(playlistData.Tracks.TrackItems))
		trackIDLookup := make(map[string]spotify.TrackDetail, len(playlistData.Tracks.TrackItems))
		for _, track := range playlistData.Tracks.TrackItems {
			trackIDsList = append(trackIDsList, track.TrackDetails.ID)
			trackIDLookup[track.TrackDetails.ID] = track.TrackDetails
			popularity.push(track.TrackDetails.ID, track.TrackDetails.Popularity)
		}

		audioFeatures, err := a.spotifyReq.GetAudioFeatures(trackIDsList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to fetch audio feature data: %s", err)
			return
		}

		var (
			energy           statsGroup
			danceability     statsGroup
			valance          statsGroup
			acousticness     statsGroup
			speechiness      statsGroup
			instrumentalness statsGroup
		)
		for _, feature := range audioFeatures {
			energy.push(feature.ID, feature.Energy)
			danceability.push(feature.ID, feature.Danceability)
			valance.push(feature.ID, feature.Valence)
			acousticness.push(feature.ID, feature.Acousticness)
			speechiness.push(feature.ID, feature.Speechiness)
			instrumentalness.push(feature.ID, feature.Instrumentalness)
		}

		energy.finish(trackIDLookup)
		danceability.finish(trackIDLookup)
		valance.finish(trackIDLookup)
		acousticness.finish(trackIDLookup)
		speechiness.finish(trackIDLookup)
		instrumentalness.finish(trackIDLookup)
		popularity.finish(trackIDLookup)

		stats := map[string]interface{}{
			"playlist_name":    playlistData.Name,
			"energy":           energy,
			"danceability":     danceability,
			"valance":          valance,
			"acousticness":     acousticness,
			"speechiness":      speechiness,
			"instrumentalness": instrumentalness,
			"popularity":       popularity,
		}

		b, err := json.Marshal(stats)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to JSON marshal audio feature data: %s", err)
			return
		}

		w.Write(b)

	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
