<template>
  <v-row>
    <v-col cols="12" md="6">
      <v-row>
        <v-col cols="12">
          <h3 class="section-heading">Audio Feature Averages</h3>
        </v-col>

        <v-col cols="10" offset="1">
          <!-- raw audio feature stat averages -->
          <PolarArea v-if="rawStatsData"
                     :chart-options="chartOptions"
                     :chart-data="chartData"
                     chart-id="polar-stats-chart"
                     dataset-id-key="polar-stats"
          />
        </v-col>
      </v-row>
    </v-col>

    <!-- stat panels -->
    <v-col md="6" cols="12">
      <div class="scrollable-container mt-md-12">
        <v-row>
          <v-col cols="12" sm="6" v-for="trackData in trackStatMappings" :key="trackData['name']">
            <TrackStatPanel
                :stat-title="trackData['name']"
                :track-name="trackData['track']['name']"
                :cover-image="trackData['track']['cover_image']"
            />
          </v-col>
        </v-row>
      </div>
    </v-col>
  </v-row>
</template>

<script>
import {PolarArea} from 'vue-chartjs/legacy'
import {ArcElement, Chart as ChartJS, Legend, RadialLinearScale, Title, Tooltip} from 'chart.js'
import {Colours} from '@/helpers/colours'
import TrackStatPanel from "@/components/TrackStatPanel";

ChartJS.register(Title, Tooltip, Legend, ArcElement, RadialLinearScale)

export default {
  name: 'AudioFeaturesPolar',
  components: {
    PolarArea,
    TrackStatPanel
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
          'Liveness',
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
              this.rawStatsData['liveness']['avg']['value'],
              this.rawStatsData['speechiness']['avg']['value'],
              this.rawStatsData['valence']['avg']['value']
            ]
          }
        ]
      }
    },
    trackStatMappings() {
      return [
        {"name": "Least Instrumental", "track": this.rawStatsData['instrumentalness']['min']},
        {"name": "Most Instrumental", "track": this.rawStatsData['instrumentalness']['max']},
        {"name": "Least Acoustic", "track": this.rawStatsData['acousticness']['min']},
        {"name": "Most Acoustic", "track": this.rawStatsData['acousticness']['max']},
        {"name": "Least Danceable", "track": this.rawStatsData['danceability']['min']},
        {"name": "Most Danceable", "track": this.rawStatsData['danceability']['max']},
        {"name": "Least Energetic", "track": this.rawStatsData['energy']['min']},
        {"name": "Most Energetic", "track": this.rawStatsData['energy']['max']},
        {"name": "Least Live", "track": this.rawStatsData['liveness']['min']},
        {"name": "Most Live", "track": this.rawStatsData['liveness']['max']},
        {"name": "Least Vocal", "track": this.rawStatsData['speechiness']['min']},
        {"name": "Most Vocal", "track": this.rawStatsData['speechiness']['max']},
        {"name": "Least Positive", "track": this.rawStatsData['valence']['min']},
        {"name": "Most Positive", "track": this.rawStatsData['valence']['max']}
      ]
    }
  }
}
</script>

<style scoped>
.scrollable-container {
  max-height: 45vw;
  width: 100%;
  overflow-x: hidden;
  padding-right: 5px;
}
</style>