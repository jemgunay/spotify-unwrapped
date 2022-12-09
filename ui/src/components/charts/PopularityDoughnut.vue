<template>
  <v-col cols="12">
    <v-row>
      <v-col cols="12" md="3">
        <h3 class="section-heading">Popularity</h3>

        <v-row>
          <v-col cols="12" sm="12">
            <TrackStatPanel
                stat-title="Least Popular"
                :track-name="rawStatsData['popularity']['min']['name']"
                :cover-image="rawStatsData['popularity']['min']['cover_image']"
                :spotify-url="rawStatsData['popularity']['min']['spotify_url']"
            />
          </v-col>
          <v-col cols="12" sm=12>
            <TrackStatPanel
                stat-title="Most Popular"
                :track-name="rawStatsData['popularity']['max']['name']"
                :cover-image="rawStatsData['popularity']['max']['cover_image']"
                :spotify-url="rawStatsData['popularity']['max']['spotify_url']"
            />
          </v-col>
        </v-row>
      </v-col>

      <v-col cols="10" offset="1" md="6" offset-md="0">
        <!-- popularity doughnut chart -->
        <h1 class="doughnut-stat" :class="$vuetify.breakpoint.smAndDown ? 'sm' : 'md'"
            :style="{'color': determineHealthColour}">
          {{ this.roundedPopularity }}%
        </h1>
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
        plugins: {
          tooltip: {
            enabled: false
          }
        }
      },
    }
  },
  computed: {
    chartData() {
      return {
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
      if (val < 30) {
        return Red;
      } else if (val < 60) {
        return Orange;
      }
      return Green;
    },
    roundedPopularity() {
      return ~~this.rawStatsData["popularity"]["avg"]["value"];
    }
  }
}
</script>

<style>
.doughnut-stat {
  height: 0;
  width: 0;
  text-shadow: 2px 2px 0px #1e7b52;
  font-size: 9vw;
  position: relative;
}

/* styles to handle dynamic scaling of the polar chart stat text based on breakpoints - its gross, but it works... */
.doughnut-stat.md {
  left: calc(33%);
  top: calc(36%);
}

.doughnut-stat.sm {
  left: calc(39%);
  top: calc(41%);
}
</style>