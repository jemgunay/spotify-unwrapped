<template>
  <Bar v-if="this.releaseDateData"
       :chart-options="chartOptions"
       :chart-data="chartData"
       chart-id="release-date-chart"
       dataset-id-key="releaseDate"
  />
</template>

<script>
import {Bar} from 'vue-chartjs/legacy'
import {BarElement, CategoryScale, Chart as ChartJS, Legend, LinearScale, Title, Tooltip} from 'chart.js'

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
        maintainAspectRatio: false
      }
    }
  },
  computed: {
    chartData() {
      return {
        labels: this.releaseDateData["keys"],
        datasets: [
          {
            backgroundColor: ['#41B883'],
            data: this.releaseDateData["values"]
          }
        ]
      }
    },
  }
}
</script>
