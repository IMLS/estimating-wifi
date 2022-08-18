<script>
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, BarElement, CategoryScale, LinearScale } from 'chart.js'


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
    dataset: {
      type: Array,
      default: () => []
    },
    labels: {
      type: Array,
      default: () => []
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

  },
  computed: {
    computedChartData: function() { 
      return {
        labels:  this.labels,
        datasets: [ { data: this.dataset } ] ,
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
