<script>

import { store } from "@/store/store.js";
import { nextTick } from 'vue'

export default {
  name: 'Fetch Data Wrapper',
  props: {
    fscsId: {
      type: String,
      default: ""
    },
    path: {
      type: String,
      default: '/'
    }, 
    queryParams: {
      type: Object,
      default: () => {}
    },
    selectedDate: {
      type: String,
      default: ''
    }, 
  },
  data() {
    return {
      store,
      fetchCount: null,
      fetchError: {},
      fetchedData: {},
      isLoading: false,
    }
  },
  watch: {
    "selectedDate": {
      async handler(newVal, oldVal)  {
        if (newVal !== oldVal) {
          await nextTick();
          this.fetchData();
        }
      },
      deep: true, 
    },
    fscsId(newVal, oldVal) {
      if (newVal !== oldVal) {
        this.fetchData();
      }
    }
  },
  async beforeMount() {
     await this.fetchData();
  },
  methods: {
    async fetchData() {
      if (this.fscsId.length !== 0) {
        this.isLoading = true;
        try {
          const response = await fetch(`${store.backendBaseUrl}${this.path}?_fscs_id=${this.fscsId}&_start=${this.selectedDate}${this.queryString}`);
          if (await !response.ok) {
            throw new Error(response.status);
          }
          this.fetchedData = await response.json();
        } catch (error) {
          this.fetchError = error;
        }
        this.isLoading = false;
      }
    },
    reduceArray(arr) {
      if (Array.isArray(arr) ) {
        const reduced = arr.reduce((previous, current) => parseInt(previous) + parseInt(current), 0)
        return this.reduceArray(reduced)
      } else {
        return arr
      }
    }
  },
  computed: {
    queryString() {
      if (this.queryParams && Object.keys(this.queryParams).length !== 0) {
        return '&' + Object.keys(this.queryParams).map(key => key + '=' + this.queryParams[key]).join('&');
      }
      return ''
    },
    responseIsOKButEmpty() {
      return (this.reduceArray(this.fetchedData) === 0)
    }
  }
};
</script>

<template>
  <div class="loading-area">
    <div v-if="this.isLoading" class="loading-indicator">
      <svg class="usa-icon usa-icon--size-9" aria-hidden="true" focusable="false" role="img">
        <use xlink:href="~uswds/img/sprite.svg#autorenew"></use>
      </svg>
    </div>
    <div class="loaded--error" v-if="this.fetchError && this.fetchError.message">
      <p>Oops! Error encountered: {{ this.fetchError.message }}</p>
    </div>
    <div class="loaded--no-data" v-if="!this.fetchedData || (this.fetchedData.length > 1 && responseIsOKButEmpty) ||  (this.fetchedData && this.fetchedData.length < 1)">
      <p>No data was found that matched your request for devices present near <b>{{ fscsId }}</b> on <b>{{ this.selectedDate }}</b>. Please choose a different date or library.</p>
    </div>
    <div class="loaded--has-data" v-else-if="this.fetchedData.length > 1">
      <slot :fetchedData="this.fetchedData" ></slot>
    </div>
  </div>
</template>

<style scoped lang="scss">
.loading-area {
  position: relative;
  width: 100%;
  min-height: 10rem;
  @media (min-width: 64em ) {
    min-height: 20rem;
  }
}
.loading-indicator {
  text-align: center;
  width: 100px;
  height: 100px;
  padding: 14px;
  background: #FFFFFFC0;
  border-radius: 100%;
  position: absolute;
  left: calc(50% - 50px);
  top: calc(5rem - 50px);
  z-index: 100;
  @media (min-width: 64em ) {
    top: calc(40% - 50px);
  }
}
.loading-indicator .usa-icon {
  animation-duration: 1s;
  animation-name: spin;
  animation-iteration-count: infinite;
  transition-timing-function: cubic-bezier(0.39, 0.58, 0.57, 1);
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(359deg);
  }
}
.loaded--no-data {
  text-align: center;
  padding: 20px;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  flex-flow: column;
  justify-content: center;
  min-height: 10rem;
  @media (min-width: 64em ) {
    min-height: 20rem;
  }

}

</style>