<template>
  <v-row class="pa-10">
    <v-col cols="4">
      <v-text-field
          label="Playlist ID"
          v-model="playlistID"
          :rules="playlistIDRules"
          :loading="loading"
          @input="getPlaylistData"></v-text-field>
    </v-col>

    <v-col cols="12">
      <p v-if="dataError">{{ dataError }}</p>

      <div v-else-if="loading">
        <p>Processing playlist {{ playlistID }}...</p>
      </div>

      <v-row>
        <v-col cols="6">
          <ExplicitChart :explicitnessData="explicitnessStats"/>
        </v-col>

        <v-col cols="6">
          <ReleaseDateChart :releaseDateData="releaseDateStats"/>
        </v-col>

        <div v-if="playlistName">
          <v-col cols="12">
            <h3>On average, your playlist is {{ generationDetails["age"] }} years old (born {{
                generationDetails["year"]
              }})! This makes it a member of the
              {{ generationDetails["name"] }} generation
              ({{ generationDetails["lower"] }} - {{ generationDetails["upper"] }})...</h3>
            <p>{{ generationDetails["summary"] }}</p>

            <h3>{{ playlistName }} by {{ playlistOwner }}</h3>
            <hr>

            <!-- stats table -->
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
                  <td>{{ stats.avg }}</td>
                </tr>
                </tbody>
              </template>
            </v-simple-table>
          </v-col>
        </div>
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
      playlistIDRules: [
        value => !!value || 'Required.',
        value => (value || '').length === 22 || 'Invalid playlist ID.',
      ],
      playlistName: null,
      playlistOwner: null,
      rawPlaylistStats: {},
      explicitnessStats: {},
      releaseDateStats: {},
      generationDetails: {}
    }
  },
  mounted() {
    this.getPlaylistData();
  },
  methods: {
    getPlaylistData() {
      this.dataError = null;

      // prevent fetching the same data for the previously searched playlist
      if (this.playlistID === this.lastPlaylistID) return;
      if ((this.playlistID || '').length !== 22) return;

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
    }
  },
}
</script>
