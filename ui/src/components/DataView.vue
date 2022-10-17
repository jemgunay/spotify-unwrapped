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

      <div v-if="dataError">
        <p>Failed to load data ({{ dataError }}).</p>
      </div>

      <div v-else-if="loading">
        <p>Processing playlist {{ playlistID }}...</p>
      </div>
      <div v-if="playlistName">
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
                v-for="(stats, statName) in playlistStats"
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
      </div>

    </v-col>
  </v-row>
</template>

<script>
import axios from "axios";

export default {
  name: 'DataView',
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
      playlistStats: {}
    }
  },
  mounted() {
    this.getPlaylistData();
  },
  methods: {
    getPlaylistData() {
      if (this.playlistID === this.lastPlaylistID) {
        return;
      }

      if ((this.playlistID || '').length !== 22) {
        return;
      }

      this.loading = true;
      axios
          .get('http://localhost:8080/api/v1/data/playlists/' + this.playlistID)
          .then(response => {
            this.playlistName = response.data["playlist_name"];
            this.playlistOwner = response.data["owner_name"];
            this.playlistStats = response.data["stats"];

            this.lastPlaylistID = this.playlistID;
            this.dataError = null;
          })
          .catch(error => {
            console.log(error);
            this.dataError = error;
            this.lastPlaylistID = "";
          })
          .finally(() => {
            this.loading = false;
          })
    }
  },
}
</script>
