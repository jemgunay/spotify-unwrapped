<template>
  <v-col md="5" sm="12">
    <h3>Common Words in Track Titles</h3>

    <!-- title word count stats table -->
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
</template>

<script>
export default {
  name: 'WordCountStats',
  props: {
    topTitleWords: {
      type: Object,
      default() {
        return null
      }
    }
  }
}
</script>

<style>
.scrollable-table {
  max-height: 60vh;
  overflow: auto;
}
</style>