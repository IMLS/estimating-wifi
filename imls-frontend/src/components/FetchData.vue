<script>

import { state } from "@/store/store.js";

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
    },
    fscsId(newVal, oldVal) {
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
      if (this.fscsId.length !== 0) await state.fetchData(this.path, `?_fscs_id=${this.fscsId}&_day=${state.selectedDate}`);
    }
  },
  computed: {
    responseIsOKButEmpty() {
      return state.fetchedData.reduce((previous, current) => previous + current, 0)
    }
  }
};
</script>

<template>
  <div class="loading-area">
    <div v-if="state.isLoading" class="loading-indicator">
      <svg class="usa-icon usa-icon--size-9" aria-hidden="true" focusable="false" role="img">
        <use xlink:href="~uswds/img/sprite.svg#autorenew"></use>
      </svg>
    </div>
    <div class="loaded--error" v-if="state.fetchError && state.fetchError.message">
      <p>Oops! Error encountered: {{ state.fetchError.message }}</p>
    </div>
    <div class="loaded--no-data" v-if="(state.fetchedData.length > 1 && responseIsOKButEmpty === 0) || state.fetchedData.length < 1">
      <p>No data was found that matched your request for devices present near <b>{{ fscsId }}</b> on <b>{{ state.selectedDate }}</b>. Please choose a different date or library.</p>
    </div>
    <div class="loaded--has-data" v-else-if="state.fetchedData.length > 1">
      <slot></slot>
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
  @media (min-width: 64em ) {
    top: calc(10rem - 50px);
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