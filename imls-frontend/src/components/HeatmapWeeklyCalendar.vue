<script>
import { state, store } from "@/store/store.js";
import { addDays, parseISO, isSameDay, format } from "date-fns";

export default {
  name: 'HeatmapWeeklyCalendar',
  props: {
    binLabels: {
      type: Array,
      default: () => []
    },
    dataset: {
      type: Array,
      default: () => []
    },
    weekStartDateISO: {
      type: String,
      default: "05-01-2022"
    },
    colorRGB: {
      type: Array,
      default: () => [0,189,227]
    }
  },  
  data() {
    return {
      state,
      store
    }
  },
  methods: {
    generateDayLabels(startingDateISO) {
      let startingDate = parseISO(startingDateISO + "T00:00");
      // 7 days = 1 week
      let dates = Array.from(Array(7), ( _, i ) => { 
          return addDays(startingDate, i)
      });
      return dates
    },
    sortArrayAscending(arr){
      return arr.sort(function(a,b){ return parseFloat(a) - parseFloat(b);});
    },
    // todo: get Quartile instead?
    getPercentile(thisVal) {
      if (thisVal === this.allValuesSorted[0] ) return 0;
      return (this.allValuesSorted.slice().filter((item) => item <= (thisVal) ).length / this.allValuesSorted.length);
    },
    isSelectedDate(dateLabel) {
      return isSameDay(parseISO(state.selectedDate + "T00:00"), dateLabel)
    },
    formatDateLabel(date, pattern) {
      return format(date, pattern);
    }
  },
  computed: {
    datesForWeek()  {
      return this.generateDayLabels(this.weekStartDateISO)
    },
    allValuesSorted() {
      return [...this.sortArrayAscending(this.dataset.slice().flat())];
    }
  }
}
</script>


<template>
  <div class="scroll-container">
    <div class="usa-sr-only">A more accessible data table follows this infographic.</div>

    <div class="weekly-calendar">

      <div class="weekly-calendar__hour-labels">
        <div class="weekly-calendar__hour-labels__label weekly-calendar__day__label weekly-calendar__cell">
          Local time
        </div>
        <!-- Each hour/bin gets a row: -->
        <div class="weekly-calendar__hour-labels__label weekly-calendar__cell" v-bind:key="header" v-for="header in binLabels">
          {{ header }} 
        </div>
      </div>

      <!-- Each day gets a column: -->
      <div class="weekly-calendar__day" v-bind:key="i" v-for="row, i in dataset" 
          :class="{ 'isSelectedDate': isSelectedDate(datesForWeek[i]) } ">
        <!-- A day column starts with a header-->
        <div class="weekly-calendar__day__label">
          <h3 class="weekly-calendar__day__label--day">{{ store.dayOfWeekLabels[i] }}</h3>
          <h4 class="weekly-calendar__day__label--date">{{ formatDateLabel(datesForWeek[i], 'M/d/yy') }}</h4>
        </div>
        <!-- A day column also has a list of values-->
        <div v-bind:key="i" v-for="cell, i in row" class="weekly-calendar__cell" 
          :data-percentile="Math.round(getPercentile(cell)*100)" :style="{ backgroundColor: 'rgba(' + colorRGB.join() + ', ' + getPercentile(cell) +')'}">
          {{ cell }}
        </div>
      </div>

    </div>

    <div class="legend-container">
      <h3 class="legend-title">
        Percentile Legend
      </h3>
      <div class="legend">
        <div class="legend__step font-mono-md text-center padding-1" v-bind:key="i" v-for="step, i in Array(11)" :style="{ backgroundColor: 'rgba(' + colorRGB.join() + ', ' + i/10 +')'}">
          {{ i*10 }}
        </div>
      </div>
    </div>


  </div>
</template>

<style scoped lang="scss">
.scroll-container {
  max-width: 100%;
  overflow-x: auto;
}
.weekly-calendar {
  display: grid;
  grid-template-columns: 10ch repeat( auto-fit, minmax(3ch, auto) );
  width: 100%;
  text-align: center;
  min-width: 50em;
}

.weekly-calendar__cell {
  padding: 5px 10px;
  height: 50px;
  display: flex;
  flex-flow: column;
  justify-content: center;
}
.weekly-calendar__day {
  border-width: 0 .25px 0 0 ;
  font-size: 20px;
  .weekly-calendar__cell {
    // border doesn't overlap background color so the stroke isn't as dark as it should be
    // use box-shadow because it renders over the background color, not instead of it
    box-shadow: inset -1px -1px 0 0px rgb(0 0 0 / 50%);
    // fake tooltips for now
    position: relative;
    &:hover {
      border-color:#fff;
      cursor: default;
      &:after {
        display: block;
      }
    }
    // todo: real tooltips someday?
    &:after {
      font-size: 14px;
      font-family: "Source Sans Pro", "Helvetica Neue", Helvetica, Arial, sans;
      display: none;
      position: absolute;
      content: 'Percentile: ' attr(data-percentile);
      padding: .5ch 1ch;
      background: #fff;
      border: 1px solid #CCC;
      border-radius: 3px;
      width: 70%;
      left: 15%;
      bottom: -3ch;
      z-index: 10;
    }
  }
}

.weekly-calendar__day__label {
  padding: 5px 10px;
  display: flex;
  flex-flow: column;
  justify-content: center;
  height: 5rem;
  position: relative;
  box-shadow: none;
  border-bottom: .25px solid rgba(0,0,0, .5);
  &:after {
    position: absolute;
    bottom: 0;
    width: 2ch;
    height: 2ch;
    right: 0;
    content: "";
    border-right: .25px solid rgba(0,0,0, .5);
  }
}
.weekly-calendar__day__label--day {
  margin: 0;
  text-transform: uppercase;
  font-weight: normal;
  font-size: 14px;
}
.weekly-calendar__day__label--date {
  margin: 0;
  font-weight: 600;
  font-size: 24px;
}
.weekly-calendar__hour-labels {
  text-align: right;
}
.weekly-calendar__hour-labels__label {
  font-size: 14px;
  position: relative;
  border-right: .25px solid rgba(0,0,0, .5);
  &:after {
    position: absolute;
    bottom: 0;
    width: 2ch;
    height: 2ch;
    right: 0;
    content: "";
    border-bottom: .25px solid rgba(0,0,0, .5);
  }
}
.weekly-calendar__hour-labels__label.weekly-calendar__day__label {
  justify-content: end;
  border-bottom: none;
  border-right: none;
  &:after {
    width: 100%;
    border-right: .25px solid rgba(0,0,0, .5);
  }
}


.isSelectedDate {
  position: relative;
  .weekly-calendar__day__label:after {
    display: none;
  }
  &:after {
    content: '';
    position: absolute;
    top: 0;
    left: -1px;
    width: calc(100% + 1px);
    height: calc(100% + 1px);
    border: 2.25px solid rgba(0,0,0.25);
    box-shadow: 0px 0px 6px 0px rgb(0 0 0 / 25%);
    display: block;
    z-index: 9;
    pointer-events: none;
  }
}
.legend-container {
  @media (min-width: 40em) {
    padding-left: 10ch;
  }
  margin-bottom: 3ch;
  min-width: 100%;
}
.legend {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  box-shadow: inset 1px 1px 0 0px rgb(0 0 0 / 50%);
}
.legend__step {
  flex: 1 1 5ch;
  box-shadow: inset -1px -1px 0 0px rgb(0 0 0 / 50%);
}
</style>