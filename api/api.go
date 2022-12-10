package api

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/jemgunay/spotify-unwrapped/config"
	"github.com/jemgunay/spotify-unwrapped/spotify"
	"github.com/jemgunay/spotify-unwrapped/stats"
)

// API is an API which also performs track data collection and aggregation.
type API struct {
	logger     config.Logger
	spotifyReq spotify.Requester
}

// New returns a Spotify API.
func New(logger config.Logger, spotifyReq spotify.Requester) API {
	return API{
		logger:     logger,
		spotifyReq: spotifyReq,
	}
}

// PlaylistsHandler processes the given Spotify playlist data used to drive visualisations.
func (a API) PlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playlistID := vars["playlistID"]

	logger := a.logger.With(zap.String("playlist", playlistID))
	logger.Info("playlist API request")

	// fetch playlist data for given playlist ID
	playlistData, err := a.spotifyReq.GetPlaylist(playlistID)
	if err != nil {
		logger.Error("failed to fetch playlist data", zap.Error(err))
		if errors.Is(err, spotify.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var (
		popularity          stats.Group
		releaseDates        stats.Group
		releaseDatesMapping = stats.NewMapping(10)
		explicitMapping     = stats.NewMapping(2)
		titleWordMapping    = stats.NewMapping(100)
		artistWordMapping   = stats.NewMapping(100)
		trackIDsList        = make([]string, 0, len(playlistData.Tracks.TrackItems))
		trackIDLookup       = make(map[string]spotify.TrackDetails, len(playlistData.Tracks.TrackItems))
	)

	for _, track := range playlistData.Tracks.TrackItems {
		trackIDsList = append(trackIDsList, track.TrackDetails.ID)
		trackIDLookup[track.TrackDetails.ID] = track.TrackDetails
		// aggregate track popularity
		popularity.Push(track.TrackDetails.ID, track.TrackDetails.Popularity)
		// aggregate by release year
		releaseDate, err := track.TrackDetails.Album.ParseReleaseDate()
		if err == nil {
			// sometimes tracks don't have release date metadata - skip them from this stat
			releaseDatesMapping.Push(strconv.Itoa(releaseDate.Year()))
			releaseDates.Push(track.TrackDetails.ID, float64(releaseDate.Unix()))
		}

		// count explicit vs explicit tracks
		explicit := "non-explicit"
		if track.TrackDetails.Explicit {
			explicit = "explicit"
		}
		explicitMapping.Push(explicit)

		// count unique sentence title words
		stats.CountWordsInSentence(track.TrackDetails.Name, titleWordMapping)
		for _, artist := range track.TrackDetails.Artists {
			artistWordMapping.Push(artist.Name)
		}
	}

	// determine the playlist age/generation
	releaseDates.Calc(trackIDLookup, stats.ToDateString())
	generation, err := stats.GetGeneration(releaseDates.Mean.DateYear())
	if err != nil {
		// don't error out - we can still display all the other data
		logger.Error("failed to determine playlist generation", zap.Error(err),
			zap.Int("avg_year", releaseDates.Mean.DateYear()))
	}

	// bulk fetch audio feature data for each track in playlist
	audioFeatures, err := a.spotifyReq.GetAudioFeatures(trackIDsList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("failed to fetch audio feature data", zap.Error(err))
		return
	}

	// aggregate track audio feature metadata for each stat
	var energy, danceability, valence, acousticness, speechiness, instrumentalness, liveness stats.Group
	var trackDuration, tempo stats.Group
	pitchKeyCounts := stats.NewMapping(12, stats.PitchKeys...)

	for _, feature := range audioFeatures {
		energy.Push(feature.ID, feature.Energy)
		danceability.Push(feature.ID, feature.Danceability)
		valence.Push(feature.ID, feature.Valence)
		acousticness.Push(feature.ID, feature.Acousticness)
		speechiness.Push(feature.ID, feature.Speechiness)
		instrumentalness.Push(feature.ID, feature.Instrumentalness)
		liveness.Push(feature.ID, feature.Liveness)
		trackDuration.Push(feature.ID, float64(feature.DurationMillis))
		tempo.Push(feature.ID, math.Round(feature.Tempo))

		// -1 is Spotify's unknown key value
		if feature.Key > -1 {
			pitchKeyCounts.Push(stats.SpotifyKeyToPitchKey(feature.Key))
		}
	}

	// perform final calculations on each stat and lookup track names
	popularity.Calc(trackIDLookup)
	trackDuration.Calc(trackIDLookup, stats.ToDurationString())
	tempo.Calc(trackIDLookup)
	// process the following stats from decimal to percentages
	toPercentage := stats.WithMultiplier(100)
	energy.Calc(trackIDLookup, toPercentage)
	danceability.Calc(trackIDLookup, toPercentage)
	valence.Calc(trackIDLookup, toPercentage)
	acousticness.Calc(trackIDLookup, toPercentage)
	speechiness.Calc(trackIDLookup, toPercentage)
	instrumentalness.Calc(trackIDLookup, toPercentage)
	liveness.Calc(trackIDLookup, toPercentage)

	// generate final response payload
	statsPayload := map[string]any{
		"metadata": map[string]any{
			"name": playlistData.Name,
			"owner": map[string]any{
				"name":        playlistData.Owner.DisplayName,
				"spotify_url": playlistData.Owner.ExternalURLs.Spotify,
			},
			"image":       playlistData.Images.First(),
			"spotify_url": playlistData.ExternalURLs.Spotify,
			"track_count": playlistData.Tracks.Total,
		},
		"stats": map[string]any{
			"raw": map[string]any{
				"popularity":       popularity,
				"energy":           energy,
				"danceability":     danceability,
				"valence":          valence,
				"acousticness":     acousticness,
				"speechiness":      speechiness,
				"instrumentalness": instrumentalness,
				"liveness":         liveness,
				"release_dates":    releaseDates,
				"track_durations":  trackDuration,
				"tempo":            tempo,
			},
			"explicitness": explicitMapping,
			"release_dates": releaseDatesMapping.OrderedLabelsAndValues(
				stats.WithSort(stats.SortKey, false),
			),
			"generation": generation,
			"top_title_words": titleWordMapping.OrderedLabelsAndValues(
				stats.WithSort(stats.SortValue, true),
				stats.WithTruncate(50),
			),
			"top_artists": artistWordMapping.OrderedLabelsAndValues(
				stats.WithSort(stats.SortValue, true),
				stats.WithTruncate(50),
			),
			"pitch_key": pitchKeyCounts.OrderedLabelsAndValues(
				stats.WithSort(stats.SortPitchKey, false),
			),
		},
	}

	respBody, err := json.Marshal(statsPayload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("failed to JSON marshal playlist API data", zap.Error(err))
		return
	}

	w.Write(respBody)
}
