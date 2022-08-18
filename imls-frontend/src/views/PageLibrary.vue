<script>
import { state } from "@/store/store.js";
import { format } from "date-fns";

import FetchData from "@/components/FetchData.vue";
import USWDSCard from "@/components/USWDSCard.vue";
import USWDSDatePicker from "@/components/USWDSDatePicker.vue";

export default {
  name: 'Single Library',
  components: {FetchData, USWDSCard, USWDSDatePicker },

  props: {
    id: {
      type: String,
      required: true,
      default: ''
    },
  },
  data() {
    return {
      state,
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
  computed: {
    chartTitle: () => {
      const localDate = state.selectedDate + "T00:00";
      return "Devices present by hour on " + format(new Date(localDate), 'PP')
    }
  }
};
</script>

<template>
  <div>
    <h1>Library {{ id }}</h1>

  <USWDSDatePicker :initialDate=state.selectedDate />

    <div class="usa-card-group margin-top-6">
      <div class="usa-card tablet:grid-col-12">
        <USWDSCard :title="chartTitle" no-footer>
          <div class="grid-row">
            <div class="grid-col">
              <FetchData :fscs-id=id :start-date="state.selectedDate"/>
            </div>
          </div>
  
          <div class="grid-row">
            <div class="grid-col maxw-tablet-lg margin-top-2">
              <p>This graph depicts all sensed wifi-enabled devices within range of the selected sensor(s), according to local time, if detected for at least 5 continuous minutes during each hour. Other explanatory text may show up here.</p>
            </div>
          </div>
        </USWDSCard>
      </div>
    </div>

  
  </div>
</template>

<style></style>
