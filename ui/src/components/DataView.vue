<template>
  <v-row class="pa-10">
    <v-col cols="4">
      <v-text-field
          label="Playlist ID"
          v-model="playlistID"
          :rules="playlistInputRules"
          :loading="loading"
          @input="getPlaylistData">
      </v-text-field>
    </v-col>

    <v-col cols="12">
      <p v-if="dataError">{{ dataError }}</p>
      <p v-else-if="loading">Unwrapping playlist <strong>{{ playlistID }}</strong>...</p>

      <v-row v-if="playlistName">

        <v-col cols="12">
          <v-divider></v-divider>
        </v-col>

        <v-col cols="12">
          <!-- playlist/owner title -->
          <h1>{{ playlistName }} by {{ playlistOwner }}</h1>
        </v-col>

        <v-col cols="12">
          <v-divider></v-divider>
        </v-col>

        <v-col cols="6">
          <!-- release date chart -->
          <ReleaseDateChart :releaseDateData="releaseDateStats"/>
        </v-col>

        <v-col cols="6">
          <!-- generations text -->
          <h4 class="mb-3">Your playlist's average age is {{ generationDetails["age"] }} years old (born {{
              generationDetails["year"]
            }})! This makes it a member of
            {{ generationDetails["name"] }}
            ({{ generationDetails["lower"] }} - {{ generationDetails["upper"] }})...</h4>
          <p>{{ generationDetails["summary"] }}</p>

          <v-divider></v-divider>

          <v-list-item three-line>
            <v-list-item-content>
              <v-list-item-title>Oldest Track</v-list-item-title>
              <v-list-item-subtitle>{{ rawPlaylistStats["releaseDates"]["min"]["name"] }}</v-list-item-subtitle>
              <v-list-item-subtitle>{{ rawPlaylistStats["releaseDates"]["min"]["date"] }}</v-list-item-subtitle>
            </v-list-item-content>

            <v-list-item-content>
              <v-list-item-title>Youngest Track</v-list-item-title>
              <v-list-item-subtitle>{{ rawPlaylistStats["releaseDates"]["max"]["name"] }}</v-list-item-subtitle>
              <v-list-item-subtitle>{{ rawPlaylistStats["releaseDates"]["max"]["date"] }}</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-col>

        <v-col cols="12">
          <v-divider></v-divider>
        </v-col>

        <v-col cols="6">
          <!-- title word count stats table -->
          <h3>Track Title Word Count</h3>

          <v-simple-table dense class="scrollable-table">
            <template v-slot:default>
              <thead>
              <tr class="text-left">
                <th>Word</th>
                <th>Count</th>
              </tr>
              </thead>
              <tbody>
              <tr
                  v-for="(key, index) in topTitleWords['keys']"
                  :key="key"
              >
                <td>{{ key }}</td>
                <td>{{ topTitleWords['values'][index] }}</td>
              </tr>
              </tbody>
            </template>
          </v-simple-table>
        </v-col>

        <v-col cols="6">
          <!-- explicit lyrics pie chart -->
          <h3>Explicit vs Non-Explicit Lyrics</h3>

          <ExplicitChart :explicitnessData="explicitnessStats"/>
        </v-col>

        <v-col cols="12">
          <v-divider></v-divider>
        </v-col>

        <v-col cols="12">
          <!-- raw stats table -->
          <v-simple-table>
            <template v-slot:default>
              <thead>
              <tr class="text-left">
                <th>Stat</th>
                <th>Min</th>
                <th>Max</th>
                <th>Avg</th>
              </tr>
              </thead>
              <tbody>
              <tr
                  v-for="(stats, statName) in rawPlaylistStats"
                  :key="statName"
              >
                <td>{{ statName }}</td>
                <td>{{ stats.min.name }} ({{ stats.min.value }})</td>
                <td>{{ stats.max.name }} ({{ stats.max.value }})</td>
                <td>{{ stats.avg.value }}</td>
              </tr>
              </tbody>
            </template>
          </v-simple-table>
        </v-col>

      </v-row>
    </v-col>
  </v-row>
</template>

<script>
import axios from "axios";

import ExplicitChart from "./charts/ExplicitChart";
import ReleaseDateChart from "@/components/charts/ReleaseDateChart";

export default {
  name: 'DataView',
  components: {
    ReleaseDateChart,
    ExplicitChart
  },
  data() {
    return {
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
          .get('http://localhost:8080/api/v1/playlists/' + this.playlistID)
          .then(response => {
            this.playlistName = response.data["playlist_name"];
            this.playlistOwner = response.data["owner_name"];
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
.scrollable-table {
  max-height: 60vh;
  overflow: auto;
}
</style>