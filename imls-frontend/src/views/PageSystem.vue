<script>
import { store } from "@/store/store.js";

import FetchData from "@/components/FetchData.vue";

export default {
  name: 'Single System',
  components: {FetchData },
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

  methods: {
  },
};
</script>

<template>
  <div>
    <h1>(FSCS) ID {{ id }}</h1>
    <FetchData :fscs-id=id />
  
  </div>
</template>

<style></style>
