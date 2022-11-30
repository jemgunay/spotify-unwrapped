<template>
  <v-card outlined>
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
        return null
      }
    }
  }
}
</script>

<style>
.list-container {
  min-height: 110px;
}
</style>