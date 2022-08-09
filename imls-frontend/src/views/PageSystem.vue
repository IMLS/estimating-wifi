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
      startDate: '2022-05-01',
      possibleStartDates: [
        '2022-05-01',
        '2022-05-02',
        '2022-05-03',
        '2022-05-11',
        '2022-05-12'
      ]
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


    <fieldset class="usa-fieldset">
      <legend class="usa-legend">Choose a date to view data for</legend>
        <div class="grid-row grid-gap-05">
          <div v-bind:key="choice" v-for="choice in possibleStartDates" class="usa-radio tablet:grid-col">
            <input
              class="usa-radio__input usa-radio__input--tile grid"
              :id="choice"
              type="radio"
              :name="choice"
              :value="choice"
              v-model="startDate"
            />
            <label class="usa-radio__label" :for="choice"
              >{{ choice }}</label
            >
          </div>
  

      </div>
    </fieldset>

    <FetchData :fscs-id=id :start-date="startDate"/>
  
  </div>
</template>

<style></style>
