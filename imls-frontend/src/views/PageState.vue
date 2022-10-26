<script>
import { store } from "@/store/store.js";

import USWDSBreadcrumb from '../components/USWDSBreadcrumb.vue';

export default {
  name: 'All Library Systems for a State',
  components: { USWDSBreadcrumb },

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
      fetchError: {},
      fetchedData: {},
      isLoading: false,
    }
  },
  watch: {
    stateInitials(newVal, oldVal) {
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
      if (this.stateInitials.length !== 0) {
        this.isLoading = true;
        try {
          const response = await fetch(`${store.backendBaseUrl}${store.backendPaths.getAllSystemsByStateInitials}?_state_code=${this.stateInitials}`);
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
    leftPadSequence(seq) {
      return (parseInt(seq) + 1000).toString().substring(1);
    },
    formatFSCSandSequence(fscsid, seq) {
      return fscsid + '-' + this.leftPadSequence(seq)
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
  beforeRouteEnter (to, from, next) {
    next(vm => {
      // access to component public instance via `vm`
      // if no matching state abbreviation is found, redirect to 404 page
      if (!vm.stateName) {
        console.log("No state found")
        vm.$router.push({name: 'NotFound'});
      }
    })
  }
};
</script>

<template>
  <div>
    <USWDSBreadcrumb :crumbs=breadcrumbs />
    <h1>{{ stateName }} Public Libraries</h1>
    <div class="loading-area">
      <div v-if="this.isLoading" class="loading-indicator">
        <svg class="usa-icon usa-icon--size-9" aria-hidden="true" focusable="false" role="img">
          <use xlink:href="~uswds/img/sprite.svg#autorenew"></use>
        </svg>
      </div>
      <div class="loaded--error" v-if="this.fetchError && this.fetchError.message">
        <p>Oops! Error encountered: {{ this.fetchError.message }}</p>
      </div> 
      <div class="loaded--has-data" v-else-if="this.fetchedData.length > 1">

        <ol class="usa-list">
          <template v-for="system in this.fetchedData">
            <li>
            <RouterLink class="usa-link" :to="{ path: '/library/' + formatFSCSandSequence(system.fscskey, system.fscs_seq) + '/' , query: $route.query}">
              {{ formatFSCSandSequence(system.fscskey, system.fscs_seq) }} - {{ system.libname }}
            </RouterLink>
            </li>
          </template>
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