<script>

export default {
  name: 'USWDSTable',
  props: {
    columnHeaders: {
      type: Array,
      default: () => []
    },
    rowHeaders: {
      type: Array,
      default: () => []
    },
    rows: {
      type: Array,
      default: () => []
    },
    caption: {
      type: String,
      default: ''
    }
  }
}
</script>


<template>
  <table class="usa-table usa-table--bordered usa-table--compact usa-table--striped">
    <caption>{{ caption }}</caption>
    <thead>
      <tr>
        <th v-if="rowHeaders.length > 0" scope="row"></th>
        <th v-for="header in columnHeaders" :key="header" scope="col">{{ header }}</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="row, headerIndex in rows" :key="headerIndex">
        <th v-if="rowHeaders.length > 0" scope="row">
          <span class="text-bold text-no-wrap">
            {{ rowHeaders[headerIndex] }}
          </span>
        </th>
        <td v-for="cell, index in row" :key="index" class="font-mono-sm text-tabular text-right">{{ cell }}</td>
      </tr>
    </tbody>
  </table>
</template>