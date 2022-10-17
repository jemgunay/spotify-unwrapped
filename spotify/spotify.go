package spotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jemgunay/spotify-unwrapped/config"
)

// Requester wraps the Spotify HTTP REST API.
type Requester struct {
	httpClient  *http.Client
	conf        config.Spotify
	accessToken string
}

// New initialises a Requester.
func New(conf config.Spotify) Requester {
	return Requester{
		conf: conf,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

type authRespBody struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// Auth requests an access token for
func (r *Requester) Auth() error {
	// create URL encoded body
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
	log.Printf("successfully created authenticated Spotify client, expires at %s", expiry.Format(time.RFC3339))
	// log.Printf("access token: %s", r.accessToken)

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
	NextURL    string      `json:"next"` // TODO: continue to read all tracks, until next is nil
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

// ErrPlaylistNotFound indicates that the requested playlist ID does not exist.
var ErrPlaylistNotFound = errors.New("playlist not found")

// GetPlaylist gets all required data for the given playlist ID.
// https://developer.spotify.com/documentation/web-api/reference/#/operations/get-playlist
func (r *Requester) GetPlaylist(id string) (Playlist, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL+"playlists/"+id, nil)
	if err != nil {
		return Playlist{}, fmt.Errorf("failed to create playlist request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.accessToken)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return Playlist{}, fmt.Errorf("failed to perform playlist request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return Playlist{}, ErrPlaylistNotFound
	default:
		return Playlist{}, fmt.Errorf("unexpected status from playlist request: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Playlist{}, fmt.Errorf("failed to read playlist resp body: %w", err)
	}

	playlist := Playlist{}
	if err := json.Unmarshal(body, &playlist); err != nil {
		return Playlist{}, fmt.Errorf("failed to JSON unmarshal playlist resp body: %w", err)
	}

	return playlist, nil
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
	idStr := strings.Join(trackIDs, ",")

	req, err := http.NewRequest(http.MethodGet, apiURL+"audio-features?ids="+idStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio features request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.accessToken)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform audio features request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status from audio features request: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio features resp body: %w", err)
	}

	audioFeaturesResult := AudioFeaturesResult{}
	if err := json.Unmarshal(body, &audioFeaturesResult); err != nil {
		return nil, fmt.Errorf("failed to JSON unmarshal audio features resp body: %w", err)
	}

	return audioFeaturesResult.Features, nil
}
