<template>
  <v-col cols="12" md="6">
    <!-- generations text -->
    <p class="mb-3 mt-md-12">Your playlist's average track age is <b>{{ generationDetails["age"] }} years old</b> ({{
        generationDetails["year"]
      }}), making it a member of
      <b>{{ generationDetails["name"] }}</b>
      ({{ generationDetails["lower"] }} - {{ generationDetails["upper"] }})...</p>
    <p class="generation-summary">"{{ generationDetails["summary"] }}"</p>

    <v-divider class="mb-3"></v-divider>

    <v-row dense>
      <v-col cols="12" sm="6">
        <TrackStatPanel
            stat-title="Oldest Track"
            :track-name="rawStatsData['release_dates']['min']['name']"
            :track-meta="rawStatsData['release_dates']['min']['value']"
            :cover-image="rawStatsData['release_dates']['min']['cover_image']"
            :spotify-url="rawStatsData['release_dates']['min']['spotify_url']"
        />
      </v-col>
      <v-col cols="12" sm="6">
        <TrackStatPanel
            stat-title="Youngest Track"
            :track-name="rawStatsData['release_dates']['max']['name']"
            :track-meta="rawStatsData['release_dates']['max']['value']"
            :cover-image="rawStatsData['release_dates']['max']['cover_image']"
            :spotify-url="rawStatsData['release_dates']['max']['spotify_url']"
        />
      </v-col>
    </v-row>

  </v-col>
</template>

<script>
import TrackStatPanel from "@/components/TrackStatPanel";

export default {
  name: 'PlaylistGeneration',
  components: {TrackStatPanel},
  props: {
    rawStatsData: {
      type: Object,
      default() {
        return null
      }
    },
    generationDetails: {
      type: Object,
      default() {
        return null
      }
    }
  }
}
</script>

<style>
.generation-summary {
  color: #3c3c3c;
  font-size: 13px;
}
</style>