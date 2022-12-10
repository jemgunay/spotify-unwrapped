<template>
  <v-app>
    <v-main>
      <v-container class="py-10 px-10" fluid>
        <Header/>
        <DataView/>
      </v-container>

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
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@500&display=swap');

h3.section-heading {
  margin-bottom: 15px;
}

.section-heading, strong {
  font-family: 'Inter', sans-serif;
  color: #3c3c3c;
}

.clickable {
  cursor: pointer;
  opacity: 1;
  transition: opacity 0.3s;
}

.clickable:hover {
  opacity: 0.8;
}

/* force show scroll bars on OSx */
.scrollable-container {
  overflow: auto;
}

.scrollable-container::-webkit-scrollbar {
  -webkit-appearance: none;
  width: 4px;
}

.scrollable-container::-webkit-scrollbar-thumb {
  border-radius: 4px;
  background-color: rgba(0, 0, 0, .5);
  box-shadow: 0 0 1px rgba(255, 255, 255, .5);
}
</style>