<template>
  <v-app>
    <NavBar/>

    <v-main>
      <DataView/>
    </v-main>
  </v-app>
</template>

<script>
import NavBar from "./components/NavBar";
import DataView from "./components/DataView";

export default {
  name: 'App',
  components: {
    NavBar,
    DataView,
  }
};
</script>
