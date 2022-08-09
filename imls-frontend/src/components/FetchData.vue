<script>
import formatISO from 'date-fns/formatISO'
import parseISO from 'date-fns/parseISO'
import format from 'date-fns/format'
import endOfDay from 'date-fns/endOfDay'


const baseUrl = 'http://127.0.0.1:3000/presences'

export default {
  name: 'Fetch Data',
  props: {
    fscsId: {
      type: String,
      required: true    
    },
    startDate: {
      type: String,
      required: false,
      default: () => '2022-05-01'
    }
  },
  data() {
    return {
      totalFound: 0,
      loadedData: {},
      loadedError: {}
    }
  },
  computed: {
    loadUrl() {
      return `${baseUrl}?limit=1000&fscs_id=eq.${this.fscsId}&start_time=gte.${formatISO(this.localStartDate)}&end_time=lt.${formatISO(endOfDay(this.localStartDate))}&order=start_time`;
    },
    localStartDate() {
      return parseISO(this.startDate, 'yyyy-MM-dd', new Date())
    }
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
    formatHumanReadableDateFromISO(dateString) {
      return format(parseISO(dateString), "bbb 'on' PPPP, zzzz");

    },
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
    }
  },
};
</script>

<template>
<div>
    <div>
      <h3>
        Load the first 1k entries from <b>{{ fscsId }}</b> for a full day starting {{ formatHumanReadableDateFromISO(startDate) }}:
      </h3>
    </div>
    <div class="margin-y-2">
      <code>{{ loadUrl }}</code>
    </div>
    <!-- <ul class="usa-button-group usa-button-group--segmented margin-bottom-3">
      <li class="usa-button-group__item" v-bind:key="fscs.id" v-for="fscs in store.fscs_ids">
        <button class="usa-button" :class="{'usa-button--outline': fscsId != fscs.id }" @click="fscsId = fscs.id">{{ fscs.id }}</button>
      </li>
    </ul> -->
    <div v-if="loadedError && loadedError.message">
      <p>Oops! Error encountered: {{ loadedError.message }}</p>
      <button @click="retry">Retry</button>
    </div>
    <div v-else-if="loadedData">
      <h3>Display {{ loadedData.length.toLocaleString('en-US') }} of {{ totalFound.toLocaleString('en-US') }} total entries found:</h3>
      <pre>{{ loadedData }}</pre>
      <div v-if="loadedData.length < 1">Request succeeded but no data was found.</div>
    </div>
    <div v-else>Loading...</div>
</div>
</template>