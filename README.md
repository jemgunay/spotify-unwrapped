# Spotify Unwrapped

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/jemgunay/spotify-unwrapped/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/jemgunay/spotify-unwrapped/tree/main)

A web app for visualising Spotify playlist data.

## Screenshots

<table class="border-collapse:collapse">
  <tr>
    <td>
        <img src="screenshots/1.png"/>
    </td>
    <td>
        <img src="screenshots/2.png"/>
    </td>
  </tr>
</table>

## Usage

```bash
go run main.go
# debug logs, e.g. full Spotify API request logs
go run main.go -debug

cd ui
npm run serve
```

## TODO

* Redo screenshots
* Add tooltips
* Add general info btn for implementation details/limitations (i.e. max tracks per query)