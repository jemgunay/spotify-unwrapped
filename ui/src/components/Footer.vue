<template>
  <v-card height="200px">
    <v-footer
        padless
        absolute
    >
      <v-card
          flat
          tile
          width="100%"
          class="lighten-1 text-center"
          :style="{'background-color':backgroundColour}"
      >
        <v-card-text>
          <!-- scroll to top of page button -->
          <v-btn
              class="mx-4"
              icon
              @click.native="scrollToTop()"
          >
            <v-icon size="24px">mdi-arrow-up-bold</v-icon>
          </v-btn>
          <!-- GitHub link -->
          <v-btn
              class="mx-4"
              icon
              href="https://github.com/jemgunay/spotify-unwrapped"
              target="_blank"
          >
            <v-icon size="24px">mdi-github</v-icon>
          </v-btn>
        </v-card-text>

        <v-divider></v-divider>

        <v-card-text class="white--text">
          <h3 class="my-1">Spotify Unwrapped by Jem Gunay</h3>
        </v-card-text>
      </v-card>
    </v-footer>
  </v-card>
</template>

<script>
import {Green} from "@/helpers/colours";

export default {
  name: 'Footer',
  data() {
    return {
      backgroundColour: Green
    }
  },
  methods: {
    scrollToTop() {
      window.scrollTo({
        top: 0,
        left: 0,
        behavior: 'smooth'
      });
    },
  }
}
</script>

<style scoped>
h3 {
  font-family: 'Pacifico', cursive;
  font-weight: 1;
}
</style>