<template>
  <v-col cols="12">
    <v-row>
      <v-col cols="12" md="6" offset-md="0">
        <h3 class="section-heading">Key Signatures</h3>

        <Bar v-if="this.pitchKeyData"
             :chart-options="chartOptions"
             :chart-data="chartData"
             chart-id="pitch-key-chart"
             dataset-id-key="pitchKeyChart"
        />
      </v-col>

      <v-col cols="12" md="6">

        <v-row dense>
          <v-col cols="12">
            <h3 class="section-heading">Tempo</h3>
            <p class="mb-0">The average track tempo is <strong>{{ rawStatsData['tempo']['avg']['value'] }} BPM</strong>.
            </p>
          </v-col>
          <v-col cols="12" sm="6">
            <TrackStatPanel
                stat-title="Lowest BPM"
                :track-name="rawStatsData['tempo']['min']['name'] + ' (' + rawStatsData['tempo']['min']['value'] + ' BPM)'"
                :cover-image="rawStatsData['tempo']['min']['cover_image']"
                :spotify-url="rawStatsData['tempo']['min']['spotify_url']"
            />
          </v-col>
          <v-col cols="12" sm=6>
            <TrackStatPanel
                stat-title="Highest BPM"
                :track-name="rawStatsData['tempo']['max']['name'] + ' (' + rawStatsData['tempo']['max']['value'] + ' BPM)'"
                :cover-image="rawStatsData['tempo']['max']['cover_image']"
                :spotify-url="rawStatsData['tempo']['max']['spotify_url']"
            />
          </v-col>
        </v-row>

        <v-row dense>
          <v-col cols="12">
            <h3 class="section-heading mt-3">Track Length</h3>
            <p class="mb-0">The average track length is <strong>{{
                rawStatsData['track_durations']['avg']['value']
              }}</strong>.</p>
          </v-col>
          <v-col cols="12" sm="6">
            <TrackStatPanel
                stat-title="Shortest Track"
                :track-name="rawStatsData['track_durations']['min']['name'] + ' (' + rawStatsData['track_durations']['min']['value'] + ')'"
                :cover-image="rawStatsData['track_durations']['min']['cover_image']"
                :spotify-url="rawStatsData['track_durations']['min']['spotify_url']"
            />
          </v-col>
          <v-col cols="12" sm=6>
            <TrackStatPanel
                stat-title="Longest Track"
                :track-name="rawStatsData['track_durations']['max']['name'] + ' (' + rawStatsData['track_durations']['max']['value'] + ')'"
                :cover-image="rawStatsData['track_durations']['max']['cover_image']"
                :spotify-url="rawStatsData['track_durations']['max']['spotify_url']"
            />
          </v-col>
        </v-row>

      </v-col>

    </v-row>
  </v-col>
</template>

<script>
import {Bar} from 'vue-chartjs/legacy'
import {BarElement, CategoryScale, Chart as ChartJS, Legend, LinearScale, Title, Tooltip} from 'chart.js'
import {Green} from '@/helpers/colours'
import TrackStatPanel from "@/components/TrackStatPanel";

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

export default {
  name: 'MusicalComposition',
  components: {
    Bar,
    TrackStatPanel
  },
  props: {
    rawStatsData: {
      type: Object,
      default() {
        return null
      }
    },
    pitchKeyData: {
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
        }
      },
    }
  },
  computed: {
    chartData() {
      return {
        labels: this.pitchKeyData["keys"],
        datasets: [
          {
            backgroundColor: [Green],
            data: this.pitchKeyData["values"]
          }
        ]
      }
    }
  }
}
</script>