<template>
  <v-col md="6" sm="12" offset-md="3">
    <h3 class="section-heading">Audio Feature Averages</h3>

    <!-- raw audio feature stat averages -->
    <PolarArea v-if="rawStatsData"
               :chart-options="chartOptions"
               :chart-data="chartData"
               chart-id="polar-stats-chart"
               dataset-id-key="polar-stats"
    />
  </v-col>
</template>

<script>
import {PolarArea} from 'vue-chartjs/legacy'
import {ArcElement, Chart as ChartJS, Legend, RadialLinearScale, Title, Tooltip} from 'chart.js'
import {Colours} from '@/components/helpers/colours'

ChartJS.register(Title, Tooltip, Legend, ArcElement, RadialLinearScale)

export default {
  name: 'PolarAudioFeatures',
  components: {
    PolarArea
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
        // maintainAspectRatio: false,
        scales: {
          r: {
            min: 0,
            max: 100
          }
        }
      }
    }
  },
  computed: {
    chartData() {
      return {
        labels: [
          'Instrumentalness',
          'Acousticness',
          'Danceability',
          'Energy',
          'Speechiness',
          'Positiveness'
        ],
        datasets: [
          {
            backgroundColor: Colours,
            data: [
              this.rawStatsData['instrumentalness']['avg']['value'],
              this.rawStatsData['acousticness']['avg']['value'],
              this.rawStatsData['danceability']['avg']['value'],
              this.rawStatsData['energy']['avg']['value'],
              this.rawStatsData['speechiness']['avg']['value'],
              this.rawStatsData['valence']['avg']['value']
            ]
          }
        ]
      }
    },
  }
}
</script>
