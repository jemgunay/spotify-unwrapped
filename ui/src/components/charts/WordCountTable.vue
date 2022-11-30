<template>
  <v-col md="5" cols="12">
    <h3 class="section-heading">Common Track Title Words</h3>

    <!-- title word count stats table -->
    <v-simple-table dense class="scrollable-container">
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
  name: 'WordCountTable',
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

<style scoped>
.scrollable-container {
  max-height: 35vw;
}
</style>