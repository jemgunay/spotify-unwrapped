<template>
  <v-col md="5" offset-md="1" sm="12">
    <h3>Explicit vs Non-Explicit Lyrics</h3>

    <!-- explicit lyrics pie chart -->
    <Pie v-if="this.explicitnessData"
         :chart-options="chartOptions"
         :chart-data="chartData"
         chart-id="explicitness-chart"
         dataset-id-key="explicitness"
         id="explicit-pie"
    />
  </v-col>
</template>

<script>
import {Pie} from 'vue-chartjs/legacy'
import {ArcElement, CategoryScale, Chart as ChartJS, Legend, Title, Tooltip} from 'chart.js'
import {Green, Red} from '@/components/helpers/helpers'

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