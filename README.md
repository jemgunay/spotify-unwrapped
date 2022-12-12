# Spotify Unwrapped

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/jemgunay/spotify-unwrapped/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/jemgunay/spotify-unwrapped/tree/main)

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
# debug logs, e.g. full egress request logs
go run main.go -debug

cd ui
npm run serve
```

## TODO

* Redo screenshots
* Add tooltips
* Add general info btn for implementation details/limitations (i.e. max tracks per query)
* Add hover for popularity vs positivity graph