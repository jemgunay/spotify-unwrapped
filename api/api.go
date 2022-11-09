package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

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

type statsMapping map[string]int

func newStatsMapping(capacity int) statsMapping {
	return make(map[string]int, capacity)
}

func (m statsMapping) push(key string) {
	m[key] = m[key] + 1
}

type orderedKVPair struct {
	Keys   []string `json:"keys"`
	Values []int    `json:"values"`
}

func (p *orderedKVPair) Len() int {
	return len(p.Keys)
}
func (p *orderedKVPair) Less(i int, j int) bool {
	return p.Keys[i] < p.Keys[j]
}
func (p *orderedKVPair) Swap(i int, j int) {
	p.Keys[i], p.Keys[j] = p.Keys[j], p.Keys[i]
	p.Values[i], p.Values[j] = p.Values[j], p.Values[i]
}

func (m statsMapping) orderedLabelsAndValues() orderedKVPair {
	pair := orderedKVPair{
		Keys:   make([]string, 0, len(m)),
		Values: make([]int, 0, len(m)),
	}
	for k, v := range m {
		pair.Keys = append(pair.Keys, k)
		pair.Values = append(pair.Values, v)
	}
	sort.Sort(&pair)
	return pair
}

// PlaylistsHandler provides data used to drive visualisations.
func (a API) PlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// fetch playlist data for given playlist ID
	playlistData, err := a.spotifyReq.GetPlaylist(vars["id"])
	if err != nil {
		log.Printf("failed to fetch playlist data: %s", err)
		if errors.Is(err, spotify.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var popularity statsGroup
	explicitMapping := newStatsMapping(2)
	releaseDatesMapping := newStatsMapping(10)
	trackIDsList := make([]string, 0, len(playlistData.Tracks.TrackItems))
	trackIDLookup := make(map[string]spotify.TrackDetails, len(playlistData.Tracks.TrackItems))

	for _, track := range playlistData.Tracks.TrackItems {
		trackIDsList = append(trackIDsList, track.TrackDetails.ID)
		trackIDLookup[track.TrackDetails.ID] = track.TrackDetails
		// aggregate track popularity
		popularity.push(track.TrackDetails.ID, track.TrackDetails.Popularity)
		// aggregate by release year
		releaseDate, err := time.Parse("2006-01-02", track.TrackDetails.Album.ReleaseDate)
		if err == nil {
			releaseDatesMapping.push(strconv.Itoa(releaseDate.Year()))
		}

		if track.TrackDetails.Explicit {
			explicitMapping.push("explicit")
		} else {
			explicitMapping.push("non-explicit")
		}
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
			"raw": map[string]interface{}{
				"popularity":       popularity,
				"energy":           energy,
				"danceability":     danceability,
				"valance":          valance,
				"acousticness":     acousticness,
				"speechiness":      speechiness,
				"instrumentalness": instrumentalness,
			},
			"explicitness":  explicitMapping,
			"release_dates": releaseDatesMapping.orderedLabelsAndValues(),
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
