package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jemgunay/spotify-unwrapped/spotify"
)

// API is an API which also performs track data collection and aggregation.
type API struct {
	spotifyReq spotify.Requester
}

// New returns a Spotify API.
func New(spotifyReq spotify.Requester) API {
	return API{
		spotifyReq: spotifyReq,
	}
}

type statsDetail struct {
	id    string
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type statsGroup struct {
	Min   statsDetail `json:"min"`
	Max   statsDetail `json:"max"`
	sum   float64
	count float64
	Mean  float64 `json:"avg"`
}

func (s *statsGroup) push(id string, val float64) {
	s.sum += val
	s.count++

	switch {
	case s.Min.id == "":
		s.Min = statsDetail{id: id, Value: val}
		s.Max = statsDetail{id: id, Value: val}
	case val > s.Max.Value:
		s.Max = statsDetail{id: id, Value: val}
	case val < s.Min.Value:
		s.Min = statsDetail{id: id, Value: val}
	}
}

func (s *statsGroup) calc(lookup map[string]spotify.TrackDetails) {
	minTrack := lookup[s.Min.id]
	s.Min.Name = minTrack.GetTrackString()
	maxTrack := lookup[s.Max.id]
	s.Max.Name = maxTrack.GetTrackString()
	s.Mean = s.sum / s.count
}

// PlaylistsHandler provides data used to drive visualisations.
func (a API) PlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// fetch playlist data for given playlist ID
	playlistData, err := a.spotifyReq.GetPlaylist(vars["id"])
	if err != nil {
		log.Printf("failed to fetch playlist data: %s", err)
		if errors.Is(err, spotify.ErrPlaylistNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var popularity statsGroup
	trackIDsList := make([]string, 0, len(playlistData.Tracks.TrackItems))
	trackIDLookup := make(map[string]spotify.TrackDetails, len(playlistData.Tracks.TrackItems))
	for _, track := range playlistData.Tracks.TrackItems {
		trackIDsList = append(trackIDsList, track.TrackDetails.ID)
		trackIDLookup[track.TrackDetails.ID] = track.TrackDetails
		// aggregate track popularity
		popularity.push(track.TrackDetails.ID, track.TrackDetails.Popularity)
	}

	// bulk fetch audio feature data for each track in playlist
	audioFeatures, err := a.spotifyReq.GetAudioFeatures(trackIDsList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to fetch audio feature data: %s", err)
		return
	}

	// aggregate track data for each stat
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

	// perform final calculations on each stat and lookup track names
	popularity.calc(trackIDLookup)
	energy.calc(trackIDLookup)
	danceability.calc(trackIDLookup)
	valance.calc(trackIDLookup)
	acousticness.calc(trackIDLookup)
	speechiness.calc(trackIDLookup)
	instrumentalness.calc(trackIDLookup)

	// generate final output payload
	stats := map[string]interface{}{
		"playlist_name": playlistData.Name,
		"owner_name":    playlistData.Owner.DisplayName,
		"stats": map[string]interface{}{
			"popularity":       popularity,
			"energy":           energy,
			"danceability":     danceability,
			"valance":          valance,
			"acousticness":     acousticness,
			"speechiness":      speechiness,
			"instrumentalness": instrumentalness,
		},
	}

	respBody, err := json.Marshal(stats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to JSON marshal playlist API data: %s", err)
		return
	}

	w.Write(respBody)
}
