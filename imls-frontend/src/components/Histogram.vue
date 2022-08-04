<script>
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, BarElement, CategoryScale, LinearScale } from 'chart.js'

const labelsWeekdays = [ "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday" ];
const labelsCalendarMonths = [ "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];


const fakeWeeklyData = [102,58,43,63,52,87,123];
const fakeMonthlyData = [432,558,623,763,987,1023,1153,925,769,656,634,402];



ChartJS.register(Title, Tooltip, BarElement, CategoryScale, LinearScale)

export default {
  name: 'BarChart',
  components: { Bar },
  props: {
    chartId: {
      type: String,
      default: 'bar-chart'
    },
    datasetIdKey: {
      type: String,
      default: 'label'
    },
    width: {
      type: Number,
      default: 400
    },
    height: {
      type: Number,
      default: 100
    },
    cssClasses: {
      default: '',
      type: String
    },
    styles: {
      type: Object,
      default: () => {}
    },
    plugins: {
      type: Object,
      default: () => {}
    },
    bins: {
      default: 'weekdays', // weekdays, months
      type: String
    },  
  },
  data() {
    return {
      chartOptions: {
        responsive: true
      }
    }
  },
  methods: {
    getLabelsByBin(bins) {
      switch (bins) {
        case 'hours':
          return [];
        case 'months':
          return labelsCalendarMonths;
        case 'weekdays':
        default:
          return labelsWeekdays;
      }
    },
    getDatasetsByBin(bins) {
      switch (bins) {
        case 'hours':
          return [];
        case 'months':
          return fakeMonthlyData;
        case 'weekdays':
        default:
          return fakeWeeklyData;
      }
    }
  },
  computed: {
    computedChartData: function() { 
      if (!!this.bins) {
        return {
          labels:  this.getLabelsByBin(this.bins),
          datasets: [ { data: this.getDatasetsByBin(this.bins) } ],
        }
      }
      return {
        labels: [],
        datasets: [ { data: [] }]
      }
    },
  }
}
</script>



<template>
  <Bar
    :chart-options="chartOptions"
    :chart-data="computedChartData"
    :chart-id="chartId"
    :dataset-id-key="datasetIdKey"
    :plugins="plugins"
    :css-classes="cssClasses"
    :styles="styles"
    :width="width"
    :height="height"
  />
</template>
