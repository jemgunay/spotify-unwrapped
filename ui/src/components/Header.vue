<template>
  <v-row>
    <v-col sm="12">
      <h1 class="header-title mt-14" :style="{'color':titleColour}">
        Spotify Unwrapped
      </h1>
    </v-col>
  </v-row>
</template>

<script>
import {Green} from "@/components/helpers/colours";

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
</style>