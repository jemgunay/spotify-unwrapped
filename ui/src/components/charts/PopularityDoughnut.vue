<template>
  <v-col md="6" offset-md="3" sm="12">
    <h3 class="section-heading">Track Average Popularity</h3>

    <!-- popularity doughnut chart -->
    <Doughnut v-if="this.rawStatsData"
              :chart-options="chartOptions"
              :chart-data="chartData"
              chart-id="popularity-chart"
              dataset-id-key="popularity"
              id="explicit-pie"
    />
  </v-col>
</template>

<script>
import {Doughnut} from 'vue-chartjs/legacy'
import {ArcElement, CategoryScale, Chart as ChartJS, Legend, Title, Tooltip} from 'chart.js'
import {Green, Orange, Red} from '@/components/helpers/colours'

ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale)

export default {
  name: 'PopularityDoughnut',
  components: {
    Doughnut
  },
  props: {
    rawStatsData: {
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
      },
    }
  },
  computed: {
    chartData() {
      return {
        // labels: ['Popularity'],
        datasets: [
          {
            backgroundColor: [this.determineHealthColour, '#d0d0d0'],
            data: [
              this.rawStatsData["popularity"]["avg"]["value"],
              100 - this.rawStatsData["popularity"]["avg"]["value"]
            ]
          }
        ]
      }
    },
    determineHealthColour() {
      let val = this.rawStatsData["popularity"]["avg"]["value"];
      console.log(val)
      if (val < 33) {
        return Red;
      } else if (val < 66) {
        return Orange;
      }
      return Green;
    }
  }
}
</script>