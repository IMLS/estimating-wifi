<script>
import { RouterView } from "vue-router";
import "uswds";

// import ComponentName from './components/ComponentName.vue'
import USWDSHeader from "./components/USWDSHeader.vue";
import USWDSFooter from "./components/USWDSFooter.vue";

export default {
  components: {
    USWDSHeader,
    USWDSFooter,
  },
  metaInfo: {
    title: "Welcome",
    description: "This tool is the result of a collaborative research project with the Institute of Museum and Library Services, the Georgia Public Library Service, and Technology Transformation Services' 10x program."
  },
  watch: {
    $route: function() {
      this.setRouteWrapperFocus();
    }
  },
  methods:{
    setRouteWrapperFocus() {
      this.$nextTick(function() {
        this.$refs.focus.focus();
      });
    },
  }
};
</script>

<template>

  <div>
    <metainfo>
      <template #title="{ content }">{{ content }} | Public Library Wifi Estimator</template>
    </metainfo>
    <USWDSHeader />

    <div class="usa-section padding-top-1">
      <div class="grid-container">
        <div class="grid-row grid-gap">
          <main
            id="main-content"
            ref="focus"
            class="grid-col"
            tabindex="-1"
          >
            <RouterView />
          </main>
        </div>
      </div>
    </div>

    <USWDSFooter />
  </div>
</template>

<style lang="sass">
@import "@/assets/index.scss"
</style>
