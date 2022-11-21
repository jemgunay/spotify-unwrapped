<template>
  <v-col md="12">
    <h3>Raw Stats</h3>

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
            v-for="(stats, statName) in rawStatsData"
            :key="statName"
        >
          <td>{{ statName }}</td>
          <td>{{ stats['min']['name'] }} ({{ stats['min']['value'] }})</td>
          <td>{{ stats['max']['name'] }} ({{ stats['max']['value'] }})</td>
          <td>{{ stats['avg']['value'] }}</td>
        </tr>
        </tbody>
      </template>
    </v-simple-table>
  </v-col>
</template>

<script>
export default {
  name: 'RawStatsTable',
  props: {
    rawStatsData: {
      type: Object,
      default() {
        return null
      }
    }
  }
}
</script>
