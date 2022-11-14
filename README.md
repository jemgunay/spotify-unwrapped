# Spotify Unwrapped

A web app for visualising Spotify playlist data. 

## Features

* Explicit lyric pie chart
* Track release year bar chart
* Playlist average age & its generation (is your playlist a Boomer, Zoomer, etc)
* General track audio feature stats table with averages
  * acousticness
  * danceability
  * energy
  * instrumentalness
  * popularity
  * speechiness
  * valence (positive/negative)

## TODO

* Implement Spotify auth store refresh on expiry
* Upper cap on playlist size, i.e. bail after 2k tracks?
* Persist playlist in URL (and cache result for a period of time server-side?)
* Images for playlist & min/max tracks
* Input tooltip
* Bar Chart animate up
* Visualisations
  * Key/tempo vs valence 
  * Get most/least energetic tune in a playlist 
  * Oldest/newest track in a playlist
  * Pie chart of artist genre frequency in a playlist
  * Languages/country
  * Words in title
  * Track duration
  * Loudness
  * Popularity & followers