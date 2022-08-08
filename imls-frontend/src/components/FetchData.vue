<script>

const baseUrl = 'http://127.0.0.1:3000/presences'

export default {
  name: 'Fetch Data',
  props: {
    fscsId: {
      type: String,
      required: true    
    },
  },
  data() {
    return {
      startDate: '2022-05-11T00:00:00',
      endDate: '',
      totalFound: 0,
      loadedData: {},
      loadedError: {}
    }
  },
  computed: {
    loadUrl() {
      return `${baseUrl}?limit=100&fscs_id=eq.${this.fscsId}&start_time=lte.${this.startDate}`;
    }
  },
  watch: {
    fscsId(newVal, oldVal) {
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
            Prefer: 'count=estimated'
          }
        })
        this.loadedData = (await response.json())
        this.totalFound = (await response.headers.get('Content-Range')).split('/')[1];
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
        Load the first hundred entries for <b>{{ fscsId }}</b> since {{ startDate }} from the backend:
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
      <h3>Display {{ loadedData.length }} of {{ totalFound }} total entries found:</h3>
      <pre>{{ loadedData }}</pre>
      <div v-if="loadedData.length < 1">Request succeeded but no data was found.</div>
    </div>
    <div v-else>Loading...</div>
</div>
</template>