package spotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jemgunay/spotify-unwrapped/config"
	"go.uber.org/zap"
)

// Requester wraps the Spotify HTTP REST API.
type Requester struct {
	httpClient  *http.Client
	conf        config.Spotify
	accessToken string
	logger      config.Logger
}

// New initialises a Requester.
func New(logger config.Logger, conf config.Spotify) Requester {
	return Requester{
		conf: conf,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		logger: logger,
	}
}

type authRespBody struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// Auth requests an access token for
func (r *Requester) Auth() error {
	formValues := url.Values{}
	formValues.Set("grant_type", "client_credentials")
	formValuesBuf := strings.NewReader(formValues.Encode())

	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", formValuesBuf)
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// base64 basic auth secrets
	base64Auth := base64.StdEncoding.EncodeToString([]byte(r.conf.ClientID + ":" + r.conf.ClientSecret))
	req.Header.Set("Authorization", "Basic "+base64Auth)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status from auth request: %s", resp.Status)
	}

	authBody := authRespBody{}
	if err := json.NewDecoder(resp.Body).Decode(&authBody); err != nil {
		return fmt.Errorf("failed to JSON decode body from auth request: %w", err)
	}

	// successfully authenticated
	r.accessToken = authBody.AccessToken
	expiry := time.Now().Add(time.Duration(authBody.ExpiresIn) * time.Second)
	r.logger.Info("successfully created authenticated Spotify client", zap.Time("expiry", expiry))

	return nil
}

const apiURL = "https://api.spotify.com/v1/"

// Playlist represents a playlist of tracks.
type Playlist struct {
	Name   string `json:"name"`
	Owner  Owner  `json:"owner"`
	Tracks Tracks `json:"tracks"`
}

// Owner represents a playlist owner.
type Owner struct {
	DisplayName string `json:"display_name"`
}

// Tracks represents a paginated set of track.
type Tracks struct {
	TrackItems []TrackItem `json:"items"`
	NextURL    string      `json:"next"`
}

// TrackItem represents the details for a given track.
type TrackItem struct {
	TrackDetails TrackDetails `json:"track"`
}

// TrackDetails represents the details of a track.
type TrackDetails struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Popularity       float64  `json:"popularity"` // 0-100
	Artists          []Artist `json:"artists"`
	Album            Album    `json:"album"`
	Explicit         bool     `json:"explicit"`
	artistsFormatted string
	trackFormatted   string
}

// GetTrackString lazy processes a track string of the format "Artist - Tracks".
func (t *TrackDetails) GetTrackString() string {
	if t.trackFormatted != "" {
		return t.trackFormatted
	}
	t.trackFormatted = t.GetArtists() + " - " + t.Name
	return t.trackFormatted
}

// GetArtists lazy processes artists into a comma separated string.
func (t *TrackDetails) GetArtists() string {
	if t.artistsFormatted != "" {
		return t.artistsFormatted
	}

	for i, artist := range t.Artists {
		t.artistsFormatted += artist.Name
		if i < len(t.Artists)-1 {
			t.artistsFormatted += ", "
		}
	}
	return t.artistsFormatted
}

// Artist represents a single artist.
type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Album represents a single album.
type Album struct {
	ReleaseDate string `json:"release_date"`
}

// ErrNotFound indicates that the requested resource does not exist.
var ErrNotFound = errors.New("not found")

// GetPlaylist gets all required data for the given playlist ID.
// https://developer.spotify.com/documentation/web-api/reference/#/operations/get-playlist
func (r *Requester) GetPlaylist(id string) (Playlist, error) {
	playlist, err := r.getPlaylist(id)
	if err != nil {
		return playlist, err
	}

	// get the rest of the paginated playlist tracks
	for {
		if playlist.Tracks.NextURL == "" {
			return playlist, nil
		}

		playlistTracks, err := r.getPlaylistTracksPage(playlist.Tracks.NextURL)
		if err != nil {
			return playlist, err
		}

		playlist.Tracks.TrackItems = append(playlist.Tracks.TrackItems, playlistTracks.TrackItems...)
		playlist.Tracks.NextURL = playlistTracks.NextURL
	}
}

func (r *Requester) getPlaylist(id string) (Playlist, error) {
	reqURL := apiURL + "playlists/" + id

	playlist := Playlist{}
	if err := r.get(reqURL, &playlist); err != nil {
		return playlist, fmt.Errorf("get playlist request failed: %w", err)
	}

	return playlist, nil
}

func (r *Requester) getPlaylistTracksPage(nextURL string) (Tracks, error) {
	tracks := Tracks{}
	if err := r.get(nextURL, &tracks); err != nil {
		return tracks, fmt.Errorf("get playlist tracks request failed: %w", err)
	}

	return tracks, nil
}

// AudioFeaturesResult represents the response body from the Spotify audio features API.
type AudioFeaturesResult struct {
	Features []AudioFeatures `json:"audio_features"`
}

// AudioFeatures represents the properties of a track.
type AudioFeatures struct {
	Danceability     float64 `json:"danceability"`
	Energy           float64 `json:"energy"`
	Key              int     `json:"key"`
	Loudness         float64 `json:"loudness"`
	Mode             int     `json:"mode"`
	Speechiness      float64 `json:"speechiness"`
	Acousticness     float64 `json:"acousticness"`
	Instrumentalness float64 `json:"instrumentalness"`
	Liveness         float64 `json:"liveness"`
	Valence          float64 `json:"valence"`
	Tempo            float64 `json:"tempo"`
	Type             string  `json:"type"`
	ID               string  `json:"id"`
	URI              string  `json:"uri"`
	TrackHref        string  `json:"track_href"`
	AnalysisURL      string  `json:"analysis_url"`
	DurationMs       int     `json:"duration_ms"`
	TimeSignature    int     `json:"time_signature"`
}

// GetAudioFeatures gets audio properties for a set of tracks.
// https://developer.spotify.com/documentation/web-api/reference/#/operations/get-several-audio-features
func (r *Requester) GetAudioFeatures(trackIDs []string) ([]AudioFeatures, error) {
	totalAudioFeatures := AudioFeaturesResult{}
	for i := 100; ; i += 100 {
		lower := i - 100
		if i > len(trackIDs) {
			i = len(trackIDs)
		}
		ids := trackIDs[lower:i]
		reqURL := apiURL + "audio-features?ids=" + strings.Join(ids, ",")

		audioFeatures := AudioFeaturesResult{}
		if err := r.get(reqURL, &audioFeatures); err != nil {
			return nil, fmt.Errorf("audio features request failed: %w", err)
		}

		if len(totalAudioFeatures.Features) == 0 {
			totalAudioFeatures = audioFeatures
		} else {
			totalAudioFeatures.Features = append(totalAudioFeatures.Features, audioFeatures.Features...)
		}

		if i == len(trackIDs) {
			return totalAudioFeatures.Features, nil
		}
	}
}

func (r *Requester) get(reqURL string, target interface{}) error {
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.accessToken)

	r.getCurl(req)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read resp body: %w", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to JSON unmarshal response body: %w", err)
	}

	return nil
}

func (r *Requester) getCurl(req *http.Request) {
	u := req.URL.String()
	method := req.Method
	auth := req.Header.Get("Authorization")
	curl := fmt.Sprintf("curl -i X%s '%s' -H 'Authorization: %s'", method, u, auth)
	r.logger.Debug(curl)
}
