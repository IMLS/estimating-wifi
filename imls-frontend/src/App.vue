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
  }
};
</script>

<template>

  <div>
    <metainfo>
      <template v-slot:title="{ content }">{{ content }} | Public Library Wifi Estimator</template>
    </metainfo>
    <USWDSHeader />

    <div class="usa-section padding-top-1">
      <div class="grid-container">
        <div class="grid-row grid-gap">
          <main
            id="main-content"
            class="grid-col"
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
