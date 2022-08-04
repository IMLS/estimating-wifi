<script>
import { store } from "@/store/store.js";
import Histogram from "@/components/Histogram.vue";
import USWDSCard from "@/components/USWDSCard.vue";

export default {
  name: 'Single Sensor',
  components: { Histogram, USWDSCard },
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
    <div class="usa-card-group margin-top-6">
      <div class="usa-card tablet:grid-col-12">
        <USWDSCard title="Devices present by weekday" no-footer>
          <Histogram bins="weekdays"/>
        </USWDSCard>
      </div>
      <div class="usa-card tablet:grid-col-12">
        <USWDSCard title="Devices present by month" no-footer>
          <Histogram bins="months"/>
        </USWDSCard>
      </div>
    </div>
  </div>
</template>

<style></style>
