# Spotify Unwrapped

A web app for visualising Spotify playlist data. 

## Features

* Explicit lyric pie chart
* Track release year bar chart
* Playlist average age & its generation (is your playlist a Boomer, Zoomer, etc)
* Average track audio features polar chart
  * Acousticness
  * Danceability
  * Energy
  * Instrumentalness
  * Speechiness
  * Liveness
  * Valence (positive/negative)
* Average track popularity doughnut chart

## Usage

```bash
go run main.go
go run main.go -debug # debug logs

cd ui
npm run serve
```

## TODO

* Fix playlist name + image display
* Fix explicit pie responsiveness
* round audio feature avgs
* API context request timeout (allow a few minutes)
* Upper cap on playlist size, i.e. bail after 2k tracks/N pages?
* Persist playlist in URL (and cache result for a period of time server-side?)
* Playlist input tooltip
* Visualisations
  * Key/tempo vs valence 
  * Get most/least energetic tune in a playlist
  * Pie chart of artist genre frequency in a playlist
  * Languages/country
  * Words in title
  * Track duration
  * Loudness
  * Popularity & artist followers
  * Most common artist
