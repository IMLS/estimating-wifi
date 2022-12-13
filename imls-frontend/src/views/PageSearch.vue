<script>
import { store } from "@/store/store.js";

export default {
  props: {
    query: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      store,
      fetchedLibraries: null, 
      fetchError: {},
      isLoading: false
    }
  },
  watch: {
    query(newVal, oldVal) {
      if (newVal !== oldVal) {
        this.searchLibraryNames();
      }
    }
  },
  async beforeMount() {
    await this.searchLibraryNames();
  },
  methods: {
    async searchLibraryNames() {
      if (this.query.length !== 0) {
        this.isLoading = true;
        try {
          const response = await fetch(`${store.backendBaseUrl}${store.backendPaths.textSearchLibraryNames}?_name=${this.query}`);
          if (await !response.ok) {
            throw new Error(response.status);
          }
          this.fetchedLibraries = await response.json();
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
    },
  },
}
</script>

<template>
  <div class="search">
    <h1>Libraries matching "{{ query }}"</h1>

    <div v-if="fetchedLibraries !== null">
      <p>Results found: {{ fetchedLibraries.length }}</p>
        <ol class="usa-list">
        <li v-for="system in fetchedLibraries" :key=system>
          <RouterLink class="usa-link" :to="{ path: '/library/' + formatFSCSandSequence(system.fscskey, system.fscs_seq) + '/' }">
              {{ formatFSCSandSequence(system.fscskey, system.fscs_seq) }} - {{ system.libname }}
            </RouterLink>
        </li>
      </ol>
    </div>
    <div v-else>
      <p>Sorry, no matching libraries found. </p>
    </div>

  </div>
</template>

<style>
</style>
