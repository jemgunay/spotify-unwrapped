# Spotify Unwrapped

A web app for visualising Spotify playlist data.

## Screenshots

<p align="center">
  <img src="screenshots/1.png" width="40%"/>
  <img src="screenshots/4.png" width="40%"/>
</p>
<p align="center">
  <img src="screenshots/2.png" width="40%"/>
  <img src="screenshots/3.png" width="40%"/>
</p>

## Usage

```bash
go run main.go
go run main.go -debug # debug logs

cd ui
npm run serve
```

## TODO

* Add tooltips to explain audio features
* Round audio feature averages
* Persist playlist ID in URL
* Cache result for a period of time server-side (memcached, TTL 5 mins)
* Visualisation ideas
    * Key/tempo vs valence
    * Pie chart of artist genre frequency in a playlist
    * Languages/country
    * Track duration
    * Loudness
