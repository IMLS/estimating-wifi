<script>

export default {
  name: 'Heatmap',
  props: {
    binLabels: {
      type: Array,
      default: () => []
    },
    datasetLabels: {
      type: Array,
      default: () => []
    },
    dataset: {
      type: Array,
      default: () => []
    },
    caption: {
      type: String,
      // todo: decide if this needs a caption or if it should direct screen reader users to the simple table that follows
      default: ''
    }
  },
  methods: {
    sortArrayAscending(arr){
      return arr.sort(function(a,b){ return parseFloat(a) - parseFloat(b);});
    },
    // todo: get Quartile instead?
    getPercentile(thisVal) {
      if (thisVal === 0 ) return 0;
      return (this.allValues.slice().filter((item) => item <= thisVal).length / this.allValues.length);
    },
    // todo: determine saturation separately from alpha channel of background color (to prevent low contrast)
  },
  computed: {
    allValues() {
      return [...this.sortArrayAscending(this.dataset.slice().flat())];
    }

  }
}
</script>


<template>
  <div class="data-grid-container">
    <table class="data-grid">
      <caption class="usa-sr-only">{{ caption }}</caption>
      <thead>
        <tr>
          <th v-if="datasetLabels.length > 0" scope="row"></th>
          <th v-bind:key="label" v-for="label in binLabels" class="data-grid__bin-label border-bottom" scope="col">
            <span class="font-sans-xs">{{ label }}</span>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-bind:key="i" v-for="row, i in dataset">
          <th v-if="datasetLabels.length > 0" class="data-grid__dataset-label padding-y-2 text-right padding-right-2 border-right" scope="row">
            <span class="text-bold text-no-wrap">
              {{ datasetLabels[i] }}
            </span>
          </th>
          <td v-bind:key="i" v-for="cell, i in row" class="data-grid__cell font-mono-md text-center padding-y-2 border" :data-percentile="Math.round(getPercentile(cell)*100)" :style="{ backgroundColor: 'rgba(120,124,206, ' + getPercentile(cell) +')'}">
            {{ cell }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <!-- TODO ADD COLOR LEGEND -->
</template>

<style scoped lang="scss">

.data-grid-container {
  max-width: 100%;
  overflow-x: auto;
}
.data-grid {
  display: grid;
  grid-template-columns: 14ch repeat( v-bind('binLabels.length'), minmax(5ch, auto));
  width: 100%;
  min-width: calc(14ch + 5ch * v-bind('binLabels.length'));
  padding-bottom: 4ch;
}
thead,
tbody, 
tr {
  display: contents;
}

th, td {
  
}
.data-grid__dataset {
}
.data-grid__dataset-label {

}
.data-grid__bin-label {
  height: 6ch;
  overflow: hidden;
  span {
    display: block;
    transform: rotate(-45deg) translateY(1.5ch) translateX(-.75ch);
  }
}
.data-grid__cell {
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
    width: 15ch;
    left: -4.5ch;
    bottom: -3ch;
    z-index: 1;
  }
}
</style>