<script>
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, BarElement, CategoryScale, LinearScale } from 'chart.js'
import ChartDataLabels from 'chartjs-plugin-datalabels';

const Y_AXIS_SCALAR = 1.1;

ChartJS.register(Title, Tooltip, BarElement, CategoryScale, LinearScale, ChartDataLabels)
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
        responsive: true,
        scales: {
          y: {
            afterDataLimits: function(axis) {
              axis.max *= Y_AXIS_SCALAR;
            }
          }
        },
        plugins: {
          datalabels: {
            // USWSDS blue-70v
            color: '#0b4778',
            anchor: 'end',
            align: 'end',
            offset: 1,
            padding: 2,
            labels: {
              title: {
                font: {
                  weight: 'bold',
                  size: 20
                }
              },
            }
          },
          tooltip: {
            enabled: false
            // displayColors: false,
            // borderWidth: 0.25,
            // borderColor: '#333',
            // backgroundColor: '#FFF',
            // // titleColor:'#333',
            // // titleAlign: 'center',
            // // titleFont: {
            // //   size: 20
            // // },
            // bodyColor:'#333',
            // bodyAlign: 'center',
            // bodyFont: {
            //   size: 20,
            //   weight: 'bold',
            //   family: 'Source Sans Pro Web, Helvetica Neue, Helvetica, Roboto, Arial, sans-serif'
            // },
            // yAlign: 'bottom',
            // padding: {
            //   left: 10,
            //   right: 10,
            //   top: 6,
            //   bottom:  6
            // },
            // caretSize: 10,
            // callbacks: {
            //   title: () => '',
            //   /* this would mean test chartjs Tooltip internals */
            //   /* c8 ignore start */
            //   label: function(context) {
            //     let label = '';
            //     if (context.parsed.y !== null) {
            //         label += context.parsed.y + " devices present"
            //     }
            //     return label;
            //     /* c8 ignore end */
            //   }
            // }
          }
        }
      }
    }
  },
  computed: {
    computedChartData: function() { 
      return {
        labels:  this.labels,
        datasets: [ 
          { 
            label: this.datasetIdKey,
            // uswds blue-70v
            backgroundColor: '#0b4778',
            data: this.dataset,
            barPercentage: .9,
            categoryPercentage: 1,
          }
        ] ,
      }
    },
  },
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
