<template>
  <v-col cols="12" md="6">
    <!-- generations text -->
    <p class="mb-3 mt-md-12">Your playlist's average track age is <b>{{ generationDetails["age"] }} years old</b> ({{
        generationDetails["year"]
      }}), making it a member of
      <b>{{ generationDetails["name"] }}</b>
      ({{ generationDetails["lower"] }} - {{ generationDetails["upper"] }})...</p>
    <p class="generation-summary">"{{ generationDetails["summary"] }}"</p>

    <v-divider class="mt-6 mb-6"></v-divider>

    <v-row>
      <v-col cols="12" sm="6">
        <TrackStatPanel
            stat-title="Oldest Track"
            :track-name="rawStatsData['releaseDates']['min']['name']"
            :track-date="rawStatsData['releaseDates']['min']['date']"
            :cover-image="rawStatsData['releaseDates']['min']['cover_image']"
            :spotify-url="rawStatsData['releaseDates']['min']['spotify_url']"
        />
      </v-col>
      <v-col cols="12" sm="6">
        <TrackStatPanel
            stat-title="Youngest Track"
            :track-name="rawStatsData['releaseDates']['max']['name']"
            :track-date="rawStatsData['releaseDates']['max']['date']"
            :cover-image="rawStatsData['releaseDates']['max']['cover_image']"
            :spotify-url="rawStatsData['releaseDates']['max']['spotify_url']"
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