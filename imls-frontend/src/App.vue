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
};
</script>

<template>
  <div>
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
