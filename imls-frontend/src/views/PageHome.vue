<script>
import { RouterLink } from "vue-router";
import { store } from "@/store/store.js";

export default {
  name: "HomePage",
  components: { RouterLink },
  data() {
    return {
      store,
    };
  },
  metaInfo: {
    title: "Homepage",
  }
}
</script>

<template>
  <div>
    <h1 id="pageTitle">Home</h1>

    <h2> All States</h2>

    <ul
      class="usa-list"
    >
      <li v-for="(stateName, stateAbbr) in store.states" :key="stateName" class="">
        <RouterLink class="usa-link" :to="{ path: '/state/' + stateAbbr + '/'}">
          {{ stateName }}
        </RouterLink>
      </li>
      
    </ul>
  </div>
</template>
