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
      <tr>
        <td v-bind:key="row" v-for="row in rows">{{ row }}</td>
      </tr>
    </tbody>
  </table>
</template>