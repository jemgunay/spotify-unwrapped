<template>
  <v-row>
    <v-col cols="12">
      <h1 class="header-title" :style="{'color':titleColour}">
        Spotify Unwrapped
        <v-icon size="48px" :style="{'color':titleColour}">mdi-gift-open-outline</v-icon>
      </h1>
    </v-col>
  </v-row>
</template>

<script>
import {Green} from "@/helpers/colours";

export default {
  name: 'Header',
  data() {
    return {
      titleColour: Green
    }
  },
}
</script>

<style scoped>
h1.header-title {
  font-family: 'Pacifico', cursive;
  font-size: 50px;
  text-align: center;
  font-weight: 1;
  text-shadow: 2px 2px 0px #1e7b52;
}

h1.header-title .v-icon {
  text-shadow: 2px 2px 0px #1e7b52;
  position: relative;
  top: -8px;
}
</style>