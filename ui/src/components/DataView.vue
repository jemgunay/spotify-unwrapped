<template>
  <v-row class="mx-10" style="min-height: 60vh">
    <v-col cols="12" md="4" offset-md="4" sm="8" offset-sm="2">
      <v-text-field
          label="Playlist ID"
          v-model="playlistID"
          :rules="playlistInputRules"
          :loading="loading"
          outlined
          clearable
          hide-details="auto"
          class="mt-6 mb-2"
          @input="getPlaylistData">
      </v-text-field>
    </v-col>

    <v-col cols="12">
      <p v-if="dataError" class="text-center">{{ dataError }}</p>
      <p v-else-if="loading" class="text-center">Unwrapping playlist <strong>{{ playlistID }}</strong>...</p>

      <v-row v-if="playlistName" id="data-container">

        <!-- TODO: remake this -->
        <v-col cols="2">
          <v-img
              :src="playlistImage"
              aspect-ratio="1"
              max-width="50"
              max-height="50"
          ></v-img>
        </v-col>
        <v-col cols="10">
          <!-- playlist/owner title -->
          <h2 class="section-heading d-inline">{{ playlistName }}</h2>
          <h3>by {{ playlistOwner }}</h3>
        </v-col>

        <v-col cols="12">
          <v-divider></v-divider>
        </v-col>

        <ReleaseDateChart :releaseDateData="releaseDateStats"/>
        <PlaylistGeneration :rawStatsData="rawPlaylistStats" :generationDetails="generationDetails"/>

        <v-col cols="12">
          <v-divider></v-divider>
        </v-col>

        <WordCountTable :topTitleWords="topTitleWords"/>
        <ExplicitPieChart :explicitnessData="explicitnessStats"/>

        <v-col cols="12">
          <v-divider></v-divider>
        </v-col>

        <AudioFeaturesPolar :rawStatsData="rawPlaylistStats"/>

        <v-col cols="12">
          <v-divider></v-divider>
        </v-col>

        <PopularityDoughnut :rawStatsData="rawPlaylistStats"/>

        <v-col cols="12">
          <v-divider></v-divider>
        </v-col>

        <RawStatsTable :rawStatsData="rawPlaylistStats"/>

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
import RawStatsTable from "@/components/charts/RawStatsTable";
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
    RawStatsTable
  },
  data() {
    return {
      apiHost: process.env.VUE_APP_API_HOST,
      dataError: null,
      loading: false,
      playlistID: "1AXy6ag2d0ag8DEdOE7kWm",
      lastPlaylistID: null,
      playlistInputRules: [
        value => !!value || 'Required.',
        value => this.isPlaylistValid(value),
      ],
      playlistName: null,
      playlistOwner: null,
      playlistImage: null,
      rawPlaylistStats: {},
      explicitnessStats: {},
      releaseDateStats: {},
      generationDetails: {},
      topTitleWords: {}
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
            this.playlistName = response.data["playlist_name"];
            this.playlistOwner = response.data["owner_name"];
            this.playlistImage = response.data["playlist_image"];
            this.rawPlaylistStats = response.data["stats"]["raw"];
            this.explicitnessStats = response.data["stats"]["explicitness"];
            this.releaseDateStats = response.data["stats"]["release_dates"];
            this.generationDetails = response.data["stats"]["generation"];
            this.topTitleWords = response.data["stats"]["top_title_words"];

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
    }
  },
}
</script>

<style>
</style>