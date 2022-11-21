<template>
  <v-col md="5" sm="12">
    <!-- generations text -->
    <h4 class="mb-3">Your playlist's average age is {{ generationDetails["age"] }} years old (born {{
        generationDetails["year"]
      }})! This makes it a member of
      {{ generationDetails["name"] }}
      ({{ generationDetails["lower"] }} - {{ generationDetails["upper"] }})...</h4>
    <p class="generation-summary">"{{ generationDetails["summary"] }}"</p>

    <v-divider class="mt-7 mb-4"></v-divider>

    <v-list-item three-line>
      <v-list-item-content>
        <v-list-item-title>Oldest Track</v-list-item-title>
        <v-list-item-subtitle>{{ rawStatsData["releaseDates"]["min"]["name"] }}</v-list-item-subtitle>
        <v-list-item-subtitle>{{ rawStatsData["releaseDates"]["min"]["date"] }}</v-list-item-subtitle>
      </v-list-item-content>

      <v-list-item-content>
        <v-list-item-title>Youngest Track</v-list-item-title>
        <v-list-item-subtitle>{{ rawStatsData["releaseDates"]["max"]["name"] }}</v-list-item-subtitle>
        <v-list-item-subtitle>{{ rawStatsData["releaseDates"]["max"]["date"] }}</v-list-item-subtitle>
      </v-list-item-content>
    </v-list-item>
  </v-col>
</template>

<script>
export default {
  name: 'PlaylistGeneration',
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
  color: rgba(0, 0, 0, 0.6);
}
</style>