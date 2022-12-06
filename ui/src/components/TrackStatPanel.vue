<template>
  <v-hover v-slot="{ hover }">
    <v-card
        outlined
        :elevation="hover ? 3 : 0"
        v-ripple
        @click="openSpotifyURL"
        class="clickable"
    >
      <v-row align="center" justify="center">
        <v-col cols="8" class="pa-2 list-container">
          <v-list-item three-line>
            <v-list-item-content>
              <v-list-item-title>{{ statTitle }}</v-list-item-title>
              <v-list-item-subtitle>{{ trackName }}</v-list-item-subtitle>
              <v-list-item-subtitle v-if="trackDate">{{ trackDate }}</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-col>

        <v-col cols="4">
          <v-img
              :src="coverImage"
              aspect-ratio="1"
          ></v-img>
        </v-col>
      </v-row>
    </v-card>
  </v-hover>
</template>

<script>
export default {
  name: 'TrackStatPanel',
  props: {
    statTitle: {
      type: String,
      default() {
        return null
      }
    },
    trackName: {
      type: String,
      default() {
        return null
      }
    },
    trackDate: {
      type: String,
      default() {
        return null
      }
    },
    coverImage: {
      type: String,
      default() {
        return "placeholder.jpg"
      }
    },
    spotifyUrl: {
      type: String,
      default() {
        return null
      }
    }
  },
  methods: {
    openSpotifyURL() {
      if (this.spotifyUrl == null) {
        return;
      }
      window.open(this.spotifyUrl, '_blank');
    }
  }
}
</script>

<style>
.list-container {
  min-height: 110px;
}
</style>