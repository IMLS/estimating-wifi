<script>
import { RouterLink } from "vue-router";
import { store } from "@/store/store.js";

export default {
  name: "USWDS Header",
  components: { RouterLink },
  data() {
    return {
      store,
    };
  },
}
</script>

<template>
  <main>
    <h1>Home</h1>
    <p>I am a homepage</p>

    <h2> All States</h2>

    <ul
      class="usa-list"
    >
      <li v-for="(stateName, stateAbbr) in store.states" :key="stateName" class="">
        <RouterLink class="usa-link" :to="{ path: '/state/' + stateAbbr + '/' , query: $route.query}">
          {{ stateName }}
        </RouterLink>
      </li>
      
    </ul>
  </main>
</template>
