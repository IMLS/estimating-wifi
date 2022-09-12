<script>
import { store } from "@/store/store.js";
import { format, formatISO, addDays, parseISO, startOfWeek, startOfYesterday } from "date-fns";

import FetchData from "@/components/FetchData.vue";
import USWDSCard from "@/components/USWDSCard.vue";
import USWDSDatePicker from "@/components/USWDSDatePicker.vue";
import Histogram from '../components/Histogram.vue';
import Heatmap from '../components/Heatmap.vue';
import HeatmapWeeklyCalendar from '../components/HeatmapWeeklyCalendar.vue';
import USWDSTable from '../components/USWDSTable.vue';


export default {
  name: 'Single Library',
  components: {FetchData, USWDSCard, USWDSDatePicker, Histogram, Heatmap, HeatmapWeeklyCalendar, USWDSTable },

  props: {
    id: {
      type: String,
      required: true,
      default: ''
    },
    selectedDate: {
      type: String,
      // load yesterday by default
      default: () => startOfYesterday().toISOString().split("T")[0]
    }, 
  },
  data() {
    return {
      store
    }
  },
  methods: {
    generateDayLabels(startingDateISO, count) {
      let startingDate = parseISO(startingDateISO + "T00:00");
      let dates = Array.from(Array(count), ( _, i ) => { 
          return format(addDays(startingDate, i), 'M/d/yy' )
      });
      return dates
    },
    navigateToSelectedDate(date) {
      this.$router.push({
        query: { ...this.$router.query, date: encodeURIComponent(date) }
      })
    }
  },
  computed: {
    activeDate() {
      return this.selectedDate;
    },
    dailyChartTitle() {
      const localDate = this.selectedDate + "T00:00";
      return "Devices present by hour on " + format(new Date(localDate), 'PP')
    },
    weeklyChartTitle() {
      const localDate = this.selectedDate + "T00:00";
      return "Devices present by hour for 7 consecutive days, starting on " + format(new Date(localDate), 'PP')
    },
    weeklyCalendarChartTitle() {
      return "Devices present by hour for a week, starting on " + format(startOfWeek(parseISO(this.selectedDate)), 'PP')
    },
    startOfWeekInISO() {
      return formatISO(startOfWeek(parseISO(this.selectedDate)), { representation: 'date' })
    }
  }
};
</script>

<template>
  <div>
    <h1>Library {{ id }}</h1>

    <USWDSDatePicker :initialDate=activeDate @date_changed="navigateToSelectedDate" />

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
              :selectedDate="selectedDate">
                <Histogram 
                :dataset="slotProps.fetchedData" 
                :labels="store.hourlyLabels" 
                :datasetIdKey="id"></Histogram>
                
                <div class="usa-accordion usa-accordion--bordered margin-top-4">
                  <h3 class="usa-accordion__heading">
                    <button
                      type="button"
                      class="usa-accordion__button bg-primary-lighter"
                      aria-expanded="false"
                      aria-controls="viewTableDaily"
                    >
                      View as table
                    </button>
                  </h3>
                  <div id="viewTableDaily" class="usa-accordion__content usa-prose" hidden>
                    <USWDSTable :columnHeaders="store.hourlyLabels" :rows="[slotProps.fetchedData]" :caption="`Devices present during each hour of the day, starting at 12am on ${this.selectedDate}`" />
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

      <!-- second graph: Binned devices by hour for one week -->
      <div class="usa-card tablet:grid-col-12">
        <USWDSCard :title="weeklyChartTitle">
          <div class="grid-row">
            <div class="grid-col">
              <FetchData 
              v-slot="slotProps"
              :fscsId=id
              :path="store.backendPaths.get24HoursBinnedByHourForNDays"
              :selectedDate="selectedDate"
              :queryParams="{ _direction: true,  _days : 7 }">
                <Heatmap 
                :dataset="slotProps.fetchedData" 
                :binLabels="store.hourlyLabels"
                :datasetLabels="generateDayLabels(this.selectedDate, 7)"></Heatmap>
                
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
                    <USWDSTable :columnHeaders="store.hourlyLabels"  :rowHeaders="generateDayLabels(this.selectedDate, 7)" :rows="slotProps.fetchedData" :caption="`Devices present during each hour of the day, starting at 12am on ${this.selectedDate}, for one week`" />
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

      <!-- third graph: Binned devices by hour for one week, calendar view -->
      <div class="usa-card tablet:grid-col-12">
        <USWDSCard :title="weeklyCalendarChartTitle">
          <div class="grid-row">
            <div class="grid-col">
              <FetchData 
              v-slot="slotProps"
              :fscsId=id
              :path="store.backendPaths.get24HoursBinnedByHourForNDays"
              :selectedDate="startOfWeekInISO"
              :queryParams="{ _direction: true,  _days : 7 }">
                <HeatmapWeeklyCalendar 
                :dataset="slotProps.fetchedData" 
                :binLabels="store.hourlyLabels"
                :weekStartDateISO="startOfWeekInISO"
                :selectedDate="selectedDate"
                ></HeatmapWeeklyCalendar>
                
                <div class="usa-accordion usa-accordion--bordered margin-top-4">
                  <h3 class="usa-accordion__heading">
                    <button
                      type="button"
                      class="usa-accordion__button"
                      aria-expanded="false"
                      aria-controls="viewTableWeeklyCalendar"
                    >
                      View as table
                    </button>
                  </h3>
                  <div id="viewTableWeeklyCalendar" class="usa-accordion__content usa-prose" hidden>
                    <USWDSTable :columnHeaders="store.hourlyLabels"  :rowHeaders="generateDayLabels(startOfWeekInISO, 7)" :rows="slotProps.fetchedData" :caption="`Devices present during each hour of the day, starting at 12am on ${startOfWeekInISO}, for one week`" />
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
