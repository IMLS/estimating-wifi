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
      fetchError: null,
      isLoading: false,
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
          const parsedResponse = await response.json();

          // the API currently returns null instead of an empty array on no matches
          this.fetchedLibraries = await parsedResponse ? parsedResponse : []
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
  metaInfo(query = this.query) {
    const pagePrefix = `Library search results for "${query}"`;
    return {
      title: pagePrefix
    }
  }
}
</script>

<template>
  <div class="search">
    <h1>Libraries matching "{{ query }}"</h1>

    <div v-if="fetchedLibraries == null || fetchedLibraries.length < 1">
      <p>Sorry, no matching libraries found. </p>
      <span v-if="fetchError">Oops! Error encountered: {{ fetchError }} </span>
    </div>
    <div v-else>
      <p>Results found: {{ fetchedLibraries.length }}</p>
        <ol class="usa-list">
        <li v-for="system in fetchedLibraries" :key=system>
          <RouterLink class="usa-link" :to="{ path: '/library/' + formatFSCSandSequence(system.fscskey, system.fscs_seq) + '/' }">
              {{ formatFSCSandSequence(system.fscskey, system.fscs_seq) }} - {{ system.libname }}
            </RouterLink>
        </li>
      </ol>
    </div>


  </div>
</template>

<style>
</style>
