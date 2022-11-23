<template>
  <v-app>
    <v-main>
      <Header/>
      <DataView/>
      <Footer/>
    </v-main>
  </v-app>
</template>

<script>
import Header from "./components/Header";
import DataView from "./components/DataView";
import Footer from "./components/Footer";

export default {
  name: 'App',
  components: {
    Header,
    DataView,
    Footer
  }
};
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Pacifico&display=swap');
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@600&display=swap');
</style>