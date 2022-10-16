# spotify-unwrapped
A web app for visualising Spotify listening data. 

## TODO

https://developer.spotify.com/documentation/web-api/reference/#/operations/get-several-audio-features

* acousticness
* danceability
* energy
* instrumentalness
* key
* tempo
* valence (positive/negative)

* Table to sort all tracks by property

1. Get averages of the above
1. Get most/least energetic tune in a playlist
1. Oldest/newest track in a playlist? Graph years playlist tracks were released in

* Artist/album popularity?

curl -H 'Authorization: Bearer BQC3djdWa31ja5TzNkAf8y2ZNj9Ry1dtZ6BFqnyoTWRBdR94z6SHp5r4O89qm1ri_v2ZM7KK6tI808RCUww' -H 'Content-Type: application/json' "https://api.spotify.com/v1/playlists/2VLA8FqcO5Oto2mACkKBOt"
curl -s -H 'Authorization: Bearer BQC3djdWa31ja5TzNkAf8y2ZNj9Ry1dtZ6BFqnyoTWRBdR94z6SHp5r4O89qm1ri_v2ZM7KK6tI808RCUww' -H 'Content-Type: application/json' "https://api.spotify.com/v1/playlists/2VLA8FqcO5Oto2mACkKBOt" > playlist.json

curl -s -H 'Authorization: Bearer BQAP28LLDysRO3pluKaywIpdtv70jeozWBc8V4WiXpVprpsPHUe0WPUKOrL4iju7HPZ4mq9QmlDmffbi5Fo' -H 'Content-Type: application/json' "https://api.spotify.com/v1/audio-features?ids=1ZD51kpXHiTNqymmz3Yy6b,7twSOmINnbYLFtIOSQbXhg" > audio_features.json
