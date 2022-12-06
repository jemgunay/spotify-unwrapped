<template>
  <v-row>
    <v-col cols="12">
      <h3 class="section-heading">Explicit vs Non-Explicit Lyrics</h3>
    </v-col>

    <v-col cols="10" offset="1">
      <!-- explicit lyrics pie chart -->
      <Pie v-if="this.explicitnessData"
           :chart-options="chartOptions"
           :chart-data="chartData"
           chart-id="explicitness-chart"
           dataset-id-key="explicitness"
           id="explicit-pie"
      />
    </v-col>
  </v-row>
</template>

<script>
import {Pie} from 'vue-chartjs/legacy'
import {ArcElement, CategoryScale, Chart as ChartJS, Legend, Title, Tooltip} from 'chart.js'
import {Green, Red} from '@/helpers/colours'

ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale)

export default {
  name: 'ExplicitPieChart',
  components: {
    Pie
  },
  props: {
    explicitnessData: {
      type: Object,
      default() {
        return null
      }
    }
  },
  data() {
    return {
      chartOptions: {
        responsive: true,
        // maintainAspectRatio: false
      }
    }
  },
  computed: {
    chartData() {
      return {
        labels: ['Explicit', 'Non-Explicit'],
        datasets: [
          {
            backgroundColor: [Red, Green],
            data: [this.explicitnessData["explicit"], this.explicitnessData["non-explicit"]]
          }
        ]
      }
    },
  }
}
</script>