<template>
  <v-col cols="12" md="6">
    <h3 class="section-heading">Track Releases by Year</h3>

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
import {Green} from '@/helpers/colours'

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
        },
        scales: {
          x: {
            title: {
              display: true,
              text: 'Year Released'
            }
          },
          y: {
            title: {
              display: true,
              text: 'Track Count'
            }
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