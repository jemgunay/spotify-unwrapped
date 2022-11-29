<template>
  <v-col cols="12">
    <v-row>
      <v-col md="6" cols="12">
        <h3 class="section-heading">Track Popularity</h3>

        <v-row>
          <v-col cols="12" sm="6">
            <TrackStatPanel
                stat-title="Least Popular"
                :track-name="rawStatsData['popularity']['min']['name']"
                :cover-image="rawStatsData['popularity']['min']['cover_image']"
            />
          </v-col>
          <v-col cols="12" sm="6">
            <TrackStatPanel
                stat-title="Most Popular"
                :track-name="rawStatsData['popularity']['max']['name']"
                :cover-image="rawStatsData['popularity']['max']['cover_image']"
            />
          </v-col>
        </v-row>
      </v-col>

      <v-col md="6" cols="12">
        <!-- popularity doughnut chart -->
        <Doughnut v-if="this.rawStatsData"
                  :chart-options="chartOptions"
                  :chart-data="chartData"
                  chart-id="popularity-chart"
                  dataset-id-key="popularity"
                  id="explicit-pie"
        />
      </v-col>
    </v-row>
  </v-col>
</template>

<script>
import {Doughnut} from 'vue-chartjs/legacy'
import {ArcElement, CategoryScale, Chart as ChartJS, Legend, Title, Tooltip} from 'chart.js'
import {Green, Orange, Red} from '@/helpers/colours'
import TrackStatPanel from "@/components/TrackStatPanel";

ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale)

export default {
  name: 'PopularityDoughnut',
  components: {
    Doughnut,
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