<script>
import Histogram from '../components/Histogram.vue';
import { store } from "@/store/store.js";

//todo: make the backend function name configurable
const baseUrl = 'http://127.0.0.1:3000/rpc/bin_devices_per_hour'

export default {
  name: 'Fetch Data',
  components: {Histogram },
  props: {
    fscsId: {
      type: String,
      required: true    
    },
    startDate: {
      type: String,
      required: false,
      default: () => '2022-05-10'
    }
  },
  data() {
    return {
      totalFound: 0,
      loadedData: {},
      loadedError: {},
      store
    }
  },
  computed: {
    loadUrl() {
      return `${baseUrl}?_fscs_id=${this.fscsId}&_day=${this.startDate}`;
    },
    getLabels(){
       return store.hourlyLabels;
    },
  },
  watch: {
    fscsId(newVal, oldVal) {
      if (newVal !== oldVal) {
        this.fetchData()
      }
    },
    startDate(newVal, oldVal) {
      if (newVal !== oldVal) {
        this.fetchData()
      }
    }
  },
  beforeMount() {
    this.fetchData();
  },
  methods: {

    async fetchData() {
      try {
        const response = await fetch(this.loadUrl, {
          headers: {
            // https://postgrest.org/en/stable/api.html#estimated-count
            Prefer: 'count=exact'
          }
        })
        this.loadedData = (await response.json())
        this.totalFound = parseInt((response.headers.get('Content-Range')).split('/')[1]);
      } catch (error) {
        this.loadedError = error
      }
    },
    formatCount(num) {
      return parseInt(num).toLocaleString('en-US')
    }
  },
};
</script>

<template>
<div>
    <div v-if="loadedError && loadedError.message">
      <p>Oops! Error encountered: {{ loadedError.message }}</p>
      <button @click="retry">Retry</button>
    </div>
    <div v-else-if="loadedData">
      <Histogram :dataset="loadedData" :labels="getLabels" ></Histogram>
      <h3>Raw output from  <code>{{ loadUrl }}</code>:</h3>
      <pre>{{ loadedData }}</pre>
      <div v-if="loadedData.length < 1">Request succeeded but no data was found.</div>
    </div>
    <div v-else>Loading...</div>

</div>
</template>