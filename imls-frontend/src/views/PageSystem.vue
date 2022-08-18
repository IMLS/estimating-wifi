<script>
import { store } from "@/store/store.js";

import FetchData from "@/components/FetchData.vue";
import USWDSCard from "@/components/USWDSCard.vue";

export default {
  name: 'Single System',
  components: {FetchData, USWDSCard },

  props: {
    id: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      store,
      startDate: '2022-05-10',
      possibleStartDates: [
        '2022-05-10',
        '2022-05-11',
        '2022-05-12',
        '2022-05-13',
        '2022-05-14'
      ]
    }
  },
  methods: {},
};
</script>

<template>
  <div>
    <h1>{{ id }}</h1>


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

    <div class="usa-card-group margin-top-6">
      <div class="usa-card tablet:grid-col-12">
        <USWDSCard title="Devices present by hour" no-footer>
          <FetchData :fscs-id=id :start-date="startDate"/>
        </USWDSCard>
      </div>
    </div>

  
  </div>
</template>

<style></style>
