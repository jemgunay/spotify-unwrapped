<template>
  <Pie v-if="this.explicitnessData"
       :chart-options="chartOptions"
       :chart-data="chartData"
       chart-id="explicitness-chart"
       dataset-id-key="explicitness"
  />
</template>

<script>
import {Pie} from 'vue-chartjs/legacy'
import {ArcElement, CategoryScale, Chart as ChartJS, Legend, Title, Tooltip} from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale)

export default {
  name: 'ExplicitChart',
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
        maintainAspectRatio: false
      }
    }
  },
  computed: {
    chartData() {
      return {
        labels: ['Explicit', 'Non-Explicit'],
        datasets: [
          {
            backgroundColor: ['#E46651', '#41B883'],
            data: [this.explicitnessData["explicit"], this.explicitnessData["non-explicit"]]
          }
        ]
      }
    },
  }
}
</script>
