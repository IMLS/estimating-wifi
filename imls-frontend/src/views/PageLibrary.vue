<script>
import { store, state } from "@/store/store.js";
import { format } from "date-fns";

import FetchData from "@/components/FetchData.vue";
import USWDSCard from "@/components/USWDSCard.vue";
import USWDSDatePicker from "@/components/USWDSDatePicker.vue";
import Histogram from '../components/Histogram.vue';
import USWDSTable from '../components/USWDSTable.vue';


export default {
  name: 'Single Library',
  components: {FetchData, USWDSCard, USWDSDatePicker, Histogram, USWDSTable },

  props: {
    id: {
      type: String,
      required: true,
      default: ''
    },
  },
  data() {
    return {
      store,
      state,
    }
  },
  methods: {},
  computed: {
    chartTitle: () => {
      const localDate = state.selectedDate + "T00:00";
      return "Devices present by hour on " + format(new Date(localDate), 'PP')
    },
      getLabels(){
       return store.hourlyLabels;
    },
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
              <FetchData 
              :fscs-id=id 
              :path="store.backendPaths.get24HoursBinnedByHour"
              :queryString="`?_fscs_id=${id}&_day=${state.selectedDate}`">
                <Histogram 
                :dataset="state.fetchedData" 
                :labels="getLabels" 
                :datasetIdKey="id"></Histogram>
                
                <div class="usa-accordion usa-accordion--bordered margin-top-4">
                  <h3 class="usa-accordion__heading">
                    <button
                      type="button"
                      class="usa-accordion__button"
                      aria-expanded="false"
                      aria-controls="viewTable"
                    >
                      View as table
                    </button>
                  </h3>
                  <div id="viewTable" class="usa-accordion__content usa-prose" hidden>
                    <USWDSTable :headers="getLabels" :rows="state.fetchedData" :caption="`Devices present during each hour of the day, starting at 12am on ${state.selectedDate}`" />
                    <div v-if="state.fetchedData.length < 1">Request succeeded but no data was found.</div>
                  </div>
                  <h3 class="usa-accordion__heading">
                    <button
                      type="button"
                      class="usa-accordion__button"
                      aria-expanded="false"
                      aria-controls="viewRaw"
                    >
                      View raw response
                    </button>
                  </h3>
                  <div id="viewRaw" class="usa-accordion__content usa-prose" hidden>
                    <pre>{{ state.fetchedData }}</pre>
                    <div v-if="state.fetchedData.length < 1">Request succeeded but no data was found.</div>
                  </div>
                </div>
              </FetchData>
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
