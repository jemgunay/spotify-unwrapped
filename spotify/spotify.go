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

	"go.uber.org/zap"

	"github.com/jemgunay/spotify-unwrapped/config"
	"github.com/jemgunay/spotify-unwrapped/spotify/auth"
)

// Requester wraps the Spotify HTTP REST API.
type Requester struct {
	conf       config.Spotify
	access     *auth.Access
	httpClient *http.Client
	logger     config.Logger
}

// New initialises a Requester.
func New(logger config.Logger, conf config.Spotify) Requester {
	r := Requester{
		conf: conf,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		logger: logger,
	}
	r.access = auth.New(r.authenticate)
	return r
}

type authRespBody struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// authenticate requests an access token for performing Spotify API requests, as well as its expiry date.
func (r *Requester) authenticate() (string, time.Time, error) {
	formValues := url.Values{}
	formValues.Set("grant_type", "client_credentials")
	formValuesBuf := strings.NewReader(formValues.Encode())

	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", formValuesBuf)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// base64 basic auth secrets
	base64Auth := base64.StdEncoding.EncodeToString([]byte(r.conf.ClientID + ":" + r.conf.ClientSecret))
	req.Header.Set("Authorization", "Basic "+base64Auth)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to perform auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("unexpected status from auth request: %s", resp.Status)
	}

	authBody := authRespBody{}
	if err := json.NewDecoder(resp.Body).Decode(&authBody); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to JSON decode body from auth request: %w", err)
	}

	// successfully authenticated
	expiry := time.Now().UTC().Add(time.Duration(authBody.ExpiresIn) * time.Second)
	r.logger.Info("successfully authenticated with Spotify", zap.Time("expiry", expiry))

	return authBody.AccessToken, expiry, nil
}

