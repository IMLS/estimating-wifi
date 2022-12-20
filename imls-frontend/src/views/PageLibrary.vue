<script>
import { store } from "@/store/store.js";
import { format, formatISO, addDays, subDays, parseISO, startOfWeek, startOfDay, startOfYesterday, endOfWeek, roundToNearestMinutes } from "date-fns";

import FetchData from "@/components/FetchData.vue";
import USWDSCard from "@/components/USWDSCard.vue";
import USWDSDatePicker from "@/components/USWDSDatePicker.vue";
import Histogram from '../components/Histogram.vue';
import Heatmap from '../components/Heatmap.vue';
import HeatmapWeeklyCalendar from '../components/HeatmapWeeklyCalendar.vue';
import USWDSTable from '../components/USWDSTable.vue';
import USWDSBreadcrumb from '../components/USWDSBreadcrumb.vue';

const DEFAULT_START_OVERRIDE = import.meta.env.VITE_DEFAULT_DATE_OVERRIDE;

export default {
  name: 'LibraryPage',
  components: {FetchData, USWDSCard, USWDSDatePicker, Histogram, Heatmap, HeatmapWeeklyCalendar, USWDSTable, USWDSBreadcrumb },

  props: {
    id: {
      type: String,
      required: true,
      default: ''
    },
    selectedDate: {
      type: String,
      // start at specific date if provided (for testing)
      default: () => {
        if (DEFAULT_START_OVERRIDE) {
          return startOfDay(parseISO(DEFAULT_START_OVERRIDE), 'PP').toISOString().split("T")[0]
        }
      return startOfYesterday().toISOString().split("T")[0]
      }
    }, 
  },
  data() {
    return {
      store,
      fetchedLibraryData: null, 
      fetchError: {},
      isLoading: false
    }
  },
  computed: {
    selectedDateUTC() {
      return new Date(this.selectedDate + "T00:00")
    },
    sixDaysAgoUTC(){
      return subDays(this.selectedDateUTC, 6)
    },
    startOfCurrentWeekUTC(){ 
      return startOfWeek(this.selectedDateUTC)
    },
    endOfCurrentWeekUTC(){ 
      return roundToNearestMinutes(endOfWeek(this.selectedDateUTC))
    },
    dailyChartTitle() {
      return "Devices present by hour on " + format(this.selectedDateUTC, 'PP')
    },
    weeklyChartTitle() {
      return "Devices present by hour for 7 consecutive days, " + format(this.sixDaysAgoUTC, 'PP') + " — " + format(this.selectedDateUTC, 'PP')
    },
    weeklyCalendarChartTitle() {
      return "Devices present by hour for the calendar week containing  " + format(this.selectedDateUTC, 'PP') + ", " + format(this.startOfCurrentWeekUTC, 'PP') + " — " + format(this.endOfCurrentWeekUTC, 'PP');
    },

    libraryName() {
      if (this.fetchedLibraryData && this.fetchedLibraryData.libname )  return this.fetchedLibraryData.libname;
      return "Library " + this.id
    },
    breadcrumbs () {
      if ( this.fetchedLibraryData == null ) return []
      return [
         { 
          name: "All States",
          link: "/" 
        },
        { 
          name: this.store.states[this.fetchedLibraryData.stabr],
          link: `/state/${this.fetchedLibraryData.stabr}/` 
        },
        {
          name: this.libraryName
        }
      ]
    }
  },
  watch: {
    id(newVal, oldVal) {
      if (newVal !== oldVal) {
        this.fetchLibraryData();
      }
    }
  },
  async beforeMount() {
    await this.fetchLibraryData();
  },
  methods: {
    toISODate(utcDate) {
      return formatISO(utcDate, { representation: 'date' })
    },
    generateDayLabels(startingDateUTC, count) {

      let dates = Array.from(Array(count), ( _, i ) => { 
          return format(addDays(startingDateUTC, i), 'M/d/yy' )
      });
      return dates
    },
    navigateToSelectedDate(date) {
      this.$router.push({
        query: { ...this.$router.query, date: encodeURIComponent(date) }
      })
    },
    async fetchLibraryData() {
      if (this.id.length !== 0) {
        this.isLoading = true;
        try {
          const response = await fetch(`${store.backendBaseUrl}${store.backendPaths.getLibraryDetailsById}?_fscs_id=${this.id}`);
          this.fetchedLibraryData = await response.json();
        } catch (error) {
          this.fetchError = error;
        }
        this.isLoading = false;
      }
    },
    leftPadSequence(seq) {
      return (parseInt(seq) + 1000).toString().substring(1);
    },
    formatFSCSandSequence(fscsid, seq) {
      return fscsid + '-' + this.leftPadSequence(seq)
    }
  }
};
</script>

