<script>
import { store } from "@/store/store.js";

import USWDSBreadcrumb from '../components/USWDSBreadcrumb.vue';

export default {
  name: 'StatePage',
  components: { USWDSBreadcrumb },
  beforeRouteEnter (to, from, next) {
    next(vm => {
      if (vm.stateName === undefined) {
        vm.$router.push({name: 'NotFound'});
      }
    })
  },

  props: {
    stateInitials: {
      type: String,
      required: true,
      default: ''
    },
  },
  data() {
    return {
      store,
      fetchCount: null,
      fetchError: null,
      fetchedLibraries: [],
      isLoading: false,
    }
  },
  computed: {
    stateName () {
      return this.store.states[this.stateInitials]
    },
    breadcrumbs () {
      return [
         { 
          name: "All States",
          link: "/" 
        },
        { 
          name: this.stateName,
          link: `/system/${this.stateInitials}` 
        }
      ]
    }
  },
  watch: {
    stateInitials(newVal, oldVal) {
      if (newVal !== oldVal) {
        this.fetchLibraries();
      }
    }
  },
  async beforeMount() {
    await this.fetchLibraries();
  },
  methods: {
    async fetchLibraries() {
      if (this.stateInitials.length !== 0) {
        this.isLoading = true;
        try {
          const response = await fetch(`${store.backendBaseUrl}${store.backendPaths.getAllSystemsByStateInitials}?_state_code=${this.stateInitials}`);
          const parsedResponse = await response.json();

          // the API currently returns null instead of an empty array on no matches
          this.fetchedLibraries = await !!parsedResponse ? parsedResponse : []
        } catch (error) {
          this.fetchError = error.message;
        }
        this.isLoading = false;
      }
    },
    leftPadSequence(seq) {
      return (parseInt(seq) + 1000).toString().substring(1);
    },
    formatFSCSandSequence(fscsid, seq) {
      return fscsid + '-' + this.leftPadSequence(seq)
    },
  },
  metaInfo(stateName = this.stateName) {
    const pagePrefix = `${stateName} Public Libraries`;
    return {
      title: pagePrefix
    }
  }
};
</script>

<template>
  <div>
    <USWDSBreadcrumb :crumbs=breadcrumbs />
    <h1>{{ stateName }} Public Libraries</h1>
    <div class="loading-area">
      <div v-if="isLoading" class="loading-indicator">
        <svg class="usa-icon usa-icon--size-9" aria-hidden="true" focusable="false" role="img">
          <use xlink:href="~uswds/img/sprite.svg#autorenew"></use>
        </svg>
      </div>
      <div v-if="fetchedLibraries == null || fetchedLibraries.length < 1" class="loaded--no-data">
        <p>Sorry, no matching libraries found. </p>
        <span v-if="fetchError">Oops! Error encountered: {{ fetchError }} </span>
      </div>
      <div v-else class="loaded--has-data">

        <ol class="usa-list">
          <li v-for="system in fetchedLibraries"  :key=system>
            <RouterLink class="usa-link" :to="{ path: '/library/' + formatFSCSandSequence(system.fscskey, system.fscs_seq) + '/' }">
              {{ formatFSCSandSequence(system.fscskey, system.fscs_seq) }} - {{ system.libname }}
            </RouterLink>
         
          </li>
        </ol>
      

      </div>
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