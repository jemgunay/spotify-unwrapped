<template>
  <v-row>
    <v-col cols="12" class="pa-10">
      <div v-if="dataError">
        <p>Failed to load data ({{ dataError }}).</p>
      </div>


      <div v-else-if="!loaded">
        <p>Loading data for playlist {{ playlistID }}...</p>
      </div>
      <div v-else>
        <h3>{{ playlistName }} by {{ playlistOwner }}</h3>
        <hr>

        <v-simple-table>
          <template v-slot:default>
            <thead>
            <tr>
              <th class="text-left">
                Stat
              </th>
              <th class="text-left">
                Min
              </th>
              <th class="text-left">
                Max
              </th>
              <th class="text-left">
                Avg
              </th>
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
      loaded: false,

      playlistID: "1AXy6ag2d0ag8DEdOE7kWm",
      playlistName: "",
      playlistOwner: "",
      playlistStats: {}
    }
  },
  mounted() {
    axios
        .get('http://localhost:8080/api/v1/data/playlists/' + this.playlistID)
        .then(response => {
          this.playlistName = response.data["playlist_name"];
          this.playlistOwner = response.data["owner_name"];
          this.playlistStats = response.data["stats"];
          this.loaded = true;
        })
        .catch(error => {
          console.log(error);
          this.dataError = error;
        })
  }
}
</script>