<template>
  <div>
    <USWDSBreadcrumb :crumbs=breadcrumbs />
    <h1>{{ libraryName }}</h1>
    <div v-if="fetchedLibraryData !== null">
      <h2>{{ formatFSCSandSequence(fetchedLibraryData.fscskey, fetchedLibraryData.fscs_seq) }}</h2>
      {{ fetchedLibraryData.address }}<br>
      {{ fetchedLibraryData.city }},   {{ fetchedLibraryData.stabr }}   {{ fetchedLibraryData.zip }}
    </div>

    <USWDSDatePicker :initial-date=toISODate(selectedDateUTC) @date_changed="navigateToSelectedDate" />

    <div class="usa-card-group margin-top-6">

      <!-- first graph: Binned devices by hour for one day -->
      <div class="usa-card tablet:grid-col-12">
        <USWDSCard :title="dailyChartTitle">
          <div class="grid-row">
            <div class="grid-col">
              <FetchData 
              v-slot="slotProps"
              :fscs-id=id
              :path="store.backendPaths.get24HoursBinnedByHour"
              :selected-date="selectedDate">
                <Histogram 
                :dataset="slotProps.fetchedData" 
                :labels="store.hourlyLabels" 
                :dataset-id-key="id"></Histogram>
                
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
                    <USWDSTable :column-headers="store.hourlyLabels" :rows="[slotProps.fetchedData]" :caption="`Devices present during each hour of the day, starting at 12am on ${selectedDate}`" />
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
              :fscs-id=id
              :path="store.backendPaths.get24HoursBinnedByHourForNDays"
              :selected-date="toISODate(sixDaysAgoUTC)"
              :query-params="{ _direction: true,  _days : 7 }">
                <Heatmap 
                :dataset="slotProps.fetchedData" 
                :bin-labels="store.hourlyLabels"
                :dataset-labels="generateDayLabels(sixDaysAgoUTC, 7)"></Heatmap>
                
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
                    <USWDSTable :column-headers="store.hourlyLabels"  :row-headers="generateDayLabels(sixDaysAgoUTC, 7)" :rows="slotProps.fetchedData" :caption="`Devices present during each hour of the day, starting at 12am on ${toISODate(sixDaysAgoUTC)}, for one week`" />
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
              :fscs-id=id
              :path="store.backendPaths.get24HoursBinnedByHourForNDays"
              :selected-date="toISODate(startOfCurrentWeekUTC)"
              :query-params="{ _direction: true,  _days : 7 }">
                <HeatmapWeeklyCalendar 
                :dataset="slotProps.fetchedData" 
                :bin-labels="store.hourlyLabels"
                :week-start-date-i-s-o="toISODate(startOfCurrentWeekUTC)"
                :selected-date="selectedDate"
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
                    <USWDSTable :column-headers="store.hourlyLabels"  :row-headers="generateDayLabels(startOfCurrentWeekUTC, 7)" :rows="slotProps.fetchedData" :caption="`Devices present during each hour of the day, starting at 12am on ${toISODate(startOfCurrentWeekUTC)}, for one week`" />
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
