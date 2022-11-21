<template>
  <v-col md="7" sm="12">
    <h3>Track Release Dates</h3>

    <!-- release date chart -->
    <Bar v-if="this.releaseDateData"
         :chart-options="chartOptions"
         :chart-data="chartData"
         chart-id="release-date-chart"
         dataset-id-key="releaseDate"
    />
  </v-col>
</template>

<script>
import {Bar} from 'vue-chartjs/legacy'
import {BarElement, CategoryScale, Chart as ChartJS, Legend, LinearScale, Title, Tooltip} from 'chart.js'
import {Green} from '@/components/helpers/helpers'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

export default {
  name: 'ReleaseDateChart',
  components: {
    Bar
  },
  props: {
    releaseDateData: {
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
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: false
          }
        }
      },
    }
  },
  computed: {
    chartData() {
      return {
        labels: this.releaseDateData["keys"],
        datasets: [
          {
            backgroundColor: [Green],
            data: this.releaseDateData["values"]
          }
        ]
      }
    },
  }
}
</script>
