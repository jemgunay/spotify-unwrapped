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

* Fix weird track title word trailing punctuation + sort alphabetically as well as by count
* round audio feature avgs
* API context request timeout (allow a few minutes)
* Persist playlist in URL
* Cache result for a period of time server-side (memcached, TTL 5 mins)
* Playlist input tooltip
* Visualisations
  * Key/tempo vs valence
  * Pie chart of artist genre frequency in a playlist
  * Languages/country
  * Track duration
  * Loudness
  * Popularity & artist followers