// Playlist represents a playlist of tracks.
type Playlist struct {
	Name         string `json:"name"`
	Owner        Owner  `json:"owner"`
	Tracks       Tracks `json:"tracks"`
	Images       Images `json:"images"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

// Owner represents a playlist owner.
type Owner struct {
	DisplayName  string `json:"display_name"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

// Tracks represents a paginated set of track.
type Tracks struct {
	TrackItems []TrackItem `json:"items"`
	NextURL    string      `json:"next"`
	Total      int         `json:"total"`
}

// TrackItem represents the details for a given track.
type TrackItem struct {
	TrackDetails TrackDetails `json:"track"`
}

// TrackDetails represents the details of a track.
type TrackDetails struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Popularity   float64  `json:"popularity"` // 0-100
	Artists      []Artist `json:"artists"`
	Album        Album    `json:"album"`
	Explicit     bool     `json:"explicit"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
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
	ReleaseDate          string `json:"release_date"`
	ReleaseDatePrecision string `json:"release_date_precision"`
	Images               Images `json:"images"`
}

// ParseReleaseDate parses the release date string variant into its equivalent time.Time.
func (a Album) ParseReleaseDate() (time.Time, error) {
	switch a.ReleaseDatePrecision {
	case "day":
		return time.Parse("2006-01-02", a.ReleaseDate)
	case "month":
		return time.Parse("2006-01", a.ReleaseDate)
	case "year":
		return time.Parse("2006", a.ReleaseDate)
	}
	return time.Time{}, errors.New("unrecognised release date precision")
}

// Images represents the set of different resolution images.
type Images []Image

// First returns the first image in the list.
func (i Images) First() string {
	if len(i) == 0 {
		return ""
	}
	return i[0].URL
}

// Image represents an album or playlist cover image.
type Image struct {
	URL string `json:"url"`
}

// ErrNotFound indicates that the requested resource does not exist.
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorised = errors.New("unauthorised")
	ErrRateLimited  = errors.New("rate limited")
)

const (
	apiURL = "https://api.spotify.com/v1/"
	// only the first ~3000 tracks of a playlist will be processed
	maxPlaylistPages = 30
)

// GetPlaylist gets all required data for the given playlist ID.
// https://developer.spotify.com/documentation/web-api/reference/#/operations/get-playlist
func (r *Requester) GetPlaylist(id string) (Playlist, error) {
	playlist, err := r.getPlaylist(id)
	if err != nil {
		return playlist, err
	}

	// get the rest of the paginated playlist tracks
	for i := 0; i < maxPlaylistPages; i++ {
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

	// we reached maximum pages
	return playlist, nil
}

func (r *Requester) getPlaylist(id string) (Playlist, error) {
	reqURL := apiURL + "playlists/" + id

	playlist := Playlist{}
	if err := r.performGetRequest(reqURL, &playlist); err != nil {
		return playlist, fmt.Errorf("get playlist request failed: %w", err)
	}

	return playlist, nil
}

func (r *Requester) getPlaylistTracksPage(nextURL string) (Tracks, error) {
	tracks := Tracks{}
	if err := r.performGetRequest(nextURL, &tracks); err != nil {
		return tracks, fmt.Errorf("get playlist tracks request failed: %w", err)
	}

	return tracks, nil
}

// AudioFeaturesResult represents the response body from the Spotify audio features API.
type AudioFeaturesResult struct {
	Features []AudioFeatures `json:"audio_features"`
}

// AudioFeatures represents the audio feature properties of a track.
type AudioFeatures struct {
	ID               string  `json:"id"`
	Danceability     float64 `json:"danceability"`
	Energy           float64 `json:"energy"`
	Speechiness      float64 `json:"speechiness"`
	Acousticness     float64 `json:"acousticness"`
	Instrumentalness float64 `json:"instrumentalness"`
	Liveness         float64 `json:"liveness"`
	Valence          float64 `json:"valence"`
	Tempo            float64 `json:"tempo"`
	Key              int     `json:"key"`
	DurationMillis   int     `json:"duration_ms"`
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
		if err := r.performGetRequest(reqURL, &audioFeatures); err != nil {
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

func (r *Requester) performGetRequest(reqURL string, target interface{}) error {
	var err error
	for i := 1; i <= 10; i++ {
		var accessToken string
		accessToken, err = r.access.Get()
		if err != nil {
			r.logger.Error("failed to refresh access token after natural token expiry", zap.Error(err),
				zap.String("url", reqURL), zap.Int("attempt", i))
			continue
		}

		err = r.get(reqURL, accessToken, target)
		switch err {
		case nil, ErrNotFound:
			return err
		case ErrUnauthorised:
			if err := r.access.Refresh(); err != nil {
				r.logger.Error("failed to refresh access token", zap.Error(err),
					zap.String("url", reqURL), zap.Int("attempt", i))
			}
		case ErrRateLimited:
			r.logger.Error("requests are being rate limited", zap.Error(err),
				zap.Error(err), zap.String("url", reqURL), zap.Int("attempt", i))
			// TODO: use API Retry-After header (this will do for now): https://stackoverflow.com/questions/30548073/spotify-web-api-rate-limits
			time.Sleep(time.Millisecond * 500)
		default:
			r.logger.Error("failed to perform get request", zap.Error(err),
				zap.Error(err), zap.String("url", reqURL), zap.Int("attempt", i))
			time.Sleep(time.Millisecond * 200)
		}
	}

	r.logger.Error("failed to perform get request after maximum attempts",
		zap.Error(err), zap.String("url", reqURL))
	return err
}

func (r *Requester) get(reqURL string, accessToken string, target interface{}) error {
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	r.logAsDebugCurl(req)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		// trigger force access token refresh
		return ErrUnauthorised
	case http.StatusTooManyRequests:
		return ErrRateLimited
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

// logAsDebugCurl is a helper func for logging
func (r *Requester) logAsDebugCurl(req *http.Request) {
	u := req.URL.String()
	method := req.Method
	auth := req.Header.Get("Authorization")
	curl := fmt.Sprintf("curl -i -X%s '%s' -H 'Authorization: %s'", method, u, auth)
	r.logger.Debug(curl)
}
