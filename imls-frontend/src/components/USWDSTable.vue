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
        <th v-bind:key="header" v-for="header in columnHeaders" scope="col">{{ header }}</th>
      </tr>
    </thead>
    <tbody>
      <tr v-bind:key="i" v-for="row, i in rows">
        <th v-if="rowHeaders.length > 0" scope="row">
          <span class="text-bold text-no-wrap">
            {{ rowHeaders[i] }}
          </span>
        </th>
        <td v-bind:key="i" v-for="cell, i in row" class="font-mono-sm text-tabular text-right">{{ cell }}</td>
      </tr>
    </tbody>
  </table>
</template>