<script>
import USWDSContent from "@/components/USWDSContent.vue";
import { RouterLink } from "vue-router";


export default {
  components: {
    USWDSContent,
  },
  data() {
    return {
      multilineContent: "I'm the first paragraph. I get special styling in USWDS Content. \nI'm a newline. Any new lines after me are parsed as regular paragraphs. I don't know markdown yet.",
    };
  },
  metaInfo: {
    title: 'About',
  }
}
</script>

<template>
  <div class="about">
    <h1>This is an about page</h1>
    <USWDSContent :multiline-content="multilineContent"/>
     <RouterLink
        to="/"
        title="Public Library Wifi Estimator"
      > Want to go back home? </RouterLink>
  </div>
</template>

<style>
</style>
