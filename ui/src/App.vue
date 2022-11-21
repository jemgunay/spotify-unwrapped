<template>
  <v-app>
    <v-main>
      <Header/>
      <DataView/>
    </v-main>
  </v-app>
</template>

<script>
import Header from "./components/Header";
import DataView from "./components/DataView";

export default {
  name: 'App',
  components: {
    Header,
    DataView,
  }
};
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Pacifico&display=swap');
</style>