<script>
import { store, state } from "@/store/store.js";
import { format, addDays, parseISO } from "date-fns";

import FetchData from "@/components/FetchData.vue";
import USWDSCard from "@/components/USWDSCard.vue";
import USWDSDatePicker from "@/components/USWDSDatePicker.vue";
import Histogram from '../components/Histogram.vue';
import Heatmap from '../components/Heatmap.vue';
import USWDSTable from '../components/USWDSTable.vue';


export default {
  name: 'Single Library',
  components: {FetchData, USWDSCard, USWDSDatePicker, Histogram, Heatmap, USWDSTable },

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
  methods: {
    generateDayLabels(startingDateISO, count) {
      let startingDate = parseISO(startingDateISO + "T00:00");
      let dates = Array.from(Array(count), ( _, i ) => { 
          return format(addDays(startingDate, i), 'EEE, MMM d' )
      });
      // let labels = dates.map(each => { format(each, 'PP' )})
      // console.log(labels)
      return dates
    }
  },
  computed: {
    dailyChartTitle: () => {
      const localDate = state.selectedDate + "T00:00";
      return "Devices present by hour on " + format(new Date(localDate), 'PP')
    },
    weeklyChartTitle: () => {
      const localDate = state.selectedDate + "T00:00";
      return "Devices present by hour for a week, starting on " + format(new Date(localDate), 'PP')
    },
  }
};
</script>

<template>
  <div>
    <h1>Library {{ id }}</h1>

    <USWDSDatePicker :initialDate=state.selectedDate />

    <div class="usa-card-group margin-top-6">

      <!-- first graph: Binned devices by hour for one day -->
      <div class="usa-card tablet:grid-col-12">
        <USWDSCard :title="dailyChartTitle">
          <div class="grid-row">
            <div class="grid-col">
              <FetchData 
              v-slot="slotProps"
              :fscsId=id
              :path="store.backendPaths.get24HoursBinnedByHour"
              :queryParams="{ _day: state.selectedDate }">
                <Histogram 
                :dataset="slotProps.fetchedData" 
                :labels="store.hourlyLabels" 
                :datasetIdKey="id"></Histogram>
                
                <div class="usa-accordion usa-accordion--bordered margin-top-4">
                  <h3 class="usa-accordion__heading">
                    <button
                      type="button"
                      class="usa-accordion__button"
                      aria-expanded="false"
                      aria-controls="viewTableDaily"
                    >
                      View as table
                    </button>
                  </h3>
                  <div id="viewTableDaily" class="usa-accordion__content usa-prose" hidden>
                    <USWDSTable :columnHeaders="store.hourlyLabels" :rows="[slotProps.fetchedData]" :caption="`Devices present during each hour of the day, starting at 12am on ${state.selectedDate}`" />
                    <div v-if="slotProps.fetchedData.length < 1">Request succeeded but no data was found.</div>
                  </div>
                  <!-- <h3 class="usa-accordion__heading">
                    <button
                      type="button"
                      class="usa-accordion__button"
                      aria-expanded="false"
                      aria-controls="viewRawDaily"
                    >
                      View raw response
                    </button>
                  </h3>
                  <div id="viewRawDaily" class="usa-accordion__content usa-prose" hidden>
                    <pre>{{ slotProps.fetchedData }}</pre>
                    <div v-if="slotProps.fetchedData.length < 1">Request succeeded but no data was found.</div>
                  </div> -->
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

      <!-- second graph: Binned devices by hour for one week -->
      <div class="usa-card tablet:grid-col-12">
        <USWDSCard :title="weeklyChartTitle">
          <div class="grid-row">
            <div class="grid-col">
              <FetchData 
              v-slot="slotProps"
              :fscsId=id
              :path="store.backendPaths.get24HoursBinnedByHourForNDays"
              :queryParams="{ _start: state.selectedDate, _direction: true,  _days : 7 }">
                <Heatmap 
                :dataset="slotProps.fetchedData" 
                :binLabels="store.hourlyLabels"
                :datasetLabels="generateDayLabels(state.selectedDate, 7)"></Heatmap>
                
                <div class="usa-accordion usa-accordion--bordered margin-top-4">
                  <h3 class="usa-accordion__heading">
                    <button
                      type="button"
                      class="usa-accordion__button"
                      aria-expanded="false"
                      aria-controls="viewTableWeekly"
                    >
                      View as table
                    </button>
                  </h3>
                  <div id="viewTableWeekly" class="usa-accordion__content usa-prose" hidden>
                    <USWDSTable :columnHeaders="store.hourlyLabels"  :rowHeaders="generateDayLabels(state.selectedDate, 7)" :rows="slotProps.fetchedData" :caption="`Devices present during each hour of the day, starting at 12am on ${state.selectedDate}, for one week`" />
                    <div v-if="slotProps.fetchedData.length < 1">Request succeeded but no data was found.</div>
                  </div>
                  <h3 class="usa-accordion__heading">
                    <button
                      type="button"
                      class="usa-accordion__button"
                      aria-expanded="false"
                      aria-controls="viewRawWeekly"
                    >
                      View raw response
                    </button>
                  </h3>
                  <div id="viewRawWeekly" class="usa-accordion__content usa-prose" hidden>
                    <pre>{{ slotProps.fetchedData }}</pre>
                    <div v-if="slotProps.fetchedData.length < 1">Request succeeded but no data was found.</div>
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
