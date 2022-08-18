<script>

import { state } from "@/store/store.js";

export default {
  name: 'Fetch Data Wrapper',
  props: {
    fscsId: {
      type: String,
      required: true    
    },
    path: {
      type: String,
      default: '/'
    },
    queryString: {
      type: String,
      default: '?'
    }
  },
  data() {
    return {
      state
    }
  },
  watch: {
    'state.selectedDate'(newVal, oldVal) {
      if (newVal !== oldVal) {
     this.fetchDataFromState();
      }
    }
  },
  async beforeMount() {
     await this.fetchDataFromState();
  },
  methods: {
    async fetchDataFromState() {
      if (this.fscsId.length !== 0) await state.fetchData(this.path, this.queryString);
    }
  },
};
</script>

<template>
  <div class="loading-area">
    <div v-if="state.isLoading" class="loading-indicator">
      <svg class="usa-icon usa-icon--size-9" aria-hidden="true" focusable="false" role="img">
        <use xlink:href="~uswds/img/sprite.svg#autorenew"></use>
      </svg>
    </div>
    <div v-if="state.fetchError && state.fetchError.message">
      <p>Oops! Error encountered: {{ state.fetchError.message }}</p>
    </div>
    <div v-else-if="state.fetchedData.length < 1">
      <p>No data was found that matched your request.</p>
    </div>
    <div v-else-if="state.fetchedData.length > 1">
      <slot></slot>
    </div>
  </div>
</template>

<style scoped>
.loading-area {
  position: relative;
  width: 100%;
  min-height: 20rem;
}
.loading-indicator {
  margin: auto auto;
  text-align: center;
  position: absolute;
  width: 100px;
  height: 100px;
  padding: 14px;
  top: calc(10rem - 50px);
  left: calc(50% - 50px);
  background: #FFFFFFC0;
  border-radius: 100%;
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


</style>