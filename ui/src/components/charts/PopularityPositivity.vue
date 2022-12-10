<template>
  <v-col cols="12">
    <v-row>
      <!-- popularity -->
      <v-col cols="12" md="6">
        <h3 class="section-heading">Popularity</h3>
        <v-row dense>
          <v-col cols="10" offset="1" md="12" offset-md="0" class="mb-3">
            <!-- popularity doughnut chart -->
            <h1 class="doughnut-stat" :class="$vuetify.breakpoint.smAndDown ? 'sm' : 'md'"
                :style="{'color': determineHealthColour}">
              {{ this.roundedPopularity }}%
            </h1>
            <Doughnut v-if="this.rawStatsData"
                      :chart-options="popularityDoughnutOptions"
                      :chart-data="popularityDoughnutChartData"
                      chart-id="popularity-doughnut"
                      dataset-id-key="popularity-doughnut"
            />
          </v-col>

          <v-col cols="12" sm="6">
            <TrackStatPanel
                stat-title="Least Popular"
                :track-name="rawStatsData['popularity']['min']['name']"
                :cover-image="rawStatsData['popularity']['min']['cover_image']"
                :spotify-url="rawStatsData['popularity']['min']['spotify_url']"
            />
          </v-col>

          <v-col cols="12" sm=6>
            <TrackStatPanel
                stat-title="Most Popular"
                :track-name="rawStatsData['popularity']['max']['name']"
                :cover-image="rawStatsData['popularity']['max']['cover_image']"
                :spotify-url="rawStatsData['popularity']['max']['spotify_url']"
            />
          </v-col>
        </v-row>

      </v-col>

      <!-- positivity -->
      <v-col cols="12" md="6">
        <v-row dense>
          <v-col cols="12">
            <h3 class="section-heading">Positivity</h3>
            <p class="mb-0">The average track positivity is <strong>{{
                rawStatsData['valence']['avg']['value']
              }}%</strong>.
              <v-icon size="34px" :style="{'color': getPositivityColor}">{{ getPositivityIcon }}</v-icon>
            </p>
          </v-col>

          <v-col cols="12" sm="6">
            <TrackStatPanel
                stat-title="Least Positive"
                :track-name="rawStatsData['valence']['min']['name']"
                :cover-image="rawStatsData['valence']['min']['cover_image']"
                :spotify-url="rawStatsData['valence']['min']['spotify_url']"
            />
          </v-col>
          <v-col cols="12" sm=6>
            <TrackStatPanel
                stat-title="Most Positive"
                :track-name="rawStatsData['valence']['max']['name']"
                :cover-image="rawStatsData['valence']['max']['cover_image']"
                :spotify-url="rawStatsData['valence']['max']['spotify_url']"
            />
          </v-col>

          <v-col cols="12" class="mt-3">
            <Bubble v-if="this.popularityGraphChartData"
                    :chart-options="popularityGraphOptions"
                    :chart-data="popularityGraphChartData"
                    chart-id="popularity-chart"
                    dataset-id-key="popularity-chart"
                    id="popularity-chart"
            />
          </v-col>
        </v-row>
      </v-col>

    </v-row>
  </v-col>
</template>

<script>
import {Bubble, Doughnut} from 'vue-chartjs/legacy'
import {ArcElement, CategoryScale, Chart as ChartJS, Legend, LinearScale, PointElement, Title, Tooltip} from 'chart.js'
import {Blue, Green, Orange, Red} from '@/helpers/colours'
import TrackStatPanel from "@/components/TrackStatPanel";

ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale, PointElement, LinearScale)

export default {
  name: 'PopularityPositivity',
  components: {
    Doughnut,
    Bubble,
    TrackStatPanel
  },
  props: {
    rawStatsData: {
      type: Object,
      default() {
        return null
      }
    },
    positivityGraphData: {
      type: Array,
      default() {
        return null
      }
    }
  },
  data() {
    return {
      popularityDoughnutOptions: {
        responsive: true,
        plugins: {
          tooltip: {
            enabled: false
          },
          title: {
            display: true,
            text: 'Average Track Popularity'
          }
        }
      },
      popularityGraphOptions: {
        responsive: true,
        plugins: {
          tooltip: {
            enabled: false
          },
          legend: {
            display: false
          },
          title: {
            display: true,
            text: 'Positivity vs Popularity vs Energy (Blob Size)'
          }
        },
        scales: {
          x: {
            title: {
              display: true,
              text: 'Positivity (%)'
            }
          },
          y: {
            title: {
              display: true,
              text: 'Popularity (%)'
            }
          }
        }
      },
    }
  },
  computed: {
    popularityDoughnutChartData() {
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
    popularityGraphChartData() {
      return {
        datasets: [
          {
            backgroundColor: [Green],
            data: this.positivityGraphData
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
    },
    getPositivityIcon() {
      let val = this.rawStatsData["valence"]["avg"]["value"];
      if (val < 25) {
        return "mdi-emoticon-sad-outline";
      } else if (val < 50) {
        return "mdi-emoticon-neutral-outline";
      } else if (val < 75) {
        return "mdi-emoticon-happy-outline";
      }
      return "mdi-emoticon-excited-outline";
    },
    getPositivityColor() {
      let val = this.rawStatsData["valence"]["avg"]["value"];
      if (val < 25) {
        return Red;
      } else if (val < 50) {
        return Orange;
      } else if (val < 75) {
        return Blue;
      }
      return Green;
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

/* handle dynamic scaling of the polar chart stat text based on breakpoints - its gross, but it works... */
.doughnut-stat.md {
  left: calc(32%);
  top: calc(38%);
}

.doughnut-stat.sm {
  left: calc(39%);
  top: calc(46%);
}
</style>