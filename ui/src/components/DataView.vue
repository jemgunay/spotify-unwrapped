<template>
  <v-row class="mx-10" style="min-height: 70vh">
    <v-col cols="12" md="4" offset-md="4" sm="8" offset-sm="2">
      <v-tooltip bottom>
        <template v-slot:activator="{ on, attrs }">
          <v-text-field
              label="Playlist ID"
              v-model="playlistID"
              :rules="playlistInputRules"
              :loading="loading"
              outlined
              clearable
              hide-details="auto"
              class="mt-6 mb-2"
              @input="getPlaylistData"
              v-bind="attrs"
              v-on="on"
          >
          </v-text-field>
        </template>
        <span>To find your playlist's ID in Spotify, select "Share", "Copy link to playlist", then paste the URL here!</span>
      </v-tooltip>
    </v-col>

    <v-col cols="12">
      <p v-if="dataError" class="text-center mt-4">{{ dataError }}</p>
      <p v-else-if="loading" class="text-center mt-4">Unwrapping playlist <strong>{{ playlistID }}</strong>...</p>

      <v-row v-if="playlistMetadata">

        <!-- playlist/owner title & playlist cover image -->
        <v-col cols="12" sm="8" class="playlist-header">
          <v-img
              :src="playlistMetadata['image']"
              aspect-ratio="1"
              max-width="65"
              max-height="65"
              v-ripple
              @click="openSpotifyURL(playlistMetadata['spotify_url'])"
              class="clickable"
          ></v-img>
          <div class="ml-3">
            <h2 class="section-heading clickable" @click="openSpotifyURL(playlistMetadata['spotify_url'])">
              {{ playlistMetadata['name'] }}</h2>
            <h3 @click="openSpotifyURL(playlistMetadata['owner']['spotify_url'])" class="clickable">
              by {{ playlistMetadata['owner']['name'] }}
            </h3>
          </div>
        </v-col>
        <v-col cols="12" sm="4" class="playlist-header-reverse">
          <h3>{{ playlistMetadata['track_count'] }} Tracks</h3>
        </v-col>

        <v-col cols="12">
          <v-divider></v-divider>
        </v-col>

        <ReleaseDateChart :releaseDateData="releaseDateStats"/>
        <PlaylistGeneration :rawStatsData="rawPlaylistStats" :generationDetails="generationDetails"/>

        <v-col cols="12">
          <v-divider></v-divider>
        </v-col>

        <v-col cols="12" md="3">
          <h3 class="section-heading">Top Artists</h3>
          <WordCountTable :topTitleWords="topArtists"/>
        </v-col>
        <v-col cols="12" md="6">
          <ExplicitPieChart :explicitnessData="explicitnessStats"/>
        </v-col>
        <v-col cols="12" md="3">
          <h3 class="section-heading">Common Track Title Words</h3>
          <WordCountTable :topTitleWords="topTitleWords"/>
        </v-col>

        <v-col cols="12">
          <v-divider></v-divider>
        </v-col>

        <AudioFeaturesPolar :rawStatsData="rawPlaylistStats"/>

        <v-col cols="12">
          <v-divider></v-divider>
        </v-col>

        <PopularityDoughnut :rawStatsData="rawPlaylistStats"/>

        <!--<v-col cols="12">-->
        <!--  <v-divider></v-divider>-->
        <!--</v-col>-->
        <!--<RawStatsTable :rawStatsData="rawPlaylistStats"/>-->

      </v-row>
    </v-col>
  </v-row>
</template>

<script>
import axios from "axios";

import ReleaseDateChart from "@/components/charts/ReleaseDateChart";
import PlaylistGeneration from "@/components/charts/PlaylistGeneration";
import WordCountTable from "@/components/charts/WordCountTable";
import ExplicitPieChart from "@/components/charts/ExplicitPieChart";
import AudioFeaturesPolar from "@/components/charts/AudioFeaturesPolar";
// import RawStatsTable from "@/components/charts/RawStatsTable";
import PopularityDoughnut from "@/components/charts/PopularityDoughnut";

export default {
  name: 'DataView',
  components: {
    ReleaseDateChart,
    PlaylistGeneration,
    WordCountTable,
    ExplicitPieChart,
    AudioFeaturesPolar,
    PopularityDoughnut,
    // RawStatsTable
  },
  data() {
    return {
      apiHost: process.env.VUE_APP_API_HOST,
      dataError: null,
      loading: false,
      playlistID: "1KnTiUzSU2HlEtejfXWPo2",
      lastPlaylistID: null,
      playlistInputRules: [
        value => !!value || 'Required.',
        value => this.isPlaylistValid(value),
      ],
      playlistMetadata: null,
      rawPlaylistStats: null,
      explicitnessStats: null,
      releaseDateStats: null,
      generationDetails: null,
      topTitleWords: null,
      topArtists: null
    }
  },
  created() {
    this.getPlaylistData();
  },
  methods: {
    getPlaylistData() {
      this.dataError = null;

      if (this.isPlaylistValid(this.playlistID) !== true) {
        return;
      }

      // if the playlist ID hasn't changed, then we don't want to request the same data
      if (this.playlistID === this.lastPlaylistID) return;

      this.loading = true;
      axios
          .get(this.apiHost + "/api/v1/playlists/" + this.playlistID)
          .then(response => {
            this.playlistMetadata = response.data["metadata"];
            this.rawPlaylistStats = response.data["stats"]["raw"];
            this.explicitnessStats = response.data["stats"]["explicitness"];
            this.releaseDateStats = response.data["stats"]["release_dates"];
            this.generationDetails = response.data["stats"]["generation"];
            this.topTitleWords = response.data["stats"]["top_title_words"];
            this.topArtists = response.data["stats"]["top_artists"];

            this.lastPlaylistID = this.playlistID;
          })
          .catch(error => {
            if (error.code === "ERR_NETWORK") {
              this.dataError = "Failed to communicate with server (network error)...";
              console.error(error);
              return;
            }
            if (error.response.status === 404) {
              this.dataError = "Playlist could not be found...";
              return;
            }
            this.lastPlaylistID = "";
            this.dataError = "Failed to load data: " + error.message;
            console.error(error);
          })
          .finally(() => {
            this.loading = false;
          })
    },
    isPlaylistValid(value) {
      const playlistIDLength = 22;
      if ((value || '').length < playlistIDLength) return "Playlist ID too short.";
      if ((value || '').length === playlistIDLength) {
        this.playlistID = value;
        return true;
      }

      // attempt to process playlist URL into raw ID
      let segments = new URL(value).pathname.split('/');
      value = segments.pop() || segments.pop();

      // validate length represents a playlist ID
      if ((value || '').length !== playlistIDLength) return "Invalid Spotify playlist URL provided.";
      this.playlistID = value;
      return true;
    },
    openSpotifyURL(spotifyURL) {
      if (spotifyURL == null) {
        return;
      }
      window.open(spotifyURL, '_blank');
    }
  },
}
</script>

<style>
.playlist-header {
  display: flex;
}

.playlist-header-reverse {
  display: flex;
  flex-direction: row-reverse;
  text-align: right;
  align-items: flex-end;
}

.playlist-header, .playlist-header-reverse h3 {
  color: rgba(0, 0, 0, 0.71);
  font-style: italic;
}
</style>