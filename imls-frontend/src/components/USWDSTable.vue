<script>

export default {
  name: 'USWDSTable',
  props: {
    headers: {
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
  <table class="usa-table usa-table--bordered">
    <caption>{{ caption }}</caption>
    <thead>
      <tr>
        <th v-bind:key="header" v-for="header in headers" scope="col">{{ header }}</th>
      </tr>
    </thead>
    <tbody>
      <tr v-bind:key="i" v-for="row, i in rows">
        <td v-bind:key="i" v-for="cell, i in row" >{{ cell }}</td>
      </tr>
    </tbody>
  </table>
</template>