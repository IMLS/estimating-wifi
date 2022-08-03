<script>
import { store } from "@/store/store.js";

export default {
  // use one or the other of these
  async beforeRouteUpdate(to, from) {
    // react to route changes...
    this.query = to.params.id;
  },
  props: {
    id: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      store,
    }
  },
  // created() {
  //   this.$watch(
  //     () => this.$route.params,
  //     (toParams, previousParams) => {
  //       // react to route changes...
  //     }
  //   )
  // },
  computed: {
    currentSensor() {
      return store.sensors.find((sensor) => sensor.id == this.id);
    },
  },
  methods: {
  },
};
</script>

<template>
  <div>
    <h1>Page for Single Sensor: {{ id }}</h1>
  </div>
</template>

<style></style>
